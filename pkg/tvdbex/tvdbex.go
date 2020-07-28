package tvdbex

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pioz/tvdb"
)

const (
	BaseEpisodesImageURL   string = "https://api.thetvdb.com/banners/episodes"
	ImageTypeKeyFanArt            = "fanart"
	ImageTypeKeyPoster            = "poster"
	ImageTypeKeySeason            = "season"
	ImageTypeKeySeasonWide        = "seasonwide"
	ImageTypeKeySeries            = "series"
)

type TVDBEx struct {
	Client      tvdb.Client
	searchCache map[string]int
	seriesCache map[int]*SeriesEx
	Language    string
	Custom      map[string]*CustomSeries
}

func NewTVDBEx(apiKey string, language string, customs ...*CustomSeries) (*TVDBEx, error) {
	anime := TVDBEx{
		Client:      tvdb.Client{Apikey: apiKey},
		Language:    language,
		searchCache: make(map[string]int),
		seriesCache: make(map[int]*SeriesEx),
		Custom:      make(map[string]*CustomSeries),
	}
	for _, c := range customs {
		if c != nil {
			anime.Custom[c.SeriesID] = c
		}
	}

	if err := anime.Client.Login(); err != nil {
		return nil, err
	}
	return &anime, nil
}

func (db *TVDBEx) Search(q string) (id int, err error) {
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

func (db *TVDBEx) GetSeriesEx(id int) (*SeriesEx, error) {
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
	custom := db.Custom[strconv.Itoa(id)]
	if custom == nil {
		custom = db.Custom[""]
	}
	s := NewSeriesEx(series, custom)
	db.seriesCache[id] = s
	return s, nil
}

func (db *TVDBEx) GetSeriesActors(series *SeriesEx) (err error) {

	if len(series.Actors) > 0 {
		return nil
	}
	db.Client.Language = db.Language
	err = db.Client.GetSeriesActors(&series.Series)
	if err != nil && !tvdb.HaveCodeError(404, err) {
		return err
	}
	return nil
}

func (db *TVDBEx) GetSeriesSummary(series *SeriesEx) (err error) {
	db.Client.Language = db.Language
	err = db.Client.GetSeriesSummary(&series.Series)
	if err != nil && !tvdb.HaveCodeError(404, err) {
		return err
	}
	return nil
}

func (db *TVDBEx) GetSeriesImages(series *SeriesEx) (err error) {
	if len(series.Images) > 0 {
		return nil
	}
	db.Client.Language = db.Language

	tmp := tvdb.Series{
		ID: series.ID,
	}
	allImageTypes := []string{
		ImageTypeKeyFanArt, ImageTypeKeyPoster, ImageTypeKeySeason, ImageTypeKeySeasonWide, ImageTypeKeySeries,
	}
	images := map[string][]tvdb.Image{}
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
			if !tvdb.HaveCodeError(404, err) {
				return err
			}
		}
	}
	series.Images = images
	return nil
}

func GetEpisodeImageUrl(episode *EpisodeEx) string {
	if episode == nil || episode.SeriesID == 0 || episode.ID == 0 {
		return ""
	}
	return fmt.Sprintf("%s/%d/%d.jpg", BaseEpisodesImageURL, episode.SeriesID, episode.ID)
}
