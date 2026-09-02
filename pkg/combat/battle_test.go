package combat

import (
	"math/rand"
	"testing"

	"terminal-roguelike/pkg/entities"
)

func TestTurnBasedBattle(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	player := entities.NewPlayer(0, 0)
	goblin := entities.NewGoblin(0, 0)

	battle := NewBattle(player, goblin, rng)

	if battle.Result != BattleOngoing {
		t.Fatalf("Expected battle to start ongoing, got %v", battle.Result)
	}

	// Test Skill (Heavy Slash)
	initialMP := player.MP
	battle.PlayerUseSkill(0) // Heavy Slash (6 MP)

	if player.MP != initialMP-6 {
		t.Errorf("Expected MP to decrease by 6, got %d", player.MP)
	}

	if goblin.HP >= goblin.MaxHP {
		t.Errorf("Expected goblin HP to decrease after Heavy Slash, got %d", goblin.HP)
	}

	// Slay the goblin
	for goblin.HP > 0 && battle.Result == BattleOngoing {
		battle.PlayerAttack()
	}

	if battle.Result != BattleVictory {
		t.Errorf("Expected BattleVictory after goblin defeated, got %v", battle.Result)
	}
}
