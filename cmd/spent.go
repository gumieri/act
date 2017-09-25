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
	"regexp"
	"strconv"
	"time"

	"github.com/gumieri/act/lib/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TimeEntryStruct is expected format of the Redmine API
type TimeEntryStruct struct {
	IssueID    int     `json:"issue_id"`
	Date       string  `json:"spent_on"`
	Time       float64 `json:"hours"`
	ActivityID int     `json:"activity_id"`
	Comment    string  `json:"comments"`
}

// PayloadStruct is the envelope to the time_entry expected by the Redmine API
type PayloadStruct struct {
	TimeEntry TimeEntryStruct `json:"time_entry"`
}

var timeEntry TimeEntryStruct
var activityAlias string

func parseMonthDay(input string) (output string, err error) {
	regexDayAndMonth := regexp.MustCompile(`^(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`)
	dayAndMonth := regexDayAndMonth.FindStringSubmatch(input)

	if len(dayAndMonth) == 0 {
		return
	}

	day, err := strconv.Atoi(dayAndMonth[2])

	if err != nil {
		return
	}

	month, err := strconv.Atoi(dayAndMonth[1])

	if err != nil {
		return
	}

	timeNow := time.Now().Local()
	monthDayDate := time.Date(timeNow.Year(), time.Month(month), day, 0, 0, 0, 0, timeNow.Location())
	output = monthDayDate.Format("2006-01-02")
	return
}

func parseRetroactiveDate(input string) (output string, err error) {
	retroactive, err := regexp.MatchString("^-[0-9]*$", input)

	if err != nil || !retroactive {
		return
	}

	daysToBack, err := strconv.Atoi(input)

	if err != nil {
		return
	}

	timeNow := time.Now().Local()
	output = timeNow.AddDate(0, 0, daysToBack).Format("2006-01-02")
	return

}

func parseDate(input string) (output string, err error) {
	complete, err := regexp.MatchString(`^([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))$`, input)

	if err != nil {
		return
	}

	if complete {
		output = input
		return
	}

	monthDay, err := parseMonthDay(input)

	if err != nil {
		return
	}

	if monthDay != "" {
		output = monthDay
		return
	}

	retroactiveDate, err := parseRetroactiveDate(input)

	if err != nil {
		return
	}

	if retroactiveDate != "" {
		output = retroactiveDate
		return
	}

	// look for a number (only the day). ex: 2
	day, err := strconv.Atoi(input)

	if err != nil {
		return
	}

	timeNow := time.Now().Local()
	output = time.Date(timeNow.Year(), timeNow.Month(), day, 0, 0, 0, 0, timeNow.Location()).Format("2006-01-02")
	return
}

func spentRun(cmd *cobra.Command, args []string) {
	var err error

	timeEntry.IssueID = getIssueID()

	// Setting the time informed (the first arg)
	timeEntry.Time, err = strconv.ParseFloat(args[0], 64)

	if err != nil {
		log.Fatal(err)
	}

	timeEntry.Date, err = parseDate(timeEntry.Date)

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

	if activityAlias != "" {
		timeEntry.ActivityID = viper.GetInt(fmt.Sprintf("activities.%s", activityAlias))
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

The Date can be informed as:
-d 2017-09-22 -- Complete
-d      09-22 -- Only the month and day. The year will be the current one.
-d         22 -- Only the day. The year and month will be the current ones.
-d         -1 -- Informing how many days back from the current date.
And if not informed, it will use the current date.

The Issue ID can be omitted if using a regex to retrieve it from the git branch.
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  spentRun,
}

func init() {
	RootCmd.AddCommand(spentCmd)

	spentCmd.Flags().IntVar(&timeEntry.ActivityID, "activity_id", 0, "The activity ID.")
	spentCmd.Flags().StringVarP(&activityAlias, "activity", "a", "", "The activity alias (alternative to activity_id).")

	currentDate := time.Now().Local().Format("2006-01-02")
	spentCmd.Flags().StringVarP(&timeEntry.Date, "date", "d", currentDate, "The date when the time was spent on.")
	spentCmd.Flags().StringVarP(&timeEntry.Comment, "comment", "m", "", "A short description of what was done.")
}
