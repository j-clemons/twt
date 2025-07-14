package tmux

import (
	"strings"

	"github.com/j-clemons/twt/internal/command"
)

func SendKeys(session string, toSend ...string) {
	args := append([]string{"send-keys", "-t", session}, toSend...)
	command.Run("tmux", args...)
}

func SetEnvironment(sessionName, key, value string) {
	command.Run("tmux", "set-environment", "-t", sessionName, key, value)
}

func GetEnvironment(sessionName, key string) string {
	out, _ := command.Run("tmux", "show-environment", "-t", sessionName, key)
	if len(out) > 0 {
		parts := strings.SplitN(out[0], "=", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}
