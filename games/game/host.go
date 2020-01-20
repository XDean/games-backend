package game

import "github.com/thoas/go-funk"

type (
	Host struct {
		Id        int
		game      Game
		clients   map[string]*Client
		ready     map[string]bool
		seat      map[string]int
		eventChan chan clientEvent
	}

	clientEvent struct {
		client *Client
		Event  interface{}
	}
)

func NewHost(id int, game Game) *Host {
	return &Host{
		Id:        id,
		game:      game,
		clients:   map[string]*Client{},
		ready:     map[string]bool{},
		seat:      map[string]int{},
		eventChan: make(chan clientEvent, 5),
	}
}

func (r *Host) Send(client *Client, event interface{}) {
	r.eventChan <- clientEvent{
		client: client,
		Event:  event,
	}
}

func (r *Host) Run() {
	go func() {
		for {
			event := <-r.eventChan
			client := event.client
			switch e := event.Event.(type) {
			case ConnectEvent:
				if _, ok := r.clients[client.id]; ok {
					client.Error("Connection already exist")
					client.Close()
				} else {
					r.clients[client.id] = client
					r.SendAll(TopicEvent{
						Topic:   "Connect",
						Payload: client.id,
					})
				}
			case DisConnectEvent:
				delete(r.clients, client.id)
				r.SendAll(TopicEvent{
					Topic:   "Disconnect",
					Payload: client.id,
				})
			case ReadyEvent:
				r.ready[client.id] = e.Ready
				r.SendAll(TopicEvent{
					Topic:   "Ready",
					Payload: client.id,
				})
			default:
				r.game.HandleEvent(client, e)
			}
		}
	}()
}
func (r *Host) SendAll(event TopicEvent) {
	for _, c := range r.clients {
		c.Send(event)
	}
}

func (r *Host) SendToSeat(event TopicEvent, seats ...int) {
	for _, c := range r.clients {
		if funk.ContainsInt(seats, r.seat[c.id]) {
			c.Send(event)
		}
	}
}

func (r *Host) SendExcludeSeat(event TopicEvent, seats ...int) {
	for _, c := range r.clients {
		if !funk.ContainsInt(seats, r.seat[c.id]) {
			c.Send(event)
		}
	}
}

func (r *Host) GetClientBySeat(seat int) *Client {
	for id, c := range r.clients {
		if r.seat[id] == seat {
			return c
		}
	}
	return nil
}

func (r *Host) Info() interface{} {
	return struct {
		Id int
	}{}
}
