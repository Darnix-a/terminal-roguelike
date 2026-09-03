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
	inv := items.NewInventory(12)

	// Starter Equipment
	starterWeapon := items.NewDagger(0, 0)
	_ = inv.Add(starterWeapon)
	inv.Equip(starterWeapon)

	starterArmor := items.NewLeatherArmor(0, 0)
	_ = inv.Add(starterArmor)
	inv.Equip(starterArmor)

	// 2 Starter Potions
	pot1 := items.NewHealthPotion(0, 0)
	pot1.HealAmount = 30
	pot1.Description = "Restores 30 HP"
	_ = inv.Add(pot1)

	pot2 := items.NewHealthPotion(0, 0)
	pot2.HealAmount = 30
	pot2.Description = "Restores 30 HP"
	_ = inv.Add(pot2)

	skills := []Skill{
		{ID: "heavy_slash", Name: "Heavy Slash", MPCost: 5, Description: "Cleave strike that pierces 50% enemy DEF"},
		{ID: "fireball", Name: "Fireball", MPCost: 7, Description: "Blast enemy with high mystical fire damage"},
		{ID: "shield_guard", Name: "Shield Guard", MPCost: 3, Description: "Brace stance that absorbs 75% incoming damage"},
		{ID: "holy_heal", Name: "Holy Heal", MPCost: 8, Description: "Chants sacred prayer restoring 35 HP"},
	}

	return &Player{
		Entity:    NewEntity("player", "Hero", '@', tcell.ColorYellow, x, y, true),
		HP:        50,
		MaxHP:     50,
		MP:        30,
		MaxMP:     30,
		BaseATK:   6,
		BaseDEF:   2,
		Level:     1,
		EXP:       0,
		MaxEXP:    35,
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
	prev := p.HP
	p.HP += amount
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	return p.HP - prev
}

func (p *Player) RestoreMP(amount int) int {
	prev := p.MP
	p.MP += amount
	if p.MP > p.MaxMP {
		p.MP = p.MaxMP
	}
	return p.MP - prev
}

func (p *Player) GainEXP(amount int) []string {
	p.EXP += amount
	messages := make([]string, 0)

	for p.EXP >= p.MaxEXP {
		p.EXP -= p.MaxEXP
		p.Level++
		p.MaxEXP = int(float64(p.MaxEXP) * 1.5)
		p.MaxHP += 10
		p.HP = p.MaxHP
		p.MaxMP += 5
		p.MP = p.MaxMP
		p.BaseATK += 2
		p.BaseDEF += 1

		messages = append(messages, fmt.Sprintf("🌟 LEVEL UP! You reached Level %d!", p.Level))
		messages = append(messages, fmt.Sprintf("Stats Increased: Max HP +10 (%d), Max MP +5 (%d), ATK +2 (%d), DEF +1 (%d)!", p.MaxHP, p.MaxMP, p.TotalATK(), p.TotalDEF()))
	}

	return messages
}
