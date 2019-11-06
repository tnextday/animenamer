package kodi

import "encoding/xml"

// https://kodi.wiki/view/NFO_files/TV_shows

type TVShow struct {
	XMLName        xml.Name       `xml:"tvshow"`
	Title          string         `xml:"title"`
	OriginalTitle  string         `xml:"originaltitle,omitempty"`
	ShowTitle      string         `xml:"showtitle,omitempty"`
	SortTitle      string         `xml:"sorttitle,omitempty"`
	Ratings        []*Rating      `xml:"ratings>rating,omitempty"`
	UserRating     int            `xml:"userrating,omitempty"`
	Top250         int            `xml:"top250,omitempty"`
	Season         int            `xml:"season,omitempty"`
	Episode        int            `xml:"episode,omitempty"`
	DisplayEpisode int            `xml:"displayepisode,omitempty"` // Not used
	DisplaySeason  int            `xml:"displayseason,omitempty"`  // Not used
	Outline        string         `xml:"outline,omitempty"`        // Not used
	Plot           string         `xml:"plot,omitempty"`
	TagLine        string         `xml:"tagline,omitempty"`
	RunTime        int            `xml:"runtime,omitempty"`
	Thumbs         []*Thumb       `xml:"thumb,omitempty"`
	Fanarts        []*Thumb       `xml:"fanart>thumb,omitempty"`
	Mpaa           string         `xml:"mpaa,omitempty"`
	PlayCount      int            `xml:"playcount,omitempty"`
	LastPlayed     string         `xml:"lastplayed,omitempty"`
	EpisodeGuide   *Url           `xml:"episodeguide>url,omitempty"`
	ID             int            `xml:"id,omitempty"`
	UniqueIDs      []*UniqueID    `xml:"uniqueid"`
	Genres         []string       `xml:"genre,omitempty"`
	Tags           []string       `xml:"tags,omitempty"`
	Premiered      string         `xml:"premiered"` // Do not use. Use <premiered> instead
	Year           string         `xml:"year,omitempty"`
	Status         string         `xml:"status,omitempty"`
	Code           string         `xml:"code,omitempty"`  // Not used
	Aired          string         `xml:"aired,omitempty"` // Not used
	Studios        []string       `xml:"studio,omitempty"`
	Trailer        string         `xml:"trailer,omitempty"`
	Actors         []*Actor       `xml:"actor,omitempty"`
	NamedSeasons   []*NamedSeason `xml:"namedseason,omitempty"`
	Resume         *Resume        `xml:"resume,omitempty"` // Not used for TV Show
	DateAdded      string         `xml:"dateadded,omitempty"`
}

type Rating struct {
	XMLName xml.Name `xml:"rating"`
	Name    string   `xml:"name,attr"`
	Max     int      `xml:"max,attr"`
	Default bool     `xml:"default,attr"`
	Value   float32  `xml:"value"`
	Votes   int      `xml:"votes"`
}

type NamedSeason struct {
	XMLName xml.Name `xml:"namedseason"`
	Number  int      `xml:"number,attr"`
	Name    string   `xml:",innerxml"`
}

type Thumb struct {
	XMLName xml.Name `xml:"thumb"`
	Aspect  string   `xml:"aspect,attr"`
	Type    string   `xml:"type,attr,omitempty"`
	Season  int      `xml:"season,attr,omitempty"`
	Preview string   `xml:"preview,attr"`
	Uri     string   `xml:",innerxml"`
}

type Url struct {
	XMLName xml.Name `xml:"url"`
	Cache   string   `xml:"cache,attr"`
	Uri     string   `xml:",innerxml"`
}

type UniqueID struct {
	XMLName xml.Name `xml:"uniqueid"`
	Type    string   `xml:"type,attr"`
	Default bool     `xml:"default,attr"`
	ID      string   `xml:",innerxml"`
}

type Actor struct {
	XMLName xml.Name `xml:"actor"`
	Name    string   `xml:"name"`
	Role    string   `xml:"role"`
	Order   int      `xml:"order"`
	Thumb   string   `xml:"thumb"`
}

type Resume struct {
	XMLName  xml.Name `xml:"resume"`
	Position float32  `xml:"position"`
	Total    float32  `xml:"total"`
}
