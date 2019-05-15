package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Settings []Setting `json:"settings"`
}

type Setting struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type Raw map[string][]string

func ReadFromFile(filename string) (Config, error) {
	if _, err := os.Stat(filename); err == nil {
		raw, readErr := fromFile(filename)
		if nil != readErr {
			return Config{}, readErr
		}
		return raw.toConfig(), nil
	} else if os.IsNotExist(err) {
		initializeInFile(filename)
		return Config{}, fmt.Errorf("Config file created: %v", filename)
	} else {
		return Config{}, err
	}
}

func fromFile(filename string) (Raw, error) {
	var raw Raw
	fileContent, fileReadErr := ioutil.ReadFile(filename)
	if fileReadErr != nil {
		return Raw{}, fmt.Errorf("Error while opening file %v: %v", filename, fileReadErr.Error())
	} else {
		var jsonErr = json.Unmarshal(fileContent, &raw)
		if jsonErr != nil {
			return Raw{}, fmt.Errorf("Error while parsing file %v: %v", filename, jsonErr.Error())
		}
	}
	return raw, nil
}

func (raw Raw) toConfig() Config {
	var config Config
	for k, v := range raw {
		setting := Setting{k,v}
		config.Settings = append(config.Settings, setting)
	}

	return config
}

func (config Config) WriteToFile(filename string) error {
	raw := config.toRaw()
	file, err := json.MarshalIndent(raw, "", " ")
	if err != nil {
		return fmt.Errorf("Error while encoding Config: %v", err.Error())
	}

	err2 := ioutil.WriteFile(filename, file, 0644)
	if err2 != nil {
		return fmt.Errorf("Error while writing file: %v", err2.Error())
	}

	return nil
}

func (config Config) toRaw() Raw {
	raw := make(Raw, len(config.Settings))
	for sID := range config.Settings {
		setting := config.Settings[sID]
		if len(setting.Key) > 0 {
			raw[setting.Key] = setting.Values
		}
	}
	return raw
}

func initializeInFile(filename string) {
	var conf Config
	conf.Settings = []Setting{{Key: "tv-shows", Values: []string{"ALF", "That 70's Show"}}}
	_ = conf.WriteToFile(filename)
}
