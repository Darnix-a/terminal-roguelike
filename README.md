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
* **Floor 0 Training Grounds (Tutorial):** Sequential 4-chamber tutorial teaching movement, line of sight, closed doors, loot pickup, equipment, and combat.
* **Dedicated Turn-Based Battle Screen (JRPG-Style):** Tactical battles with monster ASCII sprites, dynamic health bars, action menus, and scrolling battle commentary.
* **Magic Skills & Mana System:**
  * `[1] Heavy Slash` (6 MP): 1.6x ATK physical cleave, ignores 50% DEF.
  * `[2] Fireball` (8 MP): Explosive magical fire damage.
  * `[3] Shield Guard` (4 MP): Defensive stance, mitigates next incoming damage by 70%.
  * `[4] Holy Heal` (10 MP): Restores 30 HP.
* **Locked Vaults & Iron Keys (`k`, `%`):** Secret vault chambers requiring dungeon keys to unlock.
* **Treasure Chests & Hungry Mimics (`=`, `M`):** High-tier loot chests with a 15% chance of being an aggressive Mimic monster!
* **Mystic Healing Springs & Power Shrines (`0`, `&`):** Single-use springs that restore 100% HP/MP, and ancient shrines where gold sacrifices grant permanent stat blessings.
* **Elite Champion Monsters with Swarm AI:** 20% of monsters spawn with deadly affixes (`[Fiery]`, `[Vampiric]`, `[Armored]`, `[Frenzied]`) and alert nearby packs when aggroed.
* **Boss Encounter:** Descend to Floor 5 to conquer the Ancient Red Dragon.

---

## Controls

| Key | Action |
| --- | --- |
| `Arrow Keys` / `WASD` / `HJKL` | Move / Interact / Bump |
| `Y` `U` `B` `N` | Diagonal Movement |
| `g` or `,` | Pick up item from ground |
| `i` | Open Inventory (`[a-z]` to use/equip, `[d]` to drop) |
| `>` | Descend stairs to next floor |
| `Space` or `.` | Wait / Rest 1 turn |
| `q` | Quit game |

### In-Battle Commands
| Key | Action |
| --- | --- |
| `1` | Weapon Attack |
| `2` | Cast Magic Skills (`1-4`) |
| `3` | Drink Potion (`a-z`) |
| `4` | Flee / Escape Battle |

---

## Building from Source

```bash
git clone https://github.com/Darnix-a/terminal-roguelike.git
cd terminal-roguelike
make run
```
