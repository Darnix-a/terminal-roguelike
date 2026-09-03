# 🏰 Terminal Roguelike (`terminal-roguelike`)

An atmospheric, feature-packed turn-based ASCII/Unicode Dungeon Crawler RPG written in Go using `tcell`.

```text
  ██████╗  ██████╗  ██████╗ ██╗   ██╗███████╗██╗     ██╗██╗  ██╗███████╗
  ██╔══██╗██╔═══██╗██╔════╝ ██║   ██║██╔════╝██║     ██║██║ ██╔╝██╔════╝
  ██████╔╝██║   ██║██║  ███╗██║   ██║█████╗  ██║     ██║█████═╝ █████╗  
  ██╔══██╗██║   ██║██║   ██║██║   ██║██╔══╝  ██║     ██║██╔═██╗ ██╔══╝  
  ██║  ██║╚██████╔╝╚██████╔╝╚██████╔╝███████╗███████╗██║██║ ╚██╗███████╗
  ╚═╝  ╚═╝ ╚═════╝  ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚═╝╚═╝  ╚═╝╚══════╝
```

---

## ⚡ 1-Command Fast Install

### Linux & macOS
```bash
curl -fsSL https://raw.githubusercontent.com/Darnix-a/terminal-roguelike/main/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/Darnix-a/terminal-roguelike/main/install.ps1 | iex
```

---

## 🛠️ Build & Run from Source

```bash
git clone https://github.com/Darnix-a/terminal-roguelike.git
cd terminal-roguelike
go build -o roguelike ./cmd/roguelike
./roguelike
# Or simply:
make run
```

---

## 🌟 Key Features

### 1. ⚔️ Tactical Turn-Based Combat & Attack Telegraphs
* **Dedicated JRPG-Style Combat Arena:** ASCII monster silhouettes, live animated HP/MP bars, multi-menu skill dispatching, and battle action logs.
* **Heavy Attack Telegraphs:** Deadly boss and monster attacks (*Skull Crusher*, *Hellfire Hex*, *Cataclysm Roar*, *Coin Gatling*) are telegraphed 1 turn in advance with flashing warning banners.
* **Tactical Shield Guarding:** Activating **Shield Guard (`[3]`)** absorbs **75% of incoming damage**, rewarding smart tactical anticipation.

### 2. 🌟 Hero Progression & Dynamic Skill Tree Unlocking
* **Live EXP & Leveling:** Defeating dungeon denizens fills your live sidebar **EXP bar**. Leveling up restores 100% of your HP/MP and permanently boosts Max HP, Max MP, Base ATK, and Base DEF.
* **Choose 1 of 3 Abilities on Every Level Up:** Whenever you reach a new level, a Level-Up Selection Modal presents **3 randomly drawn unlearned abilities** from the master skill pool:
  * 💥 **Heavy Slash (5 MP):** Piercing cleave dealing 1.8x ATK (bypasses 50% DEF).
  * 🛡️ **Shield Guard (3 MP):** Bracing defensive stance absorbing 75% damage.
  * 🔥 **Fireball (7 MP):** Exploding fiery orb dealing high scaling magic damage (`16 + Level*4`).
  * ✨ **Holy Heal (8 MP):** Sacred restorative prayer recovering +35 HP.
  * ⚡ **Chain Lightning (10 MP):** Piercing lightning blast that **completely ignores enemy DEF**.
  * ❄️ **Frost Nova (6 MP):** Glacial blast dealing magic damage and **permanently weakening enemy ATK by -3**.
  * 🩸 **Vampiric Strike (7 MP):** Vicious slash that **converts 50% of damage dealt into player HP**.
  * 🩸 **Berserker Rage (4 MP):** Sacrifices 6 HP to unleash a colossal 2.4x physical strike.
  * 🧪 **Poison Blade (5 MP):** Infuses your weapon with toxic venom dealing +12 poison burst damage.
  * 💨 **Smoke Bomb (4 MP):** Blinds the enemy (negating their turn) and grants a **100% Guaranteed Critical Strike** on your next hit!
  * 👑 **Divine Smite (12 MP):** Celestial wrath smiting the enemy for massive 2.8x holy damage.
  * ✨ **Mana Surge (0 MP):** Sacrifices 8 HP to instantly restore +20 MP.

### 3. 🌲 Procedural Minimum Spanning Tree Dungeons
* **Branching Multi-Chamber Layouts:** Dungeons generate via spatial Minimum Spanning Trees with circular loops and dead-end side chambers.
* **Strict 1-Tile Choke Point Doorways:** Hallways remain clean and open, with at most 1–2 wooden doors placed naturally on actual choke points per floor.
* **Locked Vault Chambers (`%`, `k`):** Procedural side branches sealed by locked doors holding high-tier treasure chests (`=`) and ancient shrines.
* **Dungeon Merchant Shop (`S`):** Peaceful shopkeeper appearing on Floors 2 and 4 selling potions, scrolls, and shields. (Beware: provoking or attacking the merchant triggers a furious **Enraged Shopkeeper** boss fight!).
* **Floor 5 Final Boss:** Venture deep to face the legendary **Ancient Red Dragon (`D`)** (140 HP, 20 ATK, 7 DEF).

### 4. 🎒 Interactive Arrow-Key / WASD Inventory
* **Selection Cursor:** Intuitive cursor navigation (`▲/▼` or `WASD`) with highlighted slots.
* **Automatic Equipment Swapping:** Equipping a new weapon or armor automatically unequips previous gear without inventory lock bugs.
* **Drop Mode:** Drop unwanted gear straight to the floor with `[d]`, `[x]`, or `[Delete]`.

### 5. 💾 Main Menu, Saves & Hall of Fame Leaderboard
* **ASCII Main Menu:** Start New Runs, Continue Saved Runs, Replay Tutorial, or view High Scores.
* **Persistent Saves:** Mid-run state auto-saves on floor descent and can be saved anytime with `[Shift+S]`.
* **Top 10 High Scores Leaderboard:** Tracks your score, floor reached, hero level, kills, gold, and death/victory outcome.

---

## 🎮 Complete Controls Cheat Sheet

### 🗺️ Exploration Controls
| Key | Action |
| :--- | :--- |
| `▲ / ▼ / ◄ / ►` or `W / A / S / D` | Move Hero / Interact / Bump Attack |
| `H / J / K / L` | Classic Roguelike Vi-Keys |
| `Y / U / B / N` | Diagonal Movement |
| `g` or `,` | Pick up item on current tile |
| `i` | Open Backpack / Inventory Modal |
| `>` or `.` | Descend Stairs Down (`>`) |
| `Space` | Rest / Wait 1 turn (regens HP/MP) |
| `Shift + S` / `Esc` | Save Game & Return to Main Menu |
| `q` | Quit Run |

### 🎒 Inventory Modal Controls
| Key | Action |
| :--- | :--- |
| `▲ / ▼` or `W / S` (or `j / k`) | Move selection cursor |
| `Enter` / `Space` / `e` | **Use, Equip, or Unequip** highlighted item |
| `d` / `x` / `Delete` | **Drop** highlighted item to floor |
| `Esc` / `i` / `q` | Close Inventory |

### ⚔️ Battle Arena Controls
| Key | Action |
| :--- | :--- |
| `[1]` | Weapon Strike |
| `[2]` | Open Magic Skills Menu (Cast with `[1-9]`) |
| `[3]` | Shield Guard stance (Absorbs 75% incoming damage) |
| `[4]` | Use Healing Potion |
| `[5]` | Flee / Escape Battle (75% success) |

### 🏪 Dungeon Shop Controls
| Key | Action |
| :--- | :--- |
| `[1-5]` | Buy listed item / potion / scroll / shield |
| `[A]` | Provoke & Attack the Merchant (Triggers Boss Fight!) |
| `Esc` / `q` | Leave Shop |

---

## 🗺️ Map Legend

```text
 @  Hero (You)             .  Open Floor             #  Solid Stone Wall
 +  Wooden Door            %  Locked Vault Door      k  Iron Dungeon Key
 =  Treasure Chest         M  Hungry Mimic           0  Mystic Healing Spring
 &  Shrine of Power        S  Dungeon Merchant       >  Stairs Down
 $  Gold Pile              )  Weapon                 [  Armor
 !  Potion                 ?  Magical Scroll         g  Goblin Scout
 s  Skeletal Guard         o  Orc Berserker          w  Dark Sorcerer
 D  Ancient Red Dragon
```

---

## 📜 License

MIT License © 2026 Darnix-a.
