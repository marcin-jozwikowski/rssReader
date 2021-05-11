package reader

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

func Run(config *RuntimeConfig) []*Show {
	resultingShows := make([]*Show, len(config.Sources))
	for source := range config.Sources {
		src := &config.Sources[source]
		show := runForDataSource(src)
		src.AddResultingShow(show)
		resultingShows = append(resultingShows, show)
	}
	return resultingShows
}

func runForDataSource(data *DataSource) *Show {
	show := Show{Name: data.Name}

	html, error := GetHtmlContent(data.Url)
	if error != nil {
		panic(error)
	}

	r := regexp.MustCompile(data.RegexExtract)

	if readerPage, documentError := goquery.NewDocumentFromReader(bytes.NewReader(html)); documentError == nil {
		readerPage.Find(data.XPath).Each(func(id int, selection *goquery.Selection) {
			matches := r.FindStringSubmatch(selection.Text())
			if len(matches) < 1 {
				return
			}
			list := extractMatches(matches, r.SubexpNames())
			href, _ := selection.Attr("href")
			size, _ := strconv.ParseFloat(list["Size"], 16)
			if list["SizeName"] == "GB" {
				size *= 1000
			}
			show.AddRelease(list["Date"], list["Title"], int(size), href)
		})
	}

	return &show
}

func extractMatches(matches []string, names []string) map[string]string {
	paramsMap := make(map[string]string)
	for i, _ := range names {
		if i > 0 && i <= len(matches) {
			paramsMap[names[i]] = matches[i]
		}
	}
	return paramsMap
}
