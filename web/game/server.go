package game

import (
	"encoding/json"
	"games-backend/games/host"
	"games-backend/util"
	"github.com/gorilla/websocket"
	"github.com/thoas/go-funk"
)

type (
	Server struct {
		id        int
		host      host.Host
		clients   map[string]wsClient
		eventChan chan clientEvent
	}

	Client interface {
	}

	clientEvent struct {
		client  wsClient // can be empty
		topic   string
		payload []byte
		done    chan struct{} // can be nil
	}
)

func newServer(id int, h host.Host) Server {
	return Server{
		id:        id,
		host:      h,
		clients:   map[string]wsClient{},
		eventChan: make(chan clientEvent, 10),
	}
}

func (s Server) newClient(id string, conn *websocket.Conn) wsClient {
	return wsClient{id: id, server: s, conn: conn, eventChan: make(chan host.TopicEvent, 5)}
}

func (s Server) sendToServer(event clientEvent) {
	s.eventChan <- event
}

func (s Server) run() {
	go func() {
		for {
			event := <-s.eventChan
			client := event.client

			if client.id != "" {
				switch event.topic {
				case "connect":
					if _, ok := s.clients[client.id]; ok {
						client.error("另一个连接已经存在，你的名字可能被人占用了")
						client.close()
						continue
					}
					s.clients[client.id] = client
					s.sendAll(host.TopicEvent{
						Topic:   "connect",
						Payload: client.id,
					})
				case "disconnect":
					if s.clients[client.id].conn == client.conn {
						delete(s.clients, client.id)
						s.sendAll(host.TopicEvent{
							Topic:   "disconnect",
							Payload: client.id,
						})
					}
				case "connect-info":
					client.sendToClient(host.TopicEvent{
						Topic:   "connect-info",
						Payload: funk.Keys(s.clients),
					})
				}
			}

			err := s.host.Handle(host.NewContext(host.Context{
				ClientId: client.id,
				Topic:    event.topic,
				GetPayload: func(payload interface{}) error {
					return json.Unmarshal(event.payload, payload)
				},
				SendEvent:    s.send,
				TriggerEvent: s._triggerEvent,
			}))
			if err != nil {
				client.sendToClient(host.ErrorEvent(err.Error()))
			}
			if event.done != nil {
				close(event.done)
			}
		}
	}()
}

func (s Server) close() {
	// TODO
}

func (s Server) _triggerEvent(event host.TopicEvent) error {
	receiverFunc := util.ReflectSetReceiver(event.Payload)
	return s.host.Handle(host.NewContext(host.Context{
		ClientId: "",
		Topic:    event.Topic,
		GetPayload: func(payload interface{}) error {
			return receiverFunc(payload)
		},
		SendEvent:    s.send,
		TriggerEvent: s._triggerEvent,
	}))
}

func (s Server) send(id string, event host.TopicEvent) {
	if client, ok := s.clients[id]; ok {
		client.sendToClient(event)
	}
}

func (s Server) sendAll(event host.TopicEvent) {
	for _, c := range s.clients {
		c.sendToClient(event)
	}
}
