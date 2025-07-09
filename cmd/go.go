package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/j-clemons/twt/internal/checks"
	"github.com/j-clemons/twt/internal/command"
	"github.com/j-clemons/twt/internal/git"
	"github.com/j-clemons/twt/internal/tmux"
	"github.com/j-clemons/twt/internal/utils"
)

var goToWorktree = &cobra.Command{
	Use:   "go <branch>",
	Short: "Gets or creates a tmux session from a given branch.",
	Long: `Given a branch name, either gets or creates a new Tmux session and creates
	/ switches to that branch within that session.

	If the session already exists, switches to it regardless of if a git worktree exists
	or not. If this isn't desired, rename / delete the existing session.

	Also switches to a new session if a worktree exists (ie. the branch is checked out).
	`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
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

		flags := cmd.Flags()
		removeSession, err := flags.GetBool("remove-session")
		if err != nil {
			color.Red("Error fetching the remove session flag")
			return
		}
		currentSession, err := tmux.GetCurrentSessionName()
		if err != nil && removeSession {
			color.Red("Can't remove current session.")
		}

		// Switch to session if exists
		sessionName := utils.GenerateSessionNameFromBranch(branch)
		worktreeName := utils.GenerateWorktreeNameFromBranch(branch)
		isNewSession := !tmux.HasSession(sessionName)

		baseDir, err := git.GetBaseDir()
		if err != nil {
			color.Red(fmt.Sprint(err))
			return
		}

		if !isNewSession {
			tmux.SwitchToSession(sessionName)
			if removeSession {
				tmux.KillSession(currentSession)
			}
			return
		}

		tmux.NewSession(sessionName)
		worktreeExists := git.HasWorktree(branch)
		if worktreeExists {
			changeDirCmd := fmt.Sprintf("cd %s/%s", baseDir, worktreeName)
			tmux.SendKeys(sessionName, changeDirCmd, "Enter")
			tmux.SendKeys(sessionName, "clear", "Enter")
			tmux.SwitchToSession(sessionName)
			if removeSession {
				tmux.KillSession(currentSession)
			}
			return
		}

		// Change to worktree base to create the worktree here
		backToBaseDirCmd := fmt.Sprintf("cd %s", baseDir)
		tmux.SendKeys(sessionName, backToBaseDirCmd, "Enter")

		branchIsNew := !git.HasBranch(branch, false)
		if branchIsNew {
			newWorktreeCmd := fmt.Sprintf("git worktree add %s -b %s", worktreeName, branch)
			tmux.SendKeys(sessionName, newWorktreeCmd, "Enter")
		} else {
			newWorktreeCmd := fmt.Sprintf("git worktree add %s %s", worktreeName, branch)
			tmux.SendKeys(sessionName, newWorktreeCmd, "Enter")
		}

		changeToNewTreeCmd := fmt.Sprintf("cd %s/%s", baseDir, worktreeName)
		tmux.SendKeys(sessionName, changeToNewTreeCmd, "Enter")
		tmux.SendKeys(sessionName, "clear", "Enter")

		// Execute post init scripts
		noScripts, err := flags.GetBool("no-scripts")
		if err != nil {
			color.Red("Couldn't fetch the run scripts flag")
			return
		}
		if !noScripts {
			utils.ExecuteScriptInSession(sessionName, "go", "post.sh")
		}

		tmux.SwitchToSession(sessionName)
		if removeSession {
			tmux.KillSession(currentSession)
		}
	},
}

func init() {
	rootCmd.AddCommand(goToWorktree)

	goToWorktree.Flags().BoolP("no-scripts", "N", false, "Don't run any scripts in the common files dir if they exist for this command.")
	goToWorktree.Flags().BoolP("remove-session", "r", false, "Remove current session (not worktree) after.")
}
