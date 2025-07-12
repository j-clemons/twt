package tmux

import (
	"fmt"
)

func SwitchOrCreateSession(sessionName, currentSession string, removeCurrentSession bool) error {
	isNewSession := !HasSession(sessionName)

	if !isNewSession {
		SwitchToSession(sessionName)
		if removeCurrentSession {
			KillSession(currentSession)
		}
		return nil
	}

	return nil
}

func CreateSessionInDirectory(sessionName, directory string) error {
	NewSessionWithDirectory(sessionName, directory)
	return nil
}

func ChangeDirectoryAndClear(sessionName, directory string) {
	changeDirCmd := fmt.Sprintf("cd %s", directory)
	SendKeys(sessionName, changeDirCmd, "Enter")
	SendKeys(sessionName, "clear", "Enter")
}

func SetupWorktreeSession(sessionName, baseDir, worktreeName string) error {
	SendKeys(sessionName, "clear", "Enter")
	return nil
}

func FinalizeSession(sessionName, currentSession string, removeCurrentSession bool) {
	SwitchToSession(sessionName)
	if removeCurrentSession {
		KillSession(currentSession)
	}
}
