package tmdb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	go_tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/tnextday/animenamer/pkg/seriesdb/series"
	"github.com/tnextday/animenamer/pkg/verbose"
)

var regGroupSeason = regexp.MustCompile(`^([^\(\)]+)?\s*(\((.*)\))?$`)

type TMDB struct {
	Client              *go_tmdb.Client
	AbsoluteGroupSeason string
}

type tmdbEpisodeContext struct {
	TVEpisodeDetails *go_tmdb.TVEpisodeDetails
	TVDetails        *go_tmdb.TVDetails
}

func NewTMDB(apiKey, absoluteGroupSeason string) (*TMDB, error) {
	tmdbClient, err := go_tmdb.Init(apiKey)
	if err != nil {
		return nil, err
	}
	tmdbClient.SetClientAutoRetry()
	return &TMDB{
		Client:              tmdbClient,
		AbsoluteGroupSeason: absoluteGroupSeason,
	}, nil
}

func (t *TMDB) Search(q, language string) (id string, err error) {
	options := map[string]string{
		"language": language,
	}
	result, err := t.Client.GetSearchTVShow(q, options)
	if err != nil {
		return "", err
	}

	switch result.TotalResults {
	case 0:
		return "", series.ErrSearchNotFound
	case 1:
		return strconv.FormatInt(result.Results[0].ID, 10), nil
	default:
		return "", errors.New("returns multiple results, please use ID to match")
	}
}

type EpisodeGroup struct {
	Name string
	ID   string
}

func (t *TMDB) GetSeries(seriesId, language string, options map[string]string) (*series.SeriesDetail, error) {
	id, err := strconv.Atoi(seriesId)
	if err != nil {
		return nil, err
	}
	apiOptions := map[string]string{
		"language":           language,
		"append_to_response": "episode_groups",
	}
	tvDetail, err := t.Client.GetTVDetails(id, apiOptions)
	if err != nil {
		return nil, err
	}
	sd := &series.SeriesDetail{
		SeriesID:     seriesId,
		Name:         tvDetail.Name,
		Overview:     tvDetail.Overview,
		OriginalName: tvDetail.OriginalName,
		SeasonNames:  make(map[int]string),
		Context:      tvDetail,
	}
	for _, s := range tvDetail.Seasons {
		sd.SeasonNames[s.SeasonNumber] = s.Name
	}

	delete(apiOptions, "append_to_response")
	for _, s := range tvDetail.Seasons {
		seasonDetails, err := t.Client.GetTVSeasonDetails(id, s.SeasonNumber, apiOptions)
		if err != nil {
			return nil, err
		}
		for _, ep := range seasonDetails.Episodes {
			sep := &series.Episode{
				SeasonNumber:  ep.SeasonNumber,
				EpisodeNumber: ep.EpisodeNumber,
				Name:          ep.Name,
				Overview:      ep.Overview,
				OriginalName:  ep.Name,
				AiredDate:     ep.AirDate,
				Context: &tmdbEpisodeContext{
					TVDetails: tvDetail,
				},
			}
			sd.Episodes = append(sd.Episodes, sep)
		}
	}

	absoluteGroupSeason := t.AbsoluteGroupSeason
	if v, exists := options["absoluteGroupSeason"]; exists {
		absoluteGroupSeason = v
	}
	err = t.updateAbsoluteNumber(sd, tvDetail, absoluteGroupSeason, language)
	if err != nil {
		verbose.Printf("Series(%s) updateAbsoluteNumber %v", seriesId, err)
	}
	return sd, nil
}

func (t *TMDB) updateAbsoluteNumber(seriesDetail *series.SeriesDetail, tvDetail *go_tmdb.TVDetails, absoluteGroupSeason, language string) (err error) {
	var (
		groupName  string
		seasonName string
	)
	if absoluteGroupSeason != "" {
		matches := regGroupSeason.FindStringSubmatch(absoluteGroupSeason)
		if len(matches) == 4 {
			seasonName = strings.TrimSpace(matches[1])
			groupName = strings.TrimSpace(matches[3])
		} else {
			return fmt.Errorf("absoluteGroupSeason matches unexpected: %v", matches)
		}
	}
	var absoluteGroups []*EpisodeGroup
	for _, v := range tvDetail.EpisodeGroups.Results {
		if v.Type == 2 {
			absoluteGroups = append(absoluteGroups, &EpisodeGroup{
				Name: v.Name,
				ID:   v.ID,
			})
		}
	}
	var absoluteGroupId string
	if len(absoluteGroups) == 1 {
		absoluteGroupId = absoluteGroups[0].ID
	} else {
		for i, g := range absoluteGroups {
			if strings.TrimSpace(g.Name) == groupName {
				absoluteGroupId = absoluteGroups[i].ID
				break
			}
		}
	}
	if absoluteGroupId == "" {
		return fmt.Errorf("can't found absolute group")
	}
	apiOptions := map[string]string{
		"language": language,
	}
	episodeGroupsDetails, err := t.Client.GetTVEpisodeGroupsDetails(absoluteGroupId, apiOptions)
	if err != nil {

		return err
	}
	if episodeGroupsDetails.GroupCount == 0 {
		return err
	}

	episodesGroupIdx := 0
	for i, g := range episodeGroupsDetails.Groups {
		if seasonName != "" {
			if strings.TrimSpace(g.Name) == seasonName {
				episodesGroupIdx = i
				break
			}
		} else if g.Order == 1 {
			// 正篇的order一般为1
			episodesGroupIdx = i
			break
		}
	}
	absoluteMap := make(map[string]int)
	for _, ep := range episodeGroupsDetails.Groups[episodesGroupIdx].Episodes {
		absoluteMap[series.SeasonEpisodeNumberIndex(ep.SeasonNumber, ep.EpisodeNumber)] = ep.Order + 1
	}
	for _, ep := range seriesDetail.Episodes {
		if absoluteNumber, exists := absoluteMap[series.SeasonEpisodeNumberIndex(ep.SeasonNumber, ep.EpisodeNumber)]; exists {
			ep.AbsoluteNumber = absoluteNumber
		}
	}
	return nil
}
