package main

import (
	"fmt"
	"rssReader/src/application"
	"rssReader/src/cli"
)

func main() {
	filename := *cli.ConfigFileName
	if *cli.OneTimeUrl != "" {
		filename = cli.GetConfigFileLocation("one-time.json")
	}

	config, configErr := application.ReadRuntimeConfigFromFile(filename)
	if configErr != nil {
		if cli.IsVerbose() {
			fmt.Println(configErr.Error())
		}
		*cli.RunEditor = true // enforce config editor
	}

	if *cli.OneTimeUrl != "" {
		for _, source := range config.GetSources() {
			source.Url = *cli.OneTimeUrl
		}
	}

	if *cli.RunEditor {
		// if config.Edit() {
		//	_ = config.WriteToFile(*cli.ConfigFileName)
		// }
		return
	}

	application.RunCUI(&config)
	//publishing.Run(&config)

	//feed.Read(&config)
	//_ = config.WriteToFile(*cli.ConfigFileName)
}
