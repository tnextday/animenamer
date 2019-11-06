package kodi

import "encoding/xml"

type EpisodeDetails struct {
	XMLName         xml.Name    `xml:"episodedetails"`
	Title           string      `xml:"title"`
	OriginalTitle   string      `xml:"originaltitle,omitempty"`
	ShowTitle       string      `xml:"showtitle,omitempty"`
	Ratings         []*Rating   `xml:"ratings>rating,omitempty"`
	UserRating      int         `xml:"userrating,omitempty"`
	Top250          int         `xml:"top250,omitempty"`  // Not used
	Season          int         `xml:"season,omitempty"`  // Ignored on Import. Season is read from filename
	Episode         int         `xml:"episode,omitempty"` // Ignored on Import. Episode is read from filename
	DisplayEpisode  int         `xml:"displayepisode,omitempty"`
	DisplaySeason   int         `xml:"displayseason,omitempty"`
	Outline         string      `xml:"outline,omitempty"` // Not used
	Plot            string      `xml:"plot,omitempty"`
	TagLine         string      `xml:"tagline,omitempty"`
	RunTime         int         `xml:"runtime,omitempty"`
	Thumbs          []string    `xml:"thumb,omitempty"`
	Mpaa            string      `xml:"mpaa,omitempty"`
	PlayCount       int         `xml:"playcount,omitempty"`
	LastPlayed      string      `xml:"lastplayed,omitempty"`
	ID              int         `xml:"id,omitempty"` // Do not use as this is a Kodi generated tag.
	UniqueIDs       []*UniqueID `xml:"uniqueid"`
	Genres          []string    `xml:"genre,omitempty"`
	Credits         []string    `xml:"credits,omitempty"`
	Directors       []string    `xml:"director,omitempty"`
	Tags            []string    `xml:"tags,omitempty"`
	Premiered       string      `xml:"premiered"` // Do not use. Use <premiered> instead
	Year            string      `xml:"year,omitempty"`
	Status          string      `xml:"status,omitempty"`
	Code            string      `xml:"code,omitempty"`  // Not used
	Aired           string      `xml:"aired,omitempty"` // Not used
	Studios         []string    `xml:"studio,omitempty"`
	Trailer         string      `xml:"trailer,omitempty"`
	EpisodeBookmark string      `xml:"episodebookmark>position,omitempty"`
	Actors          []*Actor    `xml:"actor,omitempty"`
	Resume          *Resume     `xml:"resume,omitempty"`
	DateAdded       string      `xml:"dateadded,omitempty"`
}
