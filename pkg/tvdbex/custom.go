package tvdbex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/tnextday/animenamer/pkg/kodi"
	"gopkg.in/yaml.v3"
)

type CustomSeries struct {
	SeriesID    string              `json:"seriesId" yaml:"seriesId"`
	SeriesName  string              `json:"seriesName" yaml:"seriesName"`
	Overview    string              `json:"overview" yaml:"overview"`
	NamedSeason []*kodi.NamedSeason `json:"namedSeason" yaml:"namedSeason"`
	Episodes    []*CustomEpisode    `json:"episodes" yaml:"episodes"`
}

type CustomEpisode struct {
	AbsoluteNumber     int    `json:"absolute" yaml:"absolute"`
	AiredEpisodeNumber int    `json:"episode" yaml:"episode"`
	AiredSeason        int    `json:"season" yaml:"season"`
	EpisodeName        string `json:"name" yaml:"name"`
	Overview           string `json:"overview" yaml:"overview"`
}

func LoadCustomSeries(fp string) (*CustomSeries, error) {
	buf, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	ext := path.Ext(fp)
	var series CustomSeries
	switch ext {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(buf, &series)
	case ".json":
		err = json.Unmarshal(buf, &series)
	default:
		err = fmt.Errorf("format %s not support", ext)
	}
	if err != nil {
		return nil, err
	}
	return &series, nil
}
