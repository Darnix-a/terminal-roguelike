package engine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"

	"terminal-roguelike/pkg/combat"
	"terminal-roguelike/pkg/entities"
	"terminal-roguelike/pkg/items"
	"terminal-roguelike/pkg/mapgen"
)

type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StateInventory
	StateShop
	StateCombat
	StateHighScores
	StateHelp
	StateGameOver
	StateVictory
)

type Game struct {
	State            GameState
	Map              *mapgen.DungeonMap
	Player           *entities.Player
	Monsters         []*entities.Monster
	Items            []*items.Item
	Floor            int
	MaxFloors        int
	TurnCount        int
	Log              *MessageLog
	RNG              *rand.Rand
	FOVRadius        int
	InventoryIdx     int
	MapW             int
	MapH             int
	ActiveBattle     *combat.Battle
	ActiveMerchant   *entities.Monster
	LastPlayerX      int
	LastPlayerY      int
	TutorialTriggers map[string]bool
}

func NewGame(mapW, mapH int) *Game {
	startFloor := 1
	if !HasCompletedTutorial() {
		startFloor = 0
	}
	return NewGameCustom(startFloor, mapW, mapH)
}

func NewGameCustom(floor int, mapW, mapH int) *Game {
	if mapW < 50 {
		mapW = 70
	}
	if mapH < 18 {
		mapH = 22
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := &Game{
		State:            StatePlaying,
		Floor:            floor,
		MaxFloors:        5,
		TurnCount:        0,
		Log:              NewMessageLog(100),
		RNG:              rng,
		FOVRadius:        8,
		InventoryIdx:     0,
		MapW:             mapW,
		MapH:             mapH,
		TutorialTriggers: make(map[string]bool),
	}

	g.generateFloor(floor)
	return g
}

func LoadGameFromSave(mapW, mapH int) *Game {
	app := loadAppData()
	if app.ActiveSave == nil {
		return NewGame(mapW, mapH)
	}

	save := app.ActiveSave
	g := NewGameCustom(save.Floor, mapW, mapH)
	g.TurnCount = save.TurnCount

	// Restore Player Stats
	g.Player.HP = save.Player.HP
	g.Player.MaxHP = save.Player.MaxHP
	g.Player.MP = save.Player.MP
	g.Player.MaxMP = save.Player.MaxMP
	g.Player.BaseATK = save.Player.BaseATK
	g.Player.BaseDEF = save.Player.BaseDEF
	g.Player.Level = save.Player.Level
	g.Player.EXP = save.Player.EXP
	g.Player.MaxEXP = save.Player.MaxEXP
	g.Player.Gold = save.Player.Gold
	g.Player.Keys = save.Player.Keys
	g.Player.Kills = save.Player.Kills

	// Restore Inventory
	g.Player.Inventory.Items = make([]*items.Item, 0)
	for _, itm := range save.Player.Inventory {
		restoredItem := &items.Item{
			ID:          itm.ID,
			Name:        itm.Name,
			Type:        itm.Type,
			Rune:        itm.Rune,
			HealAmount:  itm.HealAmount,
			BonusATK:    itm.BonusATK,
			BonusDEF:    itm.BonusDEF,
			Value:       itm.Value,
			Description: itm.Description,
			Equipped:    itm.Equipped,
		}
		_ = g.Player.Inventory.Add(restoredItem)
	}

	g.Log.Add("=== CONTINUED SAVED GAME ===", tcell.ColorGold)
	g.Log.Add(fmt.Sprintf("Welcome back, Lv.%d Hero! Resumed at Dungeon Floor %d.", g.Player.Level, g.Floor), tcell.ColorYellow)
	return g
}

func (g *Game) generateFloor(floor int) {
	g.Floor = floor
	g.Monsters = make([]*entities.Monster, 0)
	g.Items = make([]*items.Item, 0)
	g.ActiveMerchant = nil

	if floor == 0 {
		// Handcrafted Tutorial Floor
		g.Map = mapgen.GenerateTutorial(g.MapW, g.MapH)

		if g.Player == nil {
			g.Player = entities.NewPlayer(g.Map.StartX, g.Map.StartY)
		} else {
			g.Player.MoveTo(g.Map.StartX, g.Map.StartY)
		}

		if len(g.Map.Rooms) >= 2 {
			r2 := g.Map.Rooms[1]
			cx, cy := r2.Center()
			g.Items = append(g.Items, items.NewDagger(cx-1, cy-1))
			g.Items = append(g.Items, items.NewLeatherArmor(cx+1, cy-1))
			g.Items = append(g.Items, items.NewHealthPotion(cx-1, cy+1))
			g.Items = append(g.Items, items.NewGold(cx+1, cy+1, 50))
		}

		if len(g.Map.Rooms) >= 3 {
			r3 := g.Map.Rooms[2]
			cx, cy := r3.Center()
			trainingDummy := entities.NewGoblin(cx, cy)
			trainingDummy.Name = "Training Goblin"
			trainingDummy.HP = 10
			trainingDummy.MaxHP = 10
			trainingDummy.ATK = 3
			g.Monsters = append(g.Monsters, trainingDummy)
		}

		g.Log.Add("=== TUTORIAL: TRAINING GROUNDS ===", tcell.ColorGold)
		g.Log.Add("You are '@'. Move with WASD or Arrow Keys. Darkness reveals as you walk.", tcell.ColorYellow)
		g.Log.Add("Explore the chambers to your right to learn the basics!", tcell.ColorLightCyan)

	} else {
		// Procedural Dungeon Floors (1 to 5)
		g.Map = mapgen.Generate(g.MapW, g.MapH, 14, 5, 10, g.RNG)

		if g.Player == nil {
			g.Player = entities.NewPlayer(g.Map.StartX, g.Map.StartY)
		} else {
			g.Player.MoveTo(g.Map.StartX, g.Map.StartY)
		}

		g.Log.Add(fmt.Sprintf("=== ENTERED DUNGEON FLOOR %d ===", floor), tcell.ColorGold)
		if floor == g.MaxFloors {
			g.Log.Add("WARNING: THE ANCIENT RED DRAGON LAIRS ON THIS FLOOR!", tcell.ColorRed)
		}

		// Identify candidate rooms for Vault and Shop
		numRooms := len(g.Map.Rooms)
		exitRoomIdx := numRooms - 1
		for rIdx, r := range g.Map.Rooms {
			if r.Contains(g.Map.StairsDownX, g.Map.StairsDownY) {
				exitRoomIdx = rIdx
				break
			}
		}

		// Find candidate side / dead-end rooms
		candidateRooms := make([]int, 0)
		for i := 1; i < numRooms; i++ {
			if i != exitRoomIdx {
				candidateRooms = append(candidateRooms, i)
			}
		}
		g.RNG.Shuffle(len(candidateRooms), func(i, j int) {
			candidateRooms[i], candidateRooms[j] = candidateRooms[j], candidateRooms[i]
		})

		vaultRoomIdx := -1
		shopRoomIdx := -1

		// Pick 1 random side room for Locked Vault
		if len(candidateRooms) > 0 {
			vaultRoomIdx = candidateRooms[0]
			candidateRooms = candidateRooms[1:]

			// Seal exactly 1 primary doorway with Locked Door '%'
			if len(g.Map.Doorways[vaultRoomIdx]) > 0 {
				pt := g.Map.Doorways[vaultRoomIdx][0]
				g.Map.Tiles[pt.X][pt.Y] = mapgen.NewDoorLocked()
			}

			vRoom := g.Map.Rooms[vaultRoomIdx]
			// Place chest at random interior tile
			chestX := g.RNG.Intn(vRoom.X2-vRoom.X1-1) + vRoom.X1 + 1
			chestY := g.RNG.Intn(vRoom.Y2-vRoom.Y1-1) + vRoom.Y1 + 1
			g.Map.Tiles[chestX][chestY] = mapgen.NewChest()

			// Place fountain or shrine at another random tile
			for attempt := 0; attempt < 10; attempt++ {
				fx := g.RNG.Intn(vRoom.X2-vRoom.X1-1) + vRoom.X1 + 1
				fy := g.RNG.Intn(vRoom.Y2-vRoom.Y1-1) + vRoom.Y1 + 1
				if fx != chestX || fy != chestY {
					if g.RNG.Intn(2) == 0 {
						g.Map.Tiles[fx][fy] = mapgen.NewFountain()
					} else {
						g.Map.Tiles[fx][fy] = mapgen.NewShrine()
					}
					break
				}
			}
		}

		// Pick 1 random side room for Dungeon Shop on Floors 2 and 4
		if (floor == 2 || floor == 4) && len(candidateRooms) > 0 {
			shopRoomIdx = candidateRooms[0]
			candidateRooms = candidateRooms[1:]

			sRoom := g.Map.Rooms[shopRoomIdx]
			sx, sy := sRoom.Center()

			merchant := entities.NewEntity("merchant", "Wandering Merchant", 'S', tcell.ColorGold, sx, sy, true)
			merchantMonster := &entities.Monster{
				Entity:     merchant,
				HP:         70,
				MaxHP:      70,
				ATK:        16,
				DEF:        5,
				EXP:        200,
				IsMerchant: true,
				Alerted:    false,
				Sprite: []string{
					`  [$$__$$]  `,
					`  /|====|\  `,
					`   | || |   `,
				},
				Actions: []entities.MonsterAction{
					{Name: "Ledger Strike", DamageMult: 1.2, Description: "slams you with a heavy gilded accounting tome"},
					{Name: "Coin Gatling", DamageMult: 2.0, IsTelegraphed: true, TelegraphWarning: "loads a sack of pure gold coins for a COIN GATLING assault!", Description: "unleashes a flurry of hypersonic golden coins"},
				},
			}
			g.Monsters = append(g.Monsters, merchantMonster)
			g.ActiveMerchant = merchantMonster
		}

		// Place at most 1-2 wooden doors across the entire dungeon (only on strict choke points)
		doorsPlaced := 0
		for rIdx := 0; rIdx < numRooms && doorsPlaced < 2; rIdx++ {
			if rIdx == vaultRoomIdx || rIdx == shopRoomIdx {
				continue
			}
			for _, pt := range g.Map.Doorways[rIdx] {
				if doorsPlaced < 2 && g.RNG.Intn(100) < 25 && g.Map.IsValidChokePoint(pt.X, pt.Y) {
					g.Map.Tiles[pt.X][pt.Y] = mapgen.NewDoorClosed()
					doorsPlaced++
					break
				}
			}
		}

		// Place 1 Iron Key 'k' in a random accessible room
		keyPlaced := false

		for i, room := range g.Map.Rooms {
			if i == 0 {
				continue // Skip spawn room
			}

			// Place Key
			if !keyPlaced && i != vaultRoomIdx && i != shopRoomIdx {
				kx := g.RNG.Intn(room.X2-room.X1-1) + room.X1 + 1
				ky := g.RNG.Intn(room.Y2-room.Y1-1) + room.Y1 + 1
				if !g.Map.IsBlocked(kx, ky) {
					g.Items = append(g.Items, items.NewDungeonKey(kx, ky))
					keyPlaced = true
				}
			}

			// Floor 5 Dragon Boss
			if floor == g.MaxFloors && i == exitRoomIdx {
				cx, cy := room.Center()
				boss := entities.NewDragonBoss(cx, cy)
				g.Monsters = append(g.Monsters, boss)
				continue
			}

			// Don't spawn monsters inside vault or shop
			if i == vaultRoomIdx || i == shopRoomIdx {
				continue
			}

			// 20% Chance to place a Chest outside vault
			if g.RNG.Intn(100) < 20 {
				cx := g.RNG.Intn(room.X2-room.X1-1) + room.X1 + 1
				cy := g.RNG.Intn(room.Y2-room.Y1-1) + room.Y1 + 1
				if !g.Map.IsBlocked(cx, cy) {
					g.Map.Tiles[cx][cy] = mapgen.NewChest()
				}
			}

			// 15% Chance to place a Fountain or Shrine
			if g.RNG.Intn(100) < 15 {
				fx := g.RNG.Intn(room.X2-room.X1-1) + room.X1 + 1
				fy := g.RNG.Intn(room.Y2-room.Y1-1) + room.Y1 + 1
				if !g.Map.IsBlocked(fx, fy) {
					if g.RNG.Intn(2) == 0 {
						g.Map.Tiles[fx][fy] = mapgen.NewFountain()
					} else {
						g.Map.Tiles[fx][fy] = mapgen.NewShrine()
					}
				}
			}

			// Spawn Monsters
			numMonsters := g.RNG.Intn(2) + 1
			if floor >= 3 {
				numMonsters = g.RNG.Intn(2) + 1
			}
			for m := 0; m < numMonsters; m++ {
				mx := g.RNG.Intn(room.X2-room.X1-1) + room.X1 + 1
				my := g.RNG.Intn(room.Y2-room.Y1-1) + room.Y1 + 1
				if !g.Map.IsBlocked(mx, my) && !g.isOccupied(mx, my) {
					monster := entities.GenerateRandomMonster(mx, my, floor, g.RNG)
					g.Monsters = append(g.Monsters, monster)
				}
			}

			// Spawn 1 Item (60% chance)
			if g.RNG.Intn(100) < 60 {
				ix := g.RNG.Intn(room.X2-room.X1-1) + room.X1 + 1
				iy := g.RNG.Intn(room.Y2-room.Y1-1) + room.Y1 + 1
				if !g.Map.IsBlocked(ix, iy) {
					item := items.GenerateRandomItem(ix, iy, floor, g.RNG)
					g.Items = append(g.Items, item)
				}
			}
		}

		// Auto-save on floor generation
		SaveGameProgress(g)
	}

	g.Map.ComputeFOV(g.Player.X, g.Player.Y, g.FOVRadius)
}

func (g *Game) isOccupied(x, y int) bool {
	if g.Player != nil && g.Player.X == x && g.Player.Y == y {
		return true
	}
	for _, m := range g.Monsters {
		if m.IsAlive && m.X == x && m.Y == y {
			return true
		}
	}
	return false
}

func (g *Game) getMonsterAt(x, y int) *entities.Monster {
	for _, m := range g.Monsters {
		if m.IsAlive && m.X == x && m.Y == y {
			return m
		}
	}
	return nil
}

func (g *Game) getItemAt(x, y int) *items.Item {
	for i, item := range g.Items {
		if item.X == x && item.Y == y {
			g.Items = append(g.Items[:i], g.Items[i+1:]...)
			return item
		}
	}
	return nil
}

// StartBattle initiates dedicated turn-based combat screen
func (g *Game) StartBattle(monster *entities.Monster) {
	g.LastPlayerX = g.Player.X
	g.LastPlayerY = g.Player.Y
	g.ActiveBattle = combat.NewBattle(g.Player, monster, g.RNG)
	g.State = StateCombat

	// Alert nearby monsters (Swarm Mechanic)
	g.alertNearbyMonsters(g.Player.X, g.Player.Y, 4.0)
}

func (g *Game) alertNearbyMonsters(x, y int, radius float64) {
	for _, m := range g.Monsters {
		if m.IsAlive && !m.Alerted && !m.IsMerchant {
			if entities.Distance(x, y, m.X, m.Y) <= radius {
				m.Alerted = true
			}
		}
	}
}

// ConcludeBattle processes battle resolution
func (g *Game) ConcludeBattle() {
	if g.ActiveBattle == nil {
		g.State = StatePlaying
		return
	}

	switch g.ActiveBattle.Result {
	case combat.BattleVictory:
		monster := g.ActiveBattle.Monster
		g.Log.Add(fmt.Sprintf("VICTORY! Slew %s! (+%d EXP)", monster.Name, monster.EXP), tcell.ColorGold)
		for _, msg := range g.ActiveBattle.LevelUpMsgs {
			g.Log.Add(msg, tcell.ColorAqua)
		}

		// Drop Gold
		goldReward := g.RNG.Intn(15*g.Floor+1) + 10
		if monster.IsChampion {
			goldReward *= 2
		} else if monster.IsMerchant {
			goldReward = 150
			g.Log.Add("Looted the Merchant's Heavy Vault Sack (+150 Gold)!", tcell.ColorGold)
		}
		g.Player.Gold += goldReward
		g.Log.Add(fmt.Sprintf("Looted %d Gold Coins from %s.", goldReward, monster.Name), tcell.ColorYellow)

		if monster.IsBoss {
			g.Log.Add("THE ANCIENT RED DRAGON HAS FALLEN! VICTORY IS YOURS!", tcell.ColorGold)
			g.State = StateVictory
			DeleteActiveSave()

			score := g.Player.Gold + (g.Player.Kills * 50) + 1000 + (g.Player.Level * 100)
			RecordHighScore(HighScoreEntry{
				Score:     score,
				Floor:     g.Floor,
				Level:     g.Player.Level,
				Kills:     g.Player.Kills,
				Gold:      g.Player.Gold,
				Outcome:   "Conquered Dungeon (Victory)",
				Timestamp: time.Now(),
			})
			return
		}

		g.State = StatePlaying
		if g.Floor > 0 {
			SaveGameProgress(g)
		}

	case combat.BattleDefeat:
		g.State = StateGameOver
		DeleteActiveSave()

		score := g.Player.Gold + (g.Player.Kills * 25) + (g.Floor * 100) + (g.Player.Level * 50)
		RecordHighScore(HighScoreEntry{
			Score:     score,
			Floor:     g.Floor,
			Level:     g.Player.Level,
			Kills:     g.Player.Kills,
			Gold:      g.Player.Gold,
			Outcome:   fmt.Sprintf("Slain by %s on Floor %d", g.ActiveBattle.Monster.Name, g.Floor),
			Timestamp: time.Now(),
		})

	case combat.BattleFled:
		g.Log.Add("Escaped from combat back into the dungeon shadows.", tcell.ColorYellow)
		g.Player.MoveTo(g.LastPlayerX, g.LastPlayerY)
		g.State = StatePlaying
	}

	g.ActiveBattle = nil
	g.Map.ComputeFOV(g.Player.X, g.Player.Y, g.FOVRadius)
}

// BuyShopItem handles purchasing wares from Merchant
func (g *Game) BuyShopItem(slot int) {
	switch slot {
	case 1: // Health Potion (20 Gold)
		if g.Player.Gold >= 20 {
			g.Player.Gold -= 20
			_ = g.Player.Inventory.Add(items.NewHealthPotion(0, 0))
			g.Log.Add("Bought Health Potion (+25 HP) for 20 Gold!", tcell.ColorGreen)
		} else {
			g.Log.Add("Not enough gold to buy Health Potion!", tcell.ColorRed)
		}
	case 2: // Greater Health Draught (40 Gold)
		if g.Player.Gold >= 40 {
			g.Player.Gold -= 40
			_ = g.Player.Inventory.Add(items.NewGreaterHealthPotion(0, 0))
			g.Log.Add("Bought Greater Health Draught (+50 HP) for 40 Gold!", tcell.ColorGreen)
		} else {
			g.Log.Add("Not enough gold to buy Greater Health Draught!", tcell.ColorRed)
		}
	case 3: // Scroll of Weapon Enchantment (60 Gold)
		if g.Player.Gold >= 60 {
			g.Player.Gold -= 60
			g.Player.BaseATK += 3
			g.Log.Add("✨ Read Scroll of Weapon Enchantment (-60 Gold)! Permanently gained +3 Base ATK!", tcell.ColorGold)
		} else {
			g.Log.Add("Not enough gold to buy Scroll of Weapon Enchantment!", tcell.ColorRed)
		}
	case 4: // Dragonscale Shield (75 Gold)
		if g.Player.Gold >= 75 {
			g.Player.Gold -= 75
			shield := &items.Item{
				ID: "dragon_shield", Name: "Dragonscale Shield", Type: items.TypeArmor,
				Rune: '[', Color: tcell.ColorGold, BonusDEF: 4, Description: "+4 Defense Shield",
			}
			_ = g.Player.Inventory.Add(shield)
			g.Log.Add("Bought Dragonscale Shield (+4 DEF) for 75 Gold!", tcell.ColorLightCyan)
		} else {
			g.Log.Add("Not enough gold to buy Dragonscale Shield!", tcell.ColorRed)
		}
	case 5: // Scroll of Teleportation (30 Gold)
		if g.Player.Gold >= 30 {
			g.Player.Gold -= 30
			_ = g.Player.Inventory.Add(items.NewScrollTeleport(0, 0))
			g.Log.Add("Bought Scroll of Teleportation for 30 Gold!", tcell.ColorViolet)
		} else {
			g.Log.Add("Not enough gold to buy Scroll of Teleportation!", tcell.ColorRed)
		}
	}
	SaveGameProgress(g)
}

// AttackMerchant provokes the shopkeeper into combat
func (g *Game) AttackMerchant() {
	if g.ActiveMerchant == nil {
		g.State = StatePlaying
		return
	}
	g.Log.Add("⚠️ YOU ATTACKED THE SHOPKEEPER! Prepare for his wrath!", tcell.ColorRed)
	g.StartBattle(g.ActiveMerchant)
}

// HandlePlayerAction handles exploration, interactions, and combat
func (g *Game) HandlePlayerAction(dx, dy int) {
	if g.State != StatePlaying {
		return
	}

	destX := g.Player.X + dx
	destY := g.Player.Y + dy

	if !g.Map.InBounds(destX, destY) {
		return
	}

	// 1. Check for Monster / Merchant
	targetMonster := g.getMonsterAt(destX, destY)
	if targetMonster != nil {
		if targetMonster.IsMerchant {
			g.ActiveMerchant = targetMonster
			g.State = StateShop
			return
		}
		g.StartBattle(targetMonster)
		return
	}

	// 2. Interactive Objects & Obstacles
	tile := g.Map.Tiles[destX][destY]

	switch tile.Type {
	case mapgen.TileDoorClosed:
		g.Map.Tiles[destX][destY] = mapgen.NewDoorOpen()
		g.Log.Add("You open the wooden door ('+').", tcell.ColorLightGray)
		g.EndPlayerTurn()
		return

	case mapgen.TileDoorLocked:
		if g.Player.Keys > 0 {
			g.Player.Keys--
			g.Map.Tiles[destX][destY] = mapgen.NewDoorOpen()
			g.Log.Add("🗝️ Used a Dungeon Key ('k')! Unlocked the Vault Door ('%')!", tcell.ColorGold)
			g.EndPlayerTurn()
		} else {
			g.Log.Add("🔒 The Vault Door ('%') is locked! Find an Iron Dungeon Key ('k') on this floor.", tcell.ColorRed)
		}
		return

	case mapgen.TileChest:
		// 15% Chance for Mimic!
		if g.RNG.Intn(100) < 15 {
			g.Map.Tiles[destX][destY] = mapgen.NewChestOpened()
			g.Log.Add("⚠️ THE CHEST IS ALIVE! A Hungry Mimic attacks!", tcell.ColorRed)
			mimic := entities.NewMimic(destX, destY)
			g.Monsters = append(g.Monsters, mimic)
			g.StartBattle(mimic)
			return
		}

		// Open Chest -> Rare Loot
		g.Map.Tiles[destX][destY] = mapgen.NewChestOpened()
		rewardGold := g.RNG.Intn(30) + 30
		g.Player.Gold += rewardGold
		g.Log.Add(fmt.Sprintf("✨ Opened Treasure Chest ('=')! Found %d Gold Coins and Rare Supplies!", rewardGold), tcell.ColorGold)

		rareItem := items.GenerateRandomItem(g.Player.X, g.Player.Y, g.Floor+1, g.RNG)
		_ = g.Player.Inventory.Add(rareItem)
		g.Log.Add(fmt.Sprintf("Obtained %s ('%c') from the chest!", rareItem.Name, rareItem.Rune), tcell.ColorLightCyan)
		g.EndPlayerTurn()
		return

	case mapgen.TileFountain:
		g.Map.Tiles[destX][destY] = mapgen.NewFountainUsed()
		g.Player.HP = g.Player.MaxHP
		g.Player.MP = g.Player.MaxMP
		g.Log.Add("⛲ Drank from the Mystic Healing Fountain ('0')! Fully restored HP and MP!", tcell.ColorAqua)
		g.EndPlayerTurn()
		return

	case mapgen.TileShrine:
		if g.Player.Gold >= 30 {
			g.Player.Gold -= 30
			g.Map.Tiles[destX][destY] = mapgen.NewShrineUsed()

			blessingType := g.RNG.Intn(3)
			switch blessingType {
			case 0:
				g.Player.BaseATK += 3
				g.Log.Add("⛩️ Prayed at the Shrine of Power ('&') (-30 Gold)! Blessed with +3 Permanent ATK!", tcell.ColorGold)
			case 1:
				g.Player.BaseDEF += 2
				g.Log.Add("⛩️ Prayed at the Shrine of Power ('&') (-30 Gold)! Blessed with +2 Permanent DEF!", tcell.ColorGold)
			case 2:
				g.Player.MaxHP += 15
				g.Player.HP += 15
				g.Log.Add("⛩️ Prayed at the Shrine of Power ('&') (-30 Gold)! Blessed with +15 Permanent Max HP!", tcell.ColorGold)
			}
			g.EndPlayerTurn()
		} else {
			g.Log.Add("⛩️ The Ancient Shrine ('&') requires an offering of 30 Gold Coins to grant its blessing.", tcell.ColorDarkGray)
		}
		return
	}

	if g.Map.IsBlocked(destX, destY) {
		return
	}

	// 3. Move Player
	g.LastPlayerX = g.Player.X
	g.LastPlayerY = g.Player.Y
	g.Player.Move(dx, dy)

	// Check Tutorial Room Triggers on Floor 0
	if g.Floor == 0 {
		g.checkTutorialTriggers()
	}

	// Check if stepping on an item
	for _, itm := range g.Items {
		if itm.X == g.Player.X && itm.Y == g.Player.Y {
			g.Log.Add(fmt.Sprintf("Item found: %s ('%c'). Press 'g' to pick it up!", itm.Name, itm.Rune), tcell.ColorLightCyan)
			break
		}
	}

	// Check if standing on stairs
	if g.Map.Tiles[g.Player.X][g.Player.Y].Type == mapgen.TileStairsDown {
		if g.Floor == 0 {
			g.Log.Add("Stairs Down ('>')! Press '>' to complete training and enter Floor 1.", tcell.ColorGold)
		} else {
			g.Log.Add(fmt.Sprintf("Stairs Down ('>')! Press '>' to descend to Floor %d.", g.Floor+1), tcell.ColorGold)
		}
	}

	g.EndPlayerTurn()
}

func (g *Game) checkTutorialTriggers() {
	if len(g.Map.Rooms) < 4 {
		return
	}

	px := g.Player.X
	r2 := g.Map.Rooms[1]
	r3 := g.Map.Rooms[2]
	r4 := g.Map.Rooms[3]

	if px >= r2.X1 && px <= r2.X2 && !g.TutorialTriggers["room2"] {
		g.TutorialTriggers["room2"] = true
		g.Log.Add("TUTORIAL: Armory Chamber! Stand on items and press 'g' to loot.", tcell.ColorYellow)
		g.Log.Add("Press 'i' to open inventory and equip weapons (')') & armor ('[').", tcell.ColorAqua)
	}

	if px >= r3.X1 && px <= r3.X2 && !g.TutorialTriggers["room3"] {
		g.TutorialTriggers["room3"] = true
		g.Log.Add("TUTORIAL: Combat Chamber! Training Goblin ('g') ahead.", tcell.ColorOrangeRed)
		g.Log.Add("Walk into the goblin to engage in tactical turn-based combat!", tcell.ColorYellow)
	}

	if px >= r4.X1 && !g.TutorialTriggers["room4"] {
		g.TutorialTriggers["room4"] = true
		g.Log.Add("TUTORIAL: Exit Chamber! You found the Stairs Down ('>').", tcell.ColorGold)
		g.Log.Add("Stand on the stairs and press '>' to enter Dungeon Floor 1!", tcell.ColorGold)
	}
}

// EndPlayerTurn updates FOV, processes Monster movement on map
func (g *Game) EndPlayerTurn() {
	g.TurnCount++
	g.Map.ComputeFOV(g.Player.X, g.Player.Y, g.FOVRadius)

	// Regain 1 MP every 2 exploration turns
	if g.TurnCount%2 == 0 {
		g.Player.RestoreMP(1)
	}

	// Regain 1 HP every 5 exploration turns
	if g.TurnCount%5 == 0 {
		g.Player.Heal(1)
	}

	// Process Monster AI movement on map
	for _, m := range g.Monsters {
		if !m.IsAlive || m.IsMerchant {
			continue // Merchant stays peaceful in shop
		}

		distToPlayer := entities.Distance(m.X, m.Y, g.Player.X, g.Player.Y)

		if g.Map.Tiles[m.X][m.Y].Visible || distToPlayer <= 4.0 {
			if !m.Alerted {
				m.Alerted = true
				g.alertNearbyMonsters(m.X, m.Y, 4.0)
			}
		}

		if !m.Alerted {
			continue
		}

		// If monster reaches player -> Trigger Battle Screen!
		if distToPlayer <= 1.0 {
			g.StartBattle(m)
			return
		}

		// Pathfind towards player
		stepX, stepY := entities.NextStepTowards(m.X, m.Y, g.Player.X, g.Player.Y, func(x, y int) bool {
			return g.Map.IsBlocked(x, y) || g.isOccupied(x, y)
		})

		if stepX != 0 || stepY != 0 {
			m.Move(stepX, stepY)
		}
	}
}

// PickUpItem picks up the item under the player's feet
func (g *Game) PickUpItem() {
	if g.State != StatePlaying {
		return
	}

	item := g.getItemAt(g.Player.X, g.Player.Y)
	if item == nil {
		g.Log.Add("There is nothing here to pick up.", tcell.ColorDarkGray)
		return
	}

	if item.Type == items.TypeGold {
		g.Player.Gold += item.Value
		g.Log.Add(fmt.Sprintf("Collected %d Gold Coins ('$')! (Total: %d)", item.Value, g.Player.Gold), tcell.ColorYellow)
		g.EndPlayerTurn()
		return
	}

	if item.Type == items.TypeKey {
		g.Player.Keys++
		g.Log.Add("🗝️ Found an Iron Dungeon Key ('k')! Used to unlock Vault Doors ('%').", tcell.ColorGold)
		g.EndPlayerTurn()
		return
	}

	err := g.Player.Inventory.Add(item)
	if err != nil {
		g.Items = append(g.Items, item)
		g.Log.Add(err.Error(), tcell.ColorRed)
		return
	}

	g.Log.Add(fmt.Sprintf("Picked up %s ('%c'). Press 'i' to view/equip it.", item.Name, item.Rune), tcell.ColorGreen)
	g.EndPlayerTurn()
}

// UseInventoryItem uses or equips the item at given index
func (g *Game) UseInventoryItem(index int) {
	if index < 0 || index >= len(g.Player.Inventory.Items) {
		return
	}

	item := g.Player.Inventory.Items[index]

	switch item.Type {
	case items.TypePotion:
		healed := g.Player.Heal(item.HealAmount)
		g.Log.Add(fmt.Sprintf("Drank %s ('!'). Restored %d HP! (HP: %d/%d)", item.Name, healed, g.Player.HP, g.Player.MaxHP), tcell.ColorGreen)
		_, _ = g.Player.Inventory.Remove(index)
		g.EndPlayerTurn()

	case items.TypeScroll:
		if item.ID == "scroll_teleport" {
			for {
				rx := g.RNG.Intn(g.Map.Width)
				ry := g.RNG.Intn(g.Map.Height)
				if !g.Map.IsBlocked(rx, ry) && !g.isOccupied(rx, ry) {
					g.Player.MoveTo(rx, ry)
					g.Log.Add("Scroll of Teleportation ('?') cast! Teleported to a new room!", tcell.ColorViolet)
					break
				}
			}
			_, _ = g.Player.Inventory.Remove(index)
			g.EndPlayerTurn()
		} else if item.ID == "scroll_enchant" {
			g.Player.BaseATK += 3
			g.Log.Add("✨ Read Scroll of Weapon Enchantment! Permanently gained +3 Base ATK!", tcell.ColorGold)
			_, _ = g.Player.Inventory.Remove(index)
			g.EndPlayerTurn()
		}

	case items.TypeWeapon, items.TypeArmor:
		if item.Equipped {
			msg := g.Player.Inventory.Unequip(item)
			g.Log.Add(msg, tcell.ColorLightGray)
		} else {
			msg := g.Player.Inventory.Equip(item)
			g.Log.Add(msg, tcell.ColorLightCyan)
		}
	}
}

// DropInventoryItem drops item onto floor
func (g *Game) DropInventoryItem(index int) {
	item, err := g.Player.Inventory.Remove(index)
	if err != nil {
		return
	}
	item.X = g.Player.X
	item.Y = g.Player.Y
	g.Items = append(g.Items, item)
	g.Log.Add(fmt.Sprintf("Dropped %s on the ground.", item.Name), tcell.ColorDarkGray)
	g.EndPlayerTurn()
}

// DescendStairs moves to next floor if on stairs
func (g *Game) DescendStairs() {
	if g.Map.Tiles[g.Player.X][g.Player.Y].Type != mapgen.TileStairsDown {
		g.Log.Add("There are no stairs here to descend.", tcell.ColorDarkGray)
		return
	}

	nextFloor := g.Floor + 1
	if nextFloor > g.MaxFloors {
		g.Log.Add("You have reached the bottom of the abyss!", tcell.ColorGold)
		return
	}

	if g.Floor == 0 {
		MarkTutorialCompleted()
		g.Log.Add("Training Complete! Descending into Dungeon Floor 1...", tcell.ColorGold)
		g.Log.Add("Tutorial items cleared. Entering the dungeon with starter loadout!", tcell.ColorYellow)
		g.Player = nil // Reset inventory & stats on entry to Floor 1
	} else {
		g.Log.Add(fmt.Sprintf("You descend the stairs into Floor %d...", nextFloor), tcell.ColorGold)
	}
	g.generateFloor(nextFloor)
}
