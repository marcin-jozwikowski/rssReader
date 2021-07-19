package reader

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
)

func Run(config *RuntimeConfig) []*Publishing {
	resultingPublishes := make([]*Publishing, len(config.Sources))
	for source := range config.Sources {
		src := config.Sources[source]
		publishing := runForDataSource(src)
		src.AddResultingPublishing(publishing)
		resultingPublishes = append(resultingPublishes, publishing)
	}
	return resultingPublishes
}

func runForDataSource(data *DataSource) *Publishing {
	data.SetRunning(true)
	publishing := Publishing{Name: data.Name}

	html, err := GetHtmlContent(data.Url)
	if err != nil {
		panic(err)
	}

	if readerPage, documentError := goquery.NewDocumentFromReader(bytes.NewReader(html)); documentError == nil {
		r := regexp.MustCompile(data.RegexExtract)
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
			publishing.AddRelease(list["Date"], list["Title"], int(size), href)
		})
		publishing.Sort()
	}
	data.SetRunning(false)

	return &publishing
}

func (s *DataSource) RunForRelease(release *Release) {
	defaultResult := " .... "
	dataUrl := s.InternalBaseUrl + release.Url
	release.InternalResult = defaultResult
	html, err := GetHtmlContent(dataUrl)
	if err != nil {
		release.InternalResult = "HTTP error occurred"
		return
	}

	readerPage, documentError := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if documentError != nil || readerPage == nil {
		release.InternalResult = "HTML error occurred"
		return
	}

	r := regexp.MustCompile(s.InternalRegex)
	readerPage.Find(s.InternalXPath).Each(func(id int, selection *goquery.Selection) {
		if r.MatchString(selection.Text()) {
			release.InternalResult = r.FindString(selection.Text())
			return
		}
	})

	if release.InternalResult == defaultResult {
		release.InternalResult = dataUrl
	}
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
