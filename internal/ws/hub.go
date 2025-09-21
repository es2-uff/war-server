package ws

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/labstack/echo/v4"
)

type Hub struct {
	sync.RWMutex

	Id string

	clients map[*Client]bool

	// broadcast  chan *repositories.Message
	register   chan *Client
	unregister chan *Client

	// messages []repositories.Message
	// mr       *repositories.MessageRepository
}

// func NewHub(chatroomId string, messages []repositories.Message, mr *repositories.MessageRepository) *Hub {
// 	return &Hub{
// 		Id:         chatroomId,
// 		clients:    map[*Client]bool{},
// 		broadcast:  make(chan *repositories.Message),
// 		register:   make(chan *Client),
// 		unregister: make(chan *Client),
// 		messages:   messages,
// 		mr:         mr,
// 	}
// }

// func (h *Hub) Run(c echo.Context) {
//
// 	for {
// 		select {
// 		case client := <-h.register:
// 			h.Lock()
// 			h.clients[client] = true
// 			h.Unlock()
// 			fmt.Printf("Client Registered: %s\n", client.id)
//
// 			client.send <- getMessageTemplate(c, h.messages)
//
// 		case client := <-h.unregister:
// 			h.Lock()
// 			if _, ok := h.clients[client]; ok {
// 				close(client.send)
// 				fmt.Printf("Client Unregistered: %s\n", client.id)
// 				delete(h.clients, client)
// 			}
// 			h.Unlock()
//
// 		case msg := <-h.broadcast:
// 			h.RLock()
// 			h.messages = append(h.messages, *msg)
//
// 			err := h.mr.CreateMessage(msg.UserId, msg.Content, h.Id)
//
// 			if err != nil {
// 				h.RUnlock()
// 				return
// 			}
//
// 			for client := range h.clients {
// 				select {
// 				case client.send <- getMessageTemplate(c, h.messages):
// 				default:
// 					close(client.send)
// 					delete(h.clients, client)
// 				}
// 			}
// 			h.RUnlock()
// 		}
// 	}
// }
//
// func getMessageTemplate(c echo.Context, m []repositories.Message) []byte {
//
// 	var buf bytes.Buffer
//
// 	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
//
// 	inbox_views.MessageSection(m).Render(c.Request().Context(), io.Writer(&buf))
//
// 	return buf.Bytes()
// }
