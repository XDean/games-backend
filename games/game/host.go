package game

import "github.com/thoas/go-funk"

type (
	Host struct {
		Id   int
		game Game

		clients map[string]*Client

		ready    map[string]bool
		idToSeat map[string]int
		seatToId []string

		eventChan chan clientEvent
	}

	clientEvent struct {
		client *Client
		Event  interface{}
	}

	playerInfo struct {
		Id        string
		Seat      int
		Connected bool
		Ready     bool
	}

	hostInfo struct {
		Id      int
		Players []*playerInfo
	}
)

func NewHost(id int, game Game) *Host {
	return &Host{
		Id:        id,
		game:      game,
		clients:   map[string]*Client{},
		ready:     map[string]bool{},
		idToSeat:  map[string]int{},
		seatToId:  make([]string, game.PlayerCount()),
		eventChan: make(chan clientEvent, game.PlayerCount()),
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
				if _, ok := r.idToSeat[client.id]; !ok {
					if seat, ok := r.availableSeat(); ok {
						r.seatToId[seat] = client.id
						r.idToSeat[client.id] = seat
					} else {
						client.Error("房间已满")
						client.Close()
						continue
					}
				}
				if _, ok := r.clients[client.id]; ok {
					client.Error("Connection already exist")
					client.Close()
					continue
				} else {
					r.clients[client.id] = client
					r.SendAll(TopicEvent{
						Topic:   "connect",
						Payload: client.id,
					})
					client.Send(TopicEvent{
						Topic:   "host-info",
						Payload: r.toInfo(),
					})
				}
			case DisConnectEvent:
				delete(r.clients, client.id)
				r.SendAll(TopicEvent{
					Topic:   "disconnect",
					Payload: client.id,
				})
			case ReadyEvent:
				r.ready[client.id] = e.Ready
				r.SendAll(TopicEvent{
					Topic:   "ready",
					Payload: client.id,
				})
			}
			r.game.HandleEvent(client, event.Event)
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
		if funk.ContainsInt(seats, r.idToSeat[c.id]) {
			c.Send(event)
		}
	}
}

func (r *Host) SendExcludeSeat(event TopicEvent, seats ...int) {
	for _, c := range r.clients {
		if !funk.ContainsInt(seats, r.idToSeat[c.id]) {
			c.Send(event)
		}
	}
}

func (r *Host) toInfo() hostInfo {
	players := make([]*playerInfo, 0)
	for seat, id := range r.seatToId {
		if id == "" {
			players = append(players, nil)
		} else {
			_, ok := r.clients[id]
			players = append(players, &playerInfo{
				Id:        id,
				Seat:      seat,
				Ready:     r.ready[id],
				Connected: ok,
			})
		}
	}

	return hostInfo{
		Id:      r.Id,
		Players: players,
	}
}

func (r *Host) availableSeat() (int, bool) {
	for i := 0; i < r.game.PlayerCount(); i++ {
		if r.seatToId[i] == "" {
			return i, true
		}
	}
	return 0, false
}
