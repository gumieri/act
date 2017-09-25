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
	"time"

	"github.com/gumieri/act/lib/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func stopRun(cmd *cobra.Command, args []string) {
	var activities []ActivityStruct
	var loadPath string
	var err error

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
		log.Fatal(err)
	}

	if len(activities) == 0 {
		log.Fatal(errors.New("there's no activity started"))
	}

	lastActivity := activities[len(activities)-1]
	if lastActivity == (ActivityStruct{}) {
		log.Fatal(errors.New("there's no activity started"))
	}

	if lastActivity.StoppedAt != (time.Time{}) {
		log.Fatal(errors.New("there's no activity started"))
	}

	lastActivity.StoppedAt = time.Now()

	activities[len(activities)-1] = lastActivity

	err = Save(activitiesPath, activities)

	if err != nil {
		log.Fatal(err)
	}

	duration := lastActivity.StoppedAt.Sub(lastActivity.StartedAt)
	durationHour := float64(duration) / float64(time.Hour)
	log.Printf("Activity %d stopped. Time elapsed %.2f (%s)\n", lastActivity.IssueID, durationHour, duration)
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop tracking a open activity by saving the time of it",
	Long:  ``,
	Run:   stopRun,
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
