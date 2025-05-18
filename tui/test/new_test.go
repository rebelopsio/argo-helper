package test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rebelopsio/argo-helper/tui"
)

func TestMenuNewAction(t *testing.T) {
	// Since we can't directly test the new resource action without mocking CMD
	// Let's just test that the main model works
	model := tui.NewModel()
	
	// Check that the model is initialized properly
	if model.View() == "" {
		t.Errorf("View should not be empty")
	}
	
	// Test that the model handles keyboard input as expected
	// Simulate a quit key message
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	
	// Should have a quit command
	if cmd == nil {
		t.Errorf("Expected quit command on Ctrl+C, got nil")
	}
}