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
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type ActivityStruct struct {
	IssueId    int
	ActivityId int
	StartedAt  time.Time
	StoppedAt  time.Time
	Comment    string
}

// Encode via Gob to file
func Save(path string, object interface{}) (err error) {
	file, err := os.Create(path)

	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}

	file.Close()

	return
}

// Decode Gob file
func Load(path string, object interface{}) (err error) {
	file, err := os.Open(path)

	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}

	file.Close()

	return
}

func startRun(cmd *cobra.Command, args []string) {
	var activities []ActivityStruct
	var err error

	loadPath, err := getGitRootPath()

	if err != nil {
		loadPath = filepath.Dir(os.Args[0])
	}

	activitiesPath := path.Join(loadPath, ".activities")

	_ = Load(activitiesPath, &activities)

	if len(activities) > 0 {
		lastActivity := activities[len(activities)-1]
		if lastActivity != (ActivityStruct{}) {
			if lastActivity.StoppedAt == (time.Time{}) {
				log.Fatal(errors.New("There's an activity running."))
			}
		}
	}

	issueId = getIssueId()

	activity := new(ActivityStruct)
	activity.IssueId = issueId
	activity.StartedAt = time.Now()

	activities = append(activities, *activity)
	err = Save(activitiesPath, activities)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Activity %d started.\n", issueId)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   startRun,
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
