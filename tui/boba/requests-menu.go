package boba

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type RequestMenu struct {
	cursor       int
	choices      []string
	selected     map[int]struct{}
	token        string
	refreshToken string
}

func InitialRequestMenu(token string, refreshToken string) RequestMenu {
	return RequestMenu{
		choices:      []string{"Show API Token", "Show API Refresh Token ", "Get All Users"},
		cursor:       0,
		selected:     make(map[int]struct{}),
		token:        token,
		refreshToken: refreshToken,
	}
}

func (m RequestMenu) Init() tea.Cmd {
	return tea.SetWindowTitle("Request List")
}

func (m RequestMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m RequestMenu) View() string {
	s := "\nWhat information would you like to request?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}
