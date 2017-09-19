package cmd

import (
	"errors"
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"srv-gitlab.tecnospeed.local/rafael.gumieri/act/lib/git"
)

var cfgFile string

var issueId int

func getIssueId() int {
	if issueId == 0 {
		gitPath := viper.Get("git.path")
		gitRegex := viper.Get("git.regex")

		if gitPath != nil && gitRegex != nil {
			issueId, _ = git.IssueIdFromBranch(gitPath.(string), gitRegex.(string))
		}
	}

	if issueId == 0 {
		log.Fatal(errors.New("issue_id (-i) is missing."))
	}

	return issueId
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
	RootCmd.PersistentFlags().IntVarP(&issueId, "issue_id", "i", 0, "The Issue ID.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".act" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".act")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()
}
