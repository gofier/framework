package websocket

import (
	"net/http"

	"github.com/gofier/framework/request"
	"github.com/gofier/framework/zone"
)

type Handler interface {
	DefaultChannels() []string
	OnMessage(hub Hub, msg *Msg)
	Loop(hub Hub) error

	OnPing(hub Hub, appData string)
	OnPong(hub Hub, appData string)
	OnClose(hub Hub, code int, text string)

	IConfiguration
}

type Hub interface {
	Send(msg *Msg)
	Broadcast(msg *Msg)
	BroadcastTo(channelName string, msg *Msg)

	name() string
	available() bool

	IChannel
	request.Context
}

type IConfiguration interface {
	ReadBufferSize() int
	WriteBufferSize() int
	CheckOrigin(r *http.Request) bool
	WriteTimeout() zone.Duration
	ReadTimout() zone.Duration
	MaxMessageSize() int64
}

type IChannel interface {
	JoinChannel(channelName string)
	LeaveChannel(channelName string)
}
