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

func TestAttackTelegraphAndGuard(t *testing.T) {
	rng := rand.New(rand.NewSource(100))
	player := entities.NewPlayer(0, 0)
	orc := entities.NewOrc(0, 0)

	battle := NewBattle(player, orc, rng)

	// Set a telegraphed action manually
	telegraphedMove := orc.Actions[1] // Skull Crusher
	battle.TelegraphedAction = &telegraphedMove

	// Player guards
	player.Guarding = true
	initialHP := player.HP

	battle.EnemyTurn()

	// Should have executed the telegraphed move and mitigated 70% damage
	if battle.TelegraphedAction != nil {
		t.Errorf("Expected telegraphed action to be consumed, got %+v", battle.TelegraphedAction)
	}

	damageTaken := initialHP - player.HP
	if damageTaken <= 0 {
		t.Errorf("Expected some damage taken, got %d", damageTaken)
	}
}
