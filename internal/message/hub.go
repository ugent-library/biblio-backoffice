package message

type client struct {
	id   string
	send chan []byte
}

type message struct {
	id  string
	msg []byte
}

type Hub struct {
	clients    map[string][]chan []byte
	register   chan client
	unregister chan client
	dispatch   chan message
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string][]chan []byte),
		register:   make(chan client),
		unregister: make(chan client),
		dispatch:   make(chan message),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Register(id string, send chan []byte) {
	h.register <- client{id, send}
}

func (h *Hub) Unregister(id string, send chan []byte) {
	h.unregister <- client{id, send}
}

func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

func (h *Hub) Dispatch(id string, msg []byte) {
	h.dispatch <- message{id, msg}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = append(h.clients[client.id], client.send)
		case client := <-h.unregister:
			if chans, ok := h.clients[client.id]; ok {
				for i, c := range chans {
					if c == client.send {
						chans = append(chans[:i], chans[i+1:]...)
						close(c)
					}
				}
				if len(chans) > 0 {
					h.clients[client.id] = chans
				} else {
					delete(h.clients, client.id)
				}
			}
		case msg := <-h.broadcast:
			for _, chans := range h.clients {
				for _, c := range chans {
					c <- msg
				}
				// select {
				// case send <- msg:
				// default:
				// 	close(send)
				// 	delete(h.clients, id)
				// }
			}
		case message := <-h.dispatch:
			if chans, ok := h.clients[message.id]; ok {
				for _, c := range chans {
					c <- message.msg
				}
				// select {
				// case send <- message.msg:
				// default:
				// 	close(send)
				// 	delete(h.clients, message.id)
				// }
			}
		}
	}
}
