package main

import (
	"fmt"
	"rssReader/src/cli"
	"rssReader/src/reader"
)

func main() {
	config, configErr := reader.ReadRuntimeConfigFromFile(*cli.ConfigFileName)
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

	reader.RunCUI(&config)
	//reader.Run(&config)

	//feed.Read(&config)
	//_ = config.WriteToFile(*cli.ConfigFileName)
}
