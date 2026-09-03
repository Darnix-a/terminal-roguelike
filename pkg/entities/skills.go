package entities

import (
	"math/rand"
)

type Skill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MPCost      int    `json:"mp_cost"`
	Description string `json:"description"`
}

var SkillRegistry = []Skill{
	{
		ID:          "heavy_slash",
		Name:        "Heavy Slash",
		MPCost:      5,
		Description: "Brutal cleave that pierces 50% enemy DEF (1.8x ATK)",
	},
	{
		ID:          "shield_guard",
		Name:        "Shield Guard",
		MPCost:      3,
		Description: "Defensive stance that absorbs 75% incoming damage",
	},
	{
		ID:          "fireball",
		Name:        "Fireball",
		MPCost:      7,
		Description: "Hurls an exploding fiery orb of heavy magic damage",
	},
	{
		ID:          "holy_heal",
		Name:        "Holy Heal",
		MPCost:      8,
		Description: "Chants a sacred blessing restoring 35 HP",
	},
	{
		ID:          "chain_lightning",
		Name:        "Chain Lightning",
		MPCost:      10,
		Description: "Arcs piercing lightning that completely ignores enemy DEF",
	},
	{
		ID:          "frost_nova",
		Name:        "Frost Nova",
		MPCost:      6,
		Description: "Freezes enemy for magic damage and weakens enemy ATK by 3",
	},
	{
		ID:          "vampiric_strike",
		Name:        "Vampiric Strike",
		MPCost:      7,
		Description: "Slashes enemy and converts 50% of damage dealt into HP",
	},
	{
		ID:          "berserker_rage",
		Name:        "Berserker Rage",
		MPCost:      4,
		Description: "Sacrifices 6 HP to unleash a devastating 2.4x physical blow",
	},
	{
		ID:          "poison_blade",
		Name:        "Poison Blade",
		MPCost:      5,
		Description: "Strikes with venom, inflicting +12 poison damage over 3 turns",
	},
	{
		ID:          "smoke_bomb",
		Name:        "Smoke Bomb",
		MPCost:      4,
		Description: "Blinds enemy, negating next attack and granting 100% crit next turn",
	},
	{
		ID:          "divine_smite",
		Name:        "Divine Smite",
		MPCost:      12,
		Description: "Calls down celestial fury for colossal holy damage (2.8x)",
	},
	{
		ID:          "mana_surge",
		Name:        "Mana Surge",
		MPCost:      0,
		Description: "Sacrifices 8 HP to immediately surge with +20 MP",
	},
}

// GetUnlearnedSkills returns all skills from the registry not yet known by the player
func GetUnlearnedSkills(learned []Skill) []Skill {
	learnedMap := make(map[string]bool)
	for _, s := range learned {
		learnedMap[s.ID] = true
	}

	unlearned := make([]Skill, 0)
	for _, s := range SkillRegistry {
		if !learnedMap[s.ID] {
			unlearned = append(unlearned, s)
		}
	}
	return unlearned
}

// PickRandomSkillChoices returns up to 3 random unlearned skills for level-up choice
func PickRandomSkillChoices(learned []Skill, rng *rand.Rand) []Skill {
	unlearned := GetUnlearnedSkills(learned)
	if len(unlearned) <= 3 {
		return unlearned
	}

	rng.Shuffle(len(unlearned), func(i, j int) {
		unlearned[i], unlearned[j] = unlearned[j], unlearned[i]
	})

	return unlearned[:3]
}
