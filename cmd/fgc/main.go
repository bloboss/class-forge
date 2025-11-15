package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"code.forgejo.org/forgejo/classroom/cmd/fgc/commands"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "fgc",
		Short: "Forgejo Classroom CLI",
		Long: `Forgejo Classroom is an educational assignment management system that integrates
with Forgejo to provide GitHub Classroom-like functionality for self-hosted Git platforms.

This CLI tool allows you to manage classrooms, assignments, rosters, submissions, and teams.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Global flags
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.fgc.yaml)")
	rootCmd.PersistentFlags().String("server", "", "Forgejo server URL")
	rootCmd.PersistentFlags().String("token", "", "Forgejo API token")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dry-run", false, "show what would be done without executing")

	// Bind flags to viper
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))

	// Add subcommands
	rootCmd.AddCommand(commands.NewClassroomCommand())
	rootCmd.AddCommand(commands.NewAssignmentCommand())
	rootCmd.AddCommand(commands.NewRosterCommand())
	rootCmd.AddCommand(commands.NewSubmissionCommand())
	rootCmd.AddCommand(commands.NewTeamCommand())
	rootCmd.AddCommand(commands.NewStudentCommand())

	// Initialize configuration
	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fgc" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fgc")
	}

	// Environment variables
	viper.SetEnvPrefix("FGC")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	}
}
