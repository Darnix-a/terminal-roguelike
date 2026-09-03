package items

import (
	"fmt"
	"math/rand"
)

type Inventory struct {
	Items          []*Item
	MaxCapacity    int
	EquippedWeapon *Item
	EquippedArmor  *Item
}

func NewInventory(capacity int) *Inventory {
	if capacity <= 0 || capacity > 26 {
		capacity = 26
	}
	return &Inventory{
		Items:       make([]*Item, 0, capacity),
		MaxCapacity: capacity,
	}
}

func (inv *Inventory) Add(item *Item) error {
	if len(inv.Items) >= inv.MaxCapacity {
		return fmt.Errorf("inventory is full (max %d items)", inv.MaxCapacity)
	}
	inv.Items = append(inv.Items, item)
	return nil
}

func (inv *Inventory) Remove(index int) (*Item, error) {
	if index < 0 || index >= len(inv.Items) {
		return nil, fmt.Errorf("invalid inventory slot")
	}
	item := inv.Items[index]
	// If removing equipped item, unequip first
	if item.Equipped {
		inv.Unequip(item)
	}
	inv.Items = append(inv.Items[:index], inv.Items[index+1:]...)
	return item, nil
}

func (inv *Inventory) Equip(item *Item) string {
	if item.Type == TypeWeapon {
		if inv.EquippedWeapon != nil {
			inv.EquippedWeapon.Equipped = false
		}
		// Clear any orphaned equipped flags
		for _, itm := range inv.Items {
			if itm.Type == TypeWeapon {
				itm.Equipped = false
			}
		}
		inv.EquippedWeapon = item
		item.Equipped = true
		return fmt.Sprintf("Equipped %s (+%d ATK)", item.Name, item.BonusATK)
	} else if item.Type == TypeArmor {
		if inv.EquippedArmor != nil {
			inv.EquippedArmor.Equipped = false
		}
		for _, itm := range inv.Items {
			if itm.Type == TypeArmor {
				itm.Equipped = false
			}
		}
		inv.EquippedArmor = item
		item.Equipped = true
		return fmt.Sprintf("Equipped %s (+%d DEF)", item.Name, item.BonusDEF)
	}
	return "Cannot equip this item"
}

func (inv *Inventory) Unequip(item *Item) string {
	if item.Type == TypeWeapon && (item == inv.EquippedWeapon || item.Equipped) {
		inv.EquippedWeapon = nil
		item.Equipped = false
		return fmt.Sprintf("Unequipped %s", item.Name)
	}
	if item.Type == TypeArmor && (item == inv.EquippedArmor || item.Equipped) {
		inv.EquippedArmor = nil
		item.Equipped = false
		return fmt.Sprintf("Unequipped %s", item.Name)
	}
	return ""
}

func (inv *Inventory) TotalBonusATK() int {
	if inv.EquippedWeapon != nil {
		return inv.EquippedWeapon.BonusATK
	}
	return 0
}

func (inv *Inventory) TotalBonusDEF() int {
	if inv.EquippedArmor != nil {
		return inv.EquippedArmor.BonusDEF
	}
	return 0
}

// GenerateRandomItem generates an item scaled to dungeon floor level
func GenerateRandomItem(x, y, floor int, rng *rand.Rand) *Item {
	roll := rng.Intn(100)

	// Gold Drop
	if roll < 35 {
		goldAmt := rng.Intn(10*floor) + 5*floor
		return NewGold(x, y, goldAmt)
	}

	// Potion Drop
	if roll < 65 {
		if floor >= 3 && rng.Intn(2) == 0 {
			return NewGreaterHealthPotion(x, y)
		}
		return NewHealthPotion(x, y)
	}

	// Scroll Drop
	if roll < 80 {
		if rng.Intn(3) == 0 {
			return NewScrollEnchantment(x, y)
		}
		return NewScrollTeleport(x, y)
	}

	// Weapon or Armor Drop
	if rng.Intn(2) == 0 {
		// Weapon
		switch {
		case floor >= 4:
			return NewFlamingSword(x, y)
		case floor >= 3:
			return NewBattleaxe(x, y)
		case floor >= 2:
			return NewShortsword(x, y)
		default:
			return NewDagger(x, y)
		}
	} else {
		// Armor
		switch {
		case floor >= 4:
			return NewPlateArmor(x, y)
		case floor >= 3:
			return NewChainmail(x, y)
		default:
			return NewLeatherArmor(x, y)
		}
	}
}
