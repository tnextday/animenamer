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
	"path"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/spf13/cobra"
	"github.com/tnextday/animenamer/pkg/namer"
	"github.com/tnextday/animenamer/pkg/seriesdb"
	"github.com/tnextday/animenamer/pkg/seriesdb/series"
	"github.com/tnextday/animenamer/pkg/seriesdb/tmdb"
	"github.com/tnextday/animenamer/pkg/seriesdb/tvdb"
	"github.com/tnextday/animenamer/pkg/utils"
	"github.com/tnextday/animenamer/pkg/verbose"

	"github.com/spf13/viper"
)

var (
	cfgFile           string
	DefaultTVDBApiKey string
	DefaultTMDBApiKey string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "animenamer [flags] <anime-files>",
	Short: "rename anime file with absolute number or season&episode number in the filename",
	Args:  cobra.MinimumNArgs(1),
	Run:   rootCmdFunc,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is animenamer.yml)")
	rootCmd.PersistentFlags().StringSlice("custom", []string{"animenamer.custom.yml"},
		"custom series info file.\n"+
			"(default is animenamer.custom.yml)")
	rootCmd.PersistentFlags().String("db", "tmdb", "tmdb or tvdb")
	rootCmd.PersistentFlags().String("apikey", "", "the apikey of tmdb or tvdb")
	rootCmd.PersistentFlags().String("id", "", "explicitly set the show id for db to use (applies to all files)")
	rootCmd.PersistentFlags().String("tmdb.absoluteGroupSeason", "",
		"when there are multiple tmdb absolute groups, use this name to match,\n"+
			"this value can find at group episode -> season x page,\n"+
			"the format is 'season name (group name)'")
	rootCmd.PersistentFlags().StringP("name", "n", "", "override the parsed series name with this (applies to all files)")
	rootCmd.PersistentFlags().String("mediaExt", "mkv,mp4,avi,rm,rmvb,mov,m4v,wmv", "media file extensions")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "show debugging info")

	rootCmd.Flags().Bool("renameSubtitle", true, "also rename subtitle file")
	rootCmd.Flags().String("subtitleExt", "ass,ssa,srt,sub", "subtitle file extensions")
	rootCmd.Flags().StringSliceP("pattern", "p", nil,
		"filename regex named pattern, \n"+
			"'series', 'seriesId', 'absolute', 'season', 'episode', 'ext' is the reserved name, \n"+
			"its may be override by tvdb info.\n"+
			`example: ^(?P<name>\.+)\.(?P<absolute>\d+).*\.(?P<ext>\w+)$`)
	//rootCmd.MarkFlagRequired("pattern")
	rootCmd.Flags().String("format", "{series}.S{season.2}E{episode.2}.[{absolute.3}].{ext}",
		"new filename format. variables:\n"+
			"'series', 'seriesId', 'season', 'seasonName', 'episode', 'absolute', 'date', 'episodeName', 'ext'\n"+
			"and named variables in filename pattern matched,\n"+
			"you can use {variable.n} for number padding.\n")
	rootCmd.Flags().StringP("moveToDir", "m", "",
		"move to destination directory,\n"+
			"a format string can be used in the path, for example:\n"+
			"new_path/S{season.2}-{seasonName}\n"+
			"Note that relative paths are relative to the input folder")
	rootCmd.Flags().String("replaceSpace", "", "replace the whitespace with this value in new filename")
	rootCmd.Flags().StringP("log", "l", "rename", "the rename log name for recovery")
	rootCmd.Flags().BoolP("regexpOnly", "R", false,
		"use regexp only to rename files without getting information from the web, \n"+
			"in this mode, only 'ext' is a reserved name")
	lang, _ := jibber_jabber.DetectLanguage()
	rootCmd.PersistentFlags().String("language", lang, "preferred language")
	rootCmd.PersistentFlags().BoolP("recursive", "r", true, "descend more than one level directories supplied as arguments")
	rootCmd.PersistentFlags().BoolP("dryRun", "d", false, "only print the renames")

	rootCmd.Flags().SortFlags = false
	viper.BindPFlags(rootCmd.Flags())
	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		switch flag.Name {
		case "config":
			return
		default:
			viper.BindPFlag(flag.Name, flag)
		}
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name "animenamer" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("animenamer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
	verbose.V = viper.GetBool("verbose")
}

func defaultString(s, d string) string {
	if s == "" {
		return d
	} else {
		return s
	}
}

func createSeriesDB() *seriesdb.SeriesDB {
	db := viper.GetString("db")
	var (
		provider series.Provider
		err      error
	)
	switch db {
	case series.ProviderTMDB:
		provider, err = tmdb.NewTMDB(defaultString(viper.GetString("apikey"), DefaultTMDBApiKey), viper.GetString("tmdb.absoluteGroupSeason"))
	case series.ProviderTVDB:
		provider, err = tvdb.NewTVDB(defaultString(viper.GetString("apikey"), DefaultTVDBApiKey))
	default:
		fmt.Printf("[E] unknown db type %s\n", db)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("[E] new provider: %v\n", err)
		os.Exit(1)
	}
	sdb, err := seriesdb.NewSeriesDB(provider, viper.GetString("language"), loadCustomConfigs()...)
	if err != nil {
		fmt.Printf("[E] new seriesdb error: %v\n", err)
		os.Exit(1)
	}
	return sdb
}

func rootCmdFunc(cmd *cobra.Command, args []string) {
	if viper.GetBool("verbose") {
		for k, v := range viper.AllSettings() {
			verbose.Printf("[V] %s: %v\n", k, v)
		}
	}
	regexpOnly := viper.GetBool("regexpOnly")
	var (
		sdb *seriesdb.SeriesDB
	)

	if !regexpOnly {
		sdb = createSeriesDB()
	}

	es := namer.EpisodeSearch{
		MediaExt:     namer.NewFileExtFromString(viper.GetString("mediaExt"), ","),
		SubtitlesExt: namer.NewFileExtFromString(viper.GetString("subtitleExt"), ","),
		SeriesDB:     sdb,
		SeriesName:   viper.GetString("name"),
		SeriesId:     viper.GetString("id"),
		RegexpOnly:   regexpOnly,
	}
	for _, p := range viper.GetStringSlice("pattern") {
		if err := es.AddPattern(p); err != nil {
			fmt.Printf("[E] parse pattern (%s) %v\n", p, err)
		}
	}
	if len(es.Filters) == 0 {
		fmt.Println("[E] no valid pattern")
		os.Exit(1)
	}

	recursive := viper.GetBool("recursive")
	format := viper.GetString("format")
	moveToDir := viper.GetString("moveToDir")
	renameSubtitle := viper.GetBool("renameSubtitle")
	replaceSpaceWith := viper.GetString("replaceSpace")
	dryRun := viper.GetBool("dryRun")
	logName := viper.GetString("log")
	for _, fp := range args {
		fmt.Printf("processing %s\n", fp)
		episodeFiles, err := es.ListEpisodeFile(fp, recursive)
		if err != nil {
			fmt.Printf("[E] %v\n", err)
			continue
		}
		fmt.Printf("found %d episode files\n", len(episodeFiles))
		var (
			logFile *os.File
		)
		if !dryRun {
			lfp := utils.MakeIncrementLog(fp, logName)
			if logFile, err = os.Create(lfp); err != nil {
				fmt.Printf("[E] create log file %s: %s", lfp, err)
				continue
			}
			logFile.WriteString(fmt.Sprintf("# rename files in %s\n", fp))
			defer logFile.Close()
		}
		for _, ef := range episodeFiles {
			renames := ef.Renames(format, replaceSpaceWith, renameSubtitle)
			for o, n := range renames {
				src := path.Join(ef.FileDir, o)
				var dstDir string
				if moveToDir != "" {
					dstDir = utils.NamedFormat(moveToDir, ef.Infos)
					if !filepath.IsAbs(dstDir) {
						dstDir = filepath.Join(fp, dstDir)
					}
				} else {
					dstDir = ef.FileDir
				}
				dst := path.Join(dstDir, n)
				if _, err := os.Stat(dst); err == nil {
					fmt.Printf("[E] %s exists, skipping rename %s\n", dst, src)
					continue
				}
				if dryRun {
					fmt.Printf("%s will rename(move) to %s\n", src, dst)
					continue
				}
				if err := os.MkdirAll(dstDir, 0755); err != nil {
					fmt.Printf("[E] mkdir %s error: %v\n", dstDir, err)
					continue
				}
				if err = os.Rename(src, dst); err == nil {
					fmt.Printf("%s has rename(move) to %s\n", src, dst)
					rel_src, err := filepath.Rel(fp, src)
					if err != nil {
						fmt.Printf("[E] %s\n", err)
					}
					rel_dst, err := filepath.Rel(fp, dst)
					if err != nil {
						fmt.Printf("[E] %s\n", err)
					}
					logFile.WriteString(fmt.Sprintf("R: '%s' -> '%s'\n", rel_src, rel_dst))
					logFile.Sync()
				} else {
					fmt.Printf("[E] rename %s to %s error: %v\n", src, n, err)
				}
			}
		}
	}
}

func loadCustomConfigs() []*series.CustomSeries {
	fps := viper.GetStringSlice("custom")
	if len(fps) == 0 {
		return nil
	}
	var customs []*series.CustomSeries
	for _, fp := range fps {
		c, e := series.LoadCustomSeries(fp)
		if e == nil {
			fmt.Printf("[I] use custom series info in %s\n", fp)
			customs = append(customs, c...)
		} else {
			fmt.Printf("[E] load custom series info in %s error, %v\n", fp, e)
		}
	}
	return customs
}
