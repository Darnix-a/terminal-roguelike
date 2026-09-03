package items

import (
	"testing"
)

func TestInventoryEquipAndBonuses(t *testing.T) {
	inv := NewInventory(5)

	dagger := NewDagger(0, 0)
	sword := NewShortsword(0, 0)
	armor := NewLeatherArmor(0, 0)

	_ = inv.Add(dagger)
	_ = inv.Add(sword)
	_ = inv.Add(armor)

	// Equip starter dagger & armor
	inv.Equip(dagger)
	inv.Equip(armor)

	if !dagger.Equipped || inv.EquippedWeapon != dagger {
		t.Fatalf("Expected dagger to be equipped")
	}

	if inv.TotalBonusATK() != dagger.BonusATK {
		t.Errorf("Expected bonus ATK %d, got %d", dagger.BonusATK, inv.TotalBonusATK())
	}

	// Equip sword -> should unequip dagger
	inv.Equip(sword)

	if dagger.Equipped {
		t.Errorf("Expected dagger to be unequipped when equipping sword")
	}

	if !sword.Equipped || inv.EquippedWeapon != sword {
		t.Errorf("Expected sword to be equipped")
	}

	if inv.TotalBonusATK() != sword.BonusATK {
		t.Errorf("Expected bonus ATK %d, got %d", sword.BonusATK, inv.TotalBonusATK())
	}

	// Unequip test
	inv.Unequip(sword)
	if sword.Equipped || inv.EquippedWeapon != nil {
		t.Errorf("Expected sword to be unequipped")
	}
	if inv.TotalBonusATK() != 0 {
		t.Errorf("Expected bonus ATK 0 after unequip, got %d", inv.TotalBonusATK())
	}
}
