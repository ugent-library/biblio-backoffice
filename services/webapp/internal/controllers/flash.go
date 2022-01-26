package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/context"
	"github.com/ugent-library/biblio-backend/services/webapp/internal/views"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Flashes struct {
	Context
}

func NewFlashes(c Context) *Flashes {
	return &Flashes{c}
}

func (c *Flashes) Ws(w http.ResponseWriter, r *http.Request) {
	user := context.GetUser(r.Context())

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	sendCh := make(chan []byte, 64)
	c.Engine.MessageHub.Register(user.ID, sendCh)

	go c.wsWriter(ws, sendCh)
	c.wsReader(ws, user.ID)
}

func (c *Flashes) wsReader(ws *websocket.Conn, userID string) {
	defer func() {
		ws.Close()
		c.Engine.MessageHub.Unregister(userID)
	}()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Flashes) wsWriter(ws *websocket.Conn, sendCh chan []byte) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case msg := <-sendCh:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			flash, _ := c.RenderPartial("layouts/_flash_message", views.Flash{
				Message: string(msg),
			})
			swap := `<div hx-swap-oob="beforeend:#flash-messages">` + flash + `</div>`
			if err := ws.WriteMessage(websocket.TextMessage, []byte(swap)); err != nil {
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
