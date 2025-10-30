package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewAssignmentCommand creates the assignment command and its subcommands
func NewAssignmentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assignment",
		Short: "Manage assignments",
		Long:  "Create, list, view, update, delete assignments and view statistics",
	}

	cmd.AddCommand(newAssignmentCreateCommand())
	cmd.AddCommand(newAssignmentListCommand())
	cmd.AddCommand(newAssignmentViewCommand())
	cmd.AddCommand(newAssignmentUpdateCommand())
	cmd.AddCommand(newAssignmentDeleteCommand())
	cmd.AddCommand(newAssignmentStatsCommand())

	return cmd
}

func newAssignmentCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new assignment",
		Long:  "Create a new assignment from a template repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement assignment creation
			fmt.Printf("Creating assignment: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("classroom", "c", "", "Classroom ID (required)")
	cmd.Flags().StringP("template", "t", "", "Template repository URL (required)")
	cmd.Flags().StringP("deadline", "d", "", "Assignment deadline (RFC3339 format)")
	cmd.Flags().StringP("description", "D", "", "Assignment description")
	cmd.Flags().IntP("max-teams", "m", 1, "Maximum team size (1 for individual)")
	cmd.Flags().Bool("auto-accept", false, "Automatically accept submissions")
	cmd.MarkFlagRequired("classroom")
	cmd.MarkFlagRequired("template")

	return cmd
}

func newAssignmentListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List assignments",
		Long:  "List assignments in a classroom or across all classrooms",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement assignment listing
			fmt.Println("Listing assignments...")
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("classroom", "c", "", "Filter by classroom ID")
	cmd.Flags().BoolP("active", "a", false, "Show only active assignments")
	cmd.Flags().BoolP("past", "p", false, "Show only past assignments")
	cmd.Flags().StringP("format", "f", "table", "Output format (table, json, yaml)")

	return cmd
}

func newAssignmentViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [id]",
		Short: "View assignment details",
		Long:  "Display detailed information about a specific assignment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement assignment viewing
			fmt.Printf("Viewing assignment: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("format", "f", "table", "Output format (table, json, yaml)")

	return cmd
}

func newAssignmentUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Update assignment settings",
		Long:  "Update the settings of an existing assignment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement assignment updating
			fmt.Printf("Updating assignment: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "New assignment name")
	cmd.Flags().StringP("deadline", "d", "", "New deadline (RFC3339 format)")
	cmd.Flags().StringP("description", "D", "", "New assignment description")
	cmd.Flags().IntP("max-teams", "m", 0, "New maximum team size")
	cmd.Flags().Bool("auto-accept", false, "Enable auto-accept submissions")
	cmd.Flags().Bool("no-auto-accept", false, "Disable auto-accept submissions")

	return cmd
}

func newAssignmentDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete an assignment",
		Long:  "Permanently delete an assignment and all its submissions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement assignment deletion
			fmt.Printf("Deleting assignment: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")

	return cmd
}

func newAssignmentStatsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats [id]",
		Short: "View assignment statistics",
		Long:  "Display statistics for assignment submissions and progress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement assignment statistics
			fmt.Printf("Viewing stats for assignment: %s\n", args[0])
			fmt.Println("Not yet implemented")
			return nil
		},
	}

	cmd.Flags().StringP("format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().Bool("detailed", false, "Show detailed statistics")

	return cmd
}