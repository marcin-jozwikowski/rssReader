package feed

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/headzoo/surf/errors"
	"gopkg.in/headzoo/surf.v1"
	"io/ioutil"
	"net/http"
	"os/exec"
	"rssReader/src/cli"
	"strings"
)

const readerTypeNative = "native"
const readerTypePhantomjs = "phantomJS"
const readerTypeSurf = "surf"
const readerTypeWget = "wget"

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
	params := *cli.DownloaderParams
	if strings.Contains(params, "%s") {
		params = fmt.Sprintf(params, url)
	}
	if cli.IsVerboseDebug() {
		fmt.Println("Running PhantomJS with params: " + params)
	}
	return exec.Command("phantomjs", params).Output()
}

type RssReaderWget struct {
}

func (RssReaderWget) GetXML(url string) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Running custom downloader")
	}
	cmd := exec.Command(*cli.Downloader, *cli.DownloaderParams, url)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	return out.Bytes(), nil
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
		return nil, err
	}

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	_, _ = bow.Download(foo)

	return b.Bytes(), nil
}

func GetRssReader(downloader string) RssReader {
	var reader RssReader

	switch downloader {
	case readerTypeSurf:
		reader = RssReaderSurf{}
		break

	case readerTypePhantomjs:
		reader = RssReaderPhantomJS{}
		break

	case readerTypeNative:
		reader = NativeRssReader{}
		break

	case readerTypeWget:
		reader = RssReaderWget{}
		break
	}

	return reader
}
