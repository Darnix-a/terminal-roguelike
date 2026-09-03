package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SaveData struct {
	TutorialCompleted bool `json:"tutorial_completed"`
}

func getSaveFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".terminal_roguelike.json")
}

// HasCompletedTutorial checks if tutorial was finished in a previous run
func HasCompletedTutorial() bool {
	path := getSaveFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var save SaveData
	if err := json.Unmarshal(data, &save); err != nil {
		return false
	}
	return save.TutorialCompleted
}

// MarkTutorialCompleted saves the flag so tutorial is skipped on future runs
func MarkTutorialCompleted() {
	path := getSaveFilePath()
	save := SaveData{TutorialCompleted: true}
	data, err := json.MarshalIndent(save, "", "  ")
	if err == nil {
		_ = os.WriteFile(path, data, 0644)
	}
}
