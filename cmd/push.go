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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	editor "srv-gitlab.tecnospeed.local/rafael.gumieri/act/lib/editor"
)

func commit(activity ActivityStruct) (timeEntry TimeEntryStruct, err error) {
	duration := activity.StoppedAt.Sub(activity.StartedAt)
	durationHour := float64(duration) / float64(time.Hour)

	timeEntry.IssueId = activity.IssueId
	timeEntry.ActivityId = activity.ActivityId
	timeEntry.Date = activity.StartedAt.Format("2006-01-02")
	timeEntry.Time = durationHour
	timeEntry.Comment = activity.Comment

	editorPath := viper.Get("editor")
	if editorPath != nil && timeEntry.Comment == "" {
		fileName := fmt.Sprintf("%d-comment", timeEntry.IssueId)

		helperText := fmt.Sprintf("\n\n# Issue #%d\n# Date: %s\n# Time elapsed: %.2f\n# Activity ID: %d", timeEntry.IssueId, timeEntry.Date, timeEntry.Time, timeEntry.ActivityId)

		timeEntry.Comment, err = editor.Open(editorPath.(string), fileName, helperText)

		if err != nil {
			return
		}
	}

	if timeEntry.ActivityId == 0 {
		timeEntry.ActivityId = viper.GetInt("default.activity_id")
	}

	// Validating ActivityId
	if timeEntry.ActivityId == 0 {
		err = errors.New("activity_id is missing.")
		return
	}

	// Validating ActivityId
	if strings.Trim(timeEntry.Comment, "\n ") == "" {
		err = errors.New("You must inform a comment/description to the activity.")
		return
	}

	// Sending the data to the Redmine
	payload := new(PayloadStruct)
	payload.TimeEntry = timeEntry

	marshal, err := json.Marshal(payload)

	if err != nil {
		return
	}

	url := fmt.Sprintf("http://%s/time_entries.json", viper.Get("redmine.url"))
	payloadMarshal := bytes.NewBuffer(marshal)
	request, err := http.NewRequest(http.MethodPost, url, payloadMarshal)

	if err != nil {
		return
	}

	request.Header.Add("X-Redmine-API-Key", viper.GetString("redmine.access_key"))
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		err = errors.New(fmt.Sprintf("%d", response.StatusCode))
	}

	return
}

func pushRun(cmd *cobra.Command, args []string) {
	var activities []ActivityStruct
	var err error

	loadPath, err := getGitRootPath()

	if err != nil {
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

	for index, activity := range activities {
		timeEntry, err := commit(activity)

		if err != nil {
			log.Fatal(err)
		}

		activities = append(activities[:index], activities[index+1:]...)

		err = Save(activitiesPath, activities)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Added %.2f hour(s) to the Issue #%d.", timeEntry.Time, timeEntry.IssueId)
	}
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   pushRun,
}

func init() {
	RootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
