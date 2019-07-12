package main

import (
	"fmt"
	"rssReader/src/cli"
	"rssReader/src/configuration"
	"rssReader/src/feed"
)

func main() {
	config, configErr := configuration.ReadFromFile(*cli.ConfigFileName)
	if configErr != nil {
		if cli.IsVerbose() {
			fmt.Println(configErr.Error())
		}
		*cli.RunEditor = true // enforce config editor
	}

	if *cli.RunEditor {
		config.Edit()
		_ = config.WriteToFile(*cli.ConfigFileName)
		return
	}

	config = feed.Read(config)
	_ = config.WriteToFile(*cli.ConfigFileName)
}
