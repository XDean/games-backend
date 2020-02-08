package game

import (
	"encoding/json"
	"fmt"
	"games-backend/games/host"
	"github.com/gorilla/websocket"
	"reflect"
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
						client.error("Connection already exist")
						client.close()
						continue
					}
					s.clients[client.id] = client
					s.sendAll(host.TopicEvent{
						Topic:   "connect",
						Payload: client.id,
					})
				case "disconnect":
					delete(s.clients, client.id)
					s.sendAll(host.TopicEvent{
						Topic:   "disconnect",
						Payload: client.id,
					})
				}
			}

			s.host.Handle(host.Context{
				Who:   client.id,
				Topic: event.topic,
				GetPayload: func(payload interface{}) error {
					return json.Unmarshal(event.payload, payload)
				},
				SendEvent:    s.send,
				TriggerEvent: s._triggerEvent,
			})

			if event.done != nil {
				close(event.done)
			}
		}
	}()
}

func (s Server) close() {
	// TODO
}

func (s Server) _triggerEvent(event host.TopicEvent) {
	payloadValue := reflect.ValueOf(event.Payload)
	s.host.Handle(host.Context{
		Who:   "",
		Topic: event.Topic,
		GetPayload: func(payload interface{}) error {
			rv := reflect.ValueOf(payload)
			if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
				receiveValue := rv.Elem()
				if receiveValue.Type() == payloadValue.Type() {
					receiveValue.Set(payloadValue)
					return nil
				}
			}
			return fmt.Errorf("Wrong payload type, except %T, actual %T", payload, event.Payload)
		},
		SendEvent:    s.send,
		TriggerEvent: s._triggerEvent,
	})
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
