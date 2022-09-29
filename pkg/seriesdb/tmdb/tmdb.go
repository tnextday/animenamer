package tmdb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	go_tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/tnextday/animenamer/pkg/seriesdb/series"
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
	var absoluteGroups []*EpisodeGroup

	absoluteGroupSeason := t.AbsoluteGroupSeason
	if v, exists := options["absoluteGroupSeason"]; exists {
		absoluteGroupSeason = v
	}
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
			return nil, fmt.Errorf("absoluteGroupSeason matches unexpected: %v", matches)
		}
	}
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
		return nil, fmt.Errorf("can't found absolute group")
	}
	delete(apiOptions, "append_to_response")
	episodeGroupsDetails, err := t.Client.GetTVEpisodeGroupsDetails(absoluteGroupId, apiOptions)
	if err != nil {
		return nil, err
	}
	if episodeGroupsDetails.GroupCount == 0 {
		return nil, fmt.Errorf("no episode groups")
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
	for _, ep := range episodeGroupsDetails.Groups[episodesGroupIdx].Episodes {
		sep := &series.Episode{
			SeasonNumber:   ep.SeasonNumber,
			EpisodeNumber:  ep.EpisodeNumber,
			AbsoluteNumber: ep.Order + 1,
			Name:           ep.Name,
			Overview:       ep.Overview,
			OriginalName:   ep.Name,
			AiredDate:      ep.AirDate,
			Context: &tmdbEpisodeContext{
				TVDetails: tvDetail,
			},
		}
		sd.Episodes = append(sd.Episodes, sep)

	}
	return sd, nil
}
