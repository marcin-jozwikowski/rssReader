package feed

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/headzoo/surf/errors"
	"github.com/headzoo/surf/jar"
	"github.com/laplaceon/cfbypass"
	"gopkg.in/headzoo/surf.v1"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/exec"
	"rssReader/src/cli"
	"strings"
)

const readerTypeNative = "native"
const readerTypeSurf = "surf"
const readerTypeWget = "wget"
const readerTypeCustom = "custom"

type URLReader interface {
	GetContent(string, bool) ([]byte, error)
}

type NativeURLReader struct {
}

func (NativeURLReader) GetContent(url string, protected bool) ([]byte, error) {
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

type URLReaderCustom struct {
}

func (URLReaderCustom) GetContent(url string, protected bool) ([]byte, error) {
	params := *cli.DownloaderParams
	if strings.Contains(params, "%s") {
		params = fmt.Sprintf(params, url)
	} else {
		params += " " + url
	}
	command := strings.Fields(params)[0]
	params = strings.TrimPrefix(params, command+" ")

	if cli.IsVerboseDebug() {
		fmt.Println(fmt.Sprintf("Running %s with params: %s", command, params))
	}
	return exec.Command(command, params).Output()
}

type URLReaderWget struct {
}

func (URLReaderWget) GetContent(url string, protected bool) ([]byte, error) {
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

type URLReaderSurf struct {
}

func (URLReaderSurf) GetContent(url string, protected bool) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Running SURF downloader")
	}

	bow := surf.NewBrowser()
	if protected {
		cookies := getCookieJar(url)
		bow.SetCookieJar(cookies)
	}

	err := bow.Open(url)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	_, _ = bow.Download(foo)

	return b.Bytes(), nil
}

func getCookieJar(urlAddress string) *cookiejar.Jar {
	urlData, e := url.Parse(urlAddress)
	if e != nil {
		return jar.NewMemoryCookies()
	}

	cookies := cfbypass.GetTokens(urlAddress, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36", "4")
	cookieJar := jar.NewMemoryCookies()
	cookieJar.SetCookies(urlData, cookies)

	return cookieJar
}

func GetURLReader() URLReader {
	var reader URLReader

	switch *cli.Downloader {
	default:
	case readerTypeSurf:
		reader = URLReaderSurf{}
		break

	case readerTypeNative:
		reader = NativeURLReader{}
		break

	case readerTypeWget:
		reader = URLReaderWget{}
		break

	case readerTypeCustom:
		reader = URLReaderCustom{}
		break
	}

	return reader
}
