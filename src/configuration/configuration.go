package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Settings []Setting `json:"settings"`
	Cached   []Setting `json:"cached"`
}

type Setting struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

func ReadFromFile(filename string) (Config, error) {
	var config Config
	if _, err := os.Stat(filename); err == nil {
		fileContent, fileReadErr := ioutil.ReadFile(filename)
		if fileReadErr != nil {
			return Config{}, fmt.Errorf("Error while opening file %v: %v", filename, fileReadErr.Error())
		} else {
			var jsonErr = json.Unmarshal(fileContent, &config)
			if jsonErr != nil {
				return Config{}, fmt.Errorf("Error while parsing file %v: %v", filename, jsonErr.Error())
			}
		}
	} else if os.IsNotExist(err) {
		initializeInFile(filename)
		return Config{}, fmt.Errorf("Config file created: %v", filename)
	} else {
		return Config{}, err
	}

	return config, nil
}

func (config Config) WriteToFile(filename string) error {
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return fmt.Errorf("Error while encoding Config: %v", err.Error())
	}

	err2 := ioutil.WriteFile(filename, file, 0644)
	if err2 != nil {
		return fmt.Errorf("Error while writing file: %v", err2.Error())
	}

	return nil
}

func initializeInFile(filename string) {
	var conf Config
	conf.Settings = []Setting{{Key: "tv-shows", Values: []string{"ALF", "That 70's Show"}}}
	conf.Cached = []Setting{}
	_ = conf.WriteToFile(filename)
}
