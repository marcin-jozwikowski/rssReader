package feed

import (
	"cli"
	"fmt"
	"gopkg.in/headzoo/surf.v1"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

const readerTypeNative = "native"
const readerTypePhantomjs = "phantomJS"
const readerTypeSurf = "surf"

type RssReader interface {
	GetXML(string) ([]byte, error)
}

type NativeRssReader struct {
}

func (NativeRssReader) GetXML(url string) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Running built-in downloader")
	}
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

func (RssReaderPhantomJS) GetXML(url string) ([]byte, error) {
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

type RssReaderSurf struct {
}

func (RssReaderSurf) GetXML(url string) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Running SURF downloader")
	}
	bow := surf.NewBrowser()
	err := bow.Open(url)
	if err != nil {
		panic(err)
	}

	bod := bow.Body()

	return []byte(bod), nil
}

func GetRssReader(externalCommand string) RssReader {
	var reader RssReader

	switch externalCommand {
	case readerTypeSurf:
		reader = RssReaderSurf{}
		break

	case readerTypePhantomjs:
		reader = RssReaderPhantomJS{}
		break

	case readerTypeNative:
	default:
		reader = NativeRssReader{}
		break
	}

	return reader
}
