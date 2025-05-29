package boba

import (
	tea "github.com/charmbracelet/bubbletea"
)

type AppView int

const (
	ViewLogin AppView = iota
	ViewRequestsMenu
	ViewUserInput
)

func InitialAppModel() AppModel {
	return AppModel{
		currentView: ViewLogin,
		login:       InitialLogin(),
		requestMenu: RequestMenu{},
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.login.Init(),
		m.requestMenu.Init(),
	)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case LoginSuccessMsg:
		m.requestMenu = InitialRequestMenu(msg.Token, msg.RefreshToken, msg.User)
		m.currentView = ViewRequestsMenu
		return m, nil

	case tea.KeyMsg:

		if m.currentView == ViewUserInput && msg.String() == "enter" {
			input := m.userInput.inputs[0].Value()
			m.currentView = ViewRequestsMenu

			return m, func() tea.Msg {
				return UserIDInputMsg(input)
			}
		}

		if m.currentView == ViewRequestsMenu && m.requestMenu.cursor == 4 && msg.String() == "enter" {
			m.userInput = InitialInput("user id", func(input string) (string, error) {
				return input, nil
			})
			m.currentView = ViewUserInput
			return m, m.userInput.Init()
		}
	}

	switch m.currentView {
	case ViewLogin:
		updatedLogin, cmd := m.login.Update(msg)
		m.login = updatedLogin.(Login)
		return m, cmd

	case ViewRequestsMenu:
		updatedMenu, cmd := m.requestMenu.Update(msg)
		m.requestMenu = updatedMenu.(RequestMenu)
		return m, cmd

	case ViewUserInput:
		newInput, cmd := m.userInput.Update(msg)
		m.userInput = newInput.(Input)
		return m, cmd
	}

	return m, nil
}

func (m AppModel) View() string {
	switch m.currentView {
	case ViewLogin:
		return m.login.View()
	case ViewRequestsMenu:
		return m.requestMenu.View()
	case ViewUserInput:
		return m.userInput.View()
	default:
		return "Unknown view"
	}
}
