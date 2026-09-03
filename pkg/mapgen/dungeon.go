package mapgen

import (
	"math"
	"math/rand"
	"sort"
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

func (r Rect) Contains(x, y int) bool {
	return x >= r.X1 && x <= r.X2 && y >= r.Y1 && y <= r.Y2
}

func (r Rect) InnerContains(x, y int) bool {
	return x > r.X1 && x < r.X2 && y > r.Y1 && y < r.Y2
}

type Point struct {
	X int
	Y int
}

type DungeonMap struct {
	Width       int
	Height      int
	Tiles       [][]Tile
	Rooms       []Rect
	RoomDegree  map[int]int
	Doorways    map[int][]Point
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
		Width:      width,
		Height:     height,
		Tiles:      tiles,
		Rooms:      make([]Rect, 0),
		RoomDegree: make(map[int]int),
		Doorways:   make(map[int][]Point),
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

type Edge struct {
	U int
	V int
	D float64
}

// Generate creates a rich, branching procedural dungeon with diverse rooms
func Generate(width, height, maxRooms, minRoomSize, maxRoomSize int, rng *rand.Rand) *DungeonMap {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	dungeon := NewDungeonMap(width, height)

	// 1. Generate diverse non-overlapping rooms
	for i := 0; i < maxRooms*2 && len(dungeon.Rooms) < maxRooms; i++ {
		w := rng.Intn(maxRoomSize-minRoomSize+1) + minRoomSize
		h := rng.Intn(maxRoomSize-minRoomSize+1) + minRoomSize

		// Introduce aspect ratio variety (wide rooms, tall rooms, square rooms)
		if rng.Intn(3) == 0 {
			w += rng.Intn(3)
		} else if rng.Intn(3) == 1 {
			h += rng.Intn(3)
		}

		x := rng.Intn(width-w-4) + 2
		y := rng.Intn(height-h-4) + 2

		newRoom := NewRect(x, y, w, h)

		overlap := false
		for _, otherRoom := range dungeon.Rooms {
			// Ensure 1-tile buffer
			buffered := NewRect(otherRoom.X1-1, otherRoom.Y1-1, otherRoom.X2-otherRoom.X1+2, otherRoom.Y2-otherRoom.Y1+2)
			if newRoom.Intersects(buffered) {
				overlap = true
				break
			}
		}

		if !overlap {
			dungeon.createRoom(newRoom)
			dungeon.Rooms = append(dungeon.Rooms, newRoom)
		}
	}

	// Fallback if needed
	if len(dungeon.Rooms) < 4 {
		r1 := NewRect(3, 3, 9, 7)
		r2 := NewRect(width/2-4, 3, 9, 7)
		r3 := NewRect(3, height-10, 9, 7)
		r4 := NewRect(width-14, height-10, 10, 8)
		dungeon.createRoom(r1)
		dungeon.createRoom(r2)
		dungeon.createRoom(r3)
		dungeon.createRoom(r4)
		dungeon.Rooms = []Rect{r1, r2, r3, r4}
	}

	numRooms := len(dungeon.Rooms)

	// 2. Build Spanning Network of Corridors with Branching and Loops
	// Calculate all pairwise distances
	edges := make([]Edge, 0)
	for i := 0; i < numRooms; i++ {
		cx1, cy1 := dungeon.Rooms[i].Center()
		for j := i + 1; j < numRooms; j++ {
			cx2, cy2 := dungeon.Rooms[j].Center()
			dist := math.Hypot(float64(cx1-cx2), float64(cy1-cy2))
			edges = append(edges, Edge{U: i, V: j, D: dist})
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].D < edges[j].D
	})

	// Disjoint set union (Kruskal's algorithm for minimum spanning tree)
	parent := make([]int, numRooms)
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(i int) int {
		if parent[i] == i {
			return i
		}
		parent[i] = find(parent[i])
		return parent[i]
	}
	union := func(i, j int) bool {
		rootI := find(i)
		rootJ := find(j)
		if rootI != rootJ {
			parent[rootI] = rootJ
			return true
		}
		return false
	}

	connectedEdges := make([]Edge, 0)

	// Create Minimum Spanning Tree
	for _, edge := range edges {
		if union(edge.U, edge.V) {
			connectedEdges = append(connectedEdges, edge)
		}
	}

	// Add 2-3 extra random edges for circular loops & alternate pathways
	extraAdded := 0
	for _, edge := range edges {
		if extraAdded >= 2 {
			break
		}
		// Pick short non-tree edges
		isTreeEdge := false
		for _, te := range connectedEdges {
			if (te.U == edge.U && te.V == edge.V) || (te.U == edge.V && te.V == edge.U) {
				isTreeEdge = true
				break
			}
		}
		if !isTreeEdge && rng.Intn(100) < 40 {
			connectedEdges = append(connectedEdges, edge)
			extraAdded++
		}
	}

	// Carve all connected corridors
	for _, edge := range connectedEdges {
		dungeon.RoomDegree[edge.U]++
		dungeon.RoomDegree[edge.V]++

		cx1, cy1 := dungeon.Rooms[edge.U].Center()
		cx2, cy2 := dungeon.Rooms[edge.V].Center()

		if rng.Intn(2) == 0 {
			dungeon.createHTunnel(cx1, cx2, cy1)
			dungeon.createVTunnel(cy1, cy2, cx2)
		} else {
			dungeon.createVTunnel(cy1, cy2, cx1)
			dungeon.createHTunnel(cx1, cx2, cy2)
		}
	}

	// 3. Detect True Doorways (Corridor entry tiles into room perimeters)
	for roomIdx, room := range dungeon.Rooms {
		doorways := make([]Point, 0)

		// Check Top & Bottom walls
		for x := room.X1 + 1; x < room.X2; x++ {
			// Top entry
			if dungeon.InBounds(x, room.Y1) && dungeon.Tiles[x][room.Y1].Type == TileFloor {
				doorways = append(doorways, Point{X: x, Y: room.Y1})
			}
			// Bottom entry
			if dungeon.InBounds(x, room.Y2) && dungeon.Tiles[x][room.Y2].Type == TileFloor {
				doorways = append(doorways, Point{X: x, Y: room.Y2})
			}
		}

		// Check Left & Right walls
		for y := room.Y1 + 1; y < room.Y2; y++ {
			// Left entry
			if dungeon.InBounds(room.X1, y) && dungeon.Tiles[room.X1][y].Type == TileFloor {
				doorways = append(doorways, Point{X: room.X1, Y: y})
			}
			// Right entry
			if dungeon.InBounds(room.X2, y) && dungeon.Tiles[room.X2][y].Type == TileFloor {
				doorways = append(doorways, Point{X: room.X2, Y: y})
			}
		}

		dungeon.Doorways[roomIdx] = doorways
	}

	// Start Room is Room 0
	dungeon.StartX, dungeon.StartY = dungeon.Rooms[0].Center()

	// Exit stairs in the farthest room from Start
	bestDist := -1.0
	exitRoomIdx := numRooms - 1
	for i := 1; i < numRooms; i++ {
		cx, cy := dungeon.Rooms[i].Center()
		d := math.Hypot(float64(cx-dungeon.StartX), float64(cy-dungeon.StartY))
		if d > bestDist {
			bestDist = d
			exitRoomIdx = i
		}
	}

	sx, sy := dungeon.Rooms[exitRoomIdx].Center()
	dungeon.StairsDownX = sx
	dungeon.StairsDownY = sy
	dungeon.Tiles[sx][sy] = NewStairsDown()

	return dungeon
}
