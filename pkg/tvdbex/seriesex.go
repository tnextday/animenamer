package tvdbex

import (
	"strconv"

	"github.com/tnextday/animenamer/pkg/kodi"

	"github.com/pioz/tvdb"
)

type SeriesEx struct {
	tvdb.Series
	OriginalSeriesName string
	EpisodesIndex      map[string]*EpisodeEx
	Images             map[string][]tvdb.Image
	NamedSeasons       []*kodi.NamedSeason
}

type EpisodeEx struct {
	tvdb.Episode
	OriginalEpisodeName string
}

func NewSeriesEx(s tvdb.Series, custom *CustomSeries) *SeriesEx {
	seriesEx := &SeriesEx{
		Series:             s,
		OriginalSeriesName: s.SeriesName,
		EpisodesIndex:      make(map[string]*EpisodeEx),
		Images:             map[string][]tvdb.Image{},
	}

	for i, episode := range s.Episodes {
		seId := SeasonEpisodeNumberIndex(episode.AiredSeason, episode.AiredEpisodeNumber)
		epEx := &EpisodeEx{
			Episode:             s.Episodes[i],
			OriginalEpisodeName: s.Episodes[i].EpisodeName,
		}
		seriesEx.EpisodesIndex[seId] = epEx
		if episode.AiredSeason != 0 {
			aId := AbsoluteNumberIndex(episode.AbsoluteNumber)
			seriesEx.EpisodesIndex[aId] = epEx
		}
	}
	if custom != nil && (custom.SeriesID != "" && custom.SeriesID == seriesEx.SeriesID) {
		if custom.SeriesName != "" {
			seriesEx.SeriesName = custom.SeriesName
		}
		if custom.Overview != "" {
			seriesEx.Overview = custom.Overview
		}
		seriesEx.NamedSeasons = custom.NamedSeason
		for _, ep := range custom.Episodes {
			var epex *EpisodeEx
			if ep.AiredSeason != 0 && ep.AiredEpisodeNumber != 0 {
				epex = seriesEx.GetEpisodeBySeasonEpisodeNumber(ep.AiredSeason, ep.AiredEpisodeNumber)
			} else if ep.AbsoluteNumber != 0 {
				epex = seriesEx.GetEpisodeByAbsoluteNumber(ep.AbsoluteNumber)
			}
			if epex != nil {
				if ep.EpisodeName != "" {
					epex.EpisodeName = ep.EpisodeName
				}
				if ep.Overview != "" {
					epex.Overview = ep.Overview
				}
			}
		}
	}
	return seriesEx
}

func (s *SeriesEx) GetEpisodeBySeasonEpisodeNumber(season, number int) *EpisodeEx {
	id := SeasonEpisodeNumberIndex(season, number)
	return s.EpisodesIndex[id]
}

func (s *SeriesEx) GetEpisodeByAbsoluteNumber(absoluteNumber int) *EpisodeEx {
	id := AbsoluteNumberIndex(absoluteNumber)
	return s.EpisodesIndex[id]
}

func (s *SeriesEx) GetEpisodeByIndex(index string) *EpisodeEx {
	return s.EpisodesIndex[index]
}

func AbsoluteNumberIndex(absoluteNumber int) string {
	return "a" + strconv.Itoa(absoluteNumber)
}

func SeasonEpisodeNumberIndex(season, number int) string {
	return "s" + strconv.Itoa(season) + "e" + strconv.Itoa(number)
}

//func (s *SeriesEx) GetImages(keyType string) []tvdb.Image {
//
//}
