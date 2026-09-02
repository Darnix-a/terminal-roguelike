package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"terminal-roguelike/pkg/engine"
	"terminal-roguelike/pkg/mapgen"
)

// DrawText draws a string with style at (x, y)
func DrawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	for _, r := range text {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

// DrawBox draws a border box
func DrawBox(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, title string) {
	for x := x1; x <= x2; x++ {
		screen.SetContent(x, y1, '─', nil, style)
		screen.SetContent(x, y2, '─', nil, style)
	}
	for y := y1; y <= y2; y++ {
		screen.SetContent(x1, y, '│', nil, style)
		screen.SetContent(x2, y, '│', nil, style)
	}
	screen.SetContent(x1, y1, '┌', nil, style)
	screen.SetContent(x2, y1, '┐', nil, style)
	screen.SetContent(x1, y2, '└', nil, style)
	screen.SetContent(x2, y2, '┘', nil, style)

	if title != "" {
		DrawText(screen, x1+2, y1, style.Bold(true), fmt.Sprintf(" %s ", title))
	}
}

// Render draws the entire game frame with full-height sidebar and seamless viewport
func Render(screen tcell.Screen, g *engine.Game) {
	// If in Battle, render dedicated Battle Screen
	if g.State == engine.StateCombat && g.ActiveBattle != nil {
		RenderBattleScreen(screen, g.ActiveBattle)
		return
	}

	screen.Clear()
	sw, sh := screen.Size()

	sidebarW := 32
	if sw < 85 {
		sidebarW = 28
	}

	logH := 7
	if sh < 26 {
		logH = 5
	}

	mapBoxW := sw - sidebarW - 2
	mapBoxH := sh - logH - 2

	if mapBoxW < 20 {
		mapBoxW = 20
	}
	if mapBoxH < 10 {
		mapBoxH = 10
	}

	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkSlateGray)

	// 1. Render Map Box (Top-Left)
	floorTitle := fmt.Sprintf("Dungeon Floor %d/%d", g.Floor, g.MaxFloors)
	if g.Floor == 0 {
		floorTitle = "Tutorial: Training Grounds"
	}
	DrawBox(screen, 0, 0, mapBoxW, mapBoxH, boxStyle, floorTitle)

	// 2. Render Map Tiles
	for x := 0; x < g.Map.Width && x < mapBoxW-1; x++ {
		for y := 0; y < g.Map.Height && y < mapBoxH-1; y++ {
			tile := g.Map.Tiles[x][y]
			screenX := x + 1
			screenY := y + 1

			if !tile.Explored {
				continue // Fog of war
			}

			style := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
			if tile.Visible {
				style = tcell.StyleDefault.Foreground(tcell.ColorLightGray)
			}

			var r rune
			switch tile.Type {
			case mapgen.TileWall:
				r = '#'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorDarkKhaki)
				}
			case mapgen.TileFloor:
				r = '.'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorDimGray)
				}
			case mapgen.TileDoorClosed:
				r = '+'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorBrown).Bold(true)
				}
			case mapgen.TileDoorOpen:
				r = '/'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorBrown)
				}
			case mapgen.TileDoorLocked:
				r = '%'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
				}
			case mapgen.TileStairsDown:
				r = '>'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true)
				}
			case mapgen.TileChest:
				r = '='
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true)
				}
			case mapgen.TileChestOpened:
				r = '='
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
				}
			case mapgen.TileFountain:
				r = '0'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true)
				}
			case mapgen.TileFountainUsed:
				r = '0'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorDimGray)
				}
			case mapgen.TileShrine:
				r = '&'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorMediumPurple).Bold(true)
				}
			case mapgen.TileShrineUsed:
				r = '&'
				if tile.Visible {
					style = tcell.StyleDefault.Foreground(tcell.ColorDimGray)
				}
			}

			screen.SetContent(screenX, screenY, r, nil, style)
		}
	}

	// 3. Render Items
	for _, item := range g.Items {
		if item.X < mapBoxW-1 && item.Y < mapBoxH-1 && g.Map.Tiles[item.X][item.Y].Visible {
			screen.SetContent(item.X+1, item.Y+1, item.Rune, nil, tcell.StyleDefault.Foreground(item.Color).Bold(true))
		}
	}

	// 4. Render Monsters
	for _, monster := range g.Monsters {
		if monster.IsAlive && monster.X < mapBoxW-1 && monster.Y < mapBoxH-1 && g.Map.Tiles[monster.X][monster.Y].Visible {
			screen.SetContent(monster.X+1, monster.Y+1, monster.Rune, nil, tcell.StyleDefault.Foreground(monster.Color).Bold(true))
		}
	}

	// 5. Render Player
	if g.Player.X < mapBoxW-1 && g.Player.Y < mapBoxH-1 {
		screen.SetContent(g.Player.X+1, g.Player.Y+1, g.Player.Rune, nil, tcell.StyleDefault.Foreground(g.Player.Color).Bold(true))
	}

	// 6. Render Action Log Box (Bottom-Left)
	logY := mapBoxH + 1
	DrawBox(screen, 0, logY, mapBoxW, sh-1, boxStyle, "Action Log & Guide")

	maxLines := sh - logY - 2
	if maxLines > 0 {
		recentMsgs := g.Log.GetRecent(maxLines)
		for i, msg := range recentMsgs {
			DrawText(screen, 2, logY+1+i, tcell.StyleDefault.Foreground(msg.Color), msg.Text)
		}
	}

	// 7. Render FULL-HEIGHT Right Sidebar (From y=0 to y=sh-1)
	sideX := mapBoxW + 1
	DrawBox(screen, sideX, 0, sw-1, sh-1, boxStyle, "Hero & Guide")

	DrawText(screen, sideX+2, 1, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true), fmt.Sprintf("Lv.%d %s", g.Player.Level, g.Player.Name))

	// HP Bar
	hpPct := float64(g.Player.HP) / float64(g.Player.MaxHP)
	if hpPct < 0 {
		hpPct = 0
	}
	barLen := 10
	filled := int(hpPct * float64(barLen))
	hpColor := tcell.ColorGreen
	if hpPct < 0.3 {
		hpColor = tcell.ColorRed
	} else if hpPct < 0.6 {
		hpColor = tcell.ColorYellow
	}
	hpBarStr := fmt.Sprintf("[%s%s]", strings.Repeat("=", filled), strings.Repeat(" ", barLen-filled))
	DrawText(screen, sideX+2, 2, tcell.StyleDefault.Foreground(hpColor), fmt.Sprintf("HP: %d/%d %s", g.Player.HP, g.Player.MaxHP, hpBarStr))

	// MP Bar
	mpPct := float64(g.Player.MP) / float64(g.Player.MaxMP)
	if mpPct < 0 {
		mpPct = 0
	}
	mpFilled := int(mpPct * float64(barLen))
	mpBarStr := fmt.Sprintf("[%s%s]", strings.Repeat("*", mpFilled), strings.Repeat(" ", barLen-mpFilled))
	DrawText(screen, sideX+2, 3, tcell.StyleDefault.Foreground(tcell.ColorAqua), fmt.Sprintf("MP: %d/%d %s", g.Player.MP, g.Player.MaxMP, mpBarStr))

	// Combat Stats
	DrawText(screen, sideX+2, 5, tcell.StyleDefault.Foreground(tcell.ColorLightSalmon), fmt.Sprintf("ATK: %d (%d+%d)", g.Player.TotalATK(), g.Player.BaseATK, g.Player.Inventory.TotalBonusATK()))
	DrawText(screen, sideX+2, 6, tcell.StyleDefault.Foreground(tcell.ColorSkyblue), fmt.Sprintf("DEF: %d (%d+%d)", g.Player.TotalDEF(), g.Player.BaseDEF, g.Player.Inventory.TotalBonusDEF()))
	DrawText(screen, sideX+2, 7, tcell.StyleDefault.Foreground(tcell.ColorGold), fmt.Sprintf("Gold: %d   Keys: %d", g.Player.Gold, g.Player.Keys))

	// Equipment
	weaponName := "(None)"
	if g.Player.Inventory.EquippedWeapon != nil {
		weaponName = g.Player.Inventory.EquippedWeapon.Name
	}
	DrawText(screen, sideX+2, 9, tcell.StyleDefault.Foreground(tcell.ColorLightCyan), fmt.Sprintf("Wpn: %s", weaponName))

	armorName := "(None)"
	if g.Player.Inventory.EquippedArmor != nil {
		armorName = g.Player.Inventory.EquippedArmor.Name
	}
	DrawText(screen, sideX+2, 10, tcell.StyleDefault.Foreground(tcell.ColorSilver), fmt.Sprintf("Arm: %s", armorName))

	// Controls
	DrawText(screen, sideX+2, 12, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true), "--- CONTROLS ---")
	DrawText(screen, sideX+2, 13, tcell.StyleDefault.Foreground(tcell.ColorWhite), "Move/Interact: WASD/Arrows")
	DrawText(screen, sideX+2, 14, tcell.StyleDefault.Foreground(tcell.ColorWhite), "[g] Pick Up   [i] Backpack")
	DrawText(screen, sideX+2, 15, tcell.StyleDefault.Foreground(tcell.ColorWhite), "[>] Descend   [Space] Rest")
	DrawText(screen, sideX+2, 16, tcell.StyleDefault.Foreground(tcell.ColorWhite), "[q] Quit Game")

	// Map Legend (Compact 2-column)
	DrawText(screen, sideX+2, 18, tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true), "--- MAP LEGEND ---")
	DrawText(screen, sideX+2, 19, tcell.StyleDefault.Foreground(tcell.ColorYellow), "@ Hero")
	DrawText(screen, sideX+12, 19, tcell.StyleDefault.Foreground(tcell.ColorDarkKhaki), "# Wall")
	DrawText(screen, sideX+2, 20, tcell.StyleDefault.Foreground(tcell.ColorDimGray), ". Floor")
	DrawText(screen, sideX+12, 20, tcell.StyleDefault.Foreground(tcell.ColorBrown), "+ Door")
	DrawText(screen, sideX+2, 21, tcell.StyleDefault.Foreground(tcell.ColorRed), "% Vault")
	DrawText(screen, sideX+12, 21, tcell.StyleDefault.Foreground(tcell.ColorGold), "k Key")
	DrawText(screen, sideX+2, 22, tcell.StyleDefault.Foreground(tcell.ColorGold), "= Chest")
	DrawText(screen, sideX+12, 22, tcell.StyleDefault.Foreground(tcell.ColorOrange), "M Mimic")
	DrawText(screen, sideX+2, 23, tcell.StyleDefault.Foreground(tcell.ColorAqua), "0 Spring")
	DrawText(screen, sideX+12, 23, tcell.StyleDefault.Foreground(tcell.ColorMediumPurple), "& Shrine")
	DrawText(screen, sideX+2, 24, tcell.StyleDefault.Foreground(tcell.ColorGold), "> Stairs")
	DrawText(screen, sideX+12, 24, tcell.StyleDefault.Foreground(tcell.ColorGold), "$ Gold")
	DrawText(screen, sideX+2, 25, tcell.StyleDefault.Foreground(tcell.ColorSkyblue), ") Weapon")
	DrawText(screen, sideX+12, 25, tcell.StyleDefault.Foreground(tcell.ColorSilver), "[ Armor")
	DrawText(screen, sideX+2, 26, tcell.StyleDefault.Foreground(tcell.ColorRed), "! Potion")
	DrawText(screen, sideX+12, 26, tcell.StyleDefault.Foreground(tcell.ColorViolet), "? Scroll")
	DrawText(screen, sideX+2, 27, tcell.StyleDefault.Foreground(tcell.ColorGreen), "g Goblin")
	DrawText(screen, sideX+12, 27, tcell.StyleDefault.Foreground(tcell.ColorOlive), "o Orc")
	DrawText(screen, sideX+2, 28, tcell.StyleDefault.Foreground(tcell.ColorWhite), "s Skele")
	DrawText(screen, sideX+12, 28, tcell.StyleDefault.Foreground(tcell.ColorPurple), "w Wizard")
	DrawText(screen, sideX+2, 29, tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true), "D Dragon Boss (Floor 5)")

	// 8. Modals
	switch g.State {
	case engine.StateInventory:
		RenderInventoryModal(screen, g)
	case engine.StateGameOver:
		RenderGameOverModal(screen, g)
	case engine.StateVictory:
		RenderVictoryModal(screen, g)
	}

	screen.Show()
}
