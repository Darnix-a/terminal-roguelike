package entities

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

type MonsterAction struct {
	Name             string
	DamageMult       float64
	IsMagic          bool
	IsTelegraphed    bool
	TelegraphWarning string
	Description      string
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
	IsMerchant bool
	Affix      string
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
		{Name: "Dagger Slash", DamageMult: 1.0, Description: "slashes with a small rusty dagger"},
		{Name: "Poisoned Blade", DamageMult: 1.2, Description: "stabs with a venom-tipped blade"},
		{Name: "Frenzied Flurry", DamageMult: 1.3, Description: "strikes quickly with dual daggers"},
	}
	return &Monster{
		Entity:   NewEntity("goblin", "Goblin Scout", 'g', tcell.ColorGreen, x, y, true),
		HP:       18,
		MaxHP:    18,
		ATK:      5,
		DEF:      2,
		EXP:      12,
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
		{Name: "Rusty Blade", DamageMult: 1.1, Description: "swings an ancient rusted blade"},
		{Name: "Shield Bash", DamageMult: 1.3, Description: "bashes with a cracked iron shield"},
		{Name: "Cursed Cleave", DamageMult: 1.6, IsTelegraphed: true, TelegraphWarning: "channels dark necrotic runes for a CURSED CLEAVE!", Description: "unleashes a heavy necrotic strike"},
	}
	return &Monster{
		Entity:   NewEntity("skeleton", "Skeletal Guard", 's', tcell.ColorWhite, x, y, true),
		HP:       32,
		MaxHP:    32,
		ATK:      9,
		DEF:      4,
		EXP:      22,
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
		{Name: "Skull Crusher", DamageMult: 1.9, IsTelegraphed: true, TelegraphWarning: "raises his battleaxe overhead for a devastating SKULL CRUSHER!", Description: "leaps into the air for a heavy smash"},
		{Name: "Bloodrage Frenzy", DamageMult: 1.5, Description: "unleashes a wild berserker combo"},
	}
	return &Monster{
		Entity:   NewEntity("orc", "Orc Berserker", 'o', tcell.ColorOlive, x, y, true),
		HP:       52,
		MaxHP:    52,
		ATK:      14,
		DEF:      5,
		EXP:      40,
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
		{Name: "Shadow Bolt", DamageMult: 1.3, IsMagic: true, Description: "casts a bolt of dark void magic"},
		{Name: "Life Siphon", DamageMult: 1.1, IsMagic: true, Description: "drains vitality from your soul"},
		{Name: "Hellfire Hex", DamageMult: 2.0, IsMagic: true, IsTelegraphed: true, TelegraphWarning: "begins chanting ancient incantations... HELLFIRE HEX is charging!", Description: "erupts the ground in cursed black flames"},
	}
	return &Monster{
		Entity:   NewEntity("dark_wizard", "Dark Sorcerer", 'w', tcell.ColorPurple, x, y, true),
		HP:       46,
		MaxHP:    46,
		ATK:      17,
		DEF:      4,
		EXP:      60,
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
		{Name: "Chest Chomp", DamageMult: 1.4, Description: "snaps its razor-toothed wooden jaws"},
		{Name: "Adhesive Lash", DamageMult: 1.2, Description: "whips with an adhesive sticky tongue"},
		{Name: "Gold Shrapnel", DamageMult: 1.8, IsTelegraphed: true, TelegraphWarning: "swallows deep and prepares a burst of GOLD SHRAPNEL!", Description: "blasts jagged heavy metal coins at point-blank"},
	}
	return &Monster{
		Entity:   NewEntity("mimic", "Hungry Mimic", 'M', tcell.ColorGold, x, y, true),
		HP:       45,
		MaxHP:    45,
		ATK:      14,
		DEF:      5,
		EXP:      45,
		IsBoss:   false,
		Alerted:  true,
		Sprite:   sprite,
		Actions:  actions,
	}
}

func NewEnragedShopkeeper(x, y int) *Monster {
	sprite := []string{
		`  [$$__$$]  `,
		`  /|====|\  `,
		`   | || |   `,
	}
	actions := []MonsterAction{
		{Name: "Ledger Strike", DamageMult: 1.2, Description: "slams you with a heavy gilded accounting tome"},
		{Name: "Coin Gatling", DamageMult: 2.0, IsTelegraphed: true, TelegraphWarning: "loads a sack of pure gold coins for a COIN GATLING assault!", Description: "unleashes a flurry of hypersonic golden coins"},
		{Name: "Vault Cleave", DamageMult: 1.5, Description: "swings an enchanted vault crowbar"},
	}
	return &Monster{
		Entity:   NewEntity("shopkeeper_boss", "Enraged Shopkeeper", 'S', tcell.ColorGold, x, y, true),
		HP:       90,
		MaxHP:    90,
		ATK:      17,
		DEF:      6,
		EXP:      150,
		IsBoss:   true,
		IsMerchant: true,
		Alerted:  true,
		Sprite:   sprite,
		Actions:  actions,
	}
}

func NewDragonBoss(x, y int) *Monster {
	sprite := []string{
		`  / \  //\    `,
		` <^.^> (oo)   `,
		` /|/ \ |/ | \ `,
		`   / / \ \    `,
	}
	actions := []MonsterAction{
		{Name: "Dragon Claw", DamageMult: 1.2, Description: "rakes with razor-sharp obsidian talons"},
		{Name: "Inferno Breath", DamageMult: 2.2, IsMagic: true, IsTelegraphed: true, TelegraphWarning: "inhales deeply... INFERNO BREATH incoming!", Description: "unleashes an apocalyptic torrent of scorching dragonfire"},
		{Name: "Tail Sweep", DamageMult: 1.4, Description: "swipes with its massive spiked tail"},
		{Name: "Cataclysm Roar", DamageMult: 2.5, IsMagic: true, IsTelegraphed: true, TelegraphWarning: "gathers draconic thunder for a CATACLYSM ROAR!", Description: "shatters the dungeon walls with a terrifying shockwave"},
	}
	return &Monster{
		Entity:   NewEntity("dragon", "Ancient Red Dragon", 'D', tcell.ColorRed, x, y, true),
		HP:       200,
		MaxHP:    200,
		ATK:      24,
		DEF:      9,
		EXP:      400,
		IsBoss:   true,
		Alerted:  false,
		Sprite:   sprite,
		Actions:  actions,
	}
}

// GenerateRandomMonster spawns an appropriate monster for the floor level
func GenerateRandomMonster(x, y, floor int, rng *rand.Rand) *Monster {
	var m *Monster

	switch floor {
	case 1:
		// Floor 1: Exclusively Goblin Scouts
		m = NewGoblin(x, y)
	case 2:
		// Floor 2: 70% Skeletons, 30% Goblins
		if rng.Intn(100) < 70 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}
	case 3:
		// Floor 3: 60% Orcs, 30% Skeletons, 10% Goblins
		roll := rng.Intn(100)
		if roll < 60 {
			m = NewOrc(x, y)
		} else if roll < 90 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}
	case 4:
		// Floor 4: 50% Sorcerers, 35% Orcs, 15% Skeletons
		roll := rng.Intn(100)
		if roll < 50 {
			m = NewDarkWizard(x, y)
		} else if roll < 85 {
			m = NewOrc(x, y)
		} else {
			m = NewSkeleton(x, y)
		}
	default:
		m = NewGoblin(x, y)
	}

	// Champion spawn chance (10% to 20%)
	champChance := 10 + (floor * 2)
	if champChance > 20 {
		champChance = 20
	}

	if rng.Intn(100) < champChance {
		m.IsChampion = true
		affixes := []string{"[Fiery]", "[Vampiric]", "[Armored]", "[Frenzied]"}
		m.Affix = affixes[rng.Intn(len(affixes))]
		m.Name = fmt.Sprintf("%s %s", m.Affix, m.Name)
		m.MaxHP = int(float64(m.MaxHP) * 1.3)
		m.HP = m.MaxHP
		m.EXP = int(float64(m.EXP) * 1.4)

		switch m.Affix {
		case "[Fiery]":
			m.Color = tcell.ColorOrangeRed
			m.ATK += 2
		case "[Vampiric]":
			m.Color = tcell.ColorCrimson
			m.ATK += 2
		case "[Armored]":
			m.Color = tcell.ColorSilver
			m.DEF += 2
		case "[Frenzied]":
			m.Color = tcell.ColorMediumPurple
			m.ATK += 3
		}
	}

	return m
}
