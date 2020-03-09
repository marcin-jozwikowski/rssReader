package feed

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
	PostProcess   string
	CfCookie      string
	IsHTML        bool
	IsPaginated   bool
	downloader    *URLReaderSurf
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

func (configuration *Config) AsList() []string {
	var result []string
	for _, fs := range configuration.Feeds {
		result = append(result, fs.Url)
	}

	return result
}

func (configuration *Config) WriteToFile(filename string) error {
	if file, err := json.MarshalIndent(configuration, "", " "); err != nil {
		return fmt.Errorf("error while encoding Config: %v", err.Error())
	} else {
		if err2 := ioutil.WriteFile(filename, file, 0644); err2 != nil {
			return fmt.Errorf("error while writing file: %v", err2.Error())
		}
		return nil
	}
}

func (configuration *Config) GetFeeds() *[]FeedSource {
	return &configuration.Feeds
}

func (configuration *Config) GetFeedAt(feedID int) *FeedSource {
	return &configuration.Feeds[feedID]
}

func (configuration *Config) DeleteFeetAt(keyId int) {
	configuration.Feeds = append(configuration.Feeds[:keyId], configuration.Feeds[keyId+1:]...)
}

func (configuration *Config) AddFeed(url string) {
	configuration.Feeds = append(configuration.Feeds, FeedSource{Url: url})
}

func (configuration *Config) ResetCheckedCounters() {
	for feedID := range configuration.Feeds {
		feed := &configuration.Feeds[feedID]
		feed.ResetMaxChecked()
	}
}

func (configuration *Config) CountFeeds() int {
	return len(configuration.Feeds)
}

func (feedSource *FeedSource) ResetMaxChecked() {
	feedSource.MaxChecked = 0
}

func (feedSource *FeedSource) SetMaxChecked(newMax int) {
	if newMax > 0 && newMax > feedSource.MaxChecked {
		feedSource.MaxChecked = newMax
	}
}

func (feedSource *FeedSource) AddSearchPhrase(phrase string) {
	feedSource.SearchPhrases = append(feedSource.SearchPhrases, phrase)
}

func (feedSource *FeedSource) DeleteSearchPhraseAt(key int) {
	feedSource.SearchPhrases = append(feedSource.SearchPhrases[:key-1], feedSource.SearchPhrases[key:]...)
}
