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
	fmt.Println("  C: Create new")
	fmt.Println("  X: Exit")
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
	fmt.Println("*** Keys available ***")
	for id, key := range config.parseKeys() {
		fmt.Printf("  %v: %v", strconv.Itoa(id+1), key)
		fmt.Println()
	}
	fmt.Println()
}

func (config Config) parseKeys() []string {
	if configKeys == nil {
		id := 0
		configKeys = make([]string, len(config))
		for key := range config {
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
	fmt.Println("*** Create new ***")
	r := cli.ReadString("Name new key:")
	config[r] = []string{}
	resetKeys()
}

func (config Config) editKeyValuesAction(keyId int) {
	mainKey := keyId - 1
	for {
		cli.ClearConsole()
		fmt.Println("*** Edit Key " + configKeys[mainKey])
		for key, value := range config[configKeys[mainKey]] {
			fmt.Printf("  %v: %v", strconv.Itoa(key), value)
			fmt.Println()
		}
		fmt.Println()
		fmt.Println("  D: Delete whole key")
		fmt.Println("  A: Add value")
		fmt.Println("  X: Go up")
		fmt.Println("     Name a value to remove")
		r := strings.ToLower(cli.ReadString(""))
		switch r {
		case "x":
			return

		case "d":
			delete(config, configKeys[mainKey])
			resetKeys()
			return

		case "a":
			newValue := cli.ReadString("*** Name new value for " + configKeys[mainKey])
			config[configKeys[mainKey]] = append(config[configKeys[mainKey]], newValue)
			break

		default:
			key, _ := strconv.Atoi(r)
			if len(config[configKeys[mainKey]]) > key {
				config[configKeys[mainKey]] = append(config[configKeys[mainKey]][:key], config[configKeys[mainKey]][key+1:]...)
			}
		}
	}

}
