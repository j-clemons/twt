package utils

import (
	"path/filepath"
	"strings"

	"github.com/j-clemons/twt/internal/git"
)

func GenerateSessionNameFromBranch(branchName string) string {
	projectName := getProjectName()
	sanitizedBranch := strings.Replace(branchName, "/", "__", -1)
	return projectName + "_" + sanitizedBranch
}

func GenerateWorktreeNameFromBranch(branchName string) string {
	return strings.Replace(branchName, "/", "__", -1)
}

func getProjectName() string {
	baseDir, err := git.GetBaseDir()
	if err != nil {
		return "unknown"
	}
	projectName := filepath.Base(baseDir)
	// Replace dots with underscores to avoid tmux session name conflicts
	projectName = strings.Replace(projectName, ".", "_", -1)
	return projectName
}
