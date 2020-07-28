package namer

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tnextday/animenamer/pkg/utils"

	"github.com/tnextday/animenamer/pkg/tvdbex"
)

type EpisodeFile struct {
	FileDir   string
	Filename  string
	Infos     Infos
	Subtitles []string
	Series    *tvdbex.SeriesEx
	Episode   *tvdbex.EpisodeEx
}

const (
	InfoKeySeries   = "series"
	InfoKeySeriesId = "seriesId"
	InfoKeySeason   = "season"
	InfoKeyEpisode  = "episode"
	InfoKeyAbsolute = "absolute"
	InfoKeyTitle    = "title"
	InfoKeyExt      = "ext"
	InfoKeyDate     = "date"
)

var (
	reWhitespace = regexp.MustCompile(`\s+`)
)

func (ef *EpisodeFile) Renames(format string, replaceSpaceWith string, includeSubtitles bool) map[string]string {
	ext := strings.ToLower(filepath.Ext(ef.Filename))
	oldBaseName := ef.Filename[0 : len(ef.Filename)-len(ext)]
	newName := utils.NamedFormat(format, ef.Infos)
	if replaceSpaceWith != "" {
		newName = reWhitespace.ReplaceAllString(newName, replaceSpaceWith)
	}
	var newBaseName string
	if !strings.HasSuffix(newName, ext) {
		newBaseName = newName
		newName += ext
	} else {
		newBaseName = newName[0 : len(newName)-len(ext)]
	}
	renames := make(map[string]string)
	renames[ef.Filename] = newName
	if includeSubtitles {
		for _, s := range ef.Subtitles {
			renames[s] = strings.ReplaceAll(s, oldBaseName, newBaseName)
		}
	}
	return renames
}
