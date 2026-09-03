package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"

	"terminal-roguelike/pkg/combat"
	"terminal-roguelike/pkg/engine"
	"terminal-roguelike/pkg/ui"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to initialize terminal screen: %v\n", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("Screen init error: %v\n", err)
	}
	defer screen.Fini()

	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))
	screen.Clear()

	sw, sh := screen.Size()
	sidebarW := 32
	if sw < 85 {
		sidebarW = 28
	}
	logH := 6
	if sh < 26 {
		logH = 5
	}
	mapW := sw - sidebarW - 2
	mapH := sh - logH - 3

	game := engine.NewGame(mapW, mapH)
	dropMode := false

	// Main Game Loop
	for {
		ui.Render(screen, game)

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()

		case *tcell.EventKey:
			// Global Quit
			if ev.Key() == tcell.KeyCtrlC {
				return
			}

			// Handle Input by State
			switch game.State {
			case engine.StateGameOver, engine.StateVictory:
				if ev.Rune() == 'r' || ev.Rune() == 'R' {
					game = engine.NewGame(mapW, mapH)
				} else if ev.Rune() == 'q' || ev.Rune() == 'Q' || ev.Key() == tcell.KeyEscape {
					return
				}

			case engine.StateCombat:
				battle := game.ActiveBattle
				if battle == nil {
					game.State = engine.StatePlaying
					continue
				}

				// If battle finished, any key returns to map or game over
				if battle.Result != combat.BattleOngoing {
					if ev.Key() == tcell.KeyEnter || ev.Key() == tcell.KeyEscape || ev.Rune() == ' ' || ev.Rune() == 'q' {
						game.ConcludeBattle()
					}
					continue
				}

				// Active Battle Menu Navigation
				switch battle.SubMenu {
				case combat.MenuMain:
					switch ev.Rune() {
					case '1':
						battle.PlayerAttack()
					case '2':
						battle.SubMenu = combat.MenuSkills
					case '3':
						battle.SubMenu = combat.MenuItems
					case '4':
						battle.PlayerFlee()
					}

				case combat.MenuSkills:
					if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
						battle.SubMenu = combat.MenuMain
					} else {
						r := ev.Rune()
						if r >= '1' && r <= '4' {
							skillIdx := int(r - '1')
							battle.PlayerUseSkill(skillIdx)
						}
					}

				case combat.MenuItems:
					if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
						battle.SubMenu = combat.MenuMain
					} else {
						r := ev.Rune()
						if r >= 'a' && r <= 'z' {
							itemIdx := int(r - 'a')
							battle.PlayerUseItem(itemIdx)
						}
					}
				}

			case engine.StateInventory:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'i' || ev.Rune() == 'q' {
					game.State = engine.StatePlaying
					dropMode = false
					continue
				}

				if (ev.Rune() == 'D' || ev.Rune() == 'X' || ev.Key() == tcell.KeyDelete) && !dropMode {
					dropMode = true
					game.Log.Add("Drop mode: Press [a-z] to choose item to drop.", tcell.ColorOrangeRed)
					continue
				}

				// Select slot a-z
				r := ev.Rune()
				if r >= 'a' && r <= 'z' {
					slotIdx := int(r - 'a')
					if dropMode {
						game.DropInventoryItem(slotIdx)
						dropMode = false
					} else {
						game.UseInventoryItem(slotIdx)
					}
					game.State = engine.StatePlaying
				}

			case engine.StatePlaying:
				switch ev.Key() {
				case tcell.KeyUp:
					game.HandlePlayerAction(0, -1)
				case tcell.KeyDown:
					game.HandlePlayerAction(0, 1)
				case tcell.KeyLeft:
					game.HandlePlayerAction(-1, 0)
				case tcell.KeyRight:
					game.HandlePlayerAction(1, 0)
				case tcell.KeyEscape:
					return

				default:
					switch ev.Rune() {
					// Movement: WASD & Vi-keys
					case 'w', 'W', 'k', 'K':
						game.HandlePlayerAction(0, -1)
					case 's', 'S', 'j', 'J':
						game.HandlePlayerAction(0, 1)
					case 'a', 'A', 'h', 'H':
						game.HandlePlayerAction(-1, 0)
					case 'd', 'D', 'l', 'L':
						game.HandlePlayerAction(1, 0)

					// Diagonal Movement (Vi-keys)
					case 'y', 'Y':
						game.HandlePlayerAction(-1, -1)
					case 'u', 'U':
						game.HandlePlayerAction(1, -1)
					case 'b', 'B':
						game.HandlePlayerAction(-1, 1)
					case 'n', 'N':
						game.HandlePlayerAction(1, 1)

					// Actions
					case 'g', ',':
						game.PickUpItem()

					case 'i', 'I':
						game.State = engine.StateInventory

					case '>', '.':
						if ev.Rune() == '>' {
							game.DescendStairs()
						} else {
							game.Log.Add("You wait a turn...", tcell.ColorDarkGray)
							game.EndPlayerTurn()
						}

					case ' ':
						game.Log.Add("You wait a turn...", tcell.ColorDarkGray)
						game.EndPlayerTurn()

					case 'q', 'Q':
						screen.Fini()
						fmt.Println("Thanks for playing!")
						os.Exit(0)
					}
				}
			}
		}
	}
}
