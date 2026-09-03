package engine

import (
	"os"
	"testing"
	"time"
)

func TestTutorialPersistence(t *testing.T) {
	path := getSaveFilePath()
	_ = os.Remove(path)

	if HasCompletedTutorial() {
		t.Errorf("Expected HasCompletedTutorial to be false on clean state")
	}

	MarkTutorialCompleted()

	if !HasCompletedTutorial() {
		t.Errorf("Expected HasCompletedTutorial to be true after MarkTutorialCompleted")
	}

	_ = os.Remove(path)
}

func TestHighScoresLeaderboard(t *testing.T) {
	path := getSaveFilePath()
	_ = os.Remove(path)

	RecordHighScore(HighScoreEntry{
		Score:     250,
		Floor:     2,
		Level:     2,
		Kills:     4,
		Gold:      80,
		Outcome:   "Slain by Goblin Scout on Floor 2",
		Timestamp: time.Now(),
	})

	RecordHighScore(HighScoreEntry{
		Score:     1250,
		Floor:     5,
		Level:     6,
		Kills:     15,
		Gold:      350,
		Outcome:   "Conquered Dungeon (Victory)",
		Timestamp: time.Now(),
	})

	scores := GetHighScores()
	if len(scores) != 2 {
		t.Fatalf("Expected 2 high score entries, got %d", len(scores))
	}

	if scores[0].Score != 1250 {
		t.Errorf("Expected top score to be 1250, got %d", scores[0].Score)
	}

	_ = os.Remove(path)
}
