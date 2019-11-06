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

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	configFormat string
)

// writeConfigCmd represents the writeConfig command
var writeConfigCmd = &cobra.Command{
	Use:   "writeConfig",
	Short: "Write default config to current directory",

	Run: func(cmd *cobra.Command, args []string) {

		if e := viper.WriteConfigAs("animenamer." + configFormat); e != nil {
			fmt.Println("error:", e)
		}
	},
}

func init() {
	rootCmd.AddCommand(writeConfigCmd)
	writeConfigCmd.Flags().StringVarP(&configFormat, "format", "f", "yaml", "config file format, yaml or json")
}
