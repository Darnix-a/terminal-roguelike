package mapgen

import (
	"math/rand"
	"testing"
)

func TestDungeonGeneration(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	dungeon := Generate(60, 25, 10, 5, 10, rng)

	if len(dungeon.Rooms) == 0 {
		t.Fatalf("Expected rooms to be generated, got 0")
	}

	// Verify Start Position is on a floor
	if dungeon.Tiles[dungeon.StartX][dungeon.StartY].Type != TileFloor {
		t.Errorf("Start position (%d, %d) is not a floor tile", dungeon.StartX, dungeon.StartY)
	}

	// Verify Stairs Down is placed
	if dungeon.Tiles[dungeon.StairsDownX][dungeon.StairsDownY].Type != TileStairsDown {
		t.Errorf("Stairs down (%d, %d) is not TileStairsDown", dungeon.StairsDownX, dungeon.StairsDownY)
	}
}

func TestFOVComputation(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	dungeon := Generate(60, 25, 10, 5, 10, rng)

	dungeon.ComputeFOV(dungeon.StartX, dungeon.StartY, 8)

	// Center tile must be visible and explored
	if !dungeon.Tiles[dungeon.StartX][dungeon.StartY].Visible {
		t.Errorf("Player start tile should be visible")
	}
	if !dungeon.Tiles[dungeon.StartX][dungeon.StartY].Explored {
		t.Errorf("Player start tile should be explored")
	}
}
