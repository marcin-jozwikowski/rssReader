package cli

import "os"

func GetConfigFileLocation(fileName string) string {
	home, _ := os.UserHomeDir()
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fileName = home + "/.tvReader/" + fileName
	}

	return fileName
}
