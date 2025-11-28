package ws

// HubInterface defines the common interface for both RoomHub and GameHub
type HubInterface interface {
	GetRegisterChan() chan *Client
	GetUnregisterChan() chan *Client
	GetBroadcastChan() chan []byte
}
