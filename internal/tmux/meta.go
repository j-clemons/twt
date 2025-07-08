package tmux

import (
	"github.com/j-clemons/twt/internal/command"
)

func SendKeys(session string, toSend ...string) {
	args := append([]string{"send-keys", "-t", session}, toSend...)
	command.Run("tmux", args...)
}
