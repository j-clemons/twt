package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/j-clemons/twt/internal/checks"
	"github.com/j-clemons/twt/internal/command"
	"github.com/j-clemons/twt/internal/git"
	"github.com/j-clemons/twt/internal/state"
	"github.com/j-clemons/twt/internal/tmux"
	"github.com/j-clemons/twt/internal/utils"
	"github.com/spf13/cobra"
)

var removeWorktree = &cobra.Command{
	Use:   "rm <branch>",
	Short: "Remove a git worktree, tmux session, and optionally the linked branch.",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		shouldCancel := checks.AssertReady()
		if shouldCancel {
			color.Red("Error when trying to run command, aborting.")
			return
		}

		branch := args[0]
		branch, err := command.Validate(branch)
		if err != nil {
			color.Red(err.Error())
			return
		}

		sessionName := utils.GenerateSessionNameFromBranch(branch)
		worktreeName := utils.GenerateWorktreeNameFromBranch(branch)

		flags := cmd.Flags()
		deleteBranch, err := flags.GetBool("delete-branch")
		if err != nil {
			color.Red("Couldn't check delete-branch flag")
			return
		}
		force, err := flags.GetBool("force")
		if err != nil {
			color.Red("Couldn't check force flag")
			return
		}

		confirm, err := flags.GetBool("confirm")
		if err != nil {
			color.Red("Couldn't check confirm flag")
			return
		}

		// Ask for confirmation unless --force is used
		if !force && !confirm {
			fmt.Printf("Are you sure you want to remove worktree and session for branch '%s'? (y/N): ", branch)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				color.Red("Error reading confirmation")
				return
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				color.Yellow("Operation cancelled")
				return
			}
		}
		nextBranch, err := flags.GetString("target")
		if err != nil {
			color.Red("Couldn't fetch next branch without error")
			return
		}
		var targetSession string
		if nextBranch != "" {
			targetSession = utils.GenerateSessionNameFromBranch(nextBranch)
			if !tmux.HasSession(targetSession) {
				color.Red(fmt.Sprintf("Target session '%s' doesn't exist", targetSession))
				return
			}
		}

		// Git cleanup
		branchExistsAndCheckedOut := git.HasBranch(branch, true)
		worktreeExists := git.HasWorktree(branch)
		if !branchExistsAndCheckedOut {
			color.Red(fmt.Sprintf("Branch %s doesn't exist, or isn't checked out", branch))
			return
		}
		if !worktreeExists {
			color.Red(fmt.Sprintf("Can't delete worktree %s as it doesn't exist", branch))
			return
		}
		if errs := git.RemoveWorktree(worktreeName, branch, force, deleteBranch); len(errs) > 0 {
			for _, err := range errs {
				color.Red(fmt.Sprintf("Error removing worktree: %s", err))
			}
			return
		}

		// Tmux cleanup
		existingSessions, err := tmux.ListSessions(true)
		if err != nil {
			color.Red(fmt.Sprint(err))
			return
		}
		currentSession, err := tmux.GetCurrentSessionName()
		if err != nil {
			color.Red(fmt.Sprint(err))
			return
		}
		possibleDestinations := []string{}
		for _, session := range existingSessions {
			if session != currentSession {
				possibleDestinations = append(possibleDestinations, session)
			}
		}

		// After
		if len(possibleDestinations) == 0 {
			color.Red("No available sessions to switch to")
			return
		}
		newSession := strings.ReplaceAll(possibleDestinations[0], "\"", "")
		if !tmux.HasSession(newSession) {
			color.Red("Session doesn't exist")
			return
		}
		if nextBranch != "" {
			tmux.SwitchToSession(targetSession)
		} else {
			needToSwitchSession := tmux.HasSession(sessionName) && currentSession == sessionName && len(possibleDestinations) > 0
			if needToSwitchSession {
				tmux.SwitchToSession(newSession)
			}
		}
		tmux.KillSession(sessionName)

		// Unregister session from state
		err = state.UnregisterSession(sessionName)
		if err != nil {
			color.Yellow("Warning: Failed to unregister session: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeWorktree)
	removeWorktree.Flags().BoolP("delete-branch", "d", false, "Remove branch as well as the worktree")
	removeWorktree.Flags().BoolP("force", "f", false, "Delete the worktree &| branch regardless of unstaged files")
	removeWorktree.Flags().BoolP("confirm", "y", false, "Skip confirmation prompt")
	removeWorktree.Flags().StringP("target", "t", "", "Where to go after removing session")
}
