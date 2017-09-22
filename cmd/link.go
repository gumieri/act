// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func linkRun(cmd *cobra.Command, args []string) {
	log.Printf("http://%s/issues/%d", viper.Get("redmine.url"), getIssueID())
}

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Show a link to the Issue on redmine's page",
	Long:  `Show a link to the Issue on redmine's page`,
	Run:   linkRun,
}

func init() {
	RootCmd.AddCommand(linkCmd)
}
