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
	clients    map[string]chan []byte
	register   chan client
	unregister chan string
	dispatch   chan message
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]chan []byte),
		register:   make(chan client),
		unregister: make(chan string),
		dispatch:   make(chan message),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Register(id string, send chan []byte) {
	h.register <- client{id, send}
}

func (h *Hub) Unregister(id string) {
	h.unregister <- id
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
			h.clients[client.id] = client.send
		case id := <-h.unregister:
			if send, ok := h.clients[id]; ok {
				delete(h.clients, id)
				close(send)
			}
		case msg := <-h.broadcast:
			for id, send := range h.clients {
				select {
				case send <- msg:
				default:
					close(send)
					delete(h.clients, id)
				}
			}
		case message := <-h.dispatch:
			if send, ok := h.clients[message.id]; ok {
				select {
				case send <- message.msg:
				default:
					close(send)
					delete(h.clients, message.id)
				}
			}
		}
	}
}
