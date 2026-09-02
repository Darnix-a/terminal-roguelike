package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

type MonsterAction struct {
	Name        string
	DamageMult  float64
	IsMagic     bool
	Description string
}

type Monster struct {
	*Entity
	HP         int
	MaxHP      int
	ATK        int
	DEF        int
	EXP        int
	IsBoss     bool
	IsChampion bool
	Affix      string // "[Fiery]", "[Vampiric]", "[Armored]", "[Frenzied]"
	Alerted    bool
	Guarding   bool
	Sprite     []string
	Actions    []MonsterAction
}

func NewGoblin(x, y int) *Monster {
	sprite := []string{
		`  (o.o)  `,
		`  /|><|\ `,
		`   /  \  `,
	}
	actions := []MonsterAction{
		{Name: "Dagger Slash", DamageMult: 1.0, Description: "slashes with a jagged dagger"},
		{Name: "Quick Bite", DamageMult: 1.2, Description: "lunges forward for a nasty bite"},
		{Name: "Throw Sand", DamageMult: 0.8, Description: "kicks dirty sand to distract"},
	}
	return &Monster{
		Entity:   NewEntity("goblin", "Goblin Scout", 'g', tcell.ColorGreen, x, y, true),
		HP:       14,
		MaxHP:    14,
		ATK:      5,
		DEF:      1,
		EXP:      20,
		IsBoss:   false,
		Alerted:  false,
		Sprite:   sprite,
		Actions:  actions,
	}
}

func NewSkeleton(x, y int) *Monster {
	sprite := []string{
		`  [o_o]  `,
		` --|= |--`,
		`   | |   `,
	}
	actions := []MonsterAction{
		{Name: "Rusty Blade", DamageMult: 1.0, Description: "swings an ancient rusted blade"},
		{Name: "Shield Slam", DamageMult: 1.3, Description: "bashes with a cracked iron shield"},
		{Name: "Bone Rattle", DamageMult: 0.9, Description: "unleashes a chilling necrotic aura"},
	}
	return &Monster{
		Entity:   NewEntity("skeleton", "Skeletal Guard", 's', tcell.ColorWhite, x, y, true),
		HP:       22,
		MaxHP:    22,
		ATK:      7,
		DEF:      3,
		EXP:      35,
		IsBoss:   false,
		Alerted:  false,
		Sprite:   sprite,
		Actions:  actions,
	}
}

func NewOrc(x, y int) *Monster {
	sprite := []string{
		`  (Ò_Ó)  `,
		` /(===)\ `,
		`  /   \  `,
	}
	actions := []MonsterAction{
		{Name: "Heavy Cleave", DamageMult: 1.2, Description: "swings a giant battleaxe"},
		{Name: "War Stomp", DamageMult: 1.4, Description: "stomps the ground, crushing your armor"},
		{Name: "Berserk Roar", DamageMult: 1.1, Description: "enters a bloodthirsty frenzy"},
	}
	return &Monster{
		Entity:   NewEntity("orc", "Orc Berserker", 'o', tcell.ColorOlive, x, y, true),
		HP:       34,
		MaxHP:    34,
		ATK:      10,
		DEF:      4,
		EXP:      55,
		IsBoss:   false,
		Alerted:  false,
		Sprite:   sprite,
		Actions:  actions,
	}
}

func NewDarkWizard(x, y int) *Monster {
	sprite := []string{
		`   /^\   `,
		`  (0_0)  `,
		`  <[~]>* `,
		`   / \   `,
	}
	actions := []MonsterAction{
		{Name: "Shadow Bolt", DamageMult: 1.3, IsMagic: true, Description: "casts a bolt of dark magic"},
		{Name: "Life Siphon", DamageMult: 1.0, IsMagic: true, Description: "drains vitality from your soul"},
		{Name: "Flame Hex", DamageMult: 1.5, IsMagic: true, Description: "ignites the air with cursed fire"},
	}
	return &Monster{
		Entity:   NewEntity("dark_wizard", "Dark Sorcerer", 'w', tcell.ColorPurple, x, y, true),
		HP:       30,
		MaxHP:    30,
		ATK:      13,
		DEF:      2,
		EXP:      80,
		IsBoss:   false,
		Alerted:  false,
		Sprite:   sprite,
		Actions:  actions,
	}
}

func NewMimic(x, y int) *Monster {
	sprite := []string{
		` [========] `,
		`  \VVVVVV/  `,
		`  ( O  O )  `,
		`  /^^^^^^\  `,
	}
	actions := []MonsterAction{
		{Name: "Chest Chomp", DamageMult: 1.4, Description: "snaps its razor-toothed jaws"},
		{Name: "Tongue Lash", DamageMult: 1.1, Description: "whips with an adhesive sticky tongue"},
		{Name: "Gold Spray", DamageMult: 1.2, Description: "vomits cursed heavy coins"},
	}
	return &Monster{
		Entity:   NewEntity("mimic", "Hungry Mimic", 'M', tcell.ColorGold, x, y, true),
		HP:       32,
		MaxHP:    32,
		ATK:      12,
		DEF:      4,
		EXP:      75,
		IsBoss:   false,
		Alerted:  true,
		Sprite:   sprite,
		Actions:  actions,
	}
}

func NewDragonBoss(x, y int) *Monster {
	sprite := []string{
		`    /\__/\    `,
		`   / O  O \   `,
		`  =(  --  )=  `,
		`  /| ==== |\  `,
		` / |======| \ `,
		`    /    \    `,
	}
	actions := []MonsterAction{
		{Name: "Inferno Breath", DamageMult: 1.6, IsMagic: true, Description: "breathes a catastrophic cone of dragon fire"},
		{Name: "Tail Sweep", DamageMult: 1.3, Description: "whips its colossal spiked tail"},
		{Name: "Dragon Claws", DamageMult: 1.2, Description: "rakes with razor-sharp fiery claws"},
		{Name: "Roar of Ruin", DamageMult: 1.1, Description: "shakes the entire dungeon chamber"},
	}
	return &Monster{
		Entity:   NewEntity("dragon_boss", "Ancient Red Dragon", 'D', tcell.ColorRed, x, y, true),
		HP:       120,
		MaxHP:    120,
		ATK:      18,
		DEF:      7,
		EXP:      400,
		IsBoss:   true,
		Alerted:  true,
		Sprite:   sprite,
		Actions:  actions,
	}
}

// GenerateRandomMonster spawns monsters scaled to dungeon floor with 20% Champion chance
func GenerateRandomMonster(x, y, floor int, rng *rand.Rand) *Monster {
	var m *Monster
	roll := rng.Intn(100)

	switch {
	case floor >= 4:
		if roll < 35 {
			m = NewDarkWizard(x, y)
		} else if roll < 70 {
			m = NewOrc(x, y)
		} else {
			m = NewSkeleton(x, y)
		}

	case floor >= 3:
		if roll < 45 {
			m = NewOrc(x, y)
		} else if roll < 80 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}

	case floor >= 2:
		if roll < 50 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}

	default:
		if roll < 20 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}
	}

	// 20% Chance for Monster to be an Elite Champion
	if rng.Intn(100) < 20 {
		m.IsChampion = true
		affixes := []string{"[Fiery]", "[Vampiric]", "[Armored]", "[Frenzied]"}
		m.Affix = affixes[rng.Intn(len(affixes))]
		m.Name = fmt.Sprintf("%s %s", m.Affix, m.Name)
		m.MaxHP = int(float64(m.MaxHP) * 1.4)
		m.HP = m.MaxHP
		m.EXP = int(float64(m.EXP) * 1.8)

		switch m.Affix {
		case "[Fiery]":
			m.Color = tcell.ColorOrangeRed
			m.ATK += 3
		case "[Vampiric]":
			m.Color = tcell.ColorCrimson
			m.ATK += 2
		case "[Armored]":
			m.Color = tcell.ColorSilver
			m.DEF += 4
		case "[Frenzied]":
			m.Color = tcell.ColorMediumPurple
			m.ATK += 4
		}
	}

	return m
}
