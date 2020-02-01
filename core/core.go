package core

import (
	"sort"
)

type Pages []Page

type Page struct {
	Name           string
	WebPath        string
	Filepath       string
	Markdown       string
	HTML           string
	Headers        []Header
	IndexDocuments []IndexDocument
}

type Header struct {
	Title   string
	Link    string
	Headers []Header
}

type IndexDocument struct {
	ID      string
	Link    string
	Context []string
	Content []string
	HTML    string
}

func (ps Pages) SortByName(serviceName string) Pages {
	bn := byName{serviceName: serviceName, pages: ps}
	sort.Sort(bn)
	return ps
}

type byName struct {
	serviceName string
	pages       Pages
}

func (s byName) Len() int      { return len(s.pages) }
func (s byName) Swap(i, j int) { s.pages[i], s.pages[j] = s.pages[j], s.pages[i] }
func (s byName) Less(i, j int) bool {
	if s.pages[i].Name == s.serviceName {
		return true
	}
	if s.pages[j].Name == s.serviceName {
		return false
	}
	return s.pages[i].Name < s.pages[j].Name
}

type Files []File

type File struct {
	Name        string
	Href        string
	Path        string
	ContentType string
	Content     string
}
