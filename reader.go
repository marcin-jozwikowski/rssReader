package main

import (
	"fmt"
	"rssReader/src/cli"
	"rssReader/src/feed"
)

func main() {
	config, configErr := feed.ReadConfigFromFile(*cli.ConfigFileName)
	if configErr != nil {
		if cli.IsVerbose() {
			fmt.Println(configErr.Error())
		}
		*cli.RunEditor = true // enforce config editor
	}

	if *cli.ResetChecked {
		config.ResetCheckedCounters()
	}

	if *cli.RunEditor {
		if config.Edit() {
			_ = config.WriteToFile(*cli.ConfigFileName)
		}
		return
	}

	feed.Read(&config)
	_ = config.WriteToFile(*cli.ConfigFileName)
}
