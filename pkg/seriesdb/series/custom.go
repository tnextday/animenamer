package series

import (
	"bytes"
	"io"
	"os"

	"github.com/tnextday/animenamer/pkg/kodi"
	"gopkg.in/yaml.v3"
)

type CustomSeries struct {
	SeriesID     string              `json:"seriesId" yaml:"seriesId"`
	SeriesName   string              `json:"seriesName" yaml:"seriesName"`
	Overview     string              `json:"overview" yaml:"overview"`
	NamedSeasons []*kodi.NamedSeason `json:"namedSeason" yaml:"namedSeason"`
	Episodes     []*CustomEpisode    `json:"episodes" yaml:"episodes"`
	TMDB         struct {
		AbsoluteGroupSeason string `json:"absoluteGroupSeason" yaml:"absoluteGroupSeason"`
	} `json:"tmdb" yaml:"tmdb"`
}

type CustomEpisode struct {
	AbsoluteNumber     int    `json:"absolute" yaml:"absolute"`
	AiredEpisodeNumber int    `json:"episode" yaml:"episode"`
	AiredSeason        int    `json:"season" yaml:"season"`
	EpisodeName        string `json:"name" yaml:"name"`
	Overview           string `json:"overview" yaml:"overview"`
}

func LoadCustomSeries(fp string) ([]*CustomSeries, error) {
	buf, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	var customs []*CustomSeries
	r := bytes.NewReader(buf)
	decoder := yaml.NewDecoder(r)
	for {
		var series CustomSeries
		if err := decoder.Decode(&series); err != nil {
			// Break when there are no more documents to decode
			if err != io.EOF {
				return nil, err
			}
			break
		}
		customs = append(customs, &series)
	}
	if err != nil {
		return nil, err
	}
	return customs, nil
}
