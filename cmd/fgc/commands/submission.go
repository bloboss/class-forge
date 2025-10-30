package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewSubmissionCommand creates the submission command and its subcommands
func NewSubmissionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submission",
		Short: "Manage assignment submissions",
		Long:  "List, view, and download assignment submissions",
	}

	cmd.AddCommand(newSubmissionListCommand())
	cmd.AddCommand(newSubmissionViewCommand())
	cmd.AddCommand(newSubmissionDownloadCommand())

	return cmd
}

func newSubmissionListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [assignment-id]",
		Short: "List submissions for an assignment",
		Long:  "Display all submissions for the specified assignment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement submission listing
			fmt.Printf("Listing submissions for assignment: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().StringP("status", "s", "", "Filter by status (pending, accepted, late)")
	cmd.Flags().BoolP("team-only", "t", false, "Show only team submissions")
	cmd.Flags().BoolP("individual-only", "i", false, "Show only individual submissions")

	return cmd
}

func newSubmissionViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [submission-id]",
		Short: "View submission details",
		Long:  "Display detailed information about a specific submission",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement submission viewing
			fmt.Printf("Viewing submission: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().Bool("show-commits", false, "Show commit history")

	return cmd
}

func newSubmissionDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download [assignment-id]",
		Short: "Download submissions for grading",
		Long:  "Download all submissions for an assignment as zip archives",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement submission download
			fmt.Printf("Downloading submissions for assignment: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("output", "o", "./submissions", "Output directory")
	cmd.Flags().StringP("format", "f", "zip", "Archive format (zip, tar, tar.gz)")
	cmd.Flags().Bool("latest-only", false, "Download only the latest submission from each student/team")
	cmd.Flags().StringP("at-deadline", "d", "", "Download submissions as they were at the deadline")

	return cmd
}