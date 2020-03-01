package multi_player

type (
	playerInfo struct {
		Id    string `json:"id"`
		Seat  int    `json:"seat"`
		Ready bool   `json:"ready"`
		Host  bool   `json:"host"`
	}

	watcherInfo struct {
		Id string `json:"id"`
	}

	roomInfo struct {
		Playing  bool           `json:"playing"`
		Players  []*playerInfo  `json:"players"`
		Watchers []*watcherInfo `json:"watchers"`
	}
)

func (r *Room) toInfo() roomInfo {
	players := make([]*playerInfo, 0)
	for _, player := range r.players {
		if player == nil {
			players = append(players, nil)
		} else {
			players = append(players, &playerInfo{
				Id:    player.id,
				Seat:  player.seat,
				Ready: player.ready,
				Host:  player.host,
			})
		}
	}

	watchers := make([]*watcherInfo, 0)
	for _, watcher := range r.watchers {
		watchers = append(watchers, &watcherInfo{Id: watcher.id})
	}

	return roomInfo{
		Playing:  r.playing,
		Players:  players,
		Watchers: watchers,
	}
}
