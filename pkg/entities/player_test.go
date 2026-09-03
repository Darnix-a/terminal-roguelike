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
	leveledUp, msgs := player.GainEXP(player.MaxEXP)
	if !leveledUp || len(msgs) == 0 {
		t.Errorf("Expected level up to be true with messages")
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

	// Test Learn Skill up to MaxSkills = 5
	for _, sk := range SkillRegistry {
		player.LearnSkill(sk)
	}
	if len(player.Skills) != MaxSkills {
		t.Errorf("Expected player skills to be capped at %d, got %d", MaxSkills, len(player.Skills))
	}

	// Test ReplaceSkill
	newSkill := Skill{ID: "test_skill", Name: "Test Skill", MPCost: 9, Description: "Test"}
	success := player.ReplaceSkill(0, newSkill)
	if !success || player.Skills[0].ID != "test_skill" {
		t.Errorf("Failed to replace skill at slot 0")
	}
}
