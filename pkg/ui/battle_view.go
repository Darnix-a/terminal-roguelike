package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"terminal-roguelike/pkg/combat"
	"terminal-roguelike/pkg/items"
)

// RenderBattleScreen renders the dedicated turn-based battle arena
func RenderBattleScreen(screen tcell.Screen, battle *combat.Battle) {
	screen.Clear()
	sw, sh := screen.Size()

	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	DrawBox(screen, 0, 0, sw-1, sh-1, boxStyle, fmt.Sprintf("BATTLE ARENA - VS %s", strings.ToUpper(battle.Monster.Name)))

	// 1. Top Section: Monster Sprite & HP Bar
	monster := battle.Monster
	spriteStartY := 2
	for i, line := range monster.Sprite {
		DrawText(screen, (sw-len(line))/2, spriteStartY+i, tcell.StyleDefault.Foreground(monster.Color).Bold(true), line)
	}

	monsterInfoY := spriteStartY + len(monster.Sprite) + 1
	DrawText(screen, (sw-len(monster.Name))/2, monsterInfoY, tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true), monster.Name)

	// Monster HP Bar
	mHPPct := float64(monster.HP) / float64(monster.MaxHP)
	if mHPPct < 0 {
		mHPPct = 0
	}
	mBarLen := 20
	mFilled := int(mHPPct * float64(mBarLen))
	mHPStr := fmt.Sprintf("[%s%s] %d/%d HP", strings.Repeat("=", mFilled), strings.Repeat(" ", mBarLen-mFilled), monster.HP, monster.MaxHP)
	mHPColor := tcell.ColorGreen
	if mHPPct < 0.3 {
		mHPColor = tcell.ColorRed
	} else if mHPPct < 0.6 {
		mHPColor = tcell.ColorYellow
	}
	DrawText(screen, (sw-len(mHPStr))/2, monsterInfoY+1, tcell.StyleDefault.Foreground(mHPColor).Bold(true), mHPStr)

	// Divider
	divY := monsterInfoY + 3
	for x := 1; x < sw-1; x++ {
		screen.SetContent(x, divY, '─', nil, tcell.StyleDefault.Foreground(tcell.ColorDarkGray))
	}

	// 2. Center Section: Battle Action Log
	logY := divY + 1
	logH := sh - logY - 9
	if logH < 3 {
		logH = 3
	}

	startIdx := len(battle.Log) - logH
	if startIdx < 0 {
		startIdx = 0
	}
	visibleLogs := battle.Log[startIdx:]
	for i, msg := range visibleLogs {
		DrawText(screen, 4, logY+i, tcell.StyleDefault.Foreground(msg.Color), msg.Text)
	}

	// 3. Bottom Section: Player Stats (Left) & Command Menu (Right)
	bottomY := sh - 8
	for x := 1; x < sw-1; x++ {
		screen.SetContent(x, bottomY, '─', nil, tcell.StyleDefault.Foreground(tcell.ColorDarkGray))
	}

	// Player Status (Left)
	player := battle.Player
	DrawText(screen, 4, bottomY+1, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true), fmt.Sprintf("Hero: %s (Lv.%d)", player.Name, player.Level))

	// Player HP Bar
	pHPPct := float64(player.HP) / float64(player.MaxHP)
	if pHPPct < 0 {
		pHPPct = 0
	}
	pBarLen := 14
	pFilled := int(pHPPct * float64(pBarLen))
	pHPStr := fmt.Sprintf("HP: %d/%d [%s%s]", player.HP, player.MaxHP, strings.Repeat("=", pFilled), strings.Repeat(" ", pBarLen-pFilled))
	pHPColor := tcell.ColorGreen
	if pHPPct < 0.3 {
		pHPColor = tcell.ColorRed
	}
	DrawText(screen, 4, bottomY+2, tcell.StyleDefault.Foreground(pHPColor), pHPStr)

	// Player MP Bar
	pMPPct := float64(player.MP) / float64(player.MaxMP)
	if pMPPct < 0 {
		pMPPct = 0
	}
	pMPFilled := int(pMPPct * float64(pBarLen))
	pMPStr := fmt.Sprintf("MP: %d/%d [%s%s]", player.MP, player.MaxMP, strings.Repeat("*", pMPFilled), strings.Repeat(" ", pBarLen-pMPFilled))
	DrawText(screen, 4, bottomY+3, tcell.StyleDefault.Foreground(tcell.ColorAqua), pMPStr)

	DrawText(screen, 4, bottomY+5, tcell.StyleDefault.Foreground(tcell.ColorLightSalmon), fmt.Sprintf("ATK: %d   DEF: %d", player.TotalATK(), player.TotalDEF()))

	// Command Menu (Right)
	menuX := sw/2 + 2
	screen.SetContent(menuX-2, bottomY, '┬', nil, tcell.StyleDefault.Foreground(tcell.ColorDarkGray))
	for y := bottomY + 1; y < sh-1; y++ {
		screen.SetContent(menuX-2, y, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorDarkGray))
	}

	if battle.Result == combat.BattleVictory {
		DrawText(screen, menuX, bottomY+2, tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true), "★ VICTORY ACHIEVED! ★")
		DrawText(screen, menuX, bottomY+4, tcell.StyleDefault.Foreground(tcell.ColorWhite), "Press [Enter / Space] to continue...")
	} else if battle.Result == combat.BattleDefeat {
		DrawText(screen, menuX, bottomY+2, tcell.StyleDefault.Foreground(tcell.ColorDarkRed).Bold(true), "☠ YOU WERE DEFEATED ☠")
		DrawText(screen, menuX, bottomY+4, tcell.StyleDefault.Foreground(tcell.ColorWhite), "Press [Enter / Space] to proceed...")
	} else if battle.Result == combat.BattleFled {
		DrawText(screen, menuX, bottomY+2, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true), "💨 ESCAPED SAFELY!")
		DrawText(screen, menuX, bottomY+4, tcell.StyleDefault.Foreground(tcell.ColorWhite), "Press [Enter / Space] to return to map...")
	} else {
		// Active Turn Menu
		switch battle.SubMenu {
		case combat.MenuMain:
			DrawText(screen, menuX, bottomY+1, tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true), "--- CHOOSE ACTION ---")
			DrawText(screen, menuX, bottomY+2, tcell.StyleDefault.Foreground(tcell.ColorWhite), "[1] Weapon Attack")
			DrawText(screen, menuX, bottomY+3, tcell.StyleDefault.Foreground(tcell.ColorAqua), "[2] Magic Skills (MP)")
			DrawText(screen, menuX, bottomY+4, tcell.StyleDefault.Foreground(tcell.ColorLightGreen), "[3] Use Potion")
			DrawText(screen, menuX, bottomY+5, tcell.StyleDefault.Foreground(tcell.ColorYellow), "[4] Flee Battle")

		case combat.MenuSkills:
			DrawText(screen, menuX, bottomY+1, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true), "--- MAGIC SKILLS (Press 1-4) ---")
			for i, sk := range player.Skills {
				costColor := tcell.ColorAqua
				if player.MP < sk.MPCost {
					costColor = tcell.ColorDarkGray
				}
				line := fmt.Sprintf("[%d] %s (%d MP) - %s", i+1, sk.Name, sk.MPCost, sk.Description)
				DrawText(screen, menuX, bottomY+2+i, tcell.StyleDefault.Foreground(costColor), line)
			}
			DrawText(screen, menuX, bottomY+6, tcell.StyleDefault.Foreground(tcell.ColorDarkGray), "[Esc] Back to Actions")

		case combat.MenuItems:
			DrawText(screen, menuX, bottomY+1, tcell.StyleDefault.Foreground(tcell.ColorLightGreen).Bold(true), "--- POTIONS (Press a-z) ---")
			potionsFound := 0
			for i, itm := range player.Inventory.Items {
				if itm.Type == items.TypePotion {
					letter := rune('a' + i)
					line := fmt.Sprintf("[%c] %s (+%d HP)", letter, itm.Name, itm.HealAmount)
					DrawText(screen, menuX, bottomY+2+potionsFound, tcell.StyleDefault.Foreground(tcell.ColorLightCoral), line)
					potionsFound++
					if potionsFound >= 4 {
						break
					}
				}
			}
			if potionsFound == 0 {
				DrawText(screen, menuX, bottomY+2, tcell.StyleDefault.Foreground(tcell.ColorDarkGray), "No potions in backpack!")
			}
			DrawText(screen, menuX, bottomY+6, tcell.StyleDefault.Foreground(tcell.ColorDarkGray), "[Esc] Back to Actions")
		}
	}

	screen.Show()
}
