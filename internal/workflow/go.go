package workflow

import (
	"fmt"

	"github.com/j-clemons/twt/internal/git"
	"github.com/j-clemons/twt/internal/tmux"
	"github.com/j-clemons/twt/internal/utils"
)

type GoOptions struct {
	Branch               string
	RemoveCurrentSession bool
	NoScripts            bool
	CurrentSession       string
}

func ExecuteGo(opts GoOptions) error {
	sessionName := utils.GenerateSessionNameFromBranch(opts.Branch)
	worktreeName := utils.GenerateWorktreeNameFromBranch(opts.Branch)

	baseDir, err := git.GetBaseDir()
	if err != nil {
		return err
	}

	if tmux.HasSession(sessionName) {
		tmux.SwitchToSession(sessionName)
		if opts.RemoveCurrentSession {
			tmux.KillSession(opts.CurrentSession)
		}
		return nil
	}

	worktreeExists := git.HasWorktree(opts.Branch)
	if worktreeExists {
		sessionDir := fmt.Sprintf("%s/%s", baseDir, worktreeName)
		err = tmux.CreateSessionInDirectory(sessionName, sessionDir)
		if err != nil {
			return err
		}
		err = tmux.SetupWorktreeSession(sessionName, baseDir, worktreeName)
	} else {
		err = tmux.CreateSessionInDirectory(sessionName, baseDir)
		if err != nil {
			return err
		}
		branchIsNew := !git.HasBranch(opts.Branch, false)
		err = git.CreateWorktreeInSession(sessionName, baseDir, worktreeName, opts.Branch, branchIsNew)
		if err != nil {
			return err
		}

		// Recreate session in the new worktree directory
		tmux.KillSession(sessionName)
		sessionDir := fmt.Sprintf("%s/%s", baseDir, worktreeName)
		err = tmux.CreateSessionInDirectory(sessionName, sessionDir)
		if err != nil {
			return err
		}
		tmux.SendKeys(sessionName, "clear", "Enter")
	}

	if err != nil {
		return err
	}

	return handlePostInitialization(sessionName, opts.NoScripts, opts.RemoveCurrentSession, opts.CurrentSession)
}

func handlePostInitialization(sessionName string, noScripts, removeSession bool, currentSession string) error {
	if !noScripts {
		utils.ExecuteScriptInSession(sessionName, "go", "post.sh")
	}

	tmux.FinalizeSession(sessionName, currentSession, removeSession)
	return nil
}
