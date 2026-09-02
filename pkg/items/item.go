package items

import (
	"github.com/gdamore/tcell/v2"
)

type ItemType int

const (
	TypePotion ItemType = iota
	TypeWeapon
	TypeArmor
	TypeScroll
	TypeGold
	TypeKey
)

type Item struct {
	ID          string
	Name        string
	Type        ItemType
	Rune        rune
	Color       tcell.Color
	HealAmount  int
	BonusATK    int
	BonusDEF    int
	Value       int
	X           int
	Y           int
	Description string
	Equipped    bool
}

// Preset Item Constructors

func NewHealthPotion(x, y int) *Item {
	return &Item{
		ID:          "health_potion",
		Name:        "Health Potion",
		Type:        TypePotion,
		Rune:        '!',
		Color:       tcell.ColorRed,
		HealAmount:  25,
		X:           x,
		Y:           y,
		Description: "Restores 25 HP",
	}
}

func NewGreaterHealthPotion(x, y int) *Item {
	return &Item{
		ID:          "greater_health_potion",
		Name:        "Greater Health Potion",
		Type:        TypePotion,
		Rune:        '!',
		Color:       tcell.ColorDarkRed,
		HealAmount:  50,
		X:           x,
		Y:           y,
		Description: "Restores 50 HP",
	}
}

func NewDagger(x, y int) *Item {
	return &Item{
		ID:          "dagger",
		Name:        "Iron Dagger",
		Type:        TypeWeapon,
		Rune:        ')',
		Color:       tcell.ColorLightGray,
		BonusATK:    3,
		X:           x,
		Y:           y,
		Description: "+3 Attack Power",
	}
}

func NewShortsword(x, y int) *Item {
	return &Item{
		ID:          "shortsword",
		Name:        "Steel Shortsword",
		Type:        TypeWeapon,
		Rune:        ')',
		Color:       tcell.ColorSkyblue,
		BonusATK:    6,
		X:           x,
		Y:           y,
		Description: "+6 Attack Power",
	}
}

func NewBattleaxe(x, y int) *Item {
	return &Item{
		ID:          "battleaxe",
		Name:        "Mithril Battleaxe",
		Type:        TypeWeapon,
		Rune:        ')',
		Color:       tcell.ColorDarkCyan,
		BonusATK:    10,
		X:           x,
		Y:           y,
		Description: "+10 Attack Power",
	}
}

func NewFlamingSword(x, y int) *Item {
	return &Item{
		ID:          "flaming_sword",
		Name:        "Flaming Claymore",
		Type:        TypeWeapon,
		Rune:        ')',
		Color:       tcell.ColorOrangeRed,
		BonusATK:    15,
		X:           x,
		Y:           y,
		Description: "+15 Attack Power (Infused with Fire)",
	}
}

func NewLeatherArmor(x, y int) *Item {
	return &Item{
		ID:          "leather_armor",
		Name:        "Leather Armor",
		Type:        TypeArmor,
		Rune:        '[',
		Color:       tcell.ColorSandyBrown,
		BonusDEF:    2,
		X:           x,
		Y:           y,
		Description: "+2 Defense",
	}
}

func NewChainmail(x, y int) *Item {
	return &Item{
		ID:          "chainmail",
		Name:        "Chainmail Armor",
		Type:        TypeArmor,
		Rune:        '[',
		Color:       tcell.ColorSilver,
		BonusDEF:    5,
		X:           x,
		Y:           y,
		Description: "+5 Defense",
	}
}

func NewPlateArmor(x, y int) *Item {
	return &Item{
		ID:          "plate_armor",
		Name:        "Dragonscale Plate",
		Type:        TypeArmor,
		Rune:        '[',
		Color:       tcell.ColorGold,
		BonusDEF:    9,
		X:           x,
		Y:           y,
		Description: "+9 Defense (Forged from Dragon Scales)",
	}
}

func NewScrollTeleport(x, y int) *Item {
	return &Item{
		ID:          "scroll_teleport",
		Name:        "Scroll of Teleportation",
		Type:        TypeScroll,
		Rune:        '?',
		Color:       tcell.ColorViolet,
		X:           x,
		Y:           y,
		Description: "Teleports you to a random open floor tile",
	}
}

func NewGold(x, y, amount int) *Item {
	return &Item{
		ID:          "gold",
		Name:        "Gold Coins",
		Type:        TypeGold,
		Rune:        '$',
		Color:       tcell.ColorYellow,
		Value:       amount,
		X:           x,
		Y:           y,
		Description: "Dungeon treasure",
	}
}

func NewDungeonKey(x, y int) *Item {
	return &Item{
		ID:          "dungeon_key",
		Name:        "Iron Dungeon Key",
		Type:        TypeKey,
		Rune:        'k',
		Color:       tcell.ColorGold,
		X:           x,
		Y:           y,
		Description: "Unlocks heavy iron vault doors",
	}
}
