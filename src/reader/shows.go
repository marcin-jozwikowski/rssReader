package reader

type Show struct {
	Name     string
	Episodes []*Episode
}

type Episode struct {
	Title    string
	Releases []*Release
	Show     *Show
}

type Release struct {
	Size           int
	Url            string
	Subtitle       string
	InternalResult string
	Episode        *Episode
}

func (s *Show) AddRelease(title string, subtitle string, size int, url string) {
	episode := s.getEpisodeByTitle(title)
	if episode == nil {
		episode = s.addEpisode(title)
	}
	episode.addRelease(size, url, subtitle)
}

func (s *Show) getEpisodeByTitle(title string) *Episode {
	for e := range s.Episodes {
		if s.Episodes[e].Title == title {
			return s.Episodes[e]
		}
	}
	return nil
}

func (s *Show) addEpisode(title string) *Episode {
	episode := Episode{Title: title, Show: s}
	s.Episodes = append(s.Episodes, &episode)

	return &episode
}

func (s *Show) getEpisodeByAt(id int) *Episode {
	return s.Episodes[id]
}

func (e *Episode) addRelease(size int, url string, subtitle string) {
	e.Releases = append(e.Releases, &Release{Url: url, Size: size, Subtitle: subtitle, Episode: e})
}

func (e *Episode) getReleaseAt(id int) *Release {
	return e.Releases[id]
}
