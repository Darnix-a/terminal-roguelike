package combat

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"

	"terminal-roguelike/pkg/entities"
	"terminal-roguelike/pkg/items"
)

type BattleSubMenu int

const (
	MenuMain BattleSubMenu = iota
	MenuSkills
	MenuItems
)

type BattleResult int

const (
	BattleOngoing BattleResult = iota
	BattleVictory
	BattleDefeat
	BattleFled
)

type BattleMessage struct {
	Text  string
	Color tcell.Color
}

type Battle struct {
	Player            *entities.Player
	Monster           *entities.Monster
	RNG               *rand.Rand
	SubMenu           BattleSubMenu
	Result            BattleResult
	Log               []BattleMessage
	LevelUpMsgs       []string
	LeveledUp         bool
	SelectedSlot      int
	TelegraphedAction *entities.MonsterAction
}

func NewBattle(player *entities.Player, monster *entities.Monster, rng *rand.Rand) *Battle {
	b := &Battle{
		Player:            player,
		Monster:           monster,
		RNG:               rng,
		SubMenu:           MenuMain,
		Result:            BattleOngoing,
		Log:               make([]BattleMessage, 0),
		LevelUpMsgs:       make([]string, 0),
		LeveledUp:         false,
		SelectedSlot:      0,
		TelegraphedAction: nil,
	}

	b.AddLog(fmt.Sprintf("⚔️ Engaged in combat with %s!", monster.Name), tcell.ColorGold)
	return b
}

func (b *Battle) AddLog(text string, color tcell.Color) {
	b.Log = append(b.Log, BattleMessage{Text: text, Color: color})
	if len(b.Log) > 30 {
		b.Log = b.Log[len(b.Log)-30:]
	}
}

// PlayerAttack executes standard weapon strike
func (b *Battle) PlayerAttack() {
	if b.Result != BattleOngoing {
		return
	}

	baseDmg := b.Player.TotalATK() - b.Monster.DEF
	if baseDmg < 1 {
		baseDmg = 1
	}

	// Critical check (15% base, or 100% if GuaranteedCrit)
	isCrit := b.RNG.Intn(100) < 15 || b.Player.GuaranteedCrit
	if b.Player.GuaranteedCrit {
		b.Player.GuaranteedCrit = false
	}

	dmg := baseDmg + (b.RNG.Intn(3) - 1)
	if isCrit {
		dmg = int(float64(dmg) * 1.6)
		if dmg < 3 {
			dmg = 3
		}
		b.AddLog(fmt.Sprintf("⚡ CRITICAL STRIKE! You strike %s for %d damage!", b.Monster.Name, dmg), tcell.ColorYellow)
	} else {
		b.AddLog(fmt.Sprintf("You strike %s for %d damage.", b.Monster.Name, dmg), tcell.ColorWhite)
	}

	b.Monster.HP -= dmg
	if b.checkMonsterDead() {
		return
	}

	b.EnemyTurn()
}

// PlayerUseSkill casts a special ability
func (b *Battle) PlayerUseSkill(skillIdx int) {
	if b.Result != BattleOngoing || skillIdx < 0 || skillIdx >= len(b.Player.Skills) {
		return
	}

	skill := b.Player.Skills[skillIdx]
	if b.Player.MP < skill.MPCost {
		b.AddLog(fmt.Sprintf("Not enough MP to cast %s! (Requires %d MP)", skill.Name, skill.MPCost), tcell.ColorRed)
		return
	}

	b.Player.MP -= skill.MPCost

	switch skill.ID {
	case "heavy_slash":
		dmg := int(float64(b.Player.TotalATK()) * 1.8) - (b.Monster.DEF / 2) + b.RNG.Intn(3)
		if dmg < 4 {
			dmg = 4
		}
		b.Monster.HP -= dmg
		b.AddLog(fmt.Sprintf("💥 Heavy Slash! You cleave %s for %d physical damage!", b.Monster.Name, dmg), tcell.ColorOrangeRed)

	case "shield_guard":
		b.Player.Guarding = true
		b.AddLog("🛡️ Shield Stance! You brace yourself, absorbing 75% incoming damage.", tcell.ColorSkyblue)

	case "fireball":
		dmg := 16 + (b.Player.Level * 4) + b.RNG.Intn(4)
		b.Monster.HP -= dmg
		b.AddLog(fmt.Sprintf("🔥 Fireball! Incinerating flames blast %s for %d magic damage!", b.Monster.Name, dmg), tcell.ColorOrange)

	case "holy_heal":
		healed := b.Player.Heal(35)
		b.AddLog(fmt.Sprintf("✨ Holy Light! Restored %d HP! (Current HP: %d/%d)", healed, b.Player.HP, b.Player.MaxHP), tcell.ColorGreen)

	case "chain_lightning":
		dmg := 18 + (b.Player.Level * 4) + b.RNG.Intn(5)
		b.Monster.HP -= dmg
		b.AddLog(fmt.Sprintf("⚡ Chain Lightning! Piercing electric bolts shock %s for %d TRUE damage (ignores DEF)!", b.Monster.Name, dmg), tcell.ColorAqua)

	case "frost_nova":
		dmg := 12 + (b.Player.Level * 3)
		b.Monster.HP -= dmg
		if b.Monster.ATK > 3 {
			b.Monster.ATK -= 3
		}
		b.AddLog(fmt.Sprintf("❄️ Frost Nova! Glacial ice chills %s for %d magic damage and weakens their ATK by -3!", b.Monster.Name, dmg), tcell.ColorSkyblue)

	case "vampiric_strike":
		dmg := int(float64(b.Player.TotalATK()) * 1.4) - b.Monster.DEF + b.RNG.Intn(3)
		if dmg < 3 {
			dmg = 3
		}
		b.Monster.HP -= dmg
		leech := dmg / 2
		if leech > 0 {
			b.Player.Heal(leech)
		}
		b.AddLog(fmt.Sprintf("🩸 Vampiric Strike! Slashed %s for %d damage and drained +%d HP!", b.Monster.Name, dmg, leech), tcell.ColorCrimson)

	case "berserker_rage":
		b.Player.HP -= 6
		if b.Player.HP < 1 {
			b.Player.HP = 1
		}
		dmg := int(float64(b.Player.TotalATK()) * 2.4) - b.Monster.DEF + b.RNG.Intn(4)
		if dmg < 6 {
			dmg = 6
		}
		b.Monster.HP -= dmg
		b.AddLog(fmt.Sprintf("🩸 Berserker Rage! Sacrificed 6 HP to smash %s for a massive %d damage!", b.Monster.Name, dmg), tcell.ColorRed)

	case "poison_blade":
		dmg := int(float64(b.Player.TotalATK()) * 1.2) - b.Monster.DEF + 12
		if dmg < 6 {
			dmg = 6
		}
		b.Monster.HP -= dmg
		b.AddLog(fmt.Sprintf("🧪 Poison Blade! Struck %s with deadly toxin for %d venom damage!", b.Monster.Name, dmg), tcell.ColorGreen)

	case "smoke_bomb":
		b.Player.Guarding = true
		b.Player.GuaranteedCrit = true
		b.AddLog("💨 Smoke Bomb! Blinded enemy and primed next attack for GUARANTEED CRITICAL STRIKE!", tcell.ColorViolet)

	case "divine_smite":
		dmg := int(float64(b.Player.TotalATK()) * 2.8) - (b.Monster.DEF / 3) + b.RNG.Intn(5)
		if dmg < 10 {
			dmg = 10
		}
		b.Monster.HP -= dmg
		b.AddLog(fmt.Sprintf("👑 DIVINE SMITE! Celestial light smites %s for a devastating %d holy damage!", b.Monster.Name, dmg), tcell.ColorGold)

	case "mana_surge":
		b.Player.HP -= 8
		if b.Player.HP < 1 {
			b.Player.HP = 1
		}
		b.Player.RestoreMP(20)
		b.AddLog(fmt.Sprintf("✨ Mana Surge! Sacrificed 8 HP to restore +20 MP! (MP: %d/%d)", b.Player.MP, b.Player.MaxMP), tcell.ColorAqua)
	}

	b.SubMenu = MenuMain

	if b.checkMonsterDead() {
		return
	}

	b.EnemyTurn()
}

// PlayerUseItem uses item from inventory during battle
func (b *Battle) PlayerUseItem(itemIdx int) {
	if b.Result != BattleOngoing || itemIdx < 0 || itemIdx >= len(b.Player.Inventory.Items) {
		return
	}

	item := b.Player.Inventory.Items[itemIdx]
	if item.Type == items.TypePotion {
		healed := b.Player.Heal(item.HealAmount)
		b.AddLog(fmt.Sprintf("🧪 Drank %s! Restored %d HP! (Current HP: %d/%d)", item.Name, healed, b.Player.HP, b.Player.MaxHP), tcell.ColorGreen)
		_, _ = b.Player.Inventory.Remove(itemIdx)
		b.SubMenu = MenuMain
		b.EnemyTurn()
	} else {
		b.AddLog("You can only use consumable potions in battle!", tcell.ColorDarkGray)
	}
}

// PlayerFlee attempts escape from battle
func (b *Battle) PlayerFlee() {
	if b.Result != BattleOngoing {
		return
	}

	if b.Monster.IsBoss {
		b.AddLog("There is no escaping the Dragon's lair!", tcell.ColorRed)
		return
	}

	if b.Monster.IsMerchant {
		b.AddLog("The Shopkeeper has locked the doors! Fight or perish!", tcell.ColorRed)
		return
	}

	// 75% escape chance
	if b.RNG.Intn(100) < 75 {
		b.AddLog("💨 You successfully fled from battle!", tcell.ColorYellow)
		b.Result = BattleFled
		return
	}

	b.AddLog("Failed to escape! The enemy blocks your retreat.", tcell.ColorDarkGray)
	b.EnemyTurn()
}

// EnemyTurn processes tactical monster turn with attack telegraphs
func (b *Battle) EnemyTurn() {
	if b.Monster.HP <= 0 {
		return
	}

	var action entities.MonsterAction
	isTelegraphedExecution := false

	// Check if executing a previously telegraphed heavy attack
	if b.TelegraphedAction != nil {
		action = *b.TelegraphedAction
		b.TelegraphedAction = nil
		isTelegraphedExecution = true
	} else {
		picked := b.Monster.Actions[0]
		if len(b.Monster.Actions) > 1 {
			picked = b.Monster.Actions[b.RNG.Intn(len(b.Monster.Actions))]
		}

		if picked.IsTelegraphed && b.RNG.Intn(100) < 65 {
			b.TelegraphedAction = &picked
			warningMsg := fmt.Sprintf("⚠️ WARNING: %s %s (Shield Guard next turn!)", b.Monster.Name, picked.TelegraphWarning)
			b.AddLog(warningMsg, tcell.ColorOrangeRed)
			return // Monster spends turn charging
		}

		action = picked
	}

	rawDmg := float64(b.Monster.ATK) * action.DamageMult
	var finalDmg int

	if action.IsMagic {
		finalDmg = int(rawDmg)
	} else {
		finalDmg = int(rawDmg) - b.Player.TotalDEF()
	}

	if finalDmg < 1 {
		finalDmg = 1
	}

	if b.Monster.IsChampion {
		switch b.Monster.Affix {
		case "[Fiery]":
			finalDmg += 4
			b.AddLog("🔥 Fiery Aura! Flames scorch you for +4 bonus damage!", tcell.ColorOrangeRed)
		case "[Frenzied]":
			finalDmg = int(float64(finalDmg) * 1.3)
			b.AddLog("⚡ Frenzied Rush! The champion attacks with blinding ferocity!", tcell.ColorMediumPurple)
		}
	}

	if b.Player.Guarding {
		finalDmg = int(float64(finalDmg) * 0.25)
		if finalDmg < 1 {
			finalDmg = 1
		}
		b.Player.Guarding = false
		if isTelegraphedExecution {
			b.AddLog(fmt.Sprintf("🛡️ PERFECT DEFLECTION! Shielded %s's %s! (Took %d damage)", b.Monster.Name, action.Name, finalDmg), tcell.ColorGreen)
		} else {
			b.AddLog(fmt.Sprintf("🛡️ Shield absorbed blow! %s deals only %d damage.", b.Monster.Name, finalDmg), tcell.ColorSkyblue)
		}
	} else {
		if isTelegraphedExecution {
			b.AddLog(fmt.Sprintf("💥 DIRECT HIT! %s strikes with %s! (%d damage!)", b.Monster.Name, action.Name, finalDmg), tcell.ColorRed)
		} else {
			b.AddLog(fmt.Sprintf("%s uses %s! %s (Deals %d damage)", b.Monster.Name, action.Name, action.Description, finalDmg), tcell.ColorRed)
		}
	}

	if b.Monster.IsChampion && b.Monster.Affix == "[Vampiric]" {
		leech := finalDmg / 2
		if leech > 0 {
			b.Monster.HP += leech
			if b.Monster.HP > b.Monster.MaxHP {
				b.Monster.HP = b.Monster.MaxHP
			}
			b.AddLog(fmt.Sprintf("🩸 Vampiric Drain! %s heals +%d HP!", b.Monster.Name, leech), tcell.ColorCrimson)
		}
	}

	b.Player.HP -= finalDmg
	if b.Player.HP <= 0 {
		b.Player.HP = 0
		b.AddLog("☠️ You have fallen in battle...", tcell.ColorDarkRed)
		b.Result = BattleDefeat
	}
}

func (b *Battle) checkMonsterDead() bool {
	if b.Monster.HP <= 0 {
		b.Monster.HP = 0
		b.Monster.IsAlive = false
		b.Player.Kills++
		b.Result = BattleVictory

		b.AddLog(fmt.Sprintf("🎉 VICTORY! %s was defeated! Gained +%d EXP!", b.Monster.Name, b.Monster.EXP), tcell.ColorGold)

		leveledUp, msgs := b.Player.GainEXP(b.Monster.EXP)
		b.LeveledUp = leveledUp
		b.LevelUpMsgs = msgs
		for _, msg := range b.LevelUpMsgs {
			b.AddLog(msg, tcell.ColorAqua)
		}
		return true
	}
	return false
}
