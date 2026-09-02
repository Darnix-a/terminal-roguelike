package items

import (
	"testing"
)

func TestInventoryEquipAndBonuses(t *testing.T) {
	inv := NewInventory(5)

	sword := NewShortsword(0, 0)
	shield := NewLeatherArmor(0, 0)

	if err := inv.Add(sword); err != nil {
		t.Fatalf("Failed to add sword: %v", err)
	}
	if err := inv.Add(shield); err != nil {
		t.Fatalf("Failed to add shield: %v", err)
	}

	inv.Equip(sword)
	inv.Equip(shield)

	if inv.TotalBonusATK() != sword.BonusATK {
		t.Errorf("Expected bonus ATK %d, got %d", sword.BonusATK, inv.TotalBonusATK())
	}

	if inv.TotalBonusDEF() != shield.BonusDEF {
		t.Errorf("Expected bonus DEF %d, got %d", shield.BonusDEF, inv.TotalBonusDEF())
	}

	// Unequip test
	inv.Unequip(sword)
	if inv.TotalBonusATK() != 0 {
		t.Errorf("Expected bonus ATK 0 after unequip, got %d", inv.TotalBonusATK())
	}
}
