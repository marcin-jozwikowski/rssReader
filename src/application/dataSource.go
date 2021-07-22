package application

import (
	"rssReader/src/cui"
	"rssReader/src/publishing"
	"strconv"
)

type DataSource struct {
	Name            string
	Url             string
	XPath           string
	RegexExtract    string
	GroupField      string
	InternalXPath   string
	InternalRegex   string
	InternalBaseUrl string
	publishing      *publishing.Publishing
	isRunning       bool
}

func (s DataSource) ToString() string {
	line := s.Name
	if s.GetResultingPublishing() == nil {
		if s.IsCurrentlyRunning() {
			line += " | ..."
		}
	} else {
		line += " | " + strconv.Itoa(len(s.publishing.Pieces)) + " Pieces"
	}
	return line
}

func (s *DataSource) GetListViewItems() *[]cui.ListViewItem {
	var r []cui.ListViewItem
	for _, c := range s.GetResultingPublishing().Pieces {
		r = append(r, c)
	}
	return &r
}

func (s *DataSource) AddResultingPublishing(publishing *publishing.Publishing) {
	s.publishing = publishing
}

func (s *DataSource) GetResultingPublishing() *publishing.Publishing {
	return s.publishing
}

func (s *DataSource) IsCurrentlyRunning() bool {
	return s.isRunning
}

func (s *DataSource) SetRunning(isRunning bool) {
	s.isRunning = isRunning
}
