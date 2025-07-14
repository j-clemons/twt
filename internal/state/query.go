package state

import (
	"sort"

	"github.com/j-clemons/twt/internal/git"
)

func ListSessionsForRepo(repoPath string) ([]SessionInfo, error) {
	state, err := LoadState()
	if err != nil {
		return nil, err
	}

	var sessions []SessionInfo
	for _, session := range state.Sessions {
		if session.RepoPath == repoPath {
			sessions = append(sessions, session)
		}
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	return sessions, nil
}

func ListSessionsForCurrentRepo() ([]SessionInfo, error) {
	repoPath, err := git.GetBaseDir()
	if err != nil {
		return nil, err
	}

	return ListSessionsForRepo(repoPath)
}

func ListAllSessions() ([]SessionInfo, error) {
	state, err := LoadState()
	if err != nil {
		return nil, err
	}

	var sessions []SessionInfo
	for _, session := range state.Sessions {
		sessions = append(sessions, session)
	}

	sort.Slice(sessions, func(i, j int) bool {
		if sessions[i].RepoName != sessions[j].RepoName {
			return sessions[i].RepoName < sessions[j].RepoName
		}
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	return sessions, nil
}
