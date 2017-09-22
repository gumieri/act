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
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"srv-gitlab.tecnospeed.local/labs/act/lib/editor"
)

type TimeEntryStruct struct {
	IssueID    int     `json:"issue_id"`
	Date       string  `json:"spent_on"`
	Time       float64 `json:"hours"`
	ActivityID int     `json:"activity_id"`
	Comment    string  `json:"comments"`
}

type PayloadStruct struct {
	TimeEntry TimeEntryStruct `json:"time_entry"`
}

var timeEntry TimeEntryStruct

func spentRun(cmd *cobra.Command, args []string) {
	var err error

	timeEntry.IssueID = getIssueID()

	// Setting the time informed (the first arg)
	timeEntry.Time, err = strconv.ParseFloat(args[0], 64)

	if err != nil {
		log.Fatal(err)
	}

	editorPath := viper.Get("editor")
	if editorPath != nil && timeEntry.Comment == "" {
		fileName := fmt.Sprintf("%d-comment", timeEntry.IssueID)

		helperText := fmt.Sprintf("\n\n# Issue #%d\n# Date: %s\n# Time elapsed: %.2f\n# Activity ID: %d", timeEntry.IssueID, timeEntry.Date, timeEntry.Time, timeEntry.ActivityID)

		timeEntry.Comment, err = editor.Open(editorPath.(string), fileName, helperText, true)
		if err != nil {
			log.Fatal(err)
		}
	}

	if timeEntry.Comment == "" {
		log.Fatal(errors.New("Empty note"))
	}

	if timeEntry.ActivityID == 0 {
		timeEntry.ActivityID = viper.GetInt("default.activity_id")
	}

	// Validating ActivityID
	if timeEntry.ActivityID == 0 {
		log.Fatal(errors.New("activity_id is missing"))
	}

	// Sending the data to the Redmine
	payload := new(PayloadStruct)
	payload.TimeEntry = timeEntry

	marshal, err := json.Marshal(payload)

	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://%s/time_entries.json", viper.Get("redmine.url"))
	payloadMarshal := bytes.NewBuffer(marshal)
	request, err := http.NewRequest(http.MethodPost, url, payloadMarshal)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("X-Redmine-API-Key", viper.GetString("redmine.access_key"))
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		bodyBytes, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(response.Status, "\n", string(bodyBytes))
	}

	log.Printf("Added %.2f hour(s) to the Issue #%d.", timeEntry.Time, timeEntry.IssueID)
}

// spentCmd represents the spent command
var spentCmd = &cobra.Command{
	Use:   "spent",
	Short: "Update an Issue defining the time spent on it",
	Long: `Update the Issue with the informed hours spent. The hours can be integer (ex: act spent 1) or floating point (ex: act spent 6.66).

The Activity ID can be configured with a default value (default.activity_id).

If the Date (-d) is not informed, it will use the current date.

The Issue ID can be ommited if using a regex to retrieve it from the git branch.
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  spentRun,
}

func init() {
	RootCmd.AddCommand(spentCmd)

	spentCmd.Flags().IntVar(&timeEntry.ActivityID, "activity_id", 0, "The Activity ID.")

	currentDate := time.Now().Local().Format("2006-01-02")
	spentCmd.Flags().StringVarP(&timeEntry.Date, "date", "d", currentDate, "The date when the time was spent on.")
	spentCmd.Flags().StringVarP(&timeEntry.Comment, "comment", "m", "", "A short description of what was done.")
}
