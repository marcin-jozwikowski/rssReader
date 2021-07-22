package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"rssReader/src/cui"
)

type RuntimeConfig struct {
	Sources []*DataSource
}

func (configuration RuntimeConfig) GetListViewItems() *[]cui.ListViewItem {
	var r []cui.ListViewItem
	for _, c := range configuration.Sources {
		r = append(r, c)
	}
	return &r
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

func (configuration *RuntimeConfig) GetSources() []*DataSource {
	return configuration.Sources
}

func (configuration *RuntimeConfig) GetSourceAt(feedID int) *DataSource {
	return configuration.Sources[feedID]
}

func (configuration *RuntimeConfig) DeleteSourceAt(keyId int) {
	configuration.Sources = append(configuration.Sources[:keyId], configuration.Sources[keyId+1:]...)
}

func (configuration *RuntimeConfig) AddSource(url string) {
	configuration.Sources = append(configuration.Sources, &DataSource{Url: url})
}

func (configuration *RuntimeConfig) CountSources() int {
	return len(configuration.Sources)
}