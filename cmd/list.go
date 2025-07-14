package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/j-clemons/twt/internal/git"
	"github.com/j-clemons/twt/internal/state"
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
				displayRepoSessions(sessions)
			}
		}
	},
}

func displayRepoSessions(sessions []state.SessionInfo) {
	if len(sessions) == 0 {
		color.Yellow("No TWT sessions found for this repository.")
		color.Cyan("Use 'twt go <branch>' to create a new session.")
		return
	}

	repoPath, err := git.GetBaseDir()
	if err != nil {
		color.Red("Error getting repository info: %v", err)
		return
	}
	repoName := filepath.Base(repoPath)

	fmt.Printf("Current Repository: %s (%s)\n\n", color.CyanString(repoName), repoPath)

	fmt.Printf("%-20s %-20s %-10s %-10s\n", "BRANCH", "SESSION", "STATUS", "CREATED")
	fmt.Println(strings.Repeat("-", 70))

	for _, session := range sessions {
		statusColor := color.GreenString
		if session.Status == state.StatusInactive {
			statusColor = color.RedString
		}

		fmt.Printf("%-20s %-20s %-10s %-10s\n",
			session.Branch,
			session.Name,
			statusColor(string(session.Status)),
			formatAge(session.Age()),
		)
	}

	fmt.Printf("\nUse 'twt go <branch>' to switch to a session\n")
}

func listAndDisplayAllSessions() {
	sessions, err := state.ListAllSessions()
	if err != nil {
		color.Red("Error listing sessions: %v", err)
		return
	}
	displayAllSessions(sessions)
}

func displayAllSessions(sessions []state.SessionInfo) {
	if len(sessions) == 0 {
		color.Yellow("No TWT sessions found.")
		color.Cyan("Use 'twt go <branch>' to create a new session.")
		return
	}

	fmt.Printf("All TWT Sessions:\n\n")

	fmt.Printf("%-15s %-20s %-20s %-10s %-10s\n", "REPOSITORY", "BRANCH", "SESSION", "STATUS", "CREATED")
	fmt.Println(strings.Repeat("-", 85))

	for _, session := range sessions {
		statusColor := color.GreenString
		if session.Status == state.StatusInactive {
			statusColor = color.RedString
		}

		repoName := session.RepoName
		if repoName == "" {
			repoName = filepath.Base(session.RepoPath)
		}

		fmt.Printf("%-15s %-20s %-20s %-10s %-10s\n",
			repoName,
			session.Branch,
			session.Name,
			statusColor(string(session.Status)),
			formatAge(session.Age()),
		)
	}

	fmt.Printf("\nUse 'twt go <branch>' to switch to a session\n")
}

func formatAge(duration time.Duration) string {
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%dh", int(duration.Hours()))
	} else {
		return fmt.Sprintf("%dd", int(duration.Hours()/24))
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "List sessions from all repositories")
}
