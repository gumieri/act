package cmd

import (
	"errors"
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"srv-gitlab.tecnospeed.local/labs/act/lib/git"
)

var homePath string
var cfgFile string
var issueID int

func getIssueID() int {
	if issueID == 0 {
		gitPath := viper.Get("git.path")
		gitRegex := viper.Get("git.regex")

		if gitPath != nil && gitRegex != nil {
			issueID, _ = git.IssueIDFromBranch(gitPath.(string), gitRegex.(string))
		}
	}

	if issueID == 0 {
		log.Fatal(errors.New("issue_id (-i) is missing"))
	}

	return issueID
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "act",
	Short: "act - Activity Continuous Tracking",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.act.yaml)")
	RootCmd.PersistentFlags().IntVarP(&issueID, "issue_id", "i", 0, "The Issue ID.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	// Find home directory.
	homePath, err = homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".act" (without extension).
		viper.AddConfigPath(homePath)
		viper.SetConfigName(".act")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()
}
