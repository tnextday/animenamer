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
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"github.com/tnextday/animenamer/pkg/verbose"

	"github.com/spf13/cobra"
)

var recoveryCmd = &cobra.Command{
	Use:   "recovery [flags] <log-files>",
	Short: "Recovery renamed files use log file",
	Args:  cobra.MinimumNArgs(1),
	Run:   recoveryRun,
}
var logRe = regexp.MustCompile(`^R:\s*'(.+)'\s*->\s*'(.+)'$`)

func init() {
	rootCmd.AddCommand(recoveryCmd)
	viper.BindPFlags(recoveryCmd.Flags())
}

func recoveryRun(cmd *cobra.Command, args []string) {
	if viper.GetBool("verbose") {
		for k, v := range viper.AllSettings() {
			verbose.Printf("[V] %s: %v\n", k, v)
		}
	}

	dryRun := viper.GetBool("dryRun")
	for _, logFile := range args {
		readFile, err := os.Open(logFile)
		if err != nil {
			fmt.Printf("[E] %v\n", err)
			continue
		}
		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)

		baseDir, _ := path.Split(logFile)

		for fileScanner.Scan() {
			matches := logRe.FindStringSubmatch(strings.TrimSpace(fileScanner.Text()))
			if len(matches) != 3 {
				continue
			}
			dst := path.Join(baseDir, matches[1])
			src := path.Join(baseDir, matches[2])
			if !dryRun {
				fmt.Printf("recovery %s to %s\n", src, dst)
				if _, err := os.Stat(dst); err == nil {
					fmt.Printf("[E] %s exists, skipping recovery %s\n", dst, src)
				}
				if err := os.Rename(src, dst); err != nil {
					fmt.Printf("[E] %v\n", err)
				}
			} else {
				fmt.Printf("[dry] recovery %s to %s\n", src, dst)
			}
		}
		readFile.Close()
	}
}
