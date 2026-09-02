package combat

import (
	"fmt"
	"math/rand"
)

type AttackResult struct {
	Damage      int
	IsCritical  bool
	TargetKilled bool
	Message     string
}

// CalculateAttack computes damage dealt from attacker to defender with dice rolls
func CalculateAttack(attackerName string, attackerATK int, defenderName string, defenderDEF int, defenderHP int, rng *rand.Rand) AttackResult {
	// Base damage
	baseDamage := attackerATK - defenderDEF
	if baseDamage < 1 {
		baseDamage = 1
	}

	// Damage variance: +/- 15%
	variance := rng.Intn(3) - 1 // -1, 0, +1
	finalDamage := baseDamage + variance
	if finalDamage < 1 {
		finalDamage = 1
	}

	// Critical Hit check (10% chance)
	isCrit := rng.Intn(100) < 10
	if isCrit {
		finalDamage = int(float64(finalDamage) * 1.5)
		if finalDamage < 2 {
			finalDamage = 2
		}
	}

	remainingHP := defenderHP - finalDamage
	targetKilled := remainingHP <= 0

	var msg string
	if isCrit {
		msg = fmt.Sprintf("CRITICAL HIT! %s strikes %s for %d damage!", attackerName, defenderName, finalDamage)
	} else {
		msg = fmt.Sprintf("%s hits %s for %d damage.", attackerName, defenderName, finalDamage)
	}

	if targetKilled {
		msg += fmt.Sprintf(" %s is defeated!", defenderName)
	}

	return AttackResult{
		Damage:       finalDamage,
		IsCritical:   isCrit,
		TargetKilled: targetKilled,
		Message:      msg,
	}
}
