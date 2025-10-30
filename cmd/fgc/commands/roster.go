package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewRosterCommand creates the roster command and its subcommands
func NewRosterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roster",
		Short: "Manage classroom rosters",
		Long:  "Add students to classrooms, link GitHub accounts, and import rosters",
	}

	cmd.AddCommand(newRosterAddCommand())
	cmd.AddCommand(newRosterListCommand())
	cmd.AddCommand(newRosterLinkCommand())
	cmd.AddCommand(newRosterImportCommand())

	return cmd
}

func newRosterAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [classroom-id] [student-identifier]",
		Short: "Add a student to a classroom roster",
		Long:  "Add a student to a classroom roster using their email or username",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement roster addition
			fmt.Printf("Adding student %s to classroom %s\n", args[1], args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Student's display name")
	cmd.Flags().StringP("email", "e", "", "Student's email address")
	cmd.Flags().StringP("role", "r", "student", "Student role (student, assistant, instructor)")

	return cmd
}

func newRosterListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [classroom-id]",
		Short: "List students in a classroom roster",
		Long:  "Display all students enrolled in the specified classroom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement roster listing
			fmt.Printf("Listing roster for classroom: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("format", "f", "table", "Output format (table, json, yaml, csv)")
	cmd.Flags().BoolP("linked-only", "l", false, "Show only students with linked accounts")
	cmd.Flags().BoolP("unlinked-only", "u", false, "Show only students without linked accounts")

	return cmd
}

func newRosterLinkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link [classroom-id] [student-identifier] [forgejo-username]",
		Short: "Link a student to their Forgejo account",
		Long:  "Associate a roster entry with a Forgejo username",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement account linking
			fmt.Printf("Linking student %s to Forgejo account %s in classroom %s\n", args[1], args[2], args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force linking even if account is already linked")

	return cmd
}

func newRosterImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [classroom-id] [file]",
		Short: "Import roster from CSV file",
		Long:  "Bulk import students from a CSV file with columns: name, email, identifier",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement roster import
			fmt.Printf("Importing roster from %s to classroom %s\n", args[1], args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().Bool("dry-run", false, "Show what would be imported without actually importing")
	cmd.Flags().Bool("update", false, "Update existing students instead of skipping")
	cmd.Flags().StringP("delimiter", "d", ",", "CSV delimiter character")

	return cmd
}