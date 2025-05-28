package boba

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type RequestMenu struct {
	cursor       int
	choices      []string
	selected     map[int]struct{}
	token        string
	refreshToken string
	currentUser  User
	response     string
}

func InitialRequestMenu(token string, refreshToken string, currentUser User) RequestMenu {
	return RequestMenu{
		choices:      []string{"API Token", "API Refresh Token ", "All Users", "This User"},
		cursor:       0,
		selected:     make(map[int]struct{}),
		token:        token,
		refreshToken: refreshToken,
		currentUser:  currentUser,
		response:     "",
		// response:     []string{token, refreshToken, "", fmt.Sprintf("%s (%s)", currentUser.Name, currentUser.ID)},
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
				m.response = "" // Clear response if unselecting
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
			resp, err := GenerateResponse(i, m)
			if err != nil {
				m.response = fmt.Sprintf("\nError generating response: %v\n", err)
			} else {
				m.response = fmt.Sprintf("\nResponse:\n%s\n", resp)
			}
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	if m.response != "" {
		s += "\n" + m.response
	}

	s += "\nPress q to quit.\n"

	return s
}

func GenerateResponse(selection int, m RequestMenu) (string, error) {
	switch selection {
	case 0:
		return m.token, nil
	case 1:
		return m.refreshToken, nil
	case 2:
		users, err := GetAllUsers(m.token)
		if err != nil {
			return "", fmt.Errorf("getting all users: %w", err)
		}
		response := "All Users:\n"
		for _, user := range users {
			response += fmt.Sprintf("- %s (%s)\n", user.Name, user.ID)
		}
		return response, nil
	case 3:
		return fmt.Sprintf("Current User: %s (%s)", m.currentUser.Name, m.currentUser.ID), nil
	default:
		return "", fmt.Errorf("invalid choice")
	}
}

func GetAllUsers(token string) ([]User, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/users", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	var parsed []User
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return parsed, nil
}
