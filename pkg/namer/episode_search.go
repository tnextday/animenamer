package namer

import (
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tnextday/animenamer/pkg/seriesdb"
	"github.com/tnextday/animenamer/pkg/seriesdb/series"
	"github.com/tnextday/animenamer/pkg/verbose"
)

type EpisodeSearch struct {
	Filters      []*regexp.Regexp
	MediaExt     FileExt
	SubtitlesExt FileExt
	SeriesDB     *seriesdb.SeriesDB
	SeriesName   string
	SeriesId     string
	RegexpOnly   bool
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
	if !es.RegexpOnly {
		names := re.SubexpNames()
		if es.SeriesName == "" && es.SeriesId == "" {
			if !(contains(names, InfoKeySeries) || contains(names, InfoKeySeriesId)) {
				return errors.New("series or seriesId must be defined")
			}
		}
		if !(contains(names, InfoKeyAbsolute) ||
			(contains(names, InfoKeySeason) && contains(names, InfoKeyEpisode))) {
			return errors.New("absolute or season&episode must be defined")
		}
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
			var filesInSameDir []string
			for _, fi := range fileInfos {
				if strings.HasPrefix(fi.Name(), ".") {
					continue
				}
				if fi.IsDir() && recursive {
					dirs.PushBack(path.Join(dirname, fi.Name()))
				} else {
					filesInSameDir = append(filesInSameDir, fi.Name())
				}
			}
			if fileName != "" {
				if mf, err := es.newEpisodeFile(dirname, fileName, filesInSameDir); err == nil {
					mediaFileList = append(mediaFileList, mf)
				} else {
					verbose.Printf("newEpisodeFile (%s) err: %v\n", path.Join(dirname, fileName), err)
				}
				break //当输入是一个文件时，只处理这个文件
			} else {
				for _, fn := range filesInSameDir {
					if mf, err := es.newEpisodeFile(dirname, fn, filesInSameDir); err == nil {
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

func (es *EpisodeSearch) newEpisodeFile(dirname, filename string, filesInSameDir []string) (*EpisodeFile, error) {
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
	if es.SeriesId != "" {
		infos[InfoKeySeriesId] = es.SeriesId
	}
	var (
		se      *series.Series
		episode *series.Episode
	)
	if !es.RegexpOnly && es.SeriesDB != nil {
		seriesId := infos.GetString(InfoKeySeriesId)
		if seriesId == "" {
			seriesId, err = es.SeriesDB.Search(infos.GetString(InfoKeySeries))
			if err != nil {
				verbose.Printf("search series (%s) error: %v\n", infos.GetString(InfoKeySeries), err)
			}
		}
		if seriesId != "" {
			se, err = es.SeriesDB.GetSeries(seriesId)
			if err != nil {
				return nil, fmt.Errorf("get series (%s) error: %v", seriesId, err)
			}
		}
	}
	if se != nil {
		infos[InfoKeySeries] = se.Name
		infos[InfoKeySeriesId] = se.SeriesID
		episode = se.GetEpisodeBySeasonEpisodeNumber(infos.GetInt(InfoKeySeason), infos.GetInt(InfoKeyEpisode))
		if episode == nil {
			episode = se.GetEpisodeByAbsoluteNumber(infos.GetInt(InfoKeyAbsolute))
		}
	}
	if episode != nil {
		infos[InfoKeySeason] = episode.SeasonNumber
		infos[InfoKeySeasonName] = se.SeasonNames[episode.SeasonNumber]
		infos[InfoKeyEpisode] = episode.EpisodeNumber
		infos[InfoKeyAbsolute] = episode.AbsoluteNumber
		infos[InfoKeyEpisodeName] = episode.Name
		infos[InfoKeyDate] = episode.AiredDate
	}
	mediaFile := EpisodeFile{
		FileDir:   dirname,
		Filename:  filename,
		Infos:     infos,
		Subtitles: searchSubtitles(filename, filesInSameDir, es.SubtitlesExt),
		Series:    se,
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

func searchSubtitles(mediaName string, filesInSameDir []string, subtitlesExt FileExt) []string {
	var subtitles []string
	baseName := mediaName[0 : len(mediaName)-len(filepath.Ext(mediaName))]
	for _, fn := range filesInSameDir {
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
