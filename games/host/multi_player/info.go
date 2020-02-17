package multi_player

type (
	playerInfo struct {
		Id    string `json:"id"`
		Seat  int    `json:"seat"`
		Ready bool   `json:"ready"`
	}

	watcherInfo struct {
		Id string `json:"id"`
	}

	hostInfo struct {
		Players  []*playerInfo  `json:"players"`
		Watchers []*watcherInfo `json:"watchers"`
	}
)

func (h *Host) toInfo() hostInfo {
	players := make([]*playerInfo, 0)
	for seat, id := range h.players {
		if id == "" {
			players = append(players, nil)
		} else {
			players = append(players, &playerInfo{
				Id:    id,
				Seat:  seat,
				Ready: h.ready[id],
			})
		}
	}

	watchers := make([]*watcherInfo, 0)
	for _, id := range h.watchers {
		watchers = append(watchers, &watcherInfo{Id: id})
	}

	return hostInfo{
		Players:  players,
		Watchers: watchers,
	}
}
