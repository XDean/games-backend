package game

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type (
	Client struct {
		id   string
		host *Host
		conn *websocket.Conn
		send chan TopicEvent
	}
)

func NewClient(id string, host *Host, conn *websocket.Conn) *Client {
	return &Client{id: id, host: host, conn: conn, send: make(chan TopicEvent, 5)}
}

func (c *Client) Send(event TopicEvent) {
	c.send <- event
}

func (c *Client) Error(err string) {
	c.Send(ErrorEvent(err))
}

func (c *Client) Seat() (int, bool) {
	seat, ok := c.host.idToSeat[c.id]
	return seat, ok
}

func (c *Client) Start() {
	go c.read()
	go c.write()
}

func (c *Client) Close() {
	close(c.send)
}

func (c *Client) read() {
	defer func() {
		c.host.Send(c, &DisConnectEvent{})
		_ = c.conn.Close()
	}()
	c.host.Send(c, &ConnectEvent{})
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
			c.Error(err.Error())
			continue
		}
		event := c.host.NewEvent(strings.ToLower(topic.Topic))
		if event == nil {
			c.Error("Unknown Topic: " + topic.Topic)
		} else {
			err := json.Unmarshal(topic.Payload, event)
			if err != nil {
				c.Error(err.Error())
				continue
			}
			c.host.Send(c, event)
		}
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
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

func (c *Client) writeJSON(o interface{}) error {
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
