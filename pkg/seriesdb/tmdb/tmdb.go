package tmdb

import (
	go_tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/tnextday/animenamer/pkg/seriesdb/series"
)

type TMDB struct {
	Client *go_tmdb.Client
}

func NewTMDB(apiKey string) (*TMDB, error) {
	return &TMDB{}, nil
}

func (t *TMDB) Search(q, language string) (id string, err error) {
	return "", nil
}

func (t *TMDB) GetSeries(seriesId, language string) (*series.SeriesDetail, error) {
	return nil, nil
}
