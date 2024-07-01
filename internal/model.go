package internal

import (
    "fmt"
    "strings"
    "time"

    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Todo represents a single todo item.
type Todo struct {
    ID      int
    Text    string
    Done    bool
    DueDate time.Time // New field for due date
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
    return tea.Batch(textinput.Blink, CheckDueTodos())
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
            case "m":
                if m.DeleteID > 0 && m.DeleteID <= len(m.Todos) {
                    m.Todos[m.DeleteID-1].Done = !m.Todos[m.DeleteID-1].Done
                }
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
                    m.Todos = append(m.Todos, Todo{ID: len(m.Todos) + 1, Text: text, Done: false, DueDate: time.Now().Add(24 * time.Hour)}) // Set a default due date
                    m.TextInput.SetValue("")                                                                                                // Clear the input field
                }
                m.Focused = false // Unfocus after adding todo
            case "esc":
                m.Focused = false
                m.TextInput.Blur() // Unfocus the text input
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
    case CheckDueTodosMsg:
        // Check for todos due today and delete them if not marked as done
        for i := 0; i < len(m.Todos); i++ {
            if !m.Todos[i].Done && m.Todos[i].DueDate.Before(time.Now()) {
                m.Todos = append(m.Todos[:i], m.Todos[i+1:]...)
                i-- // adjust index after deletion
            }
        }
        cmd = tea.Tick(time.Minute, func(time.Time) tea.Msg {
            return CheckDueTodosMsg{}
        })
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

    s.WriteString("\nPress 'a' to add a todo, \nPress 'l' to list todos, \nPress 'd' to delete the selected todo, \nPress 'm' to mark/unmark the selected todo, \n\nPress Ctrl+C to quit.\n")
    s.WriteString("Press 'Esc' to go back.")
    if m.Focused {
        s.WriteString("\n\nNew Todo:\n")
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

// CheckDueTodosMsg is a custom message type to trigger due todos check.
type CheckDueTodosMsg struct{}

// CheckDueTodos is a command to check for due todos.
func CheckDueTodos() tea.Cmd {
    return tea.Tick(time.Minute, func(time.Time) tea.Msg {
        return CheckDueTodosMsg{}
    })
}
