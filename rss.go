package main

import "encoding/xml"

// Rss was generated 2024-06-08 14:07:33 by https://xml-to-go.github.io/ in Ukraine.
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Atom    string   `xml:"atom,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description    string `xml:"description"`
		Language       string `xml:"language"`
		PubDate        string `xml:"pubDate"`
		LastBuildDate  string `xml:"lastBuildDate"`
		Docs           string `xml:"docs"`
		Generator      string `xml:"generator"`
		ManagingEditor string `xml:"managingEditor"`
		WebMaster      string `xml:"webMaster"`
		Items          []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
	Enclosure   struct {
		Text   string `xml:",chardata"`
		URL    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
}
