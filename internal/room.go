package internal

import (
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Room struct {
	ID      uuid.UUID
	Height  uint
	Width   uint
	players map[string]*Player
	bullets map[string]Bullet
}

func NewRoom() Room {
	return Room{
		ID:      uuid.New(),
		Height:  500,
		Width:   500,
		players: map[string]*Player{},
		bullets: map[string]Bullet{},
	}
}

func (r *Room) PlayerAdd(conn *websocket.Conn) *Player {
	player := NewPlayer(conn)

	r.players[player.ID] = player

	return r.players[player.ID]
}

func (r *Room) PlayerAction(player *Player, event Event) error {
	if event.Type == TypeShoot {
		r.PlayerShoot(player)
		return nil
	}

	if player.Direction != event.Type {
		player.Direction = event.Type
		return nil
	}

	if err := player.CheckCollision(r, event); err != nil {
		return err
	}

	r.PlayerSet(player)

	return nil
}

func (r *Room) PlayerSet(p *Player) {
	r.players[p.ID] = p
}

func (r *Room) PlayerGet(id string) *Player {
	return r.players[id]
}

func (r *Room) PlayerShoot(player *Player) {
	var (
		x uint
		y uint
	)
	switch player.Direction {
	case TypeMoveUp:
		y = player.Y + 6
		x = player.X
	case TypeMoveDown:
		y = player.Y - 6
		x = player.X
	case TypeMoveLeft:
		y = player.Y
		x = player.X + 6
	case TypeMoveRight:
		y = player.Y
		x = player.X + 6
	}

	r.bullets[uuid.New().String()] = Bullet{X: x, Y: y, Direction: player.Direction}
}

func (r *Room) broadcast() error {
	players := []*Player{}
	for _, v := range r.players {
		players = append(players, v)
	}

	bullets := []Bullet{}
	for _, v := range r.bullets {
		bullets = append(bullets, v)
	}

	e := EventReturn{Players: players, Bullets: bullets}
	for _, pr := range r.players {
		if err := pr.Conn.WriteJSON(e); err != nil {
			return err
		}
	}

	return nil
}

func (r *Room) UpdateBullets() {
	for i, b := range r.bullets {

		if b.Y >= r.Height || b.Y <= 0 || b.X <= 0 || b.X >= r.Width {
			delete(r.bullets, i)
			continue
		}

		switch b.Direction {
		case TypeMoveUp:
			if b.Y+20 >= r.Height {
				delete(r.bullets, i)
				continue
			}
			b.Y += 20
		case TypeMoveDown:
			if b.Y-20 <= 0 {
				delete(r.bullets, i)
				continue
			}
			b.Y -= 20
		case TypeMoveLeft:
			if b.X-20 <= 0 {
				delete(r.bullets, i)
				continue
			}
			b.X -= 20
		case TypeMoveRight:
			if b.X+20 >= r.Width {
				delete(r.bullets, i)
				continue
			}
			b.X += 20
		}

		deleted := false
		for i, p := range r.players {
			if b.X >= p.X-10 && b.X <= p.X+10 && b.Y >= p.Y-10 && b.Y <= p.Y+10 {
				fmt.Println("player hit")
				deleted = true

				p.Life -= 1
				if p.Life == 0 {
					delete(r.players, i)
					continue
				}

				r.players[i] = p
				continue
			}
		}

		if deleted {
			delete(r.bullets, i)
		} else {
			r.bullets[i] = b
		}
	}
}

func (r *Room) Update() {
	ticker := time.NewTicker(time.Millisecond * 100)

	for {
		r.UpdateBullets()

		r.broadcast()
		<-ticker.C
	}
}
