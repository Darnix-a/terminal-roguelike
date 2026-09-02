package combat

import (
	"math/rand"
	"testing"
)

func TestCalculateAttack(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	// Attacker 10 ATK vs Defender 3 DEF (7 base damage)
	res := CalculateAttack("Hero", 10, "Goblin", 3, 20, rng)

	if res.Damage <= 0 {
		t.Errorf("Expected damage > 0, got %d", res.Damage)
	}

	if res.TargetKilled {
		t.Errorf("Expected target not killed with 20 HP, got killed")
	}

	// Fatal blow test
	fatalRes := CalculateAttack("Hero", 30, "Goblin", 0, 10, rng)
	if !fatalRes.TargetKilled {
		t.Errorf("Expected target killed by 30 damage against 10 HP")
	}
}
