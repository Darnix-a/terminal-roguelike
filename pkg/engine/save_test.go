package engine

import (
	"os"
	"testing"
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
