package tvdb

import (
	"fmt"
	"strconv"

	go_tvdb "github.com/pioz/tvdb"
	"github.com/tnextday/animenamer/pkg/seriesdb/series"
)

const (
	BaseEpisodesImageURL   = "https://api.thetvdb.com/banners/episodes"
	ImageTypeKeyFanArt     = "fanart"
	ImageTypeKeyPoster     = "poster"
	ImageTypeKeySeason     = "season"
	ImageTypeKeySeasonWide = "seasonwide"
	ImageTypeKeySeries     = "series"
)

type TVDB struct {
	Client go_tvdb.Client
}

type tvdbEpisodeContext struct {
	series  *go_tvdb.Series
	episode *go_tvdb.Episode
}

// type SeriesEx struct {
// 	go_tvdb.Series
// 	OriginalSeriesName string
// 	NamedSeasons       []*kodi.NamedSeason
// }

func NewTVDB(apiKey string) (*TVDB, error) {
	anime := TVDB{
		Client: go_tvdb.Client{Apikey: apiKey},
	}

	if err := anime.Client.Login(); err != nil {
		return nil, err
	}
	return &anime, nil
}

func (db *TVDB) Search(q, language string) (id string, err error) {
	var searchLanguages []string
	if language != "" && language != "en" {
		searchLanguages = append(searchLanguages, language)
	}
	searchLanguages = append(searchLanguages, "")
	for _, lang := range searchLanguages {
		db.Client.Language = lang
		if series, err := db.Client.BestSearch(q); err == nil {
			return strconv.Itoa(series.ID), nil
		} else {
			if go_tvdb.HaveCodeError(404, err) {
				continue
			}
			return "", err
		}
	}
	return "", series.ErrSearchNotFound
}

func (db *TVDB) GetSeries(seriesId, lang string, options map[string]string) (*series.SeriesDetail, error) {
	id, err := strconv.Atoi(seriesId)
	if err != nil {
		return nil, err
	}
	s := go_tvdb.Series{
		ID: id,
	}
	db.Client.Language = lang
	if err := db.Client.GetSeries(&s); err != nil {
		return nil, err
	}
	if err := db.Client.GetSeriesEpisodes(&s, nil); err != nil {
		return nil, err
	}
	// custom := db.Custom[seriesId]
	// if custom == nil {
	// 	custom = db.Custom[""]
	// }

	sd := &series.SeriesDetail{
		SeriesID:     seriesId,
		Name:         s.SeriesName,
		Overview:     s.Overview,
		OriginalName: s.SeriesName,
		SeasonNames:  make(map[int]string),
		Context:      s,
	}
	for i, ep := range s.Episodes {
		sep := &series.Episode{
			SeasonNumber:   ep.AiredSeason,
			EpisodeNumber:  ep.AiredEpisodeNumber,
			AbsoluteNumber: ep.AbsoluteNumber,
			Name:           ep.EpisodeName,
			Overview:       ep.Overview,
			OriginalName:   ep.EpisodeName,
			AiredDate:      ep.FirstAired,
			Context: &tvdbEpisodeContext{
				&s,
				&s.Episodes[i],
			},
		}
		sd.Episodes = append(sd.Episodes, sep)

	}
	return sd, nil
}

func (db *TVDB) GetSeriesActors(s *go_tvdb.Series, language string) (err error) {

	if len(s.Actors) > 0 {
		return nil
	}
	db.Client.Language = language
	err = db.Client.GetSeriesActors(s)
	if err != nil && !go_tvdb.HaveCodeError(404, err) {
		return err
	}
	return nil
}

func (db *TVDB) GetSeriesSummary(s *go_tvdb.Series, language string) (err error) {
	db.Client.Language = language
	err = db.Client.GetSeriesSummary(s)
	if err != nil && !go_tvdb.HaveCodeError(404, err) {
		return err
	}
	return nil
}

func (db *TVDB) GetSeriesImages(s *go_tvdb.Series, language string) (err error) {
	if len(s.Images) > 0 {
		return nil
	}
	db.Client.Language = language

	tmp := go_tvdb.Series{
		ID: s.ID,
	}
	allImageTypes := []string{
		ImageTypeKeyFanArt, ImageTypeKeyPoster, ImageTypeKeySeason, ImageTypeKeySeasonWide, ImageTypeKeySeries,
	}
	images := map[string][]go_tvdb.Image{}
	for _, k := range allImageTypes {
		tmp.Images = nil
		var err error = nil
		switch k {
		case ImageTypeKeyFanArt:
			err = db.Client.GetSeriesFanartImages(&tmp)
		case ImageTypeKeyPoster:
			err = db.Client.GetSeriesPosterImages(&tmp)
		case ImageTypeKeySeason:
			err = db.Client.GetSeriesSeasonImages(&tmp)
		case ImageTypeKeySeasonWide:
			err = db.Client.GetSeriesSeasonwideImages(&tmp)
		case ImageTypeKeySeries:
			err = db.Client.GetSeriesSeriesImages(&tmp)
		}
		if err == nil {
			images[k] = tmp.Images
		} else {
			if !go_tvdb.HaveCodeError(404, err) {
				return err
			}
		}
	}
	// s.Images = images
	return nil
}

func GetEpisodeImageUrl(episode *go_tvdb.Episode) string {
	if episode == nil || episode.SeriesID == 0 || episode.ID == 0 {
		return ""
	}
	return fmt.Sprintf("%s/%d/%d.jpg", BaseEpisodesImageURL, episode.SeriesID, episode.ID)
}
