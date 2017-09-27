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
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/gumieri/act/lib/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func statusRun(cmd *cobra.Command, args []string) {
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

	if activities[len(activities)-1] == (ActivityStruct{}) {
		log.Fatal(errors.New("there's no activity"))
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprint(w, "\t")
	fmt.Fprint(w, "Issue")
	fmt.Fprint(w, "\t")
	fmt.Fprint(w, "Started At")
	fmt.Fprint(w, "\t")
	fmt.Fprint(w, "Stopped At")
	fmt.Fprint(w, "\t")
	fmt.Fprint(w, "Spent")
	fmt.Fprint(w, "\t")
	if logAll {
		fmt.Fprint(w, "Activity ID")
		fmt.Fprint(w, "\t")
	}
	fmt.Fprint(w, "Comment")
	fmt.Fprintln(w)

	for i, activity := range activities {

		running := activity.StoppedAt == (time.Time{})

		fmt.Fprintf(w, "{%d}", i)
		fmt.Fprint(w, "\t")

		fmt.Fprintf(w, "#%d", activity.IssueID)
		fmt.Fprint(w, "\t")

		fmt.Fprintf(w, "%s", activity.StartedAt.Format("3:04PM"))
		fmt.Fprint(w, "\t")

		if running {
			fmt.Fprint(w, "-")
		} else {
			fmt.Fprintf(w, "%s", activity.StoppedAt.Format("3:04PM"))
		}
		fmt.Fprint(w, "\t")

		if running {
			fmt.Fprintf(w, "%s", time.Now().Sub(activity.StartedAt))
		} else {
			fmt.Fprintf(w, "%s", activity.StoppedAt.Sub(activity.StartedAt))
		}
		fmt.Fprint(w, "\t")

		if logAll {
			fmt.Fprintf(w, "%d", activity.ActivityID)
			fmt.Fprint(w, "\t")
		}

		fmt.Fprintf(w, "%q", activity.Comment)
		fmt.Fprintln(w)
	}

	w.Flush()
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   statusRun,
}

func init() {
	RootCmd.AddCommand(statusCmd)

	statusCmd.Flags().BoolVarP(&logAll, "all", "a", false, "To list complete information.")
}
