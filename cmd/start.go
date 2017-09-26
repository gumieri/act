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
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gumieri/act/lib/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ActivityStruct struct {
	IssueID    int
	ActivityID int
	StartedAt  time.Time
	StoppedAt  time.Time
	Comment    string
}

// Save encode via Gob to file
func Save(path string, object interface{}) (err error) {
	file, err := os.Create(path)

	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}

	file.Close()

	return
}

// Load decode Gob file
func Load(path string, object interface{}) (err error) {
	file, err := os.Open(path)

	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}

	file.Close()

	return
}

var activityToStart ActivityStruct

func startRun(cmd *cobra.Command, args []string) {
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

	_ = Load(activitiesPath, &activities)

	if len(activities) > 0 {
		lastActivity := activities[len(activities)-1]
		if lastActivity != (ActivityStruct{}) {
			if lastActivity.StoppedAt == (time.Time{}) {
				log.Fatal(errors.New("there's an activity running"))
			}
		}
	}

	if activityAlias != "" {
		activityToStart.ActivityID = viper.GetInt(fmt.Sprintf("activities.%s", activityAlias))
	}

	if activityToStart.ActivityID == 0 {
		activityToStart.ActivityID = viper.GetInt("default.activity_id")
	}

	// Validating ActivityID
	if activityToStart.ActivityID == 0 {
		log.Fatal(errors.New("activity_id is missing"))
	}

	activityToStart.IssueID = getIssueID()
	activityToStart.StartedAt = time.Now()

	activities = append(activities, activityToStart)
	err = Save(activitiesPath, activities)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Activity %d started.\n", activityToStart.IssueID)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tracking an activity by saving the time of it",
	Long:  ``,
	Run:   startRun,
}

func init() {
	RootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVar(&activityToStart.ActivityID, "activity_id", 0, "The activity ID.")
	startCmd.Flags().StringVarP(&activityAlias, "activity", "a", "", "The activity alias (alternative to activity_id).")
}
