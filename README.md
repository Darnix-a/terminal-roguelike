# terminal-roguelike

A turn-based ASCII/Unicode Dungeon Crawler RPG written in Go with `tcell`.

## Quick Install

### Linux & macOS
```bash
curl -fsSL https://raw.githubusercontent.com/Darnix-a/terminal-roguelike/main/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/Darnix-a/terminal-roguelike/main/install.ps1 | iex
```

---

## Features

* **Main Menu, Saves & High Scores:** Title menu with local Top 10 leaderboard, run resume support, and tutorial replay.
* **Branching MST Dungeon Generation:** Procedural layouts with rooms, corridor loops, and 1-tile choke points.
* **Dedicated Turn-Based Battle Screen:** Tactical JRPG-style combat with monster ASCII art, dynamic health/mana bars, and action logs.
* **Attack Telegraphs & Guarding:** Heavy monster attacks are telegraphed 1 turn in advance. Shield Guard absorbs 75% incoming damage.
* **Hero Leveling & Skill Choices:** Gain EXP from defeated enemies to level up. Choose 1 of 3 random skills at each level up from a pool of 12 distinct abilities (max 5 active abilities, with option to skip or replace existing skills).
* **Dungeon Merchant Shop (`S`):** Spawns on floors 2 and 4 to buy potions, enchantment scrolls, and shields. Attacking the merchant triggers an Enraged Shopkeeper boss fight.
* **Locked Vaults & Keys (`k`, `%`):** Secret vault chambers requiring dungeon keys to unlock.
* **Interactive Arrow-Key Inventory:** Navigate inventory with arrow keys/WASD, equip/unequip gear with Enter, and drop items with `d` or `x`.
* **Treasure Chests & Mimics (`=`, `M`):** High-tier loot chests with a chance to encounter aggressive Mimics.
* **Healing Springs & Power Shrines (`0`, `&`):** Full recovery springs and shrines for permanent stat upgrades.
* **Elite Champions & Swarm AI:** Monsters can spawn with Champion affixes (`[Fiery]`, `[Vampiric]`, `[Armored]`, `[Frenzied]`).
* **Dragon Boss:** Descend to Floor 5 to defeat the Ancient Red Dragon.
* **Floor 0 Training Grounds:** Handcrafted tutorial covering movement, line of sight, loot pickup, equipment, and combat.

---

## Controls

### Exploration
| Key | Action |
| --- | --- |
| `Arrow Keys` / `WASD` / `HJKL` | Move / Interact / Bump Attack |
| `Y` `U` `B` `N` | Diagonal Movement |
| `g` or `,` | Pick up item |
| `i` | Open Inventory |
| `>` or `.` | Descend stairs |
| `Space` | Rest / Wait 1 turn |
| `Shift+S` / `Esc` | Save and quit to menu |
| `q` | Quit game |

### Inventory
| Key | Action |
| --- | --- |
| `Up` / `Down` (or `w` / `s`) | Move selection cursor |
| `Enter` / `Space` / `e` | Use / Equip / Unequip item |
| `d` / `x` / `Delete` | Drop item |
| `Esc` / `i` / `q` | Close inventory |

### Battle
| Key | Action |
| --- | --- |
| `1` | Weapon Attack |
| `2` | Cast Magic Skills (`1-9`) |
| `3` | Shield Guard (absorbs 75% damage) |
| `4` | Drink Potion |
| `5` | Flee / Escape Battle |

### Level-Up Screen
| Key | Action |
| --- | --- |
| `1` - `3` | Choose and learn skill (or select slot to replace) |
| `0` / `Esc` | Skip / Keep current skills |

---

## Building from Source

```bash
git clone https://github.com/Darnix-a/terminal-roguelike.git
cd terminal-roguelike
go build -o roguelike ./cmd/roguelike
./roguelike
# or: make run
```
