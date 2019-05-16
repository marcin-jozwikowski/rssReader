package configuration

import (
	"cli"
	"fmt"
	"strconv"
	"strings"
)

var configKeys []string

func (config Config) Edit() {
	canRun := true
	for {
		canRun = config.keysEditAction()
		if !canRun {
			break
		}
	}
}

func (config Config) keysEditAction() bool {
	cli.ClearConsole()
	config.printKeys()
	fmt.Println("C: Create new")
	fmt.Println("X: Exit")
	readKey := strings.ToLower(cli.ReadString(""))
	switch readKey {
	case "x":
		return false
		break

	case "c":
		config.createNewKeyAction()
		break

	default:
		keyId, _ := strconv.Atoi(readKey)
		if keyId > 0 && keyId <= len(configKeys) {
			config.editKeyValuesAction(keyId)
		}
	}
	return true
}

func (config Config) printKeys() {
	for id, key := range config.parseKeys() {
		fmt.Println(strconv.Itoa(id+1) + ": " + key)
	}
}

func (config Config) parseKeys() []string {
	if configKeys == nil {
		id := 0
		configKeys = make([]string, len(config))
		for key, _ := range config {
			configKeys[id] = key
			id++
		}
	}

	return configKeys
}

func resetKeys() {
	configKeys = nil
}

func (config Config) createNewKeyAction() {
	fmt.Println("Create new")
	r := cli.ReadString("Name new key:")
	config[r] = []string{}
	resetKeys()
}

func (config Config) editKeyValuesAction(keyId int) {
	for {
		cli.ClearConsole()
		fmt.Println("Edit Key " + configKeys[keyId-1])
		for key, value := range config[configKeys[keyId-1]] {
			fmt.Println(strconv.Itoa(key) + ": " + value)
		}
		fmt.Println("D: Delete whole key")
		fmt.Println("A: Add value")
		fmt.Println("X: Go up")
		r := strings.ToLower(cli.ReadString("Name a value to remove"))
		switch r {
		case "x":
			return

		case "d":
			delete(config, configKeys[keyId-1])
			resetKeys()
			return

		case "a":
			cli.ClearConsole()
			newValue := cli.ReadString("Name new value for " + configKeys[keyId-1])
			config[configKeys[keyId-1]] = append(config[configKeys[keyId-1]], newValue)
			break

		default:
			key, _ := strconv.Atoi(r)
			if len(config[configKeys[keyId-1]]) > key {
				config[configKeys[keyId-1]] = append(config[configKeys[keyId-1]][:key], config[configKeys[keyId-1]][key+1:]...)
			}
		}
	}

}
