package entities

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"terminal-roguelike/pkg/items"
)

type Skill struct {
	ID          string
	Name        string
	MPCost      int
	Description string
}

type Player struct {
	*Entity
	HP        int
	MaxHP     int
	MP        int
	MaxMP     int
	BaseATK   int
	BaseDEF   int
	Level     int
	EXP       int
	MaxEXP    int
	Gold      int
	Keys      int
	Kills     int
	Guarding  bool
	Skills    []Skill
	Inventory *items.Inventory
}

func NewPlayer(x, y int) *Player {
	inv := items.NewInventory(26)
	starterDagger := items.NewDagger(0, 0)
	_ = inv.Add(starterDagger)
	inv.Equip(starterDagger)

	starterPotion := items.NewHealthPotion(0, 0)
	_ = inv.Add(starterPotion)

	skills := []Skill{
		{ID: "heavy_slash", Name: "Heavy Slash", MPCost: 6, Description: "1.6x ATK physical cleave, ignores 50% DEF"},
		{ID: "fireball", Name: "Fireball", MPCost: 8, Description: "Blasts target for heavy magical fire damage"},
		{ID: "shield_guard", Name: "Shield Guard", MPCost: 4, Description: "Reduces next damage taken by 70%"},
		{ID: "holy_heal", Name: "Holy Heal", MPCost: 10, Description: "Restores 30 HP"},
	}

	return &Player{
		Entity:    NewEntity("player", "Hero", '@', tcell.ColorYellow, x, y, true),
		HP:        40,
		MaxHP:     40,
		MP:        25,
		MaxMP:     25,
		BaseATK:   6,
		BaseDEF:   2,
		Level:     1,
		EXP:       0,
		MaxEXP:    40,
		Gold:      0,
		Keys:      0,
		Kills:     0,
		Guarding:  false,
		Skills:    skills,
		Inventory: inv,
	}
}

func (p *Player) TotalATK() int {
	return p.BaseATK + p.Inventory.TotalBonusATK()
}

func (p *Player) TotalDEF() int {
	return p.BaseDEF + p.Inventory.TotalBonusDEF()
}

func (p *Player) Heal(amount int) int {
	oldHP := p.HP
	p.HP += amount
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	return p.HP - oldHP
}

func (p *Player) RestoreMP(amount int) int {
	oldMP := p.MP
	p.MP += amount
	if p.MP > p.MaxMP {
		p.MP = p.MaxMP
	}
	return p.MP - oldMP
}

func (p *Player) GainEXP(amount int) []string {
	messages := make([]string, 0)
	p.EXP += amount

	for p.EXP >= p.MaxEXP {
		p.EXP -= p.MaxEXP
		p.Level++
		p.MaxEXP = int(float64(p.MaxEXP) * 1.5)

		hpGain := 8
		mpGain := 6
		atkGain := 2
		defGain := 1

		p.MaxHP += hpGain
		p.HP += hpGain
		p.MaxMP += mpGain
		p.MP += mpGain
		p.BaseATK += atkGain
		p.BaseDEF += defGain

		messages = append(messages, fmt.Sprintf("LEVEL UP! Reached Level %d! (+%d HP, +%d MP, +%d ATK, +%d DEF)", p.Level, hpGain, mpGain, atkGain, defGain))
	}
	return messages
}
