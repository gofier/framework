package websocket

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofier/framework/log"
	request_http "github.com/gofier/framework/request/http"
	"github.com/gofier/framework/zone"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func CovertHandler(wsHandler Handler) gin.HandlerFunc {
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var pingPeriod = (wsHandler.ReadTimout() * 9) / 10
	return func(c *gin.Context) {
		goferCtx := request_http.ConvertContext(c)

		hub := newConnectionHub(goferCtx, wsHandler)

		ws, err := wsUpgrader.Upgrade(goferCtx.Writer(), goferCtx.Request(), nil)
		if err != nil {
			log.ErrorWithFields(err, map[string]interface{}{"msg": "Failed to set websocket upgrade"})
			goferCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"error": err})
			return
		}

		defer func() {
			if err := ws.Close(); err != nil {
				log.Error(err)
			}
			hub.close()
		}()

		defer func() {
			if _err := recover(); _err != nil {
				if __err, ok := _err.(error); ok {
					log.Error(__err)
					return
				}
				log.Error(errors.New(fmt.Sprint(_err)))
				return
			}
		}()

		ws.SetReadLimit(wsHandler.MaxMessageSize())

		setPingPongCloseHandler(ws, wsHandler, hub)

		go func() {
			pingTicker := zone.NewTicker(pingPeriod)
			defer pingTicker.Stop()

			for {
				if err := wsHandler.Loop(hub); err != nil {
					return
				}
				select {
				case msg, ok := <-hub.getChan():
					if !ok || msg.isDone() {
						return
					}
					if err := msg.send(ws, wsHandler); err != nil {
						log.Error(err)
						return
					}
				case <-pingTicker.C:
					ping := Msg{
						msgType: websocket.PingMessage,
						data:    &[]byte{},
						err:     nil,
					}
					if err := ping.send(ws, wsHandler); err != nil {
						log.Error(err)
						return
					}
				default:
					continue
				}
			}
		}()

		for {
			msg := &Msg{}
			if msg.scan(ws, wsHandler) != nil {
				log.Error(err)
				return
			}
			switch msg.msgType {
			case websocket.TextMessage:
				fallthrough
			case websocket.BinaryMessage:
				wsHandler.OnMessage(hub, msg)
			case websocket.PingMessage:
				pong := Msg{
					msgType: websocket.PongMessage,
					data:    &[]byte{},
					err:     nil,
				}
				if err := pong.send(ws, wsHandler); err != nil {
					log.Error(err)
					return
				}
			case websocket.PongMessage:
				if msg.err != nil {
					log.Error(err)
					return
				}
			case websocket.CloseMessage:
				closeMsg := Msg{
					msgType: websocket.CloseMessage,
					data:    &[]byte{},
					err:     nil,
				}
				if err := closeMsg.send(ws, wsHandler); err != nil {
					log.Error(err)
					return
				}
				return
			default:
				log.WarnWithFields("No websocket handler on this msgType", map[string]interface{}{"msgType": msg.msgType})
			}
		}
	}
}

func setPingPongCloseHandler(ws *websocket.Conn, wsHandler Handler, hub Hub) {
	ws.SetCloseHandler(func(code int, text string) error {
		wsHandler.OnClose(hub, code, text)
		return nil
	})
	ws.SetPingHandler(func(appData string) error {
		wsHandler.OnPing(hub, appData)
		return nil
	})
	ws.SetPongHandler(func(appData string) error {
		wsHandler.OnPong(hub, appData)
		return nil
	})
}
