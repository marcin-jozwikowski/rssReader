package configuration

import (
	"cli"
	"fmt"
	"strconv"
)

var configKeys []string

func (config Config) Edit() Config {
	canRun := true
	for {
		canRun = config.keysEditAction()
		if !canRun {
			break
		}
	}
	return config
}

func (config Config) keysEditAction() bool {
	cli.CallClear()
	config.printKeys()
	fmt.Println("C: Create new")
	fmt.Println("X: Exit")
	readKey := cli.ReadString("")
	switch readKey {
	case "X":
	case "x":
		return false
		break

	case "C":
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

func (config Config) createNewKeyAction() {
	fmt.Println("Create new")
	r := cli.ReadString("?")
	fmt.Println(r)
}


func (config Config) editKeyValuesAction(keyId int) {
	fmt.Println("Edit Key " + strconv.Itoa(keyId-1))
	r := cli.ReadString("?")
	fmt.Println(r)
}
