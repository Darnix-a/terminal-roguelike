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

	// Critical check (15%)
	isCrit := b.RNG.Intn(100) < 15
	dmg := baseDmg + (b.RNG.Intn(3) - 1)
	if isCrit {
		dmg = int(float64(dmg) * 1.5)
		if dmg < 2 {
			dmg = 2
		}
		b.AddLog(fmt.Sprintf("CRITICAL STRIKE! You hit %s for %d damage!", b.Monster.Name, dmg), tcell.ColorYellow)
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
		dmg := int(float64(b.Player.TotalATK()) * 1.8) - (b.Monster.DEF / 2)
		if dmg < 4 {
			dmg = 4
		}
		dmg += b.RNG.Intn(3)
		b.Monster.HP -= dmg
		b.AddLog(fmt.Sprintf("💥 Heavy Slash! You cleave %s for %d physical damage!", b.Monster.Name, dmg), tcell.ColorOrangeRed)

	case "fireball":
		dmg := 16 + (b.Player.Level * 4) + b.RNG.Intn(4)
		b.Monster.HP -= dmg
		b.AddLog(fmt.Sprintf("🔥 Fireball! Incinerating flames blast %s for %d magic damage!", b.Monster.Name, dmg), tcell.ColorOrange)

	case "shield_guard":
		b.Player.Guarding = true
		b.AddLog("🛡️ Shield Stance! You brace yourself, absorbing 75% incoming damage.", tcell.ColorSkyblue)

	case "holy_heal":
		healed := b.Player.Heal(35)
		b.AddLog(fmt.Sprintf("✨ Holy Light! Restored %d HP! (Current HP: %d/%d)", healed, b.Player.HP, b.Player.MaxHP), tcell.ColorGreen)
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
		// Pick action
		picked := b.Monster.Actions[0]
		if len(b.Monster.Actions) > 1 {
			picked = b.Monster.Actions[b.RNG.Intn(len(b.Monster.Actions))]
		}

		// Check if this action should be telegraphed 1 turn ahead
		if picked.IsTelegraphed && b.RNG.Intn(100) < 65 {
			b.TelegraphedAction = &picked
			warningMsg := fmt.Sprintf("⚠️ WARNING: %s %s (Shield Guard next turn!)", b.Monster.Name, picked.TelegraphWarning)
			b.AddLog(warningMsg, tcell.ColorOrangeRed)
			return // Monster spends this turn charging
		}

		action = picked
	}

	// Calculate damage
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

	// Champion Affix Effects
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

	// If player is guarding -> 75% reduction
	if b.Player.Guarding {
		finalDmg = int(float64(finalDmg) * 0.25)
		if finalDmg < 1 {
			finalDmg = 1
		}
		b.Player.Guarding = false
		if isTelegraphedExecution {
			b.AddLog(fmt.Sprintf("🛡️ PERFECT DEFLECTION! You shielded against %s's %s! (Took only %d damage)", b.Monster.Name, action.Name, finalDmg), tcell.ColorGreen)
		} else {
			b.AddLog(fmt.Sprintf("🛡️ Your shield absorbed the blow! %s deals only %d damage with %s.", b.Monster.Name, finalDmg, action.Name), tcell.ColorSkyblue)
		}
	} else {
		if isTelegraphedExecution {
			b.AddLog(fmt.Sprintf("💥 DIRECT HIT! %s unleashes %s! (Heavy %d damage!)", b.Monster.Name, action.Name, finalDmg), tcell.ColorRed)
		} else {
			b.AddLog(fmt.Sprintf("%s uses %s! %s (Deals %d damage)", b.Monster.Name, action.Name, action.Description, finalDmg), tcell.ColorRed)
		}
	}

	// Vampiric Lifesteal
	if b.Monster.IsChampion && b.Monster.Affix == "[Vampiric]" {
		leech := finalDmg / 2
		if leech > 0 {
			b.Monster.HP += leech
			if b.Monster.HP > b.Monster.MaxHP {
				b.Monster.HP = b.Monster.MaxHP
			}
			b.AddLog(fmt.Sprintf("🩸 Vampiric Drain! %s heals +%d HP from your wounds!", b.Monster.Name, leech), tcell.ColorCrimson)
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

		// Level up check
		b.LevelUpMsgs = b.Player.GainEXP(b.Monster.EXP)
		for _, msg := range b.LevelUpMsgs {
			b.AddLog(msg, tcell.ColorAqua)
		}
		return true
	}
	return false
}
