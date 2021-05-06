package reader

type Show struct {
	Name     string
	Episodes []Episode
}

type Episode struct {
	Title    string
	Releases []Release
}

type Release struct {
	Size int
	Url  string
}

func (s *Show) AddRelease(title string, size int, url string) {
	episode := s.getEpisodeByTitle(title)
	if episode == nil {
		episode = s.addEpisode(title)
	}
	episode.addRelease(size, url)
}

func (s *Show) getEpisodeByTitle(title string) *Episode {
	for e := range s.Episodes {
		if s.Episodes[e].Title == title {
			return &s.Episodes[e]
		}
	}
	return nil
}

func (s *Show) addEpisode(title string) *Episode {
	episode := Episode{Title: title}
	s.Episodes = append(s.Episodes, episode)

	return &episode
}

func (e *Episode) addRelease(size int, url string) {
	e.Releases = append(e.Releases, Release{Url: url, Size: size})
}
