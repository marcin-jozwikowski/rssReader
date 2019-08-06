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
		for id := range *config.GetFeeds() {
			fmt.Printf("  %v: %v", strconv.Itoa(id+1), config.GetFeedAt(id).Url)
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
				if !config.GetFeedAt(keyId - 1).editURLValuesAction() {
					config.DeleteFeetAt(keyId - 1)
				}
			}
		}
	}
}

func (config *Config) createNewURLAction() {
	fmt.Println("*** Create new ***")
	r := cli.ReadString("Name new URL:")
	config.AddFeed(r)
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
		if "" != feedSource.PostProcess {
			fmt.Println(fmt.Sprintf("  PostProcess: `%s`", feedSource.PostProcess))
			fmt.Println()
		}
		fmt.Println("  E: Edit URL itself")
		fmt.Println("  P: Edit PostProcess expression")
		fmt.Println("  D: Delete whole URL")
		fmt.Println("  A: Add value")
		fmt.Println("  R: Reset search counter")
		fmt.Println("  X: Go up")
		fmt.Println("     Name a value to remove")
		r := strings.ToLower(cli.ReadString(""))
		switch r {
		case "x":
			return true

		case "d":
			return false

		case "r":
			feedSource.ResetMaxChecked()
			break

		case "p":
			feedSource.PostProcess = cli.ReadString("New PostProcess expression (empty to disable)")
			break

		case "a":
			feedSource.AddSearchPhrase(cli.ReadString("*** Name new value for " + feedSource.Url))
			feedSource.ResetMaxChecked()
			break

		case "e":
			feedSource.Url = cli.ReadString("New URL")
			return true

		default:
			key, _ := strconv.Atoi(r)
			if key > 0 && len(feedSource.SearchPhrases) >= key {
				feedSource.DeleteSearchPhraseAt(key)
			}
		}
	}

}
