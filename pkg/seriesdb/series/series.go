package series

import (
	"strconv"
)

type Episode struct {
	SeasonNumber   int
	EpisodeNumber  int
	AbsoluteNumber int
	Name           string
	Overview       string
	OriginalName   string
	AiredDate      string
	Context        interface{}
}

type SeriesDetail struct {
	SeriesID     string
	Name         string
	Overview     string
	OriginalName string
	SeasonNames  map[int]string
	Episodes     []*Episode
	Context      interface{}
}

type Series struct {
	SeriesDetail
	EpisodesIndex map[string]*Episode
}

func NewSeries(sd *SeriesDetail, custom *CustomSeries) *Series {
	s := &Series{
		SeriesDetail:  *sd,
		EpisodesIndex: make(map[string]*Episode),
	}

	for _, episode := range sd.Episodes {
		seId := SeasonEpisodeNumberIndex(episode.SeasonNumber, episode.EpisodeNumber)
		s.EpisodesIndex[seId] = episode
		if episode.SeasonNumber != 0 {
			aId := AbsoluteNumberIndex(episode.AbsoluteNumber)
			s.EpisodesIndex[aId] = episode
		}
	}
	if custom != nil &&
		(custom.SeriesID == "" || (custom.SeriesID != "" && custom.SeriesID == s.SeriesID)) {
		if custom.SeriesName != "" {
			s.Name = custom.SeriesName
		}
		if custom.Overview != "" {
			s.Overview = custom.Overview
		}
		for _, namedSeason := range custom.NamedSeasons {
			s.SeasonNames[namedSeason.Number] = namedSeason.Name
		}
		for _, epc := range custom.Episodes {
			var ep *Episode
			if epc.AiredSeason != 0 && epc.AiredEpisodeNumber != 0 {
				ep = s.GetEpisodeBySeasonEpisodeNumber(epc.AiredSeason, epc.AiredEpisodeNumber)
			} else if epc.AbsoluteNumber != 0 {
				ep = s.GetEpisodeByAbsoluteNumber(epc.AbsoluteNumber)
			}
			if ep != nil {
				if epc.EpisodeName != "" {
					ep.Name = epc.EpisodeName
				}
				if epc.Overview != "" {
					ep.Overview = epc.Overview
				}
			}
		}
	}
	return s
}

func (s *Series) GetEpisodeBySeasonEpisodeNumber(season, number int) *Episode {
	id := SeasonEpisodeNumberIndex(season, number)
	return s.EpisodesIndex[id]
}

func (s *Series) GetEpisodeByAbsoluteNumber(absoluteNumber int) *Episode {
	id := AbsoluteNumberIndex(absoluteNumber)
	return s.EpisodesIndex[id]
}

func (s *Series) GetEpisodeByIndex(index string) *Episode {
	return s.EpisodesIndex[index]
}

func AbsoluteNumberIndex(absoluteNumber int) string {
	return "a" + strconv.Itoa(absoluteNumber)
}

func SeasonEpisodeNumberIndex(season, number int) string {
	return "s" + strconv.Itoa(season) + "e" + strconv.Itoa(number)
}
