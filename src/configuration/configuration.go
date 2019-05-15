package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config map[string][]string

func ReadFromFile(filename string, defaultConfig Config) (Config, error) {
	if _, err := os.Stat(filename); err == nil {
		config, readErr := fromFile(filename)
		if nil != readErr {
			return Config{}, readErr
		}
		return config, nil
	} else if os.IsNotExist(err) {
		fileWriteErr := defaultConfig.WriteToFile(filename)
		if fileWriteErr != nil {
			return Config{}, fileWriteErr
		}
		return Config{}, fmt.Errorf("Config file created: %v", filename)
	} else {
		return Config{}, err
	}
}

func fromFile(filename string) (Config, error) {
	var raw Config
	fileContent, fileReadErr := ioutil.ReadFile(filename)
	if fileReadErr != nil {
		return Config{}, fmt.Errorf("Error while opening file %v: %v", filename, fileReadErr.Error())
	} else {
		var jsonErr = json.Unmarshal(fileContent, &raw)
		if jsonErr != nil {
			return Config{}, fmt.Errorf("Error while parsing file %v: %v", filename, jsonErr.Error())
		}
	}
	return raw, nil
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
