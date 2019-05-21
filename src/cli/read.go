package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var ioReader *bufio.Reader

func init() {
	ioReader = bufio.NewReader(os.Stdin)
}

func ReadString(question string) string {
	fmt.Println(question)
	readKey, _ := ioReader.ReadString('\n')
	readKey = strings.Replace(readKey, "\n", "", -1)
	return readKey
}
