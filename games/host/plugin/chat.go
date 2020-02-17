package plugin

import (
	"games-backend/games/host"
	"time"
)

type (
	chatMessage struct {
		Id   string `json:"id"`
		Text string `json:"text"`
		Time int64  `json:"time"` // unix seconds
	}
	Chat struct {
		connected map[string]bool
		history   []chatMessage
	}
)

func NewChat() Chat {
	return Chat{connected: map[string]bool{}}
}

func (c Chat) Plug(handler host.EventHandler) host.EventHandler {
	return host.EventHandlerFunc(func(ctx host.Context) error {
		switch ctx.Topic {
		case "connect":
			c.connected[ctx.ClientId] = true
		case "disconnect":
			delete(c.connected, ctx.ClientId)
		case "chat":
			text := ""
			err := ctx.GetPayload(&text)
			if err != nil {
				return err
			}
			msg := chatMessage{
				Id:   ctx.ClientId,
				Text: text,
				Time: time.Now().Unix(),
			}
			c.history = append(c.history, msg)
			event := host.TopicEvent{
				Topic:   "chat",
				Payload: msg,
			}
			for id := range c.connected {
				ctx.SendEvent(id, event)
			}
		case "chat-history":
			ctx.SendBack(host.TopicEvent{
				Topic:   "chat-history",
				Payload: c.history,
			})
		}
		return handler.Handle(ctx)
	})
}
