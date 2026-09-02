package entities

import (
	"math/rand"
	"testing"
)

func TestMimicAndChampions(t *testing.T) {
	mimic := NewMimic(5, 5)
	if mimic.Name != "Hungry Mimic" || mimic.Rune != 'M' {
		t.Errorf("Unexpected mimic properties: %s %c", mimic.Name, mimic.Rune)
	}

	rng := rand.New(rand.NewSource(1234))
	foundChamp := false
	for i := 0; i < 50; i++ {
		m := GenerateRandomMonster(0, 0, 3, rng)
		if m.IsChampion {
			foundChamp = true
			if m.Affix == "" {
				t.Errorf("Expected champion to have an affix, got empty string")
			}
			break
		}
	}

	if !foundChamp {
		t.Errorf("Expected to generate at least 1 champion out of 50 rolls")
	}
}
