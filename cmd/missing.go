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
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/tnextday/animenamer/pkg/namer"
	"github.com/tnextday/animenamer/pkg/seriesdb/series"
	"github.com/tnextday/animenamer/pkg/verbose"

	"github.com/spf13/cobra"
)

// missingCmd represents the missing command
var missingCmd = &cobra.Command{
	Use:   "missing [flags] <anime-dirs>",
	Short: "Find missing episodes in the anime",
	Args:  cobra.MinimumNArgs(1),
	Run:   missingRun,
}

func init() {
	rootCmd.AddCommand(missingCmd)
	missingCmd.Flags().StringSliceP("missing.pattern", "p",
		[]string{
			`(?i).*s(?P<season>\d+)e(?P<episode>\d+).*`,
			`(?i).*(?P<season>\d+)x(?P<episode>\d+).*`,
		},
		"filename regex named pattern, \n"+
			"'absolute', 'season', 'episode' is the required name.\n")
	viper.BindPFlags(missingCmd.Flags())
}

func missingRun(cmd *cobra.Command, args []string) {
	if viper.GetBool("verbose") {
		for k, v := range viper.AllSettings() {
			verbose.Printf("[V] %s: %v\n", k, v)
		}
	}
	if viper.GetString("seriesId") == "" && viper.GetString("name") == "" {
		fmt.Printf("name or seriesId must be defined\n")
		os.Exit(1)
	}
	sdb := createSeriesDB()

	es := namer.EpisodeSearch{
		SeriesDB:     sdb,
		MediaExt:     namer.NewFileExtFromString(viper.GetString("mediaExt"), ","),
		SubtitlesExt: namer.NewFileExtFromString(viper.GetString("subtitleExt"), ","),
		SeriesName:   viper.GetString("name"),
		SeriesId:     viper.GetString("id"),
	}
	if es.SeriesId == "" {
		var err error
		es.SeriesId, err = sdb.Search(es.SeriesName)
		if err != nil {
			fmt.Printf("can't search series, error: %v\n", err)
			os.Exit(1)
		}
	}
	s, err := sdb.GetSeries(es.SeriesId)
	if err != nil {
		fmt.Printf("can't get series from tvdb, error: %v\n", err)
		os.Exit(1)
	}
	for _, p := range viper.GetStringSlice("missing.pattern") {
		if err = es.AddPattern(p); err != nil {
			fmt.Printf("parse pattern (%s) error: %v\n", p, err)
		}
	}
	if len(es.Filters) == 0 {
		fmt.Printf("no valid pattern")
		os.Exit(1)
	}

	recursive := viper.GetBool("recursive")
	episodeFileIndex := map[string]*namer.EpisodeFile{}
	for _, fp := range args {
		fmt.Printf("searching in %s\n", fp)
		episodeFiles, err := es.ListEpisodeFile(fp, recursive)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Printf("found %d episode files\n", len(episodeFiles))
		for _, ef := range episodeFiles {
			seId := series.SeasonEpisodeNumberIndex(ef.Episode.SeasonNumber, ef.Episode.EpisodeNumber)
			episodeFileIndex[seId] = ef
		}
	}
	fmt.Printf("\nSeries: %s\n", s.Name)
	fmt.Printf("SeriesId: %s\n", s.SeriesID)
	// if len(series.Aliases) > 0 {
	// 	fmt.Printf("Aliases: %s\n", strings.Join(series.Aliases, ", "))
	// }
	fmt.Println("")
	for _, ep := range s.Episodes {
		seId := series.SeasonEpisodeNumberIndex(ep.SeasonNumber, ep.EpisodeNumber)
		if _, exists := episodeFileIndex[seId]; !exists {
			s := fmt.Sprintf("missing S%.2dE%.2d", ep.SeasonNumber, ep.EpisodeNumber)
			if ep.SeasonNumber != 0 {
				s += fmt.Sprintf(", absolute %.3d", ep.AbsoluteNumber)
			}
			fmt.Println(s)
		}

	}
}
