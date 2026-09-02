package entities

import (
	"testing"
)

func TestPlayerLevelUpAndHealing(t *testing.T) {
	player := NewPlayer(5, 5)

	initialMaxHP := player.MaxHP
	initialATK := player.BaseATK

	// Damage player
	player.HP -= 15
	if player.HP != initialMaxHP-15 {
		t.Errorf("Expected HP %d, got %d", initialMaxHP-15, player.HP)
	}

	// Heal player
	healed := player.Heal(10)
	if healed != 10 || player.HP != initialMaxHP-5 {
		t.Errorf("Healing failed: got healed=%d, hp=%d", healed, player.HP)
	}

	// Level up
	msgs := player.GainEXP(player.MaxEXP)
	if len(msgs) == 0 {
		t.Errorf("Expected level up messages, got none")
	}

	if player.Level != 2 {
		t.Errorf("Expected level 2, got %d", player.Level)
	}

	if player.MaxHP <= initialMaxHP {
		t.Errorf("Expected MaxHP to increase after level up")
	}

	if player.BaseATK <= initialATK {
		t.Errorf("Expected BaseATK to increase after level up")
	}
}
