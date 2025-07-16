package tui

import (
	"fmt"
	"os"

	"github.com/j-clemons/twt/internal/state"
	"github.com/j-clemons/twt/internal/tui/list"
)

func RunListTui(sessions []state.SessionInfo) {
	p := list.Create(sessions)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
