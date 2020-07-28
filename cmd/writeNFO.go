/*
Copyright © 2019 tnextday <fw2k4@163.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/pioz/tvdb"

	"github.com/tnextday/animenamer/pkg/kodi"

	"github.com/spf13/viper"
	"github.com/tnextday/animenamer/pkg/namer"
	"github.com/tnextday/animenamer/pkg/tvdbex"
	"github.com/tnextday/animenamer/pkg/verbose"

	"github.com/spf13/cobra"
)

// writeNFOCmd represents the writeNFO command
var writeNFOCmd = &cobra.Command{
	Use:   "writeNFO [flags] <anime-dirs>",
	Short: "write kodi NFO files",
	Args:  cobra.MinimumNArgs(1),
	Run:   writeNFORun,
}

func init() {
	rootCmd.AddCommand(writeNFOCmd)

	writeNFOCmd.Flags().StringSliceP("writeNFO.pattern", "p",
		[]string{
			`(?i).*s(?P<season>\d+)e(?P<episode>\d+).*`,
			`(?i).*(?P<season>\d+)x(?P<episode>\d+).*`,
		},
		"filename regex named pattern, \n"+
			"'absolute', 'season', 'episode' is the required name.\n")
	writeNFOCmd.Flags().StringP("baseDir", "b", "", "series base folder to write tvshow.nfo")
	writeNFOCmd.Flags().Bool("overrideImage", false, "override exists image file")
	writeNFOCmd.Flags().Bool("downloadAllFanart", false, "download all of fanarts")
	viper.BindPFlags(writeNFOCmd.Flags())
}

func writeNFORun(cmd *cobra.Command, args []string) {
	if viper.GetBool("verbose") {
		for k, v := range viper.AllSettings() {
			verbose.Printf("[V] %s: %v\n", k, v)
		}
	}
	if viper.GetInt("seriesId") == 0 && viper.GetString("name") == "" {
		fmt.Printf("name or seriesId must be defined\n")
		os.Exit(1)
	}
	apiKey := viper.GetString("apikey")
	if apiKey == "" {
		apiKey = DefaultTvDbApiKey
	}
	tvdbEx, err := tvdbex.NewTVDBEx(apiKey, viper.GetString("language"), loadCustomConfig())
	if err != nil {
		fmt.Printf("new tvdb error: %v\n", err)
		os.Exit(1)
	}

	es := namer.EpisodeSearch{
		MediaExt:     namer.NewFileExtFromString(viper.GetString("mediaExt"), ","),
		SubtitlesExt: namer.NewFileExtFromString(viper.GetString("subtitleExt"), ","),
		TVDB:         tvdbEx,
		SeriesName:   viper.GetString("name"),
		SeriesId:     viper.GetInt("seriesId"),
	}
	fmt.Printf("get series details")
	if es.SeriesId == 0 {
		es.SeriesId, err = tvdbEx.Search(es.SeriesName)
		if err != nil {
			fmt.Printf("can't search series, error: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf(", SeriesId %d\n", es.SeriesId)

	for _, p := range viper.GetStringSlice("writeNFO.pattern") {
		if err = es.AddPattern(p); err != nil {
			fmt.Printf("parse pattern (%s) error: %v\n", p, err)
		}
	}
	if len(es.Filters) == 0 {
		fmt.Printf("no valid pattern\n")
		os.Exit(1)
	}
	series, err := tvdbEx.GetSeriesEx(es.SeriesId)
	if err != nil {
		fmt.Printf("can't get series from tvdb, error: %v\n", err)
		os.Exit(1)
	}
	if err := tvdbEx.GetSeriesActors(series); err != nil {
		fmt.Printf("can't get series actors from tvdb, error: %v\n", err)
		os.Exit(1)
	}
	if err := tvdbEx.GetSeriesImages(series); err != nil {
		fmt.Printf("can't get series images from tvdb, error: %v\n", err)
		os.Exit(1)
	}
	if err := tvdbEx.GetSeriesSummary(series); err != nil {
		fmt.Printf("can't get series summary from tvdb, error: %v\n", err)
		os.Exit(1)
	}
	recursive := viper.GetBool("recursive")
	var episodeFiles []*namer.EpisodeFile
	for _, fp := range args {
		fmt.Printf("searching episode files in %s\n", fp)
		files, err := es.ListEpisodeFile(fp, recursive)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Printf("found %d episode files\n", len(files))
		episodeFiles = append(episodeFiles, files...)
	}

	fmt.Printf("writing nfo files\n")
	overrideImage := viper.GetBool("overrideImage")
	for _, epf := range episodeFiles {
		verbose.Printf("processing %s\n",
			epf.Filename)
		if err = writeEpisodeNFO(epf, overrideImage); err != nil {
			fmt.Printf("[E]%v\n", err)
		}
	}
	baseDir := viper.GetString("baseDir")
	if baseDir != "" {
		if !path.IsAbs(baseDir) && cfgFile != "" {
			cfgDir := path.Dir(cfgFile)
			if baseDir, err = filepath.Rel(cfgDir, baseDir); err != nil {
				fmt.Printf("[E]%v\n", err)
			}
		}
		fmt.Printf("writing tvshow.nfo\n")
		if err = writeTVShowNFO(baseDir, series, overrideImage); err != nil {
			fmt.Printf("[E]%v\n", err)
		}
	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func downloadToFile(url, fp string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fp, buf, 0644)
}

func writeNFOFile(v interface{}, fp string) error {
	buf, err := xml.MarshalIndent(v, "", "\t")

	if err != nil {
		return err
	}
	return ioutil.WriteFile(fp, buf, 0644)
}

func writeEpisodeNFO(epf *namer.EpisodeFile, overrideImage bool) (err error) {
	ep := epf.Episode
	series := epf.Series

	episodeDetails := kodi.EpisodeDetails{
		Title:         ep.EpisodeName,
		OriginalTitle: ep.OriginalEpisodeName,
		ShowTitle:     ep.EpisodeName,
		Ratings: []*kodi.Rating{
			&kodi.Rating{
				Name:    "tvdb",
				Max:     10,
				Default: true,
				Value:   ep.SiteRating,
				Votes:   ep.SiteRatingCount,
			},
		},
		Season:  ep.AiredSeason,
		Episode: ep.AiredEpisodeNumber,
		Plot:    ep.Overview,
		Thumbs:  []string{tvdbex.GetEpisodeImageUrl(ep)},
		UniqueIDs: []*kodi.UniqueID{
			&kodi.UniqueID{
				Type:    "tvdb",
				Default: true,
				ID:      strconv.Itoa(ep.ID),
			},
			&kodi.UniqueID{
				Type:    "imdb",
				Default: false,
				ID:      ep.ImdbID,
			},
		},
		Genres:    series.Genre,
		Credits:   ep.Writers,
		Directors: ep.Directors,
		Premiered: ep.FirstAired,
		Studios:   []string{series.Network},
	}
	actorsMap := map[string]bool{}
	for _, a := range series.Actors {
		actorsMap[a.Name] = true
		episodeDetails.Actors = append(episodeDetails.Actors,
			&kodi.Actor{
				Name:  a.Name,
				Role:  a.Role,
				Order: a.SortOrder,
				Thumb: tvdb.ImageURL(a.Image),
			})
	}
	for _, guestStar := range ep.GuestStars {
		if _, exists := actorsMap[guestStar]; exists {
			continue
		}
		episodeDetails.Actors = append(episodeDetails.Actors,
			&kodi.Actor{
				Name:  guestStar,
				Order: len(episodeDetails.Actors),
			})
	}
	baseName := epf.Filename[0 : len(epf.Filename)-len(filepath.Ext(epf.Filename))]
	nfoPath := path.Join(epf.FileDir, baseName+".nfo")

	if err = writeNFOFile(episodeDetails, nfoPath); err != nil {
		return fmt.Errorf("write nfo %s error: %v\n", nfoPath, err)
	}
	imgExt := path.Ext(ep.Filename)
	imagePath := path.Join(epf.FileDir, baseName+"-thumb"+imgExt)
	if !overrideImage && fileExists(imagePath) {
		return nil
	}
	imageUrl := tvdb.ImageURL(ep.Filename)
	verbose.Printf("downloading image [%s](%s)\n", imagePath, imageUrl)
	if err = downloadToFile(imageUrl, imagePath); err != nil {
		return fmt.Errorf("download image [%s](%s) error: %v\n", imagePath, imageUrl, err)
	}
	return nil
}

func writeTVShowNFO(baseDir string, series *tvdbex.SeriesEx, overrideImage bool) (err error) {
	tvshowPath := path.Join(baseDir, "tvshow.nfo")
	episodeCount, _ := strconv.Atoi(series.Summary.AiredEpisodes)
	tvshow := kodi.TVShow{
		Title:         series.SeriesName,
		OriginalTitle: series.OriginalSeriesName,
		Ratings: []*kodi.Rating{
			&kodi.Rating{
				Name:    "tvdb",
				Max:     10,
				Default: true,
				Value:   series.SiteRating,
				Votes:   series.SiteRatingCount,
			},
		},
		Season:  len(series.Summary.AiredSeasons),
		Episode: episodeCount,
		Plot:    series.Overview,
		//TagLine:        "",
		//Thumbs:  nil,
		//Fanarts: nil,
		//Mpaa:           "",
		EpisodeGuide: nil,
		UniqueIDs: []*kodi.UniqueID{
			&kodi.UniqueID{
				Type:    "tvdb",
				Default: true,
				ID:      strconv.Itoa(series.ID),
			},
			&kodi.UniqueID{
				Type:    "imdb",
				Default: false,
				ID:      series.ImdbID,
			},
		},
		Genres:    series.Genre,
		Tags:      nil,
		Premiered: series.FirstAired,
		Status:    series.Status,
		Studios:   []string{series.Network},
		//Trailer:        "",
		NamedSeasons: series.NamedSeasons,
	}

	for _, a := range series.Actors {
		tvshow.Actors = append(tvshow.Actors,
			&kodi.Actor{
				Name:  a.Name,
				Role:  a.Role,
				Order: a.SortOrder,
				Thumb: tvdb.ImageURL(a.Image),
			})
	}

	imageUrls := make(map[string]string)

	addImage := func(t, fn string) {
		if _, exits := imageUrls[t]; !exits {
			imageUrls[t] = tvdb.ImageURL(fn)
		}
	}
	downloadAllFanarts := viper.GetBool("downloadAllFanart")
	for i, img := range series.Images[tvdbex.ImageTypeKeyFanArt] {
		tvshow.Fanarts = append(tvshow.Fanarts,
			&kodi.Thumb{
				Preview: tvdb.ImageURL(img.Thumbnail),
				Uri:     tvdb.ImageURL(img.FileName),
			},
		)
		if i == 0 {
			addImage("fanart", img.FileName)
		} else if downloadAllFanarts {
			addImage("fanart"+strconv.Itoa(i), img.FileName)
		}
	}
	for k, v := range series.Images {
		if k == tvdbex.ImageTypeKeyFanArt {
			continue
		}
		for _, img := range v {
			thumb := &kodi.Thumb{
				Preview: tvdb.ImageURL(img.Thumbnail),
				Uri:     tvdb.ImageURL(img.FileName),
			}
			switch img.KeyType {
			case tvdbex.ImageTypeKeyPoster:
				thumb.Aspect = "poster"
				addImage("poster", img.FileName)
			case tvdbex.ImageTypeKeySeries:
				thumb.Aspect = "banner"
				addImage("banner", img.FileName)
			case tvdbex.ImageTypeKeySeason:
				thumb.Aspect = "poster"
				thumb.Type = "season"
				thumb.Season, _ = strconv.Atoi(img.SubKey)
				if img.SubKey == "0" {
					addImage("season-specials-poster", img.FileName)
				} else {
					addImage(fmt.Sprintf("season%02d-poster", thumb.Season), img.FileName)
				}
			case tvdbex.ImageTypeKeySeasonWide:
				thumb.Aspect = "banner"
				thumb.Type = "season"
				thumb.Season, _ = strconv.Atoi(img.SubKey)
				if img.SubKey == "0" {
					addImage("season-specials-banner", img.FileName)
				} else {
					addImage(fmt.Sprintf("season%02d-banner", thumb.Season), img.FileName)
				}
			default:
				thumb.Aspect = img.KeyType
			}
			tvshow.Thumbs = append(tvshow.Thumbs, thumb)
		}
	}
	if err = writeNFOFile(tvshow, tvshowPath); err != nil {
		return fmt.Errorf("write tvshow.nfo %s error: %v\n", tvshowPath, err)
	}
	for k, u := range imageUrls {
		imgExt := path.Ext(u)
		imagePath := path.Join(baseDir, k+imgExt)
		if !overrideImage && fileExists(imagePath) {
			continue
		}
		verbose.Printf("downloading image [%s](%s)\n", imagePath, u)
		if err = downloadToFile(u, imagePath); err != nil {
			fmt.Sprintf("[E]download image [%s](%s) error: %v\n", imagePath, u, err)
		}
	}
	return nil
}
