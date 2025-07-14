package state

import (
	"time"
)

type SessionStatus string

const (
	StatusActive   SessionStatus = "active"
	StatusInactive SessionStatus = "inactive"
)

type SessionInfo struct {
	Name         string        `json:"name"`
	RepoPath     string        `json:"repo_path"`
	RepoName     string        `json:"repo_name"`
	Branch       string        `json:"branch"`
	WorktreePath string        `json:"worktree_path"`
	CreatedAt    time.Time     `json:"created_at"`
	LastAccessed time.Time     `json:"last_accessed"`
	Status       SessionStatus `json:"-"`
}

type State struct {
	Version  string                 `json:"version"`
	Sessions map[string]SessionInfo `json:"sessions"`
}

func NewState() *State {
	return &State{
		Version:  "1.0",
		Sessions: make(map[string]SessionInfo),
	}
}

func (s *SessionInfo) IsActive() bool {
	return s.Status == StatusActive
}

func (s *SessionInfo) Age() time.Duration {
	return time.Since(s.CreatedAt)
}

func (s *SessionInfo) TimeSinceAccessed() time.Duration {
	return time.Since(s.LastAccessed)
}
