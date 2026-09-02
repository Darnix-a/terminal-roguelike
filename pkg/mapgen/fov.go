package mapgen

import (
	"math"
)

// ComputeFOV calculates the field of view from (px, py) up to radius
func (m *DungeonMap) ComputeFOV(px, py, radius int) {
	// 1. Reset all visibility flags (explored remains true)
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.Tiles[x][y].Visible = false
		}
	}

	// 2. Center tile is always visible and explored
	if m.InBounds(px, py) {
		m.Tiles[px][py].Visible = true
		m.Tiles[px][py].Explored = true
	}

	// 3. 360-degree raycasting
	numRays := 360
	for i := 0; i < numRays; i++ {
		rad := float64(i) * (math.Pi / 180.0)
		dx := math.Cos(rad)
		dy := math.Sin(rad)

		curX := float64(px) + 0.5
		curY := float64(py) + 0.5

		for step := 0; step < radius; step++ {
			curX += dx
			curY += dy

			tileX := int(math.Floor(curX))
			tileY := int(math.Floor(curY))

			if !m.InBounds(tileX, tileY) {
				break
			}

			// Mark visible & explored
			m.Tiles[tileX][tileY].Visible = true
			m.Tiles[tileX][tileY].Explored = true

			// Stop ray if it hits a wall/obstacle that blocks sight
			if m.Tiles[tileX][tileY].BlockSight {
				break
			}
		}
	}
}
