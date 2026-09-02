package mapgen

// GenerateTutorial creates a handcrafted 4-room linear tutorial dungeon adapted to width/height
func GenerateTutorial(mapW, mapH int) *DungeonMap {
	if mapW < 50 {
		mapW = 50
	}
	if mapH < 18 {
		mapH = 18
	}

	dungeon := NewDungeonMap(mapW, mapH)

	// Calculate room dimensions to fill the width evenly
	numRooms := 4
	roomW := (mapW - 12) / numRooms
	if roomW > 16 {
		roomW = 16
	}
	if roomW < 8 {
		roomW = 8
	}

	roomH := mapH - 8
	if roomH > 10 {
		roomH = 10
	}
	if roomH < 6 {
		roomH = 6
	}

	centerY := mapH / 2
	roomY := centerY - (roomH / 2)

	rooms := make([]Rect, numRooms)
	spacing := (mapW - 4 - (roomW * numRooms)) / (numRooms - 1)
	if spacing < 2 {
		spacing = 2
	}

	for i := 0; i < numRooms; i++ {
		roomX := 2 + i*(roomW+spacing)
		if roomX+roomW >= mapW-1 {
			roomX = mapW - roomW - 2
		}
		r := NewRect(roomX, roomY, roomW, roomH)
		dungeon.createRoom(r)
		rooms[i] = r

		if i > 0 {
			prevCenter := rooms[i-1].X2
			currCenter := r.X1 + 1
			dungeon.createHTunnel(prevCenter-1, currCenter, centerY)
			// Place a door at the entrance of room 2 and room 3
			doorX := (prevCenter + r.X1) / 2
			if dungeon.InBounds(doorX, centerY) {
				dungeon.Tiles[doorX][centerY] = NewDoorClosed()
			}
		}
	}

	dungeon.Rooms = rooms

	// Room 1 Center = Player Start
	dungeon.StartX, dungeon.StartY = rooms[0].Center()

	// Room 4 Center = Stairs Down
	dungeon.StairsDownX, dungeon.StairsDownY = rooms[3].Center()
	dungeon.Tiles[dungeon.StairsDownX][dungeon.StairsDownY] = NewStairsDown()

	return dungeon
}
