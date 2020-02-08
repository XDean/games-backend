package game

import (
	"encoding/json"
	"games-backend/games/host"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type (
	wsClient struct {
		id        string
		server    Server
		conn      *websocket.Conn
		eventChan chan host.TopicEvent
	}
)

func (c wsClient) sendToServer(topic string, payload []byte) {
	c.server.sendToServer(clientEvent{
		client:  c,
		topic:   topic,
		payload: payload,
	})
}

func (c wsClient) sendToClient(event host.TopicEvent) {
	c.eventChan <- event
}

func (c wsClient) error(err string) {
	c.sendToClient(host.ErrorEvent(err))
}

func (c wsClient) run() {
	go c.read()
	go c.write()
}

func (c wsClient) close() {
	close(c.eventChan)
}

func (c wsClient) read() {
	defer func() {
		c.sendToServer("disconnect", nil)
		_ = c.conn.Close()
	}()
	c.sendToServer("connect", nil)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("Fail to read websocket message")
			}
			return
		}
		topic := struct {
			Topic   string          `json:"topic"`
			Payload json.RawMessage `json:"payload"`
		}{}
		err = json.Unmarshal(message, &topic)
		if err != nil {
			c.error(err.Error())
			continue
		}
		c.sendToServer(topic.Topic, topic.Payload)
	}
}

func (c wsClient) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.eventChan:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.writeJSON(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c wsClient) writeJSON(o interface{}) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(o)
	if _, err := w.Write(bs); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}
