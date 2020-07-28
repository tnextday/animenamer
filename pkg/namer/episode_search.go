package namer

import (
	"container/list"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tnextday/animenamer/pkg/verbose"

	"github.com/tnextday/animenamer/pkg/tvdbex"
)

type EpisodeSearch struct {
	Filters      []*regexp.Regexp
	MediaExt     FileExt
	SubtitlesExt FileExt
	TVDB         *tvdbex.TVDBEx
	SeriesName   string
	SeriesId     int
}

func contains(slice []string, e string) bool {
	for _, s := range slice {
		if s == e {
			return true
		}
	}
	return false
}

func (es *EpisodeSearch) AddPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	names := re.SubexpNames()
	if es.SeriesName == "" && es.SeriesId == 0 {
		if !(contains(names, InfoKeySeries) || contains(names, InfoKeySeriesId)) {
			return errors.New("series or seriesId must be defined")
		}
	}
	if !(contains(names, InfoKeyAbsolute) ||
		(contains(names, InfoKeySeason) && contains(names, InfoKeyEpisode))) {
		return errors.New("absolute or season&episode must be defined")
	}
	es.Filters = append(es.Filters, re)
	return nil
}

func (es *EpisodeSearch) ListEpisodeFile(fp string, recursive bool) ([]*EpisodeFile, error) {
	fi, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}

	baseDir := fp
	fileName := ""
	if !fi.IsDir() {
		baseDir = filepath.Dir(fp)
		fileName = fi.Name()
	}

	if fileName != "" {
		_, err := multiMatchEpisodeFile(es.Filters, fileName)
		if err != nil {
			return nil, errors.New("no matched media file found")
		}
	}
	dirs := list.New()
	dirs.PushBack(baseDir)
	var mediaFileList []*EpisodeFile
	for {
		e := dirs.Back()
		if e == nil {
			break
		}
		dirs.Remove(e)
		dirname := e.Value.(string)
		if fileInfos, err := ioutil.ReadDir(dirname); err == nil {
			var filenameList []string
			for _, fi := range fileInfos {
				if fi.IsDir() && recursive {
					dirs.PushBack(path.Join(dirname, fi.Name()))
				} else {
					filenameList = append(filenameList, fi.Name())
				}
			}
			if fileName != "" {
				if mf, err := es.newEpisodeFile(dirname, fileName, filenameList); err == nil {
					mediaFileList = append(mediaFileList, mf)
				} else {
					verbose.Printf("newEpisodeFile (%s) err: %v\n", path.Join(dirname, fileName), err)
				}
			} else {
				for _, fn := range filenameList {
					if mf, err := es.newEpisodeFile(dirname, fn, filenameList); err == nil {
						mediaFileList = append(mediaFileList, mf)
					} else {
						verbose.Printf("newEpisodeFile (%s) err: %v\n", path.Join(dirname, fn), err)
					}
				}
			}
		} else {
			verbose.Printf("readdir (%s) err: %v\n", dirname, err)
		}
	}
	return mediaFileList, nil
}

func (es *EpisodeSearch) newEpisodeFile(dirname, filename string, filenameList []string) (*EpisodeFile, error) {
	ext := filepath.Ext(filename)
	if !es.MediaExt.Contains(ext) {
		return nil, errors.New("is not media file")
	}
	infos, err := multiMatchEpisodeFile(es.Filters, filename)
	if err != nil {
		return nil, errors.New("not match")
	}
	infos[InfoKeyExt] = strings.ToLower(ext[1:])
	if es.SeriesName != "" {
		infos[InfoKeySeries] = es.SeriesName
	}
	if es.SeriesId != 0 {
		infos[InfoKeySeriesId] = es.SeriesId
	}
	var (
		series  *tvdbex.SeriesEx
		episode *tvdbex.EpisodeEx
	)
	if es.TVDB != nil {
		seriesId := infos.GetInt(InfoKeySeriesId)
		if seriesId == 0 {
			seriesId, err = es.TVDB.Search(infos.GetString(InfoKeySeries))
			if err != nil {
				verbose.Printf("search series (%s) error: %v\n", infos.GetString(InfoKeySeries), err)
			}
		}
		if seriesId != 0 {
			series, err = es.TVDB.GetSeriesEx(seriesId)
			if err != nil {
				verbose.Printf("get series (%d) error: %v\n", seriesId, err)
			}
		}
	}
	if series != nil {
		infos[InfoKeySeries] = series.SeriesName
		infos[InfoKeySeriesId] = series.ID
		episode = series.GetEpisodeBySeasonEpisodeNumber(infos.GetInt(InfoKeySeason), infos.GetInt(InfoKeyEpisode))
		if episode == nil {
			episode = series.GetEpisodeByAbsoluteNumber(infos.GetInt(InfoKeyAbsolute))
		}
	}
	if episode != nil {
		infos[InfoKeySeason] = episode.AiredSeason
		infos[InfoKeyEpisode] = episode.AiredEpisodeNumber
		infos[InfoKeyAbsolute] = episode.AbsoluteNumber
		infos[InfoKeyTitle] = episode.EpisodeName
		infos[InfoKeyDate] = episode.FirstAired
	}
	mediaFile := EpisodeFile{
		FileDir:   dirname,
		Filename:  filename,
		Infos:     infos,
		Subtitles: searchSubtitles(filename, filenameList, es.SubtitlesExt),
		Series:    series,
		Episode:   episode,
	}
	return &mediaFile, nil
}

func matchEpisodeFile(filter *regexp.Regexp, filename string) (infos Infos, err error) {
	matches := filter.FindStringSubmatch(filename)
	if matches == nil {
		return nil, errors.New("not match")
	}
	infos = make(Infos)
	for i, v := range filter.SubexpNames() {
		switch v {
		case "":
			break
		case InfoKeySeriesId, InfoKeySeason, InfoKeyEpisode, InfoKeyAbsolute:
			if i, err := strconv.Atoi(strings.TrimSpace(matches[i])); err == nil {
				infos[v] = i
			} else {
				verbose.Printf("%s, convert %s(%s) to int err: %v\n", filename, v, matches[i], err)
			}
		default:
			infos[v] = strings.TrimSpace(matches[i])
		}
	}
	return infos, nil

}

func multiMatchEpisodeFile(filters []*regexp.Regexp, filename string) (infos Infos, err error) {
	for _, filter := range filters {
		infos, err = matchEpisodeFile(filter, filename)
		if err == nil {
			return
		}
	}
	return
}

func searchSubtitles(mediaName string, filenameList []string, subtitlesExt FileExt) []string {
	var subtitles []string
	baseName := mediaName[0 : len(mediaName)-len(filepath.Ext(mediaName))]
	for _, fn := range filenameList {
		if !subtitlesExt.Is(fn) {
			continue
		}
		if strings.HasPrefix(fn, baseName) {
			subtitles = append(subtitles, fn)
		}
	}
	return subtitles
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
