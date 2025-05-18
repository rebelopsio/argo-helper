package test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rebelopsio/argo-helper/tui"
)

func TestNewModel(t *testing.T) {
	model := tui.NewModel()

	// Test that the model is initialized properly
	// This is very basic, but shows the test structure
	if model.View() == "" {
		t.Errorf("View should not be empty")
	}
}

// This is a simple test that checks that the model updates correctly
func TestModelUpdate(t *testing.T) {
	model := tui.NewModel()

	// Simulate a quit key message
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	// Check that the quit command was returned
	if cmd == nil {
		t.Errorf("Expected quit command on Ctrl+C, got nil")
	}

	// We'd typically check more model state here in a real test
}
