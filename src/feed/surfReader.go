package feed

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/headzoo/surf/browser"
	"github.com/headzoo/surf/jar"
	"gopkg.in/headzoo/surf.v1"
	"net/http"
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
	return r.getContentBytes(pagedUrl, feed.CfCookie)
}

func (r *URLReaderSurf) GetContent(feed *FeedSource) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Running SURF downloader")
	}
	return r.getContentBytes(feed.Url, feed.CfCookie)
}

func (r *URLReaderSurf) getContentBytes(url string, cfCookie string) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Running SURF downloader for " + url)
	}
	if !r.initiated {
		r.surfBrowser = surf.NewBrowser()
		if cfCookie != "" {
			r.surfBrowser.SetCookieJar(r.getCookieJarForFeed(url, cfCookie))
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

func (r URLReaderSurf) getCookieJarForFeed(feedUrl string, cookieValue string) *cookiejar.Jar {
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

	a := strings.Split(urlData.Host, ".")
	if len(a) > 2 {
		// wildcard cookie if not main domain
		_, a = a[0], a[1:]
	}
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name: "__cfduid",
		Domain: "." + strings.Join(a, "."),
		Value: cookieValue,
	})
	cookieJar := jar.NewMemoryCookies()
	cookieJar.SetCookies(urlData, cookies)

	return cookieJar
}

func GetURLReader() *URLReaderSurf {
	return &URLReaderSurf{
		initiated: false,
	}
}
