// internal/model.go
package internal

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
)

// Import the Todo struct from todo.go
// Todo struct definition
type Todo struct {
    ID   int
    Text string
    Done bool
}

// Model struct
type Model struct {
    Todos []Todo // Capitalize Todos to export it
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
    // You can perform any initialization logic here
    return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages and update the state here
    return m, nil
}

// View renders the TUI using LipGLOSS
func (m Model) View() string {
    // Define how to render your TUI here using LipGLOSS
    return fmt.Sprintf("Todo App\n\nNumber of todos: %d", len(m.Todos))
}
