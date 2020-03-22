package plugin

import (
	"games-backend/games/host"
	"time"
)

const (
	TopicChat        = "chat-message"
	TopicChatHistory = "chat-history"
)

type (
	chatMessage struct {
		Id   string `json:"id"`
		Text string `json:"text"`
		Time int64  `json:"time"` // unix seconds
	}

	Chat struct {
		Connect *Connect `inject:""`
		history []chatMessage
	}
)

func NewChat() *Chat {
	return &Chat{history: []chatMessage{}}
}

func (c *Chat) Plug(handler host.EventHandler) host.EventHandler {
	return host.EventHandlerFunc(func(ctx host.Context) error {
		switch ctx.Topic {
		case TopicChat:
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
				Topic:   TopicChat,
				Payload: msg,
			}
			for _, id := range c.Connect.GetAll() {
				ctx.SendEvent(id, event)
			}
		case TopicChatHistory:
			ctx.SendBack(host.TopicEvent{
				Topic:   TopicChatHistory,
				Payload: c.history,
			})
		}
		return handler.Handle(ctx)
	})
}
