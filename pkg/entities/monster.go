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
		{Name: "Dagger Slash", DamageMult: 1.0, Description: "slashes with a jagged dagger"},
		{Name: "Poisoned Blade", DamageMult: 1.3, Description: "stabs with a venom-tipped blade"},
		{Name: "Frenzied Flurry", DamageMult: 1.5, Description: "strikes wildly with dual daggers"},
	}
	return &Monster{
		Entity:   NewEntity("goblin", "Goblin Scout", 'g', tcell.ColorGreen, x, y, true),
		HP:       18,
		MaxHP:    18,
		ATK:      7,
		DEF:      2,
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
		{Name: "Shield Bash", DamageMult: 1.4, Description: "bashes hard with a spiked iron shield"},
		{Name: "Cursed Cleave", DamageMult: 1.7, IsTelegraphed: true, TelegraphWarning: "channels dark necrotic runes into its blade for a CURSED CLEAVE!", Description: "unleashes a brutal necrotic strike"},
	}
	return &Monster{
		Entity:   NewEntity("skeleton", "Skeletal Guard", 's', tcell.ColorWhite, x, y, true),
		HP:       30,
		MaxHP:    30,
		ATK:      10,
		DEF:      4,
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
		{Name: "Heavy Cleave", DamageMult: 1.3, Description: "swings a giant battleaxe with crushing power"},
		{Name: "Skull Crusher", DamageMult: 2.2, IsTelegraphed: true, TelegraphWarning: "raises his colossal battleaxe overhead for a devastating SKULL CRUSHER!", Description: "leaps into the air and delivers a bone-shattering smash"},
		{Name: "Bloodrage Frenzy", DamageMult: 1.8, Description: "unleashes a bloodthirsty berserker combo"},
	}
	return &Monster{
		Entity:   NewEntity("orc", "Orc Berserker", 'o', tcell.ColorOlive, x, y, true),
		HP:       48,
		MaxHP:    48,
		ATK:      15,
		DEF:      5,
		EXP:      75,
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
		{Name: "Shadow Bolt", DamageMult: 1.4, IsMagic: true, Description: "casts a piercing bolt of dark void magic"},
		{Name: "Life Siphon", DamageMult: 1.2, IsMagic: true, Description: "drains vitality directly from your soul"},
		{Name: "Hellfire Hex", DamageMult: 2.3, IsMagic: true, IsTelegraphed: true, TelegraphWarning: "begins chanting ancient incantations... HELLFIRE HEX is charging!", Description: "erupts the ground in catastrophic cursed black flames"},
	}
	return &Monster{
		Entity:   NewEntity("dark_wizard", "Dark Sorcerer", 'w', tcell.ColorPurple, x, y, true),
		HP:       42,
		MaxHP:    42,
		ATK:      18,
		DEF:      3,
		EXP:      100,
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
		{Name: "Chest Chomp", DamageMult: 1.5, Description: "snaps its razor-toothed wooden jaws"},
		{Name: "Adhesive Lash", DamageMult: 1.3, Description: "whips with an adhesive sticky tongue"},
		{Name: "Gold Shrapnel", DamageMult: 2.0, IsTelegraphed: true, TelegraphWarning: "swallows deep and prepares a burst of GOLD SHRAPNEL!", Description: "blasts jagged heavy metal coins at point-blank"},
	}
	return &Monster{
		Entity:   NewEntity("mimic", "Hungry Mimic", 'M', tcell.ColorGold, x, y, true),
		HP:       45,
		MaxHP:    45,
		ATK:      16,
		DEF:      5,
		EXP:      90,
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
		{Name: "Coin Gatling", DamageMult: 2.2, IsTelegraphed: true, TelegraphWarning: "loads a sack of pure gold coins for a COIN GATLING assault!", Description: "unleashes a relentless flurry of hypersonic golden coins"},
		{Name: "Vault Cleave", DamageMult: 1.6, Description: "swings an enchanted vault crowbar"},
	}
	return &Monster{
		Entity:   NewEntity("shopkeeper_boss", "Enraged Shopkeeper", 'S', tcell.ColorGold, x, y, true),
		HP:       85,
		MaxHP:    85,
		ATK:      19,
		DEF:      6,
		EXP:      250,
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
		{Name: "Dragon Claws", DamageMult: 1.4, Description: "rakes with razor-sharp fiery claws"},
		{Name: "Tail Slam", DamageMult: 1.5, Description: "whips its colossal spiked tail against you"},
		{Name: "Inferno Breath", DamageMult: 2.5, IsMagic: true, IsTelegraphed: true, TelegraphWarning: "takes a massive breath, glowing with searing infernal MAGMA!", Description: "breathes a catastrophic inferno across the arena"},
		{Name: "Cataclysm Roar", DamageMult: 2.2, IsTelegraphed: true, TelegraphWarning: "rears back to unleash a world-shattering CATACLYSM ROAR!", Description: "unleashes a roar of pure annihilation"},
	}
	return &Monster{
		Entity:   NewEntity("dragon_boss", "Ancient Red Dragon", 'D', tcell.ColorRed, x, y, true),
		HP:       180,
		MaxHP:    180,
		ATK:      24,
		DEF:      9,
		EXP:      500,
		IsBoss:   true,
		Alerted:  true,
		Sprite:   sprite,
		Actions:  actions,
	}
}

// GenerateRandomMonster spawns monsters scaled to dungeon floor with 25% Champion chance
func GenerateRandomMonster(x, y, floor int, rng *rand.Rand) *Monster {
	var m *Monster
	roll := rng.Intn(100)

	switch {
	case floor >= 4:
		if roll < 40 {
			m = NewDarkWizard(x, y)
		} else if roll < 75 {
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
		if roll < 55 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}

	default:
		if roll < 25 {
			m = NewSkeleton(x, y)
		} else {
			m = NewGoblin(x, y)
		}
	}

	// 25% Chance for Monster to be an Elite Champion
	if rng.Intn(100) < 25 {
		m.IsChampion = true
		affixes := []string{"[Fiery]", "[Vampiric]", "[Armored]", "[Frenzied]"}
		m.Affix = affixes[rng.Intn(len(affixes))]
		m.Name = fmt.Sprintf("%s %s", m.Affix, m.Name)
		m.MaxHP = int(float64(m.MaxHP) * 1.5)
		m.HP = m.MaxHP
		m.EXP = int(float64(m.EXP) * 2.0)

		switch m.Affix {
		case "[Fiery]":
			m.Color = tcell.ColorOrangeRed
			m.ATK += 4
		case "[Vampiric]":
			m.Color = tcell.ColorCrimson
			m.ATK += 3
		case "[Armored]":
			m.Color = tcell.ColorSilver
			m.DEF += 5
		case "[Frenzied]":
			m.Color = tcell.ColorMediumPurple
			m.ATK += 5
		}
	}

	return m
}
