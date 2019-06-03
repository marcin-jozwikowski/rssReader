package feed

import (
	"cli"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

const readerTypeNative = "native"
const readerTypePhantomjs = "phantomJS"

type RssReader interface {
	GetXML(string) ([]byte, error)
}

type NativeRssReader struct {
}

func (_ NativeRssReader) GetXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err.Error())
	}

	return data, nil
}

type RssReaderPhantomJS struct {
}

func (extRR RssReaderPhantomJS) GetXML(url string) ([]byte, error) {
	params := *cli.DownloadParams
	if strings.Contains(params, "%s") {
		params = fmt.Sprintf(params, url)
	}
	if *cli.ProxyAddr != "" {
		params = fmt.Sprintf("--proxy=%s ", *cli.ProxyAddr) + params
		if *cli.ProxyType != "" {
			params = fmt.Sprintf("--proxy-type=%s ", *cli.ProxyType) + params
		}
		if *cli.ProxyAuth != "" {
			params = fmt.Sprintf("--proxy-auth=\"%s\" ", *cli.ProxyAuth) + params
		}
	}
	if cli.IsVerboseDebug() {
		fmt.Println("Running PhantomJS with params: " + params)
	}
	return exec.Command("phantomjs", params).Output()
}

func GetRssReader(externalCommand string) RssReader {
	var reader RssReader

	switch externalCommand {
	case readerTypePhantomjs:
		reader = RssReaderPhantomJS{}

	case readerTypeNative:
	default:
		reader = NativeRssReader{}

	}

	return reader
}
