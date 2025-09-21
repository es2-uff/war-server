package ws

import (
	"github.com/labstack/echo/v4"
)

type WSRoomServer struct {
	rooms map[*Hub]bool
	mr    *repositories.MessageRepository
}

func NewWsSever(mr *repositories.MessageRepository) *WsServer {
	return &WsServer{
		mr:        mr,
		chatrooms: make(map[*Hub]bool),
	}
}

func (w *WsServer) NewHub(c echo.Context, chatroomId string) *Hub {

	for k := range w.chatrooms {
		if k.Id == chatroomId {
			return k
		}
	}

	messages, err := w.mr.GetChatroomMessages(chatroomId)

	if err != nil {

	}

	hub := NewHub(chatroomId, messages, w.mr)

	w.chatrooms[hub] = true

	go hub.Run(c)

	return hub
}
