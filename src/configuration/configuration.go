package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Feeds []FeedSource
}

type FeedSource struct {
	Url           string
	SearchPhrases []string
	MaxChecked    int
}

func ReadConfigFromFile(filename string) (Config, error) {
	if _, err := os.Stat(filename); err == nil {
		config, readErr := fromFile(filename)
		if nil != readErr {
			return Config{}, readErr
		}
		return config, nil
	} else if os.IsNotExist(err) {
		newConfig := Config{}
		fileWriteErr := newConfig.WriteToFile(filename)
		if fileWriteErr != nil {
			return Config{}, fileWriteErr
		}
		return Config{}, fmt.Errorf("config file created: %v", filename)
	} else {
		return Config{}, err
	}
}

func fromFile(filename string) (Config, error) {
	var raw Config
	fileContent, fileReadErr := ioutil.ReadFile(filename)
	if fileReadErr != nil {
		return Config{}, fmt.Errorf("error while opening file %v: %v", filename, fileReadErr.Error())
	} else {
		var jsonErr = json.Unmarshal(fileContent, &raw)
		if jsonErr != nil {
			return Config{}, fmt.Errorf("error while parsing file %v: %v", filename, jsonErr.Error())
		}
	}
	return raw, nil
}

func (config *Config) WriteToFile(filename string) error {
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return fmt.Errorf("error while encoding Config: %v", err.Error())
	}

	err2 := ioutil.WriteFile(filename, file, 0644)
	if err2 != nil {
		return fmt.Errorf("error while writing file: %v", err2.Error())
	}

	return nil
}

func (config *Config) GetFeeds() *[]FeedSource {
	return &config.Feeds
}

func (config *Config) GetFeedAt(feedID int) *FeedSource {
	return &config.Feeds[feedID]
}

func (config *Config) DeleteFeetAt(keyId int) {
	config.Feeds = append(config.Feeds[:keyId], config.Feeds[keyId+1:]...)
}

func (config *Config) AddFeed(url string) {
	config.Feeds = append(config.Feeds, FeedSource{Url: url})
}

func (feedSource *FeedSource) ResetMaxChecked() {
	feedSource.MaxChecked = 0
}

func (feedSource *FeedSource) AddSearchPhrase(phrase string) {
	feedSource.SearchPhrases = append(feedSource.SearchPhrases, phrase)
}

func (feedSource *FeedSource) SetMaxChecked(newMax int) {
	feedSource.MaxChecked = newMax
}

func (feedSource *FeedSource) DeleteSearchPhraseAt(key int) {
	feedSource.SearchPhrases = append(feedSource.SearchPhrases[:key-1], feedSource.SearchPhrases[key:]...)
}
