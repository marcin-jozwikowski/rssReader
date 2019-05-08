package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Settings []Setting `json:"settings"`
	Cached []Setting `json:"cached"`
}

type Setting struct {
	Key string `json:"key"`
	Values []string `json:"values"`
}

func ReadFromFile(filename string) (Config, error) {
	var config Config
	var fileContent, err = ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("Error while opening file %v: %v", filename, err.Error())
	} else {
		var jsonErr = json.Unmarshal(fileContent, &config)
		if jsonErr != nil {
			return Config{}, fmt.Errorf("Error while parsing file %v: %v", filename, jsonErr.Error())
		}
	}

	return config, nil
}