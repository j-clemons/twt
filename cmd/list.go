package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/j-clemons/twt/internal/git"
	"github.com/j-clemons/twt/internal/state"
	"github.com/j-clemons/twt/internal/tui"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List TWT-managed tmux sessions",
	Long: `List all tmux sessions created and managed by TWT.

By default, shows sessions for the current git repository.
Use --all to show sessions from all repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")

		if all {
			listAndDisplayAllSessions()
		} else {
			sessions, err := state.ListSessionsForCurrentRepo()
			if err != nil {
				_, ok := err.(*git.NotInGitDirError)
				if ok {
					listAndDisplayAllSessions()
				} else {
					color.Red("Error listing sessions: %v", err)
					return
				}
			} else {
				tui.RunListTui(sessions)
			}
		}
	},
}

func listAndDisplayAllSessions() {
	sessions, err := state.ListAllSessions()
	if err != nil {
		color.Red("Error listing sessions: %v", err)
		return
	}
	tui.RunListTui(sessions)
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "List sessions from all repositories")
}
