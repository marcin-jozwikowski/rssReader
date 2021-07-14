package reader

import "sort"

type Publishing struct {
	Name     string
	Episodes []*Episode
}

type Episode struct {
	Title      string
	Releases   []*Release
	Publishing *Publishing
}

type ByEpisodeTitle []*Episode

func (a ByEpisodeTitle) Len() int           { return len(a) }
func (a ByEpisodeTitle) Less(i, j int) bool { return a[i].Title > a[j].Title } //desc
func (a ByEpisodeTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type Release struct {
	Size           int
	Url            string
	Subtitle       string
	InternalResult string
	Episode        *Episode
}

type ByReleaseSize []*Release

func (a ByReleaseSize) Len() int           { return len(a) }
func (a ByReleaseSize) Less(i, j int) bool { return a[i].Size < a[j].Size }
func (a ByReleaseSize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (s *Publishing) AddRelease(title string, subtitle string, size int, url string) {
	episode := s.getEpisodeByTitle(title)
	if episode == nil {
		episode = s.addEpisode(title)
	}
	episode.addRelease(size, url, subtitle)
}

func (s *Publishing) getEpisodeByTitle(title string) *Episode {
	for e := range s.Episodes {
		if s.Episodes[e].Title == title {
			return s.Episodes[e]
		}
	}
	return nil
}

func (s *Publishing) addEpisode(title string) *Episode {
	episode := Episode{Title: title, Publishing: s}
	s.Episodes = append(s.Episodes, &episode)

	return &episode
}

func (s *Publishing) getEpisodeByAt(id int) *Episode {
	return s.Episodes[id]
}

func (s *Publishing) Sort() {
	sort.Sort(ByEpisodeTitle(s.Episodes))
}

func (e *Episode) addRelease(size int, url string, subtitle string) {
	e.Releases = append(e.Releases, &Release{Url: url, Size: size, Subtitle: subtitle, Episode: e})
	sort.Sort(ByReleaseSize(e.Releases))
}

func (e *Episode) getReleaseAt(id int) *Release {
	return e.Releases[id]
}
