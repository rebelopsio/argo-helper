package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b8bb26")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#b8bb26")).
			Padding(1, 2)

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d5c4a1")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#a89984")).
			Padding(1, 2)

	checkboxChecked   = "✓"
	checkboxUnchecked = "□"
)

type initModel struct {
	projectInput textinput.Model
	pathInput    textinput.Model
	withExamples bool
	focusIndex   int
	submitted    bool
	err          error
}

func initialInitModel() initModel {
	projectInput := textinput.New()
	projectInput.Placeholder = "Enter project name"
	projectInput.Focus()
	projectInput.CharLimit = 30
	projectInput.Width = 40

	pathInput := textinput.New()
	cwd, err := os.Getwd()
	if err == nil {
		pathInput.Placeholder = cwd
	} else {
		pathInput.Placeholder = "Enter repository path"
	}
	pathInput.CharLimit = 100
	pathInput.Width = 40

	return initModel{
		projectInput: projectInput,
		pathInput:    pathInput,
		withExamples: false,
		focusIndex:   0,
		submitted:    false,
		err:          nil,
	}
}

func (m initModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return NewModel(), nil
		case "tab", "shift+tab", "up", "down":
			// Handle input focus
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = 2 // Cycle back to examples toggle
				}
			} else {
				m.focusIndex++
				if m.focusIndex > 2 {
					m.focusIndex = 0 // Cycle back to project input
				}
			}

			// Handle focus changes
			if m.focusIndex == 0 {
				m.projectInput.Focus()
				m.pathInput.Blur()
			} else if m.focusIndex == 1 {
				m.projectInput.Blur()
				m.pathInput.Focus()
			} else {
				m.projectInput.Blur()
				m.pathInput.Blur()
			}

			return m, cmd

		case "enter":
			if m.focusIndex == 2 {
				// Toggle examples
				m.withExamples = !m.withExamples
				return m, nil
			}

			// Validate and submit if project name is provided
			if m.projectInput.Value() == "" {
				m.err = fmt.Errorf("Project name is required")
				return m, nil
			}

			// Form is complete, run init command
			m.submitted = true

			// Build args for the command
			path := m.pathInput.Value()
			if path == "" {
				var err error
				path, err = os.Getwd()
				if err != nil {
					m.err = fmt.Errorf("Failed to get current directory: %w", err)
					return m, nil
				}
			} else {
				// Convert relative path to absolute if needed
				if !filepath.IsAbs(path) {
					cwd, err := os.Getwd()
					if err != nil {
						m.err = fmt.Errorf("Failed to get current directory: %w", err)
						return m, nil
					}
					path = filepath.Join(cwd, path)
				}
			}

			// We'll implement this properly later
			_ = path

			// Instead of using the actual commands directly,
			// for now just simulate success
			// TODO: Properly integrate with cmd package

			// Simulate success
			fmt.Printf("Creating repository structure at %s\n", path)
			err := fmt.Errorf("") // No error
			if err != nil {
				m.err = err
				m.submitted = false
				return m, nil
			}

			// Return to the main menu after successful submission
			return NewModel(), nil

		case "space":
			if m.focusIndex == 2 {
				// Toggle examples
				m.withExamples = !m.withExamples
			}
		}
	}

	// Handle text input updates
	if m.focusIndex == 0 {
		m.projectInput, cmd = m.projectInput.Update(msg)
		return m, cmd
	} else if m.focusIndex == 1 {
		m.pathInput, cmd = m.pathInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m initModel) View() string {
	if m.submitted {
		return appStyle.Render("Initializing repository...")
	}

	var projectInputStyle, pathInputStyle lipgloss.Style
	if m.focusIndex == 0 {
		projectInputStyle = focusedStyle
		pathInputStyle = blurredStyle
	} else if m.focusIndex == 1 {
		projectInputStyle = blurredStyle
		pathInputStyle = focusedStyle
	} else {
		projectInputStyle = blurredStyle
		pathInputStyle = blurredStyle
	}

	// Checkbox styling for examples
	checkboxStyle := blurredStyle
	checkbox := checkboxUnchecked
	if m.withExamples {
		checkbox = checkboxChecked
	}
	if m.focusIndex == 2 {
		checkboxStyle = focusedStyle
	}

	// Error display
	errorText := ""
	if m.err != nil {
		errorText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fb4934")).
			Render(fmt.Sprintf("Error: %v", m.err))
	}

	title := titleStyle.Render("Initialize ArgoCD Repository")
	projectInput := fmt.Sprintf("Project Name (required):\n%s", projectInputStyle.Render(m.projectInput.View()))
	pathInput := fmt.Sprintf("Repository Path (optional):\n%s", pathInputStyle.Render(m.pathInput.View()))
	examplesToggle := fmt.Sprintf("Include examples:\n%s", checkboxStyle.Render(checkbox))

	help := "\nTab/Shift+Tab: Navigate • Space: Toggle checkbox • Enter: Submit • Esc: Cancel"

	return appStyle.Render(
		fmt.Sprintf(
			"%s\n\n%s\n\n%s\n\n%s\n\n%s\n%s",
			title,
			projectInput,
			pathInput,
			examplesToggle,
			errorText,
			help,
		),
	)
}
