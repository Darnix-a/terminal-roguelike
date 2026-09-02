package mapgen

type TileType int

const (
	TileWall TileType = iota
	TileFloor
	TileDoorClosed
	TileDoorOpen
	TileDoorLocked
	TileStairsDown
	TileStairsUp
	TileChest
	TileChestOpened
	TileFountain
	TileFountainUsed
	TileShrine
	TileShrineUsed
)

type Tile struct {
	Type       TileType
	Blocked    bool
	BlockSight bool
	Explored   bool
	Visible    bool
}

func NewWall() Tile {
	return Tile{Type: TileWall, Blocked: true, BlockSight: true}
}

func NewFloor() Tile {
	return Tile{Type: TileFloor, Blocked: false, BlockSight: false}
}

func NewDoorClosed() Tile {
	return Tile{Type: TileDoorClosed, Blocked: true, BlockSight: true}
}

func NewDoorOpen() Tile {
	return Tile{Type: TileDoorOpen, Blocked: false, BlockSight: false}
}

func NewDoorLocked() Tile {
	return Tile{Type: TileDoorLocked, Blocked: true, BlockSight: true}
}

func NewStairsDown() Tile {
	return Tile{Type: TileStairsDown, Blocked: false, BlockSight: false}
}

func NewStairsUp() Tile {
	return Tile{Type: TileStairsUp, Blocked: false, BlockSight: false}
}

func NewChest() Tile {
	return Tile{Type: TileChest, Blocked: true, BlockSight: false}
}

func NewChestOpened() Tile {
	return Tile{Type: TileChestOpened, Blocked: false, BlockSight: false}
}

func NewFountain() Tile {
	return Tile{Type: TileFountain, Blocked: true, BlockSight: false}
}

func NewFountainUsed() Tile {
	return Tile{Type: TileFountainUsed, Blocked: false, BlockSight: false}
}

func NewShrine() Tile {
	return Tile{Type: TileShrine, Blocked: true, BlockSight: false}
}

func NewShrineUsed() Tile {
	return Tile{Type: TileShrineUsed, Blocked: false, BlockSight: false}
}
