package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"terminal-roguelike/pkg/items"
)

type SavedItem struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        items.ItemType `json:"type"`
	Rune        rune           `json:"rune"`
	HealAmount  int            `json:"heal_amount"`
	BonusATK    int            `json:"bonus_atk"`
	BonusDEF    int            `json:"bonus_def"`
	Value       int            `json:"value"`
	Description string         `json:"description"`
	Equipped    bool           `json:"equipped"`
}

type SavedPlayer struct {
	HP        int         `json:"hp"`
	MaxHP     int         `json:"max_hp"`
	MP        int         `json:"mp"`
	MaxMP     int         `json:"max_mp"`
	BaseATK   int         `json:"base_atk"`
	BaseDEF   int         `json:"base_def"`
	Level     int         `json:"level"`
	EXP       int         `json:"exp"`
	MaxEXP    int         `json:"max_exp"`
	Gold      int         `json:"gold"`
	Keys      int         `json:"keys"`
	Kills     int         `json:"kills"`
	Inventory []SavedItem `json:"inventory"`
}

type SavedGameState struct {
	Floor     int         `json:"floor"`
	TurnCount int         `json:"turn_count"`
	Player    SavedPlayer `json:"player"`
}

type HighScoreEntry struct {
	Score     int       `json:"score"`
	Floor     int       `json:"floor"`
	Level     int       `json:"level"`
	Kills     int       `json:"kills"`
	Gold      int       `json:"gold"`
	Outcome   string    `json:"outcome"`
	Timestamp time.Time `json:"timestamp"`
}

type AppData struct {
	TutorialCompleted bool            `json:"tutorial_completed"`
	ActiveSave        *SavedGameState `json:"active_save,omitempty"`
	HighScores        []HighScoreEntry `json:"high_scores"`
}

func getSaveFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".terminal_roguelike.json")
}

func loadAppData() AppData {
	path := getSaveFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return AppData{HighScores: make([]HighScoreEntry, 0)}
	}

	var app AppData
	if err := json.Unmarshal(data, &app); err != nil {
		return AppData{HighScores: make([]HighScoreEntry, 0)}
	}
	if app.HighScores == nil {
		app.HighScores = make([]HighScoreEntry, 0)
	}
	return app
}

func saveAppData(app AppData) {
	path := getSaveFilePath()
	data, err := json.MarshalIndent(app, "", "  ")
	if err == nil {
		_ = os.WriteFile(path, data, 0644)
	}
}

// HasCompletedTutorial checks if tutorial was finished
func HasCompletedTutorial() bool {
	return loadAppData().TutorialCompleted
}

// MarkTutorialCompleted saves the flag
func MarkTutorialCompleted() {
	app := loadAppData()
	app.TutorialCompleted = true
	saveAppData(app)
}

// HasActiveSave checks if a saved game exists to continue
func HasActiveSave() bool {
	return loadAppData().ActiveSave != nil
}

// DeleteActiveSave removes saved game on death or victory
func DeleteActiveSave() {
	app := loadAppData()
	app.ActiveSave = nil
	saveAppData(app)
}

// SaveGameProgress saves current player and floor state
func SaveGameProgress(g *Game) {
	if g.Floor == 0 || g.Player == nil {
		return // Don't save tutorial
	}

	savedInv := make([]SavedItem, len(g.Player.Inventory.Items))
	for i, itm := range g.Player.Inventory.Items {
		savedInv[i] = SavedItem{
			ID:          itm.ID,
			Name:        itm.Name,
			Type:        itm.Type,
			Rune:        itm.Rune,
			HealAmount:  itm.HealAmount,
			BonusATK:    itm.BonusATK,
			BonusDEF:    itm.BonusDEF,
			Value:       itm.Value,
			Description: itm.Description,
			Equipped:    itm.Equipped,
		}
	}

	savedPlayer := SavedPlayer{
		HP:        g.Player.HP,
		MaxHP:     g.Player.MaxHP,
		MP:        g.Player.MP,
		MaxMP:     g.Player.MaxMP,
		BaseATK:   g.Player.BaseATK,
		BaseDEF:   g.Player.BaseDEF,
		Level:     g.Player.Level,
		EXP:       g.Player.EXP,
		MaxEXP:    g.Player.MaxEXP,
		Gold:      g.Player.Gold,
		Keys:      g.Player.Keys,
		Kills:     g.Player.Kills,
		Inventory: savedInv,
	}

	app := loadAppData()
	app.ActiveSave = &SavedGameState{
		Floor:     g.Floor,
		TurnCount: g.TurnCount,
		Player:    savedPlayer,
	}
	saveAppData(app)
}

// RecordHighScore adds a finished run to high score list
func RecordHighScore(entry HighScoreEntry) {
	app := loadAppData()
	app.HighScores = append(app.HighScores, entry)

	// Sort high scores descending
	sort.Slice(app.HighScores, func(i, j int) bool {
		return app.HighScores[i].Score > app.HighScores[j].Score
	})

	// Keep Top 10
	if len(app.HighScores) > 10 {
		app.HighScores = app.HighScores[:10]
	}
	saveAppData(app)
}

// GetHighScores returns the sorted high score list
func GetHighScores() []HighScoreEntry {
	return loadAppData().HighScores
}
