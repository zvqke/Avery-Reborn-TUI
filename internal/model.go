package internal

import (
    "fmt"
    "strings"

    "github.com/charmbracelet/bubbles/textinput"
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
    Todos     []Todo
    TextInput textinput.Model
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
    return textinput.Blink
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "enter": // Add a todo
            if text := strings.TrimSpace(m.TextInput.Value()); text != "" {
                m.Todos = append(m.Todos, Todo{ID: len(m.Todos) + 1, Text: text, Done: false})
                m.TextInput.SetValue("") // Clear the input field
            }
        case "d": // Remove the last todo
            if len(m.Todos) > 0 {
                m.Todos = m.Todos[:len(m.Todos)-1]
            }
        }
    }

    m.TextInput, _ = m.TextInput.Update(msg)
    cmds = append(cmds, m.TextInput.Update(msg))

    return m, tea.Batch(cmds...)
}

// View renders the model.
func (m Model) View() string {
    var s strings.Builder
    s.WriteString("Todos:\n")
    for _, todo := range m.Todos {
        status := "[ ]"
        if todo.Done {
            status = "[x]"
        }
        s.WriteString(fmt.Sprintf("%s %s\n", status, todo.Text))
    }

    s.WriteString("\nPress Enter to add a todo, 'd' to remove the last todo, and Ctrl+C to quit.\n")
    s.WriteString("\nNew Todo:\n")
    s.WriteString(m.TextInput.View())

    return lipgloss.NewStyle().
        Background(lipgloss.Color("#242424")).
        Foreground(lipgloss.Color("#FAFAFA")).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#FFA500")).
        Render(s.String())
}

// NewModel creates a new model instance with initialized text input.
func NewModel(todos []Todo) Model {
    ti := textinput.New()
    ti.Placeholder = "Enter a new todo"
    ti.Focus()
    ti.CharLimit = 156
    ti.Width = 20

    return Model{
        Todos:     todos,
        TextInput: ti,
    }
}
