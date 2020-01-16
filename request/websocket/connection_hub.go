package websocket

import (
	"github.com/gofier/framework/request"

	"github.com/gofrs/uuid"
)

type connectionHub struct {
	_name    string
	msgChan  chan *Msg
	isClosed bool
	request.Context
}

func (ch *connectionHub) JoinChannel(channelName string) {
	channelMap.Join(channelName, ch)
}

func (ch *connectionHub) LeaveChannel(channelName string) {
	channelMap.Leave(channelName, ch)
}

func newConnectionHub(c request.Context, handler Handler) *connectionHub {
	ch := &connectionHub{}

	ch._name = uuid.Must(uuid.NewV4()).String()
	ch.msgChan = make(chan *Msg, 256)
	ch.isClosed = false
	ch.Context = c

	channelMap.Join(goferDefaultChannelName, ch)

	for _, channelName := range handler.DefaultChannels() {
		channelMap.Join(channelName, ch)
	}

	return ch
}

func (ch *connectionHub) name() string {
	return ch._name
}

func (ch *connectionHub) Send(msg *Msg) {
	ch.msgChan <- msg
}

func (ch *connectionHub) Broadcast(msg *Msg) {
	ch.BroadcastTo(goferDefaultChannelName, msg)
}

func (ch *connectionHub) BroadcastTo(channelName string, msg *Msg) {
	for _, hub := range channelMap.Hubs(channelName).All() {
		if !hub.available() {
			continue
		}
		hub.Send(msg)
	}
}

func (ch *connectionHub) getChan() chan *Msg {
	return ch.msgChan
}

func (ch *connectionHub) close() {
	close(ch.msgChan)
	ch.isClosed = true

	channelMap.LeaveAll(ch)
}

func (ch *connectionHub) available() bool {
	return !ch.isClosed
}
