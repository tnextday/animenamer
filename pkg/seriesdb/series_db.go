package seriesdb

import (
	"fmt"
	"strings"

	"github.com/tnextday/animenamer/pkg/seriesdb/series"
	"github.com/tnextday/animenamer/pkg/seriesdb/tmdb"
	"github.com/tnextday/animenamer/pkg/seriesdb/tvdb"
)

type SeriesDB struct {
	provider    series.Provider
	searchCache map[string]string
	seriesCache map[string]*series.Series
	Language    string
	Custom      map[string]*series.CustomSeries
}

func NewSeriesDB(db, apiKey, language string, customs ...*series.CustomSeries) (*SeriesDB, error) {
	var (
		provider series.Provider
		err      error
	)
	switch db {
	case series.ProviderTVDB:
		provider, err = tvdb.NewTVDB(apiKey)
	case series.ProviderTMDB:
		provider, err = tmdb.NewTMDB(apiKey)
	default:
		err = fmt.Errorf("unknow provider type (%s)", db)
	}
	if err != nil {
		return nil, err
	}
	sdb := SeriesDB{
		provider:    provider,
		Language:    language,
		searchCache: make(map[string]string),
		seriesCache: make(map[string]*series.Series),
		Custom:      make(map[string]*series.CustomSeries),
	}
	for _, c := range customs {
		if c != nil {
			sdb.Custom[c.SeriesID] = c
		}
	}

	return &sdb, nil
}

func (db *SeriesDB) Search(q string) (seriesId string, err error) {
	q = strings.ToLower(strings.TrimSpace(q))
	if id, exist := db.searchCache[q]; exist {
		return id, nil
	}
	if seriesId, err := db.provider.Search(q, db.Language); err == nil {
		db.searchCache[q] = seriesId
		return seriesId, nil
	} else {
		return "", err
	}
}

func (db *SeriesDB) GetSeries(seriesId string) (*series.Series, error) {
	if series, exists := db.seriesCache[seriesId]; exists {
		return series, nil
	}
	seriesDetail, err := db.provider.GetSeries(seriesId, db.Language)
	if err != nil {
		return nil, err
	}
	custom := db.Custom[seriesId]
	if custom == nil {
		custom = db.Custom[""]
	}
	series := series.NewSeries(seriesDetail, custom)
	db.seriesCache[seriesId] = series
	return series, nil
}
