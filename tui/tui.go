package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3c3836")).
		Width(60)

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#83a598")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)
)

type menuItem struct {
	title       string
	description string
	action      func() (tea.Model, tea.Cmd)
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

// Model represents the main TUI model
type Model struct {
	List list.Model // Exported for testing
}

// Expose the Model type for testing
var _ = Model{} // This ensures the type is exported

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles updates to the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			selectedItem, ok := m.List.SelectedItem().(menuItem)
			if ok && selectedItem.action != nil {
				return selectedItem.action()
			}
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

// View renders the model
func (m Model) View() string {
	title := titleStyle.Render("ArgoCD Helper")
	content := m.List.View()
	return appStyle.Render(fmt.Sprintf("%s\n%s", title, content))
}

// NewModel creates a new TUI model
func NewModel() Model {
	items := []list.Item{
		menuItem{
			title:       "Initialize Repository",
			description: "Create a new ArgoCD repository structure",
			action:      menuInitAction,
		},
		menuItem{
			title:       "New Resource",
			description: "Create a new ArgoCD resource (ApplicationSet, etc.)",
			action:      menuNewAction,
		},
		menuItem{
			title:       "Quit",
			description: "Exit the application",
			action: func() (tea.Model, tea.Cmd) {
				return Model{}, tea.Quit
			},
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Choose an option:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	return Model{List: l}
}

// Run starts the TUI application
func Run() error {
	p := tea.NewProgram(NewModel())
	_, err := p.Run()
	return err
}

// menuInitAction creates the init form model
func menuInitAction() (tea.Model, tea.Cmd) {
	return initialInitModel(), textinput.Blink
}

// menuNewAction creates the new resource form model
func menuNewAction() (tea.Model, tea.Cmd) {
	return initialNewModel(), textinput.Blink
}

// ExportedMenuNewAction is a wrapper for tests to access menuNewAction
func ExportedMenuNewAction() (tea.Model, tea.Cmd) {
	return menuNewAction()
}