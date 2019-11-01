package tvdbex

import (
	"strconv"

	"github.com/pioz/tvdb"
)

type Series struct {
	tvdb.Series
	EpisodesIndex map[string]*tvdb.Episode
}

func NewSeries(s tvdb.Series) *Series {
	series := &Series{
		Series:        s,
		EpisodesIndex: make(map[string]*tvdb.Episode),
	}
	for i, episode := range s.Episodes {
		seId := SeasonEpisodeNumberIndex(episode.AiredSeason, episode.AiredEpisodeNumber)
		series.EpisodesIndex[seId] = &s.Episodes[i]
		if episode.AiredSeason != 0 {
			aId := AbsoluteNumberIndex(episode.AbsoluteNumber)
			series.EpisodesIndex[aId] = &s.Episodes[i]
		}
	}
	return series
}

func (s *Series) GetEpisodeBySeasonEpisodeNumber(season, number int) *tvdb.Episode {
	id := SeasonEpisodeNumberIndex(season, number)
	return s.EpisodesIndex[id]
}

func (s *Series) GetEpisodeByAbsoluteNumber(absoluteNumber int) *tvdb.Episode {
	id := AbsoluteNumberIndex(absoluteNumber)
	return s.EpisodesIndex[id]
}

func (s *Series) GetEpisodeByIndex(index string) *tvdb.Episode {
	return s.EpisodesIndex[index]
}

func AbsoluteNumberIndex(absoluteNumber int) string {
	return "a" + strconv.Itoa(absoluteNumber)
}

func SeasonEpisodeNumberIndex(season, number int) string {
	return "s" + strconv.Itoa(season) + "e" + strconv.Itoa(number)
}
