package tmux

func CreateSessionInDirectory(sessionName, directory string) error {
	NewSessionWithDirectory(sessionName, directory)
	return nil
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
