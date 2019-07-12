package configuration

import (
	"fmt"
	"rssReader/src/cli"
	"strconv"
	"strings"
)

func (config *Config) Edit() bool {
	for {
		cli.ClearConsole()
		fmt.Println("*** Entries available ***")
		for id, key := range config.Feeds {
			fmt.Printf("  %v: %v", strconv.Itoa(id+1), key.Url)
			fmt.Println()
		}
		fmt.Println()
		fmt.Println("  C: Create new")
		fmt.Println("  S: Save and Exit")
		fmt.Println("  X: Discard changes and Exit")
		readKey := strings.ToLower(cli.ReadString(""))
		switch readKey {
		case "s":
			return true

		case "x":
			return false

		case "c":
			config.createNewURLAction()
			break

		default:
			keyId, _ := strconv.Atoi(readKey)
			if keyId > 0 && keyId <= len(config.Feeds) {
				if !config.Feeds[keyId-1].editURLValuesAction() {
					config.Feeds = append(config.Feeds[:keyId-1], config.Feeds[keyId:]...)
				}
			}
		}
	}
}

func (config *Config) createNewURLAction() {
	fmt.Println("*** Create new ***")
	r := cli.ReadString("Name new URL:")
	config.Feeds = append(config.Feeds, FeedSource{Url: r})
}

func (feedSource *FeedSource) editURLValuesAction() bool {
	for {
		cli.ClearConsole()
		fmt.Println("*** Edit URL " + feedSource.Url)
		for key, value := range feedSource.SearchPhrases {
			fmt.Printf("  %v: %v", strconv.Itoa(key+1), value)
			fmt.Println()
		}
		fmt.Println()
		fmt.Println("  E: Edit URL itself")
		fmt.Println("  D: Delete whole URL")
		fmt.Println("  A: Add value")
		fmt.Println("  X: Go up")
		fmt.Println("     Name a value to remove")
		r := strings.ToLower(cli.ReadString(""))
		switch r {
		case "x":
			return true

		case "d":
			return false

		case "a":
			newValue := cli.ReadString("*** Name new value for " + feedSource.Url)
			feedSource.SearchPhrases = append(feedSource.SearchPhrases, newValue)
			break

		case "e":
			newUrl := cli.ReadString("New URL")
			feedSource.Url = newUrl
			return true

		default:
			key, _ := strconv.Atoi(r)
			if key > 0 && len(feedSource.SearchPhrases) >= key {
				feedSource.SearchPhrases = append(feedSource.SearchPhrases[:key-1], feedSource.SearchPhrases[key:]...)
			}
		}
	}

}
