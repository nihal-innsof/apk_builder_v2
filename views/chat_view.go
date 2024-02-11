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
	recvChan  chan wsMessage
	msgs      []string
}

type wsMessage struct {
	Msg string
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
	recvChan := make(chan wsMessage)
	go func() {
		for {
			msg := wsMessage{}
			msgBytes := make([]byte, 1024)
			if _, err := conn.Read(msgBytes); err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			msg.Msg = string(msgBytes)
			recvChan <- msg
		}
	}()
	return &ChatView{
		viewport:  vp,
		chatInput: ci,
		conn:      conn,
		recvChan:  recvChan,
	}
}

func (m ChatView) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, getNextMsg(m.recvChan))
}

func (m ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		ciCmd tea.Cmd
		vpCmd tea.Cmd
		mCmd  tea.Cmd
	)
	m.chatInput, ciCmd = m.chatInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	switch msg := msg.(type) {
	case wsMessage:
		m.msgs = append(m.msgs, msg.Msg)
		m.viewport.SetContent(m.viewMessages())
		m.viewport.GotoBottom()
		mCmd = getNextMsg(m.recvChan)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			message := m.chatInput.Value()
			if message != "" {
				writeMessage(m.conn, m.chatInput.Value())
			}
			m.chatInput.Reset()
		}
	}
	return m, tea.Batch(ciCmd, vpCmd, mCmd)
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

func (m ChatView) viewMessages() string {
	s := ""
	for i := range m.msgs {
		s += fmt.Sprintf("%s\n", m.msgs[i])
	}
	return s
}

func getNextMsg(c <-chan wsMessage) tea.Cmd {
	return func() tea.Msg {
		return <-c
	}
}

func writeMessage(conn *websocket.Conn, message string) {
	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Println("Error writing message:", err)
	}
}
