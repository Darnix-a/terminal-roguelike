package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"terminal-roguelike/pkg/engine"
)

var TitleBanner = []string{
	` _______ ______ _____  __  __ _____ _   _          _      `,
	`|__   __|  ____|  __ \|  \/  |_   _| \ | |   /\   | |     `,
	`   | |  | |__  | |__) | \  / | | | |  \| |  /  \  | |     `,
	`   | |  |  __| |  _  /| |\/| | | | | . ` + "`" + ` | / /\ \ | |     `,
	`   | |  | |____| | \ \| |  | |_| |_| |\  |/ ____ \| |____ `,
	`   |_|  |______|_|  \_\_|  |_|_____|_| \_/_/    \_\______|`,
	`               === D U N G E O N   C R A W L E R ===      `,
}

// RenderMainMenu renders the primary main menu screen
func RenderMainMenu(screen tcell.Screen, hasSave bool) {
	screen.Clear()
	sw, sh := screen.Size()

	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true)
	DrawBox(screen, 0, 0, sw-1, sh-1, boxStyle, "TERMINAL ROGUELIKE")

	// Render Title Banner
	bannerStartY := 2
	for i, line := range TitleBanner {
		color := tcell.ColorGold
		if i == len(TitleBanner)-1 {
			color = tcell.ColorYellow
		}
		DrawText(screen, (sw-len(line))/2, bannerStartY+i, tcell.StyleDefault.Foreground(color).Bold(true), line)
	}

	// Menu Options
	menuY := bannerStartY + len(TitleBanner) + 2
	boxW := 48
	boxX := (sw - boxW) / 2

	DrawBox(screen, boxX, menuY, boxX+boxW, menuY+10, tcell.StyleDefault.Foreground(tcell.ColorSkyblue), "MAIN MENU")

	optStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	dimStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)

	DrawText(screen, boxX+4, menuY+2, optStyle, "[1] New Game (Descent into Floor 1)")

	if hasSave {
		DrawText(screen, boxX+4, menuY+3, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true), "[2] Continue Saved Run")
	} else {
		DrawText(screen, boxX+4, menuY+3, dimStyle, "[2] Continue (No Save Found)")
	}

	DrawText(screen, boxX+4, menuY+4, optStyle, "[3] Training Grounds (Tutorial Floor)")
	DrawText(screen, boxX+4, menuY+5, optStyle, "[4] Hall of Fame (High Scores)")
	DrawText(screen, boxX+4, menuY+6, optStyle, "[5] Guide & Keybindings")
	DrawText(screen, boxX+4, menuY+8, tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true), "[Q] Exit to Terminal")

	// Footer
	footText := "Select an option [1-5] or [Q]uit"
	DrawText(screen, (sw-len(footText))/2, sh-3, tcell.StyleDefault.Foreground(tcell.ColorDarkGray), footText)

	screen.Show()
}

// RenderHighScoresScreen renders the Hall of Fame leaderboard
func RenderHighScoresScreen(screen tcell.Screen, highScores []engine.HighScoreEntry) {
	screen.Clear()
	sw, sh := screen.Size()

	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true)
	DrawBox(screen, 0, 0, sw-1, sh-1, boxStyle, "HALL OF FAME - TOP 10 DUNGEON RUNS")

	header := "RANK  SCORE   FLOOR  LEVEL  KILLS  GOLD  OUTCOME"
	DrawText(screen, 4, 3, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true), header)

	divider := strings.Repeat("─", sw-8)
	DrawText(screen, 4, 4, tcell.StyleDefault.Foreground(tcell.ColorDarkGray), divider)

	if len(highScores) == 0 {
		DrawText(screen, 6, 6, tcell.StyleDefault.Foreground(tcell.ColorWhite), "No recorded runs yet! Venture into the dungeon and etch your name into glory.")
	} else {
		for i, entry := range highScores {
			rankColor := tcell.ColorWhite
			if i == 0 {
				rankColor = tcell.ColorGold
			} else if i == 1 {
				rankColor = tcell.ColorSilver
			} else if i == 2 {
				rankColor = tcell.ColorPeru
			}

			line := fmt.Sprintf("#%-3d %-7d Floor %-2d Lv.%-2d  %-5d  %-4d  %s",
				i+1, entry.Score, entry.Floor, entry.Level, entry.Kills, entry.Gold, entry.Outcome)

			DrawText(screen, 4, 5+i, tcell.StyleDefault.Foreground(rankColor), line)
		}
	}

	backPrompt := "Press [Esc], [q], or [Enter] to return to Main Menu"
	DrawText(screen, (sw-len(backPrompt))/2, sh-3, tcell.StyleDefault.Foreground(tcell.ColorGold), backPrompt)

	screen.Show()
}

// RenderHelpScreen renders the full gameplay guide & controls
func RenderHelpScreen(screen tcell.Screen) {
	screen.Clear()
	sw, sh := screen.Size()

	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorSkyblue).Bold(true)
	DrawBox(screen, 0, 0, sw-1, sh-1, boxStyle, "HOW TO PLAY & CONTROLS")

	y := 2
	sectionColor := tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true)
	bodyColor := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	DrawText(screen, 3, y, sectionColor, "1. EXPLORATION CONTROLS:")
	y++
	DrawText(screen, 5, y, bodyColor, "• Move / Step:        WASD or Arrow Keys / vi-keys (hjkl)")
	y++
	DrawText(screen, 5, y, bodyColor, "• Wait 1 Turn:        Spacebar or '.'")
	y++
	DrawText(screen, 5, y, bodyColor, "• Pick Up Item:       'g' or ',' (Picks up items, keys 'k', and gold '$')")
	y++
	DrawText(screen, 5, y, bodyColor, "• Inventory:          'i' (Use/Equip items with [a-z], Drop with [Shift+D])")
	y++
	DrawText(screen, 5, y, bodyColor, "• Descend Stairs:     '>' or '.' when standing on Stairs Down ('>')")
	y++
	DrawText(screen, 5, y, bodyColor, "• Save & Quit:        'S' (Saves your run to resume later)")
	y += 2

	DrawText(screen, 3, y, sectionColor, "2. TACTICAL TURN-BASED COMBAT:")
	y++
	DrawText(screen, 5, y, bodyColor, "• Attack:             [1] Standard weapon strike (15% Critical Chance)")
	y++
	DrawText(screen, 5, y, bodyColor, "• Heavy Slash:        [2]->[1] (6 MP) Cleaves enemy armor")
	y++
	DrawText(screen, 5, y, bodyColor, "• Fireball:           [2]->[2] (8 MP) High magic burst damage")
	y++
	DrawText(screen, 5, y, bodyColor, "• Shield Guard:       [2]->[3] (4 MP) Crucial! Mitigates 70% incoming damage")
	y++
	DrawText(screen, 5, y, bodyColor, "• Holy Heal:          [2]->[4] (10 MP) Restores 30 HP")
	y++
	DrawText(screen, 5, y, bodyColor, "• Attack Telegraphs:  When an enemy is charging a heavy move, USE SHIELD GUARD!")
	y += 2

	DrawText(screen, 3, y, sectionColor, "3. DUNGEON OBJECTS & NPC:")
	y++
	DrawText(screen, 5, y, bodyColor, "• Merchant ('S'):     Peaceful shopkeeper on Floors 2 & 4 to buy potions/gear")
	y++
	DrawText(screen, 5, y, bodyColor, "• Locked Vaults ('%'): Requires an Iron Key ('k') found on that floor")
	y++
	DrawText(screen, 5, y, bodyColor, "• Chests ('='):       Contains rare loot (beware of Hungry Mimics 'M'!)")
	y++
	DrawText(screen, 5, y, bodyColor, "• Fountains ('0'):    Fully restores HP and MP")
	y++
	DrawText(screen, 5, y, bodyColor, "• Shrines ('&'):      Offering 30 Gold grants permanent ATK, DEF, or Max HP")

	backPrompt := "Press [Esc], [q], or [Enter] to return to Main Menu"
	DrawText(screen, (sw-len(backPrompt))/2, sh-2, tcell.StyleDefault.Foreground(tcell.ColorGold), backPrompt)

	screen.Show()
}
