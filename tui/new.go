package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cmdPkg "github.com/rebelopsio/argo-helper/cmd"
)

type newModel struct {
	resourceTypeInput textinput.Model
	resourceNameInput textinput.Model
	outputPathInput   textinput.Model
	focusIndex        int
	submitted         bool
	err               error
}

func initialNewModel() newModel {
	// Resource Type Input (defaults to applicationset)
	resourceTypeInput := textinput.New()
	resourceTypeInput.Placeholder = "applicationset"
	resourceTypeInput.Focus()
	resourceTypeInput.CharLimit = 30
	resourceTypeInput.Width = 40
	resourceTypeInput.SetValue("applicationset")

	// Resource Name Input
	resourceNameInput := textinput.New()
	resourceNameInput.Placeholder = "Enter resource name"
	resourceNameInput.CharLimit = 30
	resourceNameInput.Width = 40

	// Output Path Input (defaults to templates/apps)
	outputPathInput := textinput.New()
	outputPathInput.Placeholder = "templates/apps"
	outputPathInput.CharLimit = 100
	outputPathInput.Width = 40

	return newModel{
		resourceTypeInput: resourceTypeInput,
		resourceNameInput: resourceNameInput,
		outputPathInput:   outputPathInput,
		focusIndex:        0,
		submitted:         false,
		err:               nil,
	}
}

func (m newModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m newModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					m.focusIndex = 2 // Cycle back to output path input
				}
			} else {
				m.focusIndex++
				if m.focusIndex > 2 {
					m.focusIndex = 0 // Cycle back to resource type input
				}
			}

			// Handle focus changes
			switch m.focusIndex {
			case 0:
				m.resourceTypeInput.Focus()
				m.resourceNameInput.Blur()
				m.outputPathInput.Blur()
			case 1:
				m.resourceTypeInput.Blur()
				m.resourceNameInput.Focus()
				m.outputPathInput.Blur()
			case 2:
				m.resourceTypeInput.Blur()
				m.resourceNameInput.Blur()
				m.outputPathInput.Focus()
			}

			return m, cmd

		case "enter":
			// Validate and submit if all required fields are provided
			resourceType := m.resourceTypeInput.Value()
			if resourceType == "" {
				resourceType = "applicationset" // Default
			}

			if resourceType != "applicationset" {
				m.err = fmt.Errorf("unsupported resource type: %s (only 'applicationset' is currently supported)", resourceType)
				return m, nil
			}

			if m.resourceNameInput.Value() == "" {
				m.err = fmt.Errorf("resource name is required")
				return m, nil
			}

			// Form is complete, run 'new' command
			m.submitted = true

			// Get the values
			resourceName := m.resourceNameInput.Value()
			outputPath := m.outputPathInput.Value()
			if outputPath == "" {
				outputPath = "templates/apps" // Default
			}

			// Make sure output path exists
			if !filepath.IsAbs(outputPath) {
				cwd, err := os.Getwd()
				if err != nil {
					m.err = fmt.Errorf("failed to get current directory: %w", err)
					return m, nil
				}
				outputPath = filepath.Join(cwd, outputPath)
			}

			// Call the CLI handler function
			// We need to use actual command, create a minimal one
			cmdPkg.SetNewFlags(nil, resourceType, resourceName, outputPath)

			// Execute without an actual cobra command
			err := generateNewResource(resourceType, resourceName, outputPath)

			if err != nil {
				m.err = err
				m.submitted = false
				return m, nil
			}

			// Return to the main menu after successful submission
			return NewModel(), nil
		}
	}

	// Handle text input updates
	switch m.focusIndex {
	case 0:
		m.resourceTypeInput, cmd = m.resourceTypeInput.Update(msg)
		return m, cmd
	case 1:
		m.resourceNameInput, cmd = m.resourceNameInput.Update(msg)
		return m, cmd
	case 2:
		m.outputPathInput, cmd = m.outputPathInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// generateNewResource creates a new resource file
// This is a simplified version of the command to avoid dependency issues
func generateNewResource(resourceType, resourceName, outputPath string) error {
	// Validate resource type
	if resourceType != "applicationset" {
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	// If resourceName is not provided, error
	if resourceName == "" {
		return fmt.Errorf("resource name is required")
	}

	// Set default output path if not provided
	if outputPath == "" {
		outputPath = "templates/apps"
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate the content
	content := generateApplicationSetContent(resourceName)

	// Write the file
	filename := fmt.Sprintf("%s-%s.yaml", resourceType, resourceName)
	filePath := filepath.Join(outputPath, filename)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}

	fmt.Printf("Created file: %s\n", filePath)
	return nil
}

// generateApplicationSetContent generates an ApplicationSet template
func generateApplicationSetContent(name string) string {
	return fmt.Sprintf(`apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: %s
  namespace: argocd
spec:
  generators:
    - git:
        repoURL: {{ .Values.global.repoURL }}
        revision: {{ .Values.global.targetRevision }}
        directories:
          - path: apps/*
  template:
    metadata:
      name: '{{ "{{ path.basename }}" }}'
      labels:
        {{- include "common.labels" . | nindent 8 }}
    spec:
      project: {{ include "common.projectName" . }}
      source:
        repoURL: {{ .Values.global.repoURL }}
        targetRevision: {{ .Values.global.targetRevision }}
        path: '{{ "{{ path }}" }}'
      destination:
        server: {{ .Values.destination.server | default "https://kubernetes.default.svc" }}
        namespace: '{{ "{{ path.basename }}" }}'
      syncPolicy:
        {{- toYaml .Values.applications.defaults.syncPolicy | nindent 8 }}
`, name)
}

func (m newModel) View() string {
	if m.submitted {
		return appStyle.Render("Creating resource...")
	}

	var resourceTypeStyle, resourceNameStyle, outputPathStyle lipgloss.Style

	switch m.focusIndex {
	case 0:
		resourceTypeStyle = focusedStyle
		resourceNameStyle = blurredStyle
		outputPathStyle = blurredStyle
	case 1:
		resourceTypeStyle = blurredStyle
		resourceNameStyle = focusedStyle
		outputPathStyle = blurredStyle
	case 2:
		resourceTypeStyle = blurredStyle
		resourceNameStyle = blurredStyle
		outputPathStyle = focusedStyle
	}

	// Error display
	errorText := ""
	if m.err != nil {
		errorText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fb4934")).
			Render(fmt.Sprintf("Error: %v", m.err))
	}

	title := titleStyle.Render("Create New ArgoCD Resource")
	resourceTypeInput := fmt.Sprintf("Resource Type (default: applicationset):\n%s", resourceTypeStyle.Render(m.resourceTypeInput.View()))
	resourceNameInput := fmt.Sprintf("Resource Name (required):\n%s", resourceNameStyle.Render(m.resourceNameInput.View()))
	outputPathInput := fmt.Sprintf("Output Path (default: templates/apps):\n%s", outputPathStyle.Render(m.outputPathInput.View()))

	help := "\nTab/Shift+Tab: Navigate • Enter: Submit • Esc: Cancel"

	return appStyle.Render(
		fmt.Sprintf(
			"%s\n\n%s\n\n%s\n\n%s\n\n%s\n%s",
			title,
			resourceTypeInput,
			resourceNameInput,
			outputPathInput,
			errorText,
			help,
		),
	)
}
