package git

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/j-clemons/twt/internal/command"
)

func CreateWorktreeInSession(sessionName, baseDir, worktreeName, branch string, isNewBranch bool) error {
	// Change to base directory
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(baseDir)
	if err != nil {
		return err
	}

	// Run git worktree add command directly (synchronously)
	var args []string
	if isNewBranch {
		args = []string{"worktree", "add", worktreeName, "-b", branch}
	} else {
		args = []string{"worktree", "add", worktreeName, branch}
	}

	_, stderr := command.Run("git", args...)
	if len(stderr) > 0 {
		return fmt.Errorf("git worktree add failed: %v", stderr)
	}

	return nil
}

func VerifyWorktreeReady(baseDir, worktreeName, branch string) error {
	worktreePath := filepath.Join(baseDir, worktreeName)

	// Check if directory exists and is accessible
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree directory does not exist: %s", worktreePath)
	}

	// Verify .git file exists (indicates proper worktree setup)
	gitFile := filepath.Join(worktreePath, ".git")
	if _, err := os.Stat(gitFile); os.IsNotExist(err) {
		return fmt.Errorf("worktree .git file missing: %s", gitFile)
	}

	// Verify we can read the current branch
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(worktreePath); err != nil {
		return fmt.Errorf("cannot access worktree directory: %v", err)
	}

	// Check git status to ensure repository is in good state
	_, stderr := command.Run("git", "status", "--porcelain")
	if len(stderr) > 0 {
		return fmt.Errorf("git status failed in worktree: %v", stderr)
	}

	return nil
}

func WaitForWorktreeReady(baseDir, worktreeName, branch string, timeout time.Duration) error {
	start := time.Now()
	for {
		if err := VerifyWorktreeReady(baseDir, worktreeName, branch); err == nil {
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for worktree to be ready after %v", timeout)
		}

		time.Sleep(50 * time.Millisecond)
	}
}
