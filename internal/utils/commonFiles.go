package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/j-clemons/twt/internal/git"
)

func GetCommonFilesDirPath() (string, error) {
	baseDir, err := git.GetBaseDir()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/common", baseDir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", errors.New("Common files dir doesn't exist.")
	}
	return path, nil
}

func GetScriptsDirPath() (string, error) {
	commonFilesDir, err := GetCommonFilesDirPath()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/scripts", commonFilesDir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", errors.New("Scripts dir doesn't exist.")
	}
	return path, nil
}

func GetScriptsDirPathForCommand(command string) (string, error) {
	scriptsDir, err := GetScriptsDirPath()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", scriptsDir, command)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", errors.New("Scripts dir doesn't exist.")
	}
	return path, nil
}
