package backend

import (
	"sync"
	"time"
)

type Coordinate struct {
	X int32
	Y int32
}

type Direction int32

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
	DirectionStop
)

type Player struct {
	Position  Coordinate
	Name      string
	Direction Direction
	Icon      rune
	Mux       sync.Mutex
}

type Game struct {
	Players          map[string]*Player
	Mux              sync.Mutex
	OnPositionChange func(*Player)
}

func NewGame() Game {
	game := Game{
		Players: make(map[string]*Player),
	}
	return game
}

func (game *Game) Start() {
	go func() {
		lastmove := map[string]time.Time{}
		for {
			game.Mux.Lock()
			for name, player := range game.Players {
				player.Mux.Lock()
				if player.Direction == DirectionStop || lastmove[name].After(time.Now().Add(-50*time.Millisecond)) {
					player.Direction = DirectionStop
					player.Mux.Unlock()
					continue
				}
				switch player.Direction {
				case DirectionUp:
					player.Position.Y--
				case DirectionDown:
					player.Position.Y++
				case DirectionLeft:
					player.Position.X--
				case DirectionRight:
					player.Position.X++
				}
				player.Direction = DirectionStop
				lastmove[name] = time.Now()
				if game.OnPositionChange != nil {
					game.OnPositionChange(player)
				}
				player.Mux.Unlock()
			}
			game.Mux.Unlock()
		}
	}()
}