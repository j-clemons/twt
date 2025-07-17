package list

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	var s strings.Builder

	highlightStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	s.WriteString("TWT Sessions:\n\n")

	s.WriteString(fmt.Sprintf("  %-15s %-25s %-10s\n", "REPOSITORY", "BRANCH", "CREATED"))
	s.WriteString(strings.Repeat("-", 65) + "\n")

	for i, session := range m.sessions {
		repoName := session.RepoName
		if repoName == "" {
			repoName = filepath.Base(session.RepoPath)
		}

		sessionStr := fmt.Sprintf("%-15s %-25s %-10s",
			repoName,
			session.Branch,
			formatAge(session.Age()),
		)

		if m.cursor == i {
			s.WriteString(highlightStyle.Render(fmt.Sprintf("> %s", sessionStr)))
		} else {
			s.WriteString(fmt.Sprintf("  %s", sessionStr))
		}

		s.WriteString("\n")
	}

	s.WriteString("\nPress 'enter' or 'space' to switch to selected session\n")
	s.WriteString("Press 'q' to quit\n")

	return s.String()
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
