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
    Focused   bool
    DeleteID  int
    Listing   bool
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
    return textinput.Blink
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        if m.Listing {
            switch msg.String() {
            case "d":
                if m.DeleteID > 0 && m.DeleteID <= len(m.Todos) {
                    m.Todos = append(m.Todos[:m.DeleteID-1], m.Todos[m.DeleteID:]...)
                    m.DeleteID = 0
                }
                m.Listing = false
            case "esc":
                m.Listing = false
            case "up":
                if m.DeleteID > 1 {
                    m.DeleteID--
                }
            case "down":
                if m.DeleteID < len(m.Todos) {
                    m.DeleteID++
                }
            }
        } else if m.Focused {
            switch msg.String() {
            case "enter":
                if text := strings.TrimSpace(m.TextInput.Value()); text != "" {
                    m.Todos = append(m.Todos, Todo{ID: len(m.Todos) + 1, Text: text, Done: false})
                    m.TextInput.SetValue("") // Clear the input field
                }
                m.Focused = false // Unfocus after adding todo
            default:
                m.TextInput, cmd = m.TextInput.Update(msg)
            }
        } else {
            switch msg.String() {
            case "ctrl+c":
                return m, tea.Quit
            case "a":
                m.Focused = true
                m.TextInput.Focus()
            case "l":
                m.Listing = true
                m.DeleteID = 1
            }
        }
    }

    return m, cmd
}

// View renders the model.
func (m Model) View() string {
    var s strings.Builder
    maxWidth := 50 // Define a max width for todos

    s.WriteString("Todos:\n")
    for i, todo := range m.Todos {
        status := "[ ]"
        if todo.Done {
            status = "[x]"
        }
        gradientStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", (i*15)%256, 0xFF-(i*15)%256, 0x80)))
        todoText := fmt.Sprintf("%s %s", status, todo.Text)
        if m.Listing && i+1 == m.DeleteID {
            gradientStyle = gradientStyle.Bold(true).Background(lipgloss.Color("#FF0000"))
        }
        s.WriteString(gradientStyle.Render(lipgloss.NewStyle().Width(maxWidth).Render(todoText)))
        s.WriteString("\n")
    }

    s.WriteString("\nPress 'a' to add a todo, 'l' to list todos, 'd' to delete the selected todo, and Ctrl+C to quit.\n")
    if m.Focused {
        s.WriteString("\nNew Todo:\n")
        s.WriteString(m.TextInput.View())
    }

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
    ti.CharLimit = 156
    ti.Width = 20

    return Model{
        Todos:     todos,
        TextInput: ti,
        Focused:   false,
        DeleteID:  0,
        Listing:   false,
    }
}
