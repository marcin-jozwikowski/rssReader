package feed

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/headzoo/surf/browser"
	"github.com/headzoo/surf/jar"
	"github.com/laplaceon/cfbypass"
	"gopkg.in/headzoo/surf.v1"
	"net/http/cookiejar"
	"net/url"
	"rssReader/src/cli"
	"strconv"
	"strings"
)

type URLReaderSurf struct {
	surfBrowser *browser.Browser
	cookieJar   *cookiejar.Jar
	initiated   bool
}

func (r *URLReaderSurf) GetContentPaginated(feed *FeedSource, page int) ([]byte, error) {
	pagedUrl := strings.Replace(feed.Url, "{page}", strconv.Itoa(page), -1)
	return r.getContentBytes(pagedUrl, feed.IsProtected)
}

func (r *URLReaderSurf) GetContent(feed *FeedSource) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Running SURF downloader")
	}
	return r.getContentBytes(feed.Url, feed.IsProtected)
}

func (r *URLReaderSurf) getContentBytes(url string, isProtected bool) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Running SURF downloader for " + url)
	}
	if !r.initiated {
		r.surfBrowser = surf.NewBrowser()
		if isProtected {
			r.surfBrowser.SetCookieJar(r.getCookieJarForFeed(url))
		}
		r.initiated = true
	}

	err := r.surfBrowser.Open(url)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	_, _ = r.surfBrowser.Download(foo)

	return b.Bytes(), nil
}

func (r URLReaderSurf) getCookieJarForFeed(feedUrl string) *cookiejar.Jar {
	urlData, e := url.Parse(feedUrl)
	if e != nil {
		if cli.IsVerboseDebug() {
			fmt.Println("URL error " + e.Error())
		}
		return jar.NewMemoryCookies()
	}
	if cli.IsVerboseDebug() {
		fmt.Println("Setting cookies for " + feedUrl)
	}

	cookies := cfbypass.GetTokens(feedUrl, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36", "4")
	cookieJar := jar.NewMemoryCookies()
	cookieJar.SetCookies(urlData, cookies)

	return cookieJar
}

func GetURLReader() *URLReaderSurf {
	return &URLReaderSurf{
		initiated: false,
	}
}
