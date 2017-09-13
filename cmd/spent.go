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
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TimeEntryStruct struct {
	IssueId    int     `json:"issue_id"`
	Date       string  `json:"spent_on"`
	Time       float64 `json:"hours"`
	ActivityId int     `json:"activity_id"`
	Comment    string  `json:"comments"`
}

type PayloadStruct struct {
	TimeEntry TimeEntryStruct `json:"time_entry"`
}

var timeEntry TimeEntryStruct

func typeOnEditor(editorCommand string) (text string, err error) {
	filePath := fmt.Sprintf("%s/%d-comment", os.TempDir(), timeEntry.IssueId)

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	vimcmd := exec.Command(editorCommand, filePath)
	vimcmd.Stdin = os.Stdin
	vimcmd.Stdout = os.Stdout
	vimcmd.Stderr = os.Stderr

	err = vimcmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = vimcmd.Wait()
	if err != nil {

	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	text = string(content)

	return
}

func getIssueFromGitBranch(gitPath string) (issueId int, err error) {
	out, _ := exec.Command(gitPath, "rev-parse", "--abbrev-ref", "HEAD").Output()

	regexC, err := regexp.Compile("[0-9]*")

	if err != nil {
		return
	}

	issueId, err = strconv.Atoi(regexC.FindString(string(out)))

	return
}

func spentRun(cmd *cobra.Command, args []string) {
	// Setting values who require viper loaded
	if timeEntry.IssueId == 0 {
		gitPath := viper.Get("git.path")
		if gitPath != nil {
			timeEntry.IssueId, _ = getIssueFromGitBranch(gitPath.(string))
		}
	}

	// Validating IssueId
	if timeEntry.IssueId == 0 {
		log.Fatal(errors.New("issue_id (-i) is missing."))
	}

	var err error

	editor := viper.Get("editor")
	if editor != nil && timeEntry.Comment == "" {
		timeEntry.Comment, err = typeOnEditor(editor.(string))
		if err != nil {
			log.Fatal(err)
		}
	}

	if timeEntry.ActivityId == 0 {
		timeEntry.ActivityId = viper.GetInt("default.activity_id")
	}

	// Validating ActivityId
	if timeEntry.ActivityId == 0 {
		log.Fatal(errors.New("activity_id is missing."))
	}

	// Setting the time informed (the first arg)
	timeEntry.Time, err = strconv.ParseFloat(args[0], 64)

	if err != nil {
		log.Fatal(err)
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

	log.Printf("Added %.2f hour(s) to the Issue #%d.", timeEntry.Time, timeEntry.IssueId)
	fmt.Println()
}

// spentCmd represents the spent command
var spentCmd = &cobra.Command{
	Use:   "spent",
	Short: "Update an Issue defining the time spent on it.",
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

	spentCmd.Flags().IntVarP(&timeEntry.IssueId, "issue_id", "i", 0, "The Issue ID.")
	spentCmd.Flags().IntVar(&timeEntry.ActivityId, "activity_id", 0, "The Activity ID.")

	current_date := time.Now().Local().Format("2006-01-02")
	spentCmd.Flags().StringVarP(&timeEntry.Date, "date", "d", current_date, "The date when the time was spent on.")
	spentCmd.Flags().StringVarP(&timeEntry.Comment, "comment", "m", "", "A short description of what was done.")
}
