package main

import (
	"fmt"
	"rssReader/src/application"
	"rssReader/src/cli"
)

func main() {
	config, configErr := application.ReadRuntimeConfigFromFile(*cli.ConfigFileName)
	if configErr != nil {
		if cli.IsVerbose() {
			fmt.Println(configErr.Error())
		}
		*cli.RunEditor = true // enforce config editor
	}

	if *cli.RunEditor {
		//if config.Edit() {
		//	_ = config.WriteToFile(*cli.ConfigFileName)
		//}
		return
	}

	application.RunCUI(&config)
	//publishing.Run(&config)

	//feed.Read(&config)
	//_ = config.WriteToFile(*cli.ConfigFileName)
}
