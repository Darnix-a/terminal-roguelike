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
		HP:       14,
		MaxHP:    14,
		ATK:      5,
		DEF:      1,
		EXP:      25,
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
		HP:       22,
		MaxHP:    22,
		ATK:      8,
		DEF:      3,
		EXP:      45,
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
		HP:       38,
		MaxHP:    38,
		ATK:      12,
		DEF:      4,
		EXP:      70,
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
		HP:       36,
		MaxHP:    36,
		ATK:      15,
		DEF:      3,
		EXP:      95,
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
		HP:       38,
		MaxHP:    38,
		ATK:      13,
		DEF:      4,
		EXP:      80,
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
		HP:       70,
		MaxHP:    70,
		ATK:      16,
		DEF:      5,
		EXP:      200,
		IsBoss:   false,
		IsMerchant: true,
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
		{Name: "Dragon Claws", DamageMult: 1.3, Description: "rakes with fiery dragon claws"},
		{Name: "Tail Slam", DamageMult: 1.4, Description: "whips its colossal spiked tail against you"},
		{Name: "Inferno Breath", DamageMult: 2.2, IsMagic: true, IsTelegraphed: true, TelegraphWarning: "takes a massive breath, glowing with searing infernal MAGMA!", Description: "breathes a catastrophic inferno across the arena"},
		{Name: "Cataclysm Roar", DamageMult: 2.0, IsTelegraphed: true, TelegraphWarning: "rears back to unleash a world-shattering CATACLYSM ROAR!", Description: "unleashes a roar of pure annihilation"},
	}
	return &Monster{
		Entity:   NewEntity("dragon_boss", "Ancient Red Dragon", 'D', tcell.ColorRed, x, y, true),
		HP:       140,
		MaxHP:    140,
		ATK:      20,
		DEF:      7,
		EXP:      500,
		IsBoss:   true,
		Alerted:  true,
		Sprite:   sprite,
		Actions:  actions,
	}
}

// GenerateRandomMonster spawns monsters scaled smoothly to dungeon floor
func GenerateRandomMonster(x, y, floor int, rng *rand.Rand) *Monster {
	var m *Monster
	roll := rng.Intn(100)

	switch {
	case floor >= 4:
		if roll < 45 {
			m = NewDarkWizard(x, y)
		} else if roll < 80 {
			m = NewOrc(x, y)
		} else {
			m = NewSkeleton(x, y)
		}

	case floor >= 3:
		if roll < 50 {
			m = NewOrc(x, y)
		} else if roll < 80 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}

	case floor >= 2:
		if roll < 60 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}

	default: // Floor 1: Exclusively Goblin Scouts for a balanced onboarding
		m = NewGoblin(x, y)
	}

	// Smooth progressive floor stat bonuses
	if floor > 1 {
		floorBonusHP := (floor - 1) * 3
		floorBonusATK := (floor - 1) * 1
		floorBonusDEF := (floor - 1) / 2
		floorBonusEXP := (floor - 1) * 8

		m.MaxHP += floorBonusHP
		m.HP = m.MaxHP
		m.ATK += floorBonusATK
		m.DEF += floorBonusDEF
		m.EXP += floorBonusEXP
	}

	// Champion spawn chance scales smoothly (10% on F1 up to 25% on F4+)
	champChance := 10 + (floor * 3)
	if champChance > 25 {
		champChance = 25
	}

	if rng.Intn(100) < champChance {
		m.IsChampion = true
		affixes := []string{"[Fiery]", "[Vampiric]", "[Armored]", "[Frenzied]"}
		m.Affix = affixes[rng.Intn(len(affixes))]
		m.Name = fmt.Sprintf("%s %s", m.Affix, m.Name)
		m.MaxHP = int(float64(m.MaxHP) * 1.3)
		m.HP = m.MaxHP
		m.EXP = int(float64(m.EXP) * 1.5)

		switch m.Affix {
		case "[Fiery]":
			m.Color = tcell.ColorOrangeRed
			m.ATK += 2
		case "[Vampiric]":
			m.Color = tcell.ColorCrimson
			m.ATK += 2
		case "[Armored]":
			m.Color = tcell.ColorSilver
			m.DEF += 3
		case "[Frenzied]":
			m.Color = tcell.ColorMediumPurple
			m.ATK += 3
		}
	}

	return m
}
