package series

import (
	"errors"
)

type Provider interface {
	Search(query, language string) (seriesId string, err error)
	GetSeries(seriesId, language string) (series *SeriesDetail, err error)
}

var ErrSearchNotFound = errors.New("search not found")

var (
	ProviderTMDB string = "tmdb"
	ProviderTVDB string = "tvdb"
)
