package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"terminal-roguelike/pkg/engine"
	"terminal-roguelike/pkg/entities"
	"terminal-roguelike/pkg/items"
)

// RenderInventoryModal renders the interactive keyboard/arrow navigation inventory
func RenderInventoryModal(screen tcell.Screen, g *engine.Game) {
	sw, sh := screen.Size()
	w, h := 58, 18
	x1 := (sw - w) / 2
	y1 := (sh - h) / 2
	x2 := x1 + w
	y2 := y1 + h

	// Fill background
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}

	DrawBox(screen, x1, y1, x2, y2, tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true), "INVENTORY")

	inv := g.Player.Inventory
	if len(inv.Items) == 0 {
		DrawText(screen, x1+4, y1+4, tcell.StyleDefault.Foreground(tcell.ColorDarkGray), "Your backpack is empty.")
	} else {
		if g.InventoryIdx < 0 {
			g.InventoryIdx = 0
		}
		if g.InventoryIdx >= len(inv.Items) {
			g.InventoryIdx = len(inv.Items) - 1
		}

		for i, itm := range inv.Items {
			isSelected := (i == g.InventoryIdx)
			style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

			status := ""
			if itm.Equipped {
				status = " [EQUIPPED]"
			}

			prefix := "  "
			if isSelected {
				prefix = "> "
				style = tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorNavy).Bold(true)
			} else if itm.Equipped {
				style = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
			} else if itm.Type == items.TypePotion {
				style = tcell.StyleDefault.Foreground(tcell.ColorLightCoral)
			} else if itm.Type == items.TypeScroll {
				style = tcell.StyleDefault.Foreground(tcell.ColorViolet)
			}

			line := fmt.Sprintf("%s%s - %s%s", prefix, itm.Name, itm.Description, status)
			if isSelected && len(line) < w-4 {
				line += fmt.Sprintf("%*s", w-4-len(line), "")
			}
			DrawText(screen, x1+2, y1+2+i, style, line)
		}
	}

	footer := "▲/▼ Move | [Enter] Use/Equip | [d/x] Drop | [Esc/i] Close"
	DrawText(screen, (sw-len(footer))/2, y2-2, tcell.StyleDefault.Foreground(tcell.ColorYellow), footer)
}

// RenderGameOverModal renders the death screen
func RenderGameOverModal(screen tcell.Screen, g *engine.Game) {
	sw, sh := screen.Size()
	w, h := 46, 14
	x1 := (sw - w) / 2
	y1 := (sh - h) / 2
	x2 := x1 + w
	y2 := y1 + h

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}

	DrawBox(screen, x1, y1, x2, y2, tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true), "YOU DIED")

	score := g.Player.Gold + (g.Player.Kills * 25) + (g.Floor * 100) + (g.Player.Level * 50)

	DrawText(screen, x1+4, y1+3, tcell.StyleDefault.Foreground(tcell.ColorDarkRed).Bold(true), "Your journey ends in the dark abyss...")
	DrawText(screen, x1+4, y1+5, tcell.StyleDefault.Foreground(tcell.ColorWhite), fmt.Sprintf("Dungeon Floor: %d / %d", g.Floor, g.MaxFloors))
	DrawText(screen, x1+4, y1+6, tcell.StyleDefault.Foreground(tcell.ColorWhite), fmt.Sprintf("Hero Level:    %d", g.Player.Level))
	DrawText(screen, x1+4, y1+7, tcell.StyleDefault.Foreground(tcell.ColorWhite), fmt.Sprintf("Monsters Slain: %d", g.Player.Kills))
	DrawText(screen, x1+4, y1+8, tcell.StyleDefault.Foreground(tcell.ColorGold), fmt.Sprintf("Gold Treasure: %d", g.Player.Gold))
	DrawText(screen, x1+4, y1+10, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true), fmt.Sprintf("FINAL SCORE:   %d", score))

	DrawText(screen, x1+4, y2-2, tcell.StyleDefault.Foreground(tcell.ColorDarkGray), "Press [r] to Restart or [q] to Quit")
}

// RenderVictoryModal renders the game victory screen
func RenderVictoryModal(screen tcell.Screen, g *engine.Game) {
	sw, sh := screen.Size()
	w, h := 50, 15
	x1 := (sw - w) / 2
	y1 := (sh - h) / 2
	x2 := x1 + w
	y2 := y1 + h

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}

	DrawBox(screen, x1, y1, x2, y2, tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true), "VICTORY!")

	score := g.Player.Gold + (g.Player.Kills * 50) + 1000 + (g.Player.Level * 100)

	DrawText(screen, x1+4, y1+3, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true), "THE ANCIENT DRAGON HAS BEEN DEFEATED!")
	DrawText(screen, x1+4, y1+4, tcell.StyleDefault.Foreground(tcell.ColorGreen), "You have conquered the Dungeon of Shadows!")
	DrawText(screen, x1+4, y1+6, tcell.StyleDefault.Foreground(tcell.ColorWhite), fmt.Sprintf("Hero Level:     %d", g.Player.Level))
	DrawText(screen, x1+4, y1+7, tcell.StyleDefault.Foreground(tcell.ColorWhite), fmt.Sprintf("Monsters Slain: %d", g.Player.Kills))
	DrawText(screen, x1+4, y1+8, tcell.StyleDefault.Foreground(tcell.ColorGold), fmt.Sprintf("Gold Treasure:  %d", g.Player.Gold))
	DrawText(screen, x1+4, y1+10, tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true), fmt.Sprintf("VICTORY SCORE:  %d", score))

	DrawText(screen, x1+4, y2-2, tcell.StyleDefault.Foreground(tcell.ColorDarkGray), "Press [r] to Play Again or [q] to Quit")
}

// RenderShopModal renders the interactive dungeon merchant shop
func RenderShopModal(screen tcell.Screen, g *engine.Game) {
	sw, sh := screen.Size()
	w, h := 56, 17
	x1 := (sw - w) / 2
	y1 := (sh - h) / 2
	x2 := x1 + w
	y2 := y1 + h

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}

	title := fmt.Sprintf("DUNGEON MERCHANT (Your Gold: %d)", g.Player.Gold)
	DrawBox(screen, x1, y1, x2, y2, tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true), title)

	DrawText(screen, x1+3, y1+2, tcell.StyleDefault.Foreground(tcell.ColorYellow), "\"Welcome, traveler! Spend your gold on fine dungeon wares!\"")

	DrawText(screen, x1+3, y1+4, tcell.StyleDefault.Foreground(tcell.ColorLightCoral), "[1] Health Potion (+25 HP)              -  20 Gold")
	DrawText(screen, x1+3, y1+5, tcell.StyleDefault.Foreground(tcell.ColorRed), "[2] Greater Health Draught (+50 HP)     -  40 Gold")
	DrawText(screen, x1+3, y1+6, tcell.StyleDefault.Foreground(tcell.ColorGold), "[3] Scroll of Weapon Enchant (+3 ATK)   -  60 Gold")
	DrawText(screen, x1+3, y1+7, tcell.StyleDefault.Foreground(tcell.ColorSilver), "[4] Dragonscale Shield (+4 DEF)         -  75 Gold")
	DrawText(screen, x1+3, y1+8, tcell.StyleDefault.Foreground(tcell.ColorViolet), "[5] Scroll of Teleportation             -  30 Gold")

	DrawText(screen, x1+3, y1+11, tcell.StyleDefault.Foreground(tcell.ColorOrangeRed), "[A] Provoke / Attack the Merchant")
	DrawText(screen, x1+3, y1+13, tcell.StyleDefault.Foreground(tcell.ColorYellow), "Press [1-5] to Buy | [A] Attack | [Esc/q] Leave Shop")
}

// RenderLevelUpModal renders the skill choice / replacement level-up modal
func RenderLevelUpModal(screen tcell.Screen, g *engine.Game) {
	sw, sh := screen.Size()
	w, h := 66, 18
	x1 := (sw - w) / 2
	y1 := (sh - h) / 2
	x2 := x1 + w
	y2 := y1 + h

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
		}
	}

	if g.PendingSkill != nil {
		// Replacement selection screen (5/5 capacity reached)
		title := fmt.Sprintf("REPLACE SKILL WITH [%s]", g.PendingSkill.Name)
		DrawBox(screen, x1, y1, x2, y2, tcell.StyleDefault.Foreground(tcell.ColorOrangeRed).Bold(true), title)

		DrawText(screen, x1+3, y1+2, tcell.StyleDefault.Foreground(tcell.ColorYellow), "Ability capacity reached (5/5)! Choose which skill to forget:")

		for i, sk := range g.Player.Skills {
			rowY := y1 + 4 + (i * 2)
			line := fmt.Sprintf("[%d] %s (%d MP) - %s", i+1, sk.Name, sk.MPCost, sk.Description)
			if len(line) > w-6 {
				line = line[:w-9] + "..."
			}
			DrawText(screen, x1+4, rowY, tcell.StyleDefault.Foreground(tcell.ColorLightCoral), line)
		}

		prompt := fmt.Sprintf("Press [1-%d] to Replace | [0/Esc] Cancel & Keep Current", len(g.Player.Skills))
		DrawText(screen, (sw-len(prompt))/2, y2-2, tcell.StyleDefault.Foreground(tcell.ColorGold), prompt)
		return
	}

	// Offered skills screen
	skillCount := len(g.Player.Skills)
	statusStr := fmt.Sprintf("[%d/%d Abilities]", skillCount, entities.MaxSkills)
	title := fmt.Sprintf("🌟 LEVEL UP! (Level %d) %s", g.Player.Level, statusStr)
	DrawBox(screen, x1, y1, x2, y2, tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true), title)

	if skillCount >= entities.MaxSkills {
		DrawText(screen, x1+3, y1+2, tcell.StyleDefault.Foreground(tcell.ColorOrange), "Select an ability to REPLACE an existing skill, or press [0] to keep current:")
	} else {
		DrawText(screen, x1+3, y1+2, tcell.StyleDefault.Foreground(tcell.ColorYellow), "Select a new ability to learn, or press [0] to skip:")
	}

	for i, sk := range g.OfferedSkills {
		rowY := y1 + 4 + (i * 3)
		nameLine := fmt.Sprintf("[%d] %s (%d MP)", i+1, sk.Name, sk.MPCost)
		DrawText(screen, x1+4, rowY, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true), nameLine)
		DrawText(screen, x1+8, rowY+1, tcell.StyleDefault.Foreground(tcell.ColorLightGray), sk.Description)
	}

	prompt := fmt.Sprintf("Press [1-%d] to Select | [0] Skip / Don't Take Skill", len(g.OfferedSkills))
	DrawText(screen, (sw-len(prompt))/2, y2-2, tcell.StyleDefault.Foreground(tcell.ColorGold), prompt)
}
