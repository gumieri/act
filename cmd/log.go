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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RedmineGetTimeEntriesStruct struct {
	TimeEntries []RedmineTimeEntryStruct `json:"time_entries"`
}

type RedmineTimeEntryStruct struct {
	ID        int                   `json:"id"`
	Time      float64               `json:"hours"`
	Comment   string                `json:"comments"`
	Date      string                `json:"spent_on"`
	CreatedOn string                `json:"created_on"`
	UpdatedOn string                `json:"updated_on"`
	Project   RedmineProjectSctruct `json:"project"`
	Issue     RedmineIssueStruct    `json:"issue"`
	User      RedmineUserStruct     `json:"user"`
	Activity  RedmineActivityStruct `json:"activity"`
}

type RedmineProjectSctruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RedmineIssueStruct struct {
	ID int `json:"id"`
}

type RedmineUserStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RedmineActivityStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func logRun(cmd *cobra.Command, args []string) {
	issueID = getIssueID()

	url := fmt.Sprintf("http://%s/issues/%d/time_entries.json", viper.Get("redmine.url"), issueID)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("X-Redmine-API-Key", viper.GetString("redmine.access_key"))

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	payload := new(RedmineGetTimeEntriesStruct)
	err = json.NewDecoder(response.Body).Decode(payload)

	if err != nil {
		log.Fatal(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, "Date\tTime\tUser Name\tAct. ID\tActivity Name\tComment")
	for _, timeEntry := range payload.TimeEntries {
		fmt.Fprintf(w, "%s", timeEntry.Date)
		fmt.Fprint(w, "\t")
		fmt.Fprintf(w, "%.2f", timeEntry.Time)
		fmt.Fprint(w, "\t")
		fmt.Fprintf(w, "%s", timeEntry.User.Name)
		fmt.Fprint(w, "\t")
		fmt.Fprintf(w, "%d", timeEntry.Activity.ID)
		fmt.Fprint(w, "\t")
		fmt.Fprintf(w, "%s", timeEntry.Activity.Name)
		fmt.Fprintf(w, "%q", timeEntry.Comment)
		fmt.Fprintln(w)
	}
	w.Flush()

}

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show the appointments of a specific issue",
	Long:  ``,
	Run:   logRun,
}

func init() {
	RootCmd.AddCommand(logCmd)
}
