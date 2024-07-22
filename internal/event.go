package internal

type Type string

const (
	TypeMoveUp    Type = "up"
	TypeMoveDown  Type = "down"
	TypeMoveLeft  Type = "left"
	TypeMoveRight Type = "right"
	TypeJoin      Type = "join"
	TypeShoot     Type = "shoot"
)

type Event struct {
	Type Type
}

type EventReturn struct {
	Players []*Player
	Bullets []Bullet
}
