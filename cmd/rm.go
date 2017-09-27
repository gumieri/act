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
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gumieri/act/lib/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func rmRun(cmd *cobra.Command, args []string) {
	var activities []ActivityStruct
	var loadPath string
	var err error

	index, err := strconv.Atoi(args[0])

	if err != nil {
		log.Fatal(errors.New("The argument must be a valid index number from status list"))
	}

	gitPath := viper.Get("git.path")
	if gitPath != nil {
		loadPath, _ = git.TopLevelPath(gitPath.(string))
	}

	if loadPath == "" {
		loadPath = filepath.Dir(os.Args[0])
	}

	activitiesPath := path.Join(loadPath, ".activities")

	err = Load(activitiesPath, &activities)

	if err != nil {
		return
	}

	if len(activities) == 0 {
		return
	}

	if index > (len(activities) - 1) {
		log.Fatal(errors.New("The argument must be a valid index number from status list"))
	}

	activities = append(activities[:index], activities[index+1:]...)

	err = Save(activitiesPath, activities)

	if err != nil {
		log.Fatal(err)
	}
}

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove activiries without push it",
	Long:  `Use the status command to use the {index} as reference to be removed by this command.`,
	Args:  cobra.MinimumNArgs(1),
	Run:   rmRun,
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
