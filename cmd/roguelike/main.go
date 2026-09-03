package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"

	"terminal-roguelike/pkg/combat"
	"terminal-roguelike/pkg/engine"
	"terminal-roguelike/pkg/ui"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tcell screen: %v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing tcell screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	screen.Clear()

	sw, sh := screen.Size()
	mapW := sw - 33
	mapH := sh - 9
	if mapW < 40 {
		mapW = 40
	}
	if mapH < 15 {
		mapH = 15
	}

	game := engine.NewGame(mapW, mapH)
	game.State = engine.StateMainMenu

	for {
		// Render current state
		switch game.State {
		case engine.StateMainMenu:
			ui.RenderMainMenu(screen, engine.HasActiveSave())

		case engine.StateHighScores:
			ui.RenderHighScoresScreen(screen, engine.GetHighScores())

		case engine.StateHelp:
			ui.RenderHelpScreen(screen)

		case engine.StateCombat:
			if game.ActiveBattle != nil {
				ui.RenderBattleScreen(screen, game.ActiveBattle)
			} else {
				game.State = engine.StatePlaying
				ui.Render(screen, game)
			}

		default:
			ui.Render(screen, game)
		}

		// Event polling
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
			sw, sh = screen.Size()
			mapW = sw - 33
			mapH = sh - 9
			if mapW < 40 {
				mapW = 40
			}
			if mapH < 15 {
				mapH = 15
			}
			game.MapW = mapW
			game.MapH = mapH

		case *tcell.EventKey:
			switch game.State {

			// 1. MAIN MENU
			case engine.StateMainMenu:
				switch ev.Key() {
				case tcell.KeyEscape:
					return
				default:
					switch ev.Rune() {
					case '1': // New Game
						game = engine.NewGameCustom(1, mapW, mapH)
						game.State = engine.StatePlaying

					case '2': // Continue Game
						if engine.HasActiveSave() {
							game = engine.LoadGameFromSave(mapW, mapH)
							game.State = engine.StatePlaying
						}

					case '3': // Replay Tutorial
						game = engine.NewGameCustom(0, mapW, mapH)
						game.State = engine.StatePlaying

					case '4': // High Scores
						game.State = engine.StateHighScores

					case '5', '?': // Help
						game.State = engine.StateHelp

					case 'q', 'Q':
						return
					}
				}

			// 2. HIGH SCORES
			case engine.StateHighScores:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyEnter || ev.Rune() == 'q' || ev.Rune() == 'Q' || ev.Rune() == ' ' {
					game.State = engine.StateMainMenu
				}

			// 3. HELP SCREEN
			case engine.StateHelp:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyEnter || ev.Rune() == 'q' || ev.Rune() == 'Q' || ev.Rune() == ' ' {
					game.State = engine.StateMainMenu
				}

			// 4. COMBAT ARENA
			case engine.StateCombat:
				battle := game.ActiveBattle
				if battle == nil {
					game.State = engine.StatePlaying
					continue
				}

				if battle.Result != combat.BattleOngoing {
					game.ConcludeBattle()
					continue
				}

				switch battle.SubMenu {
				case combat.MenuMain:
					switch ev.Rune() {
					case '1':
						battle.PlayerAttack()
					case '2':
						battle.SubMenu = combat.MenuSkills
					case '3':
						battle.PlayerUseSkill(2) // Shield Guard shortcut
					case '4':
						battle.SubMenu = combat.MenuItems
					case '5':
						battle.PlayerFlee()
					}

				case combat.MenuSkills:
					if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
						battle.SubMenu = combat.MenuMain
					} else {
						r := ev.Rune()
						if r >= '1' && r <= '9' {
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

			// 5. LEVEL UP SKILL CHOICE MODAL
			case engine.StateLevelUp:
				r := ev.Rune()
				if r >= '1' && r <= '3' {
					choiceIdx := int(r - '1')
					game.SelectLevelUpSkill(choiceIdx)
				} else if ev.Key() == tcell.KeyEscape {
					game.State = engine.StatePlaying
				}

			// 6. INVENTORY MODAL (Arrow Key / WASD navigation)
			case engine.StateInventory:
				invCount := len(game.Player.Inventory.Items)

				switch ev.Key() {
				case tcell.KeyUp:
					if game.InventoryIdx > 0 {
						game.InventoryIdx--
					}
				case tcell.KeyDown:
					if game.InventoryIdx < invCount-1 {
						game.InventoryIdx++
					}
				case tcell.KeyEnter:
					if invCount > 0 {
						game.UseInventoryItem(game.InventoryIdx)
						game.State = engine.StatePlaying
					}
				case tcell.KeyDelete, tcell.KeyBackspace, tcell.KeyBackspace2:
					if invCount > 0 {
						game.DropInventoryItem(game.InventoryIdx)
						game.State = engine.StatePlaying
					}
				case tcell.KeyEscape:
					game.State = engine.StatePlaying

				default:
					switch ev.Rune() {
					case 'w', 'k', 'W', 'K':
						if game.InventoryIdx > 0 {
							game.InventoryIdx--
						}
					case 's', 'j', 'S', 'J':
						if game.InventoryIdx < invCount-1 {
							game.InventoryIdx++
						}
					case ' ', 'e', 'E':
						if invCount > 0 {
							game.UseInventoryItem(game.InventoryIdx)
							game.State = engine.StatePlaying
						}
					case 'd', 'D', 'x', 'X':
						if invCount > 0 {
							game.DropInventoryItem(game.InventoryIdx)
							game.State = engine.StatePlaying
						}
					case 'i', 'I', 'q', 'Q':
						game.State = engine.StatePlaying
					}
				}

			// 7. DUNGEON SHOP
			case engine.StateShop:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Rune() == 'Q' {
					game.State = engine.StatePlaying
					continue
				}

				if ev.Rune() == 'a' || ev.Rune() == 'A' {
					game.AttackMerchant()
					continue
				}

				r := ev.Rune()
				if r >= '1' && r <= '5' {
					slotIdx := int(r - '0')
					game.BuyShopItem(slotIdx)
				}

			// 8. PLAYING (DUNGEON EXPLORATION)
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
					engine.SaveGameProgress(game)
					game.State = engine.StateMainMenu

				default:
					switch ev.Rune() {
					case 'w', 'k':
						game.HandlePlayerAction(0, -1)
					case 's', 'j':
						game.HandlePlayerAction(0, 1)
					case 'a', 'h':
						game.HandlePlayerAction(-1, 0)
					case 'd', 'l':
						game.HandlePlayerAction(1, 0)
					case 'g', ',':
						game.PickUpItem()
					case 'i':
						game.State = engine.StateInventory
					case '>', '.':
						game.DescendStairs()
					case ' ':
						game.EndPlayerTurn() // Wait 1 turn
					case 'S':
						engine.SaveGameProgress(game)
						game.State = engine.StateMainMenu
					case 'q', 'Q':
						engine.SaveGameProgress(game)
						game.State = engine.StateMainMenu
					}
				}

			// 8. GAME OVER & VICTORY
			case engine.StateGameOver, engine.StateVictory:
				switch ev.Rune() {
				case 'r', 'R':
					game = engine.NewGameCustom(1, mapW, mapH)
					game.State = engine.StatePlaying
				case 'm', 'M':
					game.State = engine.StateMainMenu
				case 'q', 'Q':
					game.State = engine.StateMainMenu
				}
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyEnter {
					game.State = engine.StateMainMenu
				}
			}
		}
	}
}
