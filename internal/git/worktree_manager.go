package git

import (
	"fmt"
	"os"

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
