package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/net/websocket"

	"nihal/apk_builder_v2/views"
)

const wsUrl = "ws://localhost:8080/ws"

func main() {
	origin := "http://ws.postman-echo.com/"
	conn, err := websocket.Dial(wsUrl, "", origin)
	if err != nil {
		log.Println("Error dialing websocket:", err)
		os.Exit(1)
	}
	p := tea.NewProgram(views.NewChatView(conn))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
