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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"srv-gitlab.tecnospeed.local/rafael.gumieri/act/lib/editor"
)

type IssuePayloadStruct struct {
	Issue IssueStruct `json:"issue"`
}

type IssueStruct struct {
	Note string `json:"notes"`
}

func noteRun(cmd *cobra.Command, args []string) {
	var err error

	issueId := getIssueId()
	note := args[0]

	editorPath := viper.Get("editor")
	if editorPath != nil && note == "" {
		fileName := fmt.Sprintf("%d-note", issueId)

		note, err = editor.Open(editorPath.(string), fileName, "")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Sending the data to the Redmine
	payload := new(IssuePayloadStruct)
	payload.Issue.Note = note

	marshal, err := json.Marshal(payload)

	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://%s/issues/%d.json", viper.Get("redmine.url"), issueId)
	payloadMarshal := bytes.NewBuffer(marshal)
	request, err := http.NewRequest(http.MethodPut, url, payloadMarshal)

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

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(response.Status, "\n", string(bodyBytes))
	}

	log.Printf("Added the note to the Issue #%d.", issueId)
}

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Add a note to the Issue",
	Long: `The informed argument is sent as note to the Issue.

The Issue ID can be ommited if using a regex to retrieve it from the git branch.
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  noteRun,
}

func init() {
	RootCmd.AddCommand(noteCmd)
}