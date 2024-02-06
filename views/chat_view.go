package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/net/websocket"
)

type ChatView struct {
	viewport  viewport.Model
	chatInput textinput.Model
	conn      *websocket.Conn
}

func NewChatView(conn *websocket.Conn) *ChatView {
	vp := viewport.New(60, 8)
	vp.SetContent(`Welcome to chat room!
    Type a message and press Enter to send.
    `)
	ci := textinput.New()
	ci.Placeholder = "Type a message..."
	ci.CharLimit = 255
	ci.Width = 60
	ci.Focus()
	return &ChatView{
		viewport:  vp,
		chatInput: ci,
		conn:      conn,
	}
}

func (m ChatView) Init() tea.Cmd {
	return textinput.Blink
}

func (m ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		ciCmd tea.Cmd
		vpCmd tea.Cmd
	)
	m.chatInput, ciCmd = m.chatInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			writeMessage(m.conn, m.chatInput.Value())
			m.chatInput.Reset()
		}
	}
	return m, tea.Batch(ciCmd, vpCmd)
}

func (m ChatView) View() string {
	s := fmt.Sprintf(`
%s
%s
    `,
		m.viewport.View(),
		m.chatInput.View(),
	)
	return s
}

func writeMessage(conn *websocket.Conn, message string) {
	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Println("Error writing message:", err)
	}
}
