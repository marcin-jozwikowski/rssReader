package reader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type RuntimeConfig struct {
	Sources []DataSource
}

type DataSource struct {
	Name            string
	Url             string
	XPath           string
	RegexExtract    string
	GroupField      string
	InternalXPath   string
	InternalRegex   string
	InternalBaseUrl string
	show            *Show
}

func (s *DataSource) AddResultingShow(show *Show) {
	s.show = show
}

func (s *DataSource) GetResultingShow() *Show {
	return s.show
}

func ReadRuntimeConfigFromFile(filename string) (RuntimeConfig, error) {
	if _, err := os.Stat(filename); err == nil {
		if config, readErr := fromFile(filename); nil == readErr {
			return config, nil
		} else {
			return RuntimeConfig{}, readErr
		}
	} else if os.IsNotExist(err) {
		newRuntimeConfig := RuntimeConfig{}
		if fileWriteErr := newRuntimeConfig.WriteToFile(filename); fileWriteErr != nil {
			return RuntimeConfig{}, fileWriteErr
		}
		return RuntimeConfig{}, fmt.Errorf("config file created: %v", filename)
	} else {
		return RuntimeConfig{}, err
	}
}

func fromFile(filename string) (RuntimeConfig, error) {
	if fileContent, fileReadErr := ioutil.ReadFile(filename); fileReadErr != nil {
		return RuntimeConfig{}, fmt.Errorf("error while opening file %v: %v", filename, fileReadErr.Error())
	} else {
		var raw RuntimeConfig
		if jsonErr := json.Unmarshal(fileContent, &raw); jsonErr != nil {
			return RuntimeConfig{}, fmt.Errorf("error while parsing file %v: %v", filename, jsonErr.Error())
		}
		return raw, nil
	}
}

func (configuration *RuntimeConfig) SourcesAsList() []string {
	var result []string
	for _, fs := range configuration.Sources {
		result = append(result, fs.Name)
	}

	return result
}

func (configuration *RuntimeConfig) WriteToFile(filename string) error {
	if file, err := json.MarshalIndent(configuration, "", " "); err != nil {
		return fmt.Errorf("error while encoding RuntimeConfig: %v", err.Error())
	} else {
		if err2 := ioutil.WriteFile(filename, file, 0644); err2 != nil {
			return fmt.Errorf("error while writing file: %v", err2.Error())
		}
		return nil
	}
}

func (configuration *RuntimeConfig) GetSources() *[]DataSource {
	return &configuration.Sources
}

func (configuration *RuntimeConfig) GetSourceAt(feedID int) *DataSource {
	return &configuration.Sources[feedID]
}

func (configuration *RuntimeConfig) DeleteSourceAt(keyId int) {
	configuration.Sources = append(configuration.Sources[:keyId], configuration.Sources[keyId+1:]...)
}

func (configuration *RuntimeConfig) AddSource(url string) {
	configuration.Sources = append(configuration.Sources, DataSource{Url: url})
}

func (configuration *RuntimeConfig) CountSources() int {
	return len(configuration.Sources)
}
