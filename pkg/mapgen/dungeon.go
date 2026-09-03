package mapgen

import (
	"math/rand"
	"time"
)

type Rect struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

func NewRect(x, y, w, h int) Rect {
	return Rect{
		X1: x,
		Y1: y,
		X2: x + w,
		Y2: y + h,
	}
}

func (r Rect) Center() (int, int) {
	return (r.X1 + r.X2) / 2, (r.Y1 + r.Y2) / 2
}

func (r Rect) Intersects(other Rect) bool {
	return r.X1 <= other.X2 && r.X2 >= other.X1 &&
		r.Y1 <= other.Y2 && r.Y2 >= other.Y1
}

type DungeonMap struct {
	Width       int
	Height      int
	Tiles       [][]Tile
	Rooms       []Rect
	StairsDownX int
	StairsDownY int
	StartX      int
	StartY      int
}

func NewDungeonMap(width, height int) *DungeonMap {
	tiles := make([][]Tile, width)
	for x := 0; x < width; x++ {
		tiles[x] = make([]Tile, height)
		for y := 0; y < height; y++ {
			tiles[x][y] = NewWall()
		}
	}

	return &DungeonMap{
		Width:  width,
		Height: height,
		Tiles:  tiles,
		Rooms:  make([]Rect, 0),
	}
}

func (m *DungeonMap) InBounds(x, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height
}

func (m *DungeonMap) IsBlocked(x, y int) bool {
	if !m.InBounds(x, y) {
		return true
	}
	return m.Tiles[x][y].Blocked
}

func (m *DungeonMap) createRoom(room Rect) {
	for x := room.X1 + 1; x < room.X2; x++ {
		for y := room.Y1 + 1; y < room.Y2; y++ {
			if m.InBounds(x, y) {
				m.Tiles[x][y] = NewFloor()
			}
		}
	}
}

func (m *DungeonMap) createHTunnel(x1, x2, y int) {
	minX, maxX := x1, x2
	if x1 > x2 {
		minX, maxX = x2, x1
	}
	for x := minX; x <= maxX; x++ {
		if m.InBounds(x, y) {
			m.Tiles[x][y] = NewFloor()
		}
	}
}

func (m *DungeonMap) createVTunnel(y1, y2, x int) {
	minY, maxY := y1, y2
	if y1 > y2 {
		minY, maxY = y2, y1
	}
	for y := minY; y <= maxY; y++ {
		if m.InBounds(x, y) {
			m.Tiles[x][y] = NewFloor()
		}
	}
}

// Generate creates a procedural dungeon with guaranteed room connectivity
func Generate(width, height, maxRooms, minRoomSize, maxRoomSize int, rng *rand.Rand) *DungeonMap {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	dungeon := NewDungeonMap(width, height)

	for i := 0; i < maxRooms; i++ {
		w := rng.Intn(maxRoomSize-minRoomSize+1) + minRoomSize
		h := rng.Intn(maxRoomSize-minRoomSize+1) + minRoomSize
		x := rng.Intn(width-w-3) + 2
		y := rng.Intn(height-h-3) + 2

		newRoom := NewRect(x, y, w, h)

		overlap := false
		for _, otherRoom := range dungeon.Rooms {
			// Ensure 1-tile buffer between rooms
			buffered := NewRect(otherRoom.X1-1, otherRoom.Y1-1, otherRoom.X2-otherRoom.X1+2, otherRoom.Y2-otherRoom.Y1+2)
			if newRoom.Intersects(buffered) {
				overlap = true
				break
			}
		}

		if !overlap {
			dungeon.createRoom(newRoom)

			newCenterX, newCenterY := newRoom.Center()

			if len(dungeon.Rooms) == 0 {
				dungeon.StartX = newCenterX
				dungeon.StartY = newCenterY
			} else {
				// Connect directly to the previous room center
				prevCenterX, prevCenterY := dungeon.Rooms[len(dungeon.Rooms)-1].Center()

				if rng.Intn(2) == 1 {
					dungeon.createHTunnel(prevCenterX, newCenterX, prevCenterY)
					dungeon.createVTunnel(prevCenterY, newCenterY, newCenterX)
				} else {
					dungeon.createVTunnel(prevCenterY, newCenterY, prevCenterX)
					dungeon.createHTunnel(prevCenterX, newCenterX, newCenterY)
				}
			}

			dungeon.Rooms = append(dungeon.Rooms, newRoom)
		}
	}

	// Fallback if too few rooms placed
	if len(dungeon.Rooms) < 3 {
		r1 := NewRect(4, 4, 10, 8)
		r2 := NewRect(width-16, height-12, 10, 8)
		dungeon.createRoom(r1)
		dungeon.createRoom(r2)
		dungeon.createHTunnel(9, width-11, 8)
		dungeon.createVTunnel(8, height-8, width-11)
		dungeon.StartX = 9
		dungeon.StartY = 8
		dungeon.Rooms = []Rect{r1, r2}
	}

	// Place stairs down in the last room
	lastRoom := dungeon.Rooms[len(dungeon.Rooms)-1]
	sx, sy := lastRoom.Center()
	dungeon.StairsDownX = sx
	dungeon.StairsDownY = sy
	dungeon.Tiles[sx][sy] = NewStairsDown()

	return dungeon
}
