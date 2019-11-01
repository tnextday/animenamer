package tvdbex

import (
	"errors"
	"strings"

	"github.com/pioz/tvdb"
)

type TVDB struct {
	Client      tvdb.Client
	searchCache map[string]int
	seriesCache map[int]*Series
	Language    string
}

func NewTVDB(apiKey string, language string) (*TVDB, error) {
	anime := TVDB{
		Client:      tvdb.Client{Apikey: apiKey},
		Language:    language,
		searchCache: make(map[string]int),
		seriesCache: make(map[int]*Series),
	}
	if err := anime.Client.Login(); err != nil {
		return nil, err
	}
	return &anime, nil
}

func (db *TVDB) Search(q string) (id int, err error) {
	q = strings.ToLower(strings.TrimSpace(q))
	if id, exist := db.searchCache[q]; exist {
		return id, nil
	}
	var searchLanguages []string
	if db.Language != "" && db.Language != "en" {
		searchLanguages = append(searchLanguages, db.Language)
	}
	searchLanguages = append(searchLanguages, "")
	for _, lang := range searchLanguages {
		db.Client.Language = lang
		if series, err := db.Client.BestSearch(q); err == nil {
			db.searchCache[q] = series.ID
			return series.ID, nil
		} else {
			if tvdb.HaveCodeError(404, err) {
				continue
			}
			return 0, err
		}
	}
	return 0, errors.New("not found")
}

func (db *TVDB) GetSeries(id int) (*Series, error) {
	if series, exists := db.seriesCache[id]; exists {
		return series, nil
	}
	db.Client.Language = db.Language
	series := tvdb.Series{
		ID: id,
	}
	if err := db.Client.GetSeries(&series); err != nil {
		return nil, err
	}
	if err := db.Client.GetSeriesEpisodes(&series, nil); err != nil {
		return nil, err
	}
	s := NewSeries(series)
	db.seriesCache[id] = s
	return s, nil
}
