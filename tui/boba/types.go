package boba

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/lib/pq"
)

type AppModel struct {
	currentView AppView
	login       Login
	requestMenu RequestMenu
	userInput   Input
}

type User struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Online   bool           `json:"online"`
	Channels pq.StringArray `json:"channels" sql:"type:text[]"`
	Created  int64          `json:"created"`
	Updated  int64          `json:"updated"`
}

type Input struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	inoutLabel string
	inputFunc  InputFunc
}

type InputFunc func(string) (string, error)

type UserIDInputMsg string

type Login struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

type LoginResponse struct {
	Message      string `json:"message"`
	RefreshToken string `json:"refreshToken"`
	Token        string `json:"token"`
	User         User   `json:"user"`
}

type LoginSuccessMsg struct {
	Token        string
	RefreshToken string
	User         User
}

type RequestMenu struct {
	cursor       int
	choices      []string
	selected     map[int]struct{}
	token        string
	refreshToken string
	currentUser  User
	response     string
	tempUserID   string
}
