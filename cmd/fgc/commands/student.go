package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewStudentCommand creates the student command and its subcommands
func NewStudentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "student",
		Short: "Student operations",
		Long:  "Commands for students to interact with assignments",
	}

	cmd.AddCommand(newStudentAcceptCommand())

	return cmd
}

func newStudentAcceptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept [assignment-id]",
		Short: "Accept an assignment",
		Long:  "Accept an assignment and create a repository for submission",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement assignment acceptance
			fmt.Printf("Accepting assignment: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("team", "t", "", "Join or create a team (for team assignments)")

	return cmd
}