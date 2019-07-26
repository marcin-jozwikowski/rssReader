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
		if config, readErr := fromFile(filename); nil == readErr {
			return config, nil
		} else {
			return Config{}, readErr
		}
	} else if os.IsNotExist(err) {
		newConfig := Config{}
		if fileWriteErr := newConfig.WriteToFile(filename); fileWriteErr != nil {
			return Config{}, fileWriteErr
		}
		return Config{}, fmt.Errorf("config file created: %v", filename)
	} else {
		return Config{}, err
	}
}

func fromFile(filename string) (Config, error) {
	if fileContent, fileReadErr := ioutil.ReadFile(filename); fileReadErr != nil {
		return Config{}, fmt.Errorf("error while opening file %v: %v", filename, fileReadErr.Error())
	} else {
		var raw Config
		if jsonErr := json.Unmarshal(fileContent, &raw); jsonErr != nil {
			return Config{}, fmt.Errorf("error while parsing file %v: %v", filename, jsonErr.Error())
		}
		return raw, nil
	}
}

func (config *Config) WriteToFile(filename string) error {
	if file, err := json.MarshalIndent(config, "", " "); err != nil {
		return fmt.Errorf("error while encoding Config: %v", err.Error())
	} else {
		if err2 := ioutil.WriteFile(filename, file, 0644); err2 != nil {
			return fmt.Errorf("error while writing file: %v", err2.Error())
		}
		return nil
	}
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

func (config *Config) ResetCheckedCounters() {
	for feedID := range config.Feeds {
		feed := &config.Feeds[feedID]
		feed.ResetMaxChecked()
	}
}

func (config *Config) CountFeeds() int {
	return len(config.Feeds)
}

func (feedSource *FeedSource) ResetMaxChecked() {
	feedSource.MaxChecked = 0
}

func (feedSource *FeedSource) AddSearchPhrase(phrase string) {
	feedSource.SearchPhrases = append(feedSource.SearchPhrases, phrase)
}

func (feedSource *FeedSource) SetMaxChecked(newMax int) {
	if newMax > 0 && newMax > feedSource.MaxChecked {
		feedSource.MaxChecked = newMax
	}
}

func (feedSource *FeedSource) DeleteSearchPhraseAt(key int) {
	feedSource.SearchPhrases = append(feedSource.SearchPhrases[:key-1], feedSource.SearchPhrases[key:]...)
}
