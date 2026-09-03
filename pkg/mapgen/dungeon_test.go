package mapgen

import (
	"math/rand"
	"testing"
)

func TestDungeonGenerationRobustness(t *testing.T) {
	rng := rand.New(rand.NewSource(999))

	for i := 0; i < 50; i++ {
		dungeon := Generate(70, 24, 14, 5, 10, rng)

		if len(dungeon.Rooms) < 2 {
			t.Errorf("Iteration %d: expected >= 2 rooms, got %d", i, len(dungeon.Rooms))
		}

		if !dungeon.InBounds(dungeon.StartX, dungeon.StartY) {
			t.Errorf("Iteration %d: StartX/StartY out of bounds", i)
		}

		if !dungeon.InBounds(dungeon.StairsDownX, dungeon.StairsDownY) {
			t.Errorf("Iteration %d: StairsDown out of bounds", i)
		}

		if dungeon.Tiles[dungeon.StairsDownX][dungeon.StairsDownY].Type != TileStairsDown {
			t.Errorf("Iteration %d: expected StairsDown tile at exit", i)
		}
	}
}
