package boba

import tea "github.com/charmbracelet/bubbletea"

type AppView int

const (
	ViewLogin AppView = iota
	ViewRequestsMenu
	ViewOther // Add more views as needed
)

type AppModel struct {
	currentView AppView
	login       Login
	requestMenu RequestMenu
}

func InitialAppModel() AppModel {
	return AppModel{
		currentView: ViewLogin,
		login:       InitialLogin(),
		requestMenu: RequestMenu{}, // initialized later after login
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.login.Init(),
		m.requestMenu.Init(), // Initialize the request menu
	)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case LoginSuccessMsg:
		m.requestMenu = InitialRequestMenu(msg.Token, msg.RefreshToken)
		m.currentView = ViewRequestsMenu
		return m, nil
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
	}

	return m, nil
}

func (m AppModel) View() string {
	switch m.currentView {
	case ViewLogin:
		return m.login.View()
	case ViewRequestsMenu:
		return m.requestMenu.View()
	default:
		return "Unknown view"
	}
}
