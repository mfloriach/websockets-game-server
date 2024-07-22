package internal

import (
	"errors"

	"github.com/gofiber/contrib/websocket"
)

type Player struct {
	ID        string
	X         uint
	Y         uint
	Direction Type
	Width     uint
	Height    uint
	Life      int
	Conn      *websocket.Conn `json:"-"`
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{
		ID:        conn.Params("name"),
		X:         10,
		Y:         10,
		Direction: TypeMoveUp,
		Life:      3,
		Conn:      conn,
	}
}

func (p *Player) CheckCollision(r *Room, event Event) error {
	if err := p.checkCollisionScenario(r, event); err != nil {
		return err
	}

	return p.checkCollisionOtherPlayers(r, event)
}

func (player *Player) checkCollisionScenario(r *Room, event Event) error {
	switch event.Type {
	case TypeMoveUp: // up
		if player.Y+20 >= r.Height {
			return nil
		}
		player.Y += 20
	case TypeMoveDown: // down
		if player.Y-20 <= 0 {
			return nil
		}
		player.Y -= 20
	case TypeMoveLeft: // left
		if player.X-1 <= 0 {
			return nil
		}
		player.X -= 20
	case TypeMoveRight: // right
		if player.X+1 >= r.Width {
			return nil
		}
		player.X += 20
	default:
		return errors.New("invalid direction")
	}

	return nil
}

func (player *Player) checkCollisionOtherPlayers(r *Room, event Event) error {
	for _, p := range r.players {
		clone := *player

		switch event.Type {
		case TypeMoveUp: // up
			clone.Y += 20
		case TypeMoveDown: // down
			clone.Y -= 20
		case TypeMoveLeft: // left
			clone.X -= 20
		case TypeMoveRight: // right
			clone.X += 20
		default:
			return errors.New("invalid direction")
		}

		if clone.X == p.X && clone.Y == p.Y {
			return nil
		}
	}

	return nil
}
