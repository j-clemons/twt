package list

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/j-clemons/twt/internal/state"
	"github.com/j-clemons/twt/internal/tmux"
)

type model struct {
	sessions []state.SessionInfo
	cursor   int
}

func CreateModel(sessions []state.SessionInfo) model {
	return model{
		sessions: sessions,
	}
}

func Create(sessions []state.SessionInfo) tea.Program {
	return *tea.NewProgram(CreateModel(sessions))
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.sessions)-1 {
				m.cursor++
			}

		case "enter", " ":
			tmux.SwitchToSession(m.sessions[m.cursor].Name)
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if len(m.sessions) == 0 {
		color.Yellow("No TWT sessions found.")
		color.Cyan("Use 'twt go <branch>' to create a new session.")
		return ""
	}

	s := "All TWT Sessions:\n\n"

	s += fmt.Sprintf(" %-15s %-20s %-20s %-10s %-10s\n", "REPOSITORY", "BRANCH", "SESSION", "STATUS", "CREATED")
	s += strings.Repeat("-", 85) + "\n"

	for i, session := range m.sessions {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		statusColor := color.GreenString
		if session.Status == state.StatusInactive {
			statusColor = color.RedString
		}

		repoName := session.RepoName
		if repoName == "" {
			repoName = filepath.Base(session.RepoPath)
		}

		s += fmt.Sprintf("%s%-15s %-20s %-20s %-10s %-10s\n",
			cursor,
			repoName,
			session.Branch,
			session.Name,
			statusColor(string(session.Status)),
			formatAge(session.Age()),
		)
	}

	s += "\nPress 'enter' or 'space' to switch to selected session\n"
	s += "Press 'q' to quit\n"

	return s
}

func formatAge(duration time.Duration) string {
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%dh", int(duration.Hours()))
	} else {
		return fmt.Sprintf("%dd", int(duration.Hours()/24))
	}
}
