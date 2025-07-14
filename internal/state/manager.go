package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/j-clemons/twt/internal/tmux"
)

const (
	StateFileName = "sessions.json"
	ConfigDirName = "twt"
)

func getStateFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	twtConfigDir := filepath.Join(configDir, ConfigDirName)
	if err := os.MkdirAll(twtConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(twtConfigDir, StateFileName), nil
}

func LoadState() (*State, error) {
	stateFile, err := getStateFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return NewState(), nil
	}

	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	for name, session := range state.Sessions {
		if tmux.HasSession(name) {
			session.Status = StatusActive
		} else {
			session.Status = StatusInactive
		}
		state.Sessions[name] = session
	}

	return &state, nil
}

func SaveState(state *State) error {
	stateFile, err := getStateFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

func RegisterSession(sessionName, repoPath, repoName, branch, worktreePath string) error {
	state, err := LoadState()
	if err != nil {
		return err
	}

	now := time.Now()
	session := SessionInfo{
		Name:         sessionName,
		RepoPath:     repoPath,
		RepoName:     repoName,
		Branch:       branch,
		WorktreePath: worktreePath,
		CreatedAt:    now,
		LastAccessed: now,
		Status:       StatusActive,
	}

	state.Sessions[sessionName] = session

	tmux.SetEnvironment(sessionName, "TWT_REPO_PATH", repoPath)
	tmux.SetEnvironment(sessionName, "TWT_BRANCH", branch)
	tmux.SetEnvironment(sessionName, "TWT_MANAGED", "true")

	return SaveState(state)
}

func UnregisterSession(sessionName string) error {
	state, err := LoadState()
	if err != nil {
		return err
	}

	delete(state.Sessions, sessionName)
	return SaveState(state)
}

func UpdateLastAccessed(sessionName string) error {
	state, err := LoadState()
	if err != nil {
		return err
	}

	if session, exists := state.Sessions[sessionName]; exists {
		session.LastAccessed = time.Now()
		state.Sessions[sessionName] = session
		return SaveState(state)
	}

	return nil
}
