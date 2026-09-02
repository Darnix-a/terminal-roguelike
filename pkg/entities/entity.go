package entities

import (
	"github.com/gdamore/tcell/v2"
)

type Entity struct {
	ID             string
	Name           string
	Rune           rune
	Color          tcell.Color
	X              int
	Y              int
	BlocksMovement bool
	IsAlive        bool
}

func NewEntity(id, name string, r rune, color tcell.Color, x, y int, blocks bool) *Entity {
	return &Entity{
		ID:             id,
		Name:           name,
		Rune:           r,
		Color:          color,
		X:              x,
		Y:              y,
		BlocksMovement: blocks,
		IsAlive:        true,
	}
}

func (e *Entity) Move(dx, dy int) {
	e.X += dx
	e.Y += dy
}

func (e *Entity) MoveTo(x, y int) {
	e.X = x
	e.Y = y
}
