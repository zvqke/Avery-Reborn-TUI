package internal

import (
    "fmt"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Todo represents a single todo item.
type Todo struct {
    ID   int
    Text string
    Done bool
}

// Model represents the Bubble Tea model.
type Model struct {
    Todos []Todo
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
    return nil
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "a": // Add a todo
            m.Todos = append(m.Todos, Todo{ID: len(m.Todos) + 1, Text: "New Task", Done: false})
        case "d": // Remove the last todo
            if len(m.Todos) > 0 {
                m.Todos = m.Todos[:len(m.Todos)-1]
            }
        }
    }
    return m, nil
}

// View renders the model.
func (m Model) View() string {
    var s string
    for _, todo := range m.Todos {
        status := "[ ]"
        if todo.Done {
            status = "[x]"
        }
        s += fmt.Sprintf("%s %s\n", status, todo.Text)
    }
    return lipgloss.NewStyle().
        Background(lipgloss.Color("#242424")).
        Foreground(lipgloss.Color("#FAFAFA")).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#FFA500")).
        Render(s)
}
