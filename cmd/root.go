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
	"regexp"
	"strconv"

	"github.com/spf13/pflag"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/spf13/cobra"
	"github.com/tnextday/animenamer/pkg/namer"
	"github.com/tnextday/animenamer/pkg/tvdbex"
	"github.com/tnextday/animenamer/pkg/verbose"

	"github.com/spf13/viper"
)

var (
	cfgFile           string
	DefaultTvDbApiKey string
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
	rootCmd.PersistentFlags().String("custom", "animenamer.custom.yml",
		"custom series info file.\n"+
			"its support yaml or json.\n"+
			"(default is animenamer.custom.yml)")
	rootCmd.PersistentFlags().String("apikey", "", "the tvdb apikey")
	rootCmd.PersistentFlags().String("seriesId", "", "explicitly set the show id for TVdb to use (applies to all files)")
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
			"'series', 'seriesId', 'season', 'episode', 'absolute', 'date', 'title', 'ext'\n"+
			"and named variables in filename pattern matched,\n"+
			"you can use {variable.n} for number padding.\n")
	rootCmd.Flags().String("replaceSpace", "", "replace the whitespace with this value in new filename")
	rootCmd.Flags().StringP("log", "l", "rename", "the rename log name for recovery")
	lang, _ := jibber_jabber.DetectLanguage()
	rootCmd.PersistentFlags().String("language", lang, "preferred language")
	rootCmd.PersistentFlags().BoolP("recursive", "r", true, "descend more than one level directories supplied as arguments")
	rootCmd.Flags().BoolP("dryRun", "d", false, "only print the renames")

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

func rootCmdFunc(cmd *cobra.Command, args []string) {
	if viper.GetBool("verbose") {
		for k, v := range viper.AllSettings() {
			verbose.Printf("[V] %s: %v\n", k, v)
		}
	}
	apiKey := viper.GetString("apikey")
	if apiKey == "" {
		apiKey = DefaultTvDbApiKey
	}

	tvdb, err := tvdbex.NewTVDBEx(apiKey, viper.GetString("language"), loadCustomConfig())
	if err != nil {
		fmt.Printf("new tvdb error: %v\n", err)
		os.Exit(1)
	}

	es := namer.EpisodeSearch{
		MediaExt:     namer.NewFileExtFromString(viper.GetString("mediaExt"), ","),
		SubtitlesExt: namer.NewFileExtFromString(viper.GetString("subtitleExt"), ","),
		TVDB:         tvdb,
		SeriesName:   viper.GetString("name"),
		SeriesId:     viper.GetInt("seriesId"),
	}
	for _, p := range viper.GetStringSlice("pattern") {
		if err = es.AddPattern(p); err != nil {
			fmt.Printf("parse pattern (%s) error: %v\n", p, err)
		}
	}
	if len(es.Filters) == 0 {
		fmt.Println("no valid pattern")
		os.Exit(1)
	}

	recursive := viper.GetBool("recursive")
	format := viper.GetString("format")
	renameSubtitle := viper.GetBool("renameSubtitle")
	replaceSpaceWith := viper.GetString("replaceSpace")
	dryRun := viper.GetBool("dryRun")
	logName := viper.GetString("log")
	for _, fp := range args {
		fmt.Printf("processing %s\n", fp)
		episodeFiles, err := es.ListEpisodeFile(fp, recursive)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Printf("found %d episode files\n", len(episodeFiles))
		var logFile *os.File
		if !dryRun {
			lfp := makeLogName(fp, logName)
			if logFile, err = os.Create(lfp); err != nil {
				fmt.Printf("error create log file %s: %s", lfp, err)
				continue
			}
			logFile.WriteString(fmt.Sprintf("# rename files in %s\n", fp))
			defer logFile.Close()
		}
		for _, ef := range episodeFiles {
			renames := ef.Renames(format, replaceSpaceWith, renameSubtitle)
			if dryRun {
				for o, n := range renames {
					fmt.Printf("%s will be rename to %s\n", o, n)
				}
			} else {
				for o, n := range renames {
					src := path.Join(ef.FileDir, o)
					dst := path.Join(ef.FileDir, n)
					if err := os.Rename(src, dst); err == nil {
						fmt.Printf("%s has rename to %s\n", o, n)
						logFile.WriteString(fmt.Sprintf("R: '%s' -> '%s'\n", src, dst))
						logFile.Sync()
					} else {
						fmt.Printf("rename %s to %s error: %v\n", o, n, err)
					}
				}
			}
		}
	}
}

func loadCustomConfig() *tvdbex.CustomSeries {
	fp := viper.GetString("custom")
	if fp == "" {
		return nil
	}
	c, e := tvdbex.LoadCustomSeries(fp)
	if e == nil {
		fmt.Printf("[I] use custom series info in %s\n", fp)
		return c
	} else {
		fmt.Printf("[E] load custom series info in %s error, %v\n", fp, e)
		return nil
	}
}

func makeLogName(dir, name string) string {
	logIndex := 0
	iRe := regexp.MustCompile(fmt.Sprintf(`%s\.(\d+)\.log`, name))
	if files, err := filepath.Glob(fmt.Sprintf("%s/%s.*.log", dir, name)); err != nil {
		for _, f := range files {
			_, fn := filepath.Split(f)
			matchs := iRe.FindStringSubmatch(fn)

			if len(matchs) != 2 {
				continue
			}
			if i, err := strconv.Atoi(matchs[1]); err == nil {
				if i > logIndex {
					logIndex = i
				}
			}
		}
	}
	return path.Join(dir, fmt.Sprintf("%s.%d.log", name, logIndex))
}
