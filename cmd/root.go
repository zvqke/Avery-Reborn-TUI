package cmd

import (
	"fmt"
	"os"

	"Avery-Reborn-TUI/internal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todoapp",
	Short: "A CLI todo app with TUI",
	Long:  "A command-line todo application with a text user interface (TUI) built with LipGLOSS and Bubble Tea.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Todo App...")

		// Initialize your todos slice
		todos := []internal.Todo{
			{ID: 1, Text: "Task 1", Done: false},
			{ID: 2, Text: "Task 2", Done: true},
			// Add more todos as needed
		}

		// Initialize your Bubble Tea model
		model := internal.Model{
			Todos: todos,
		}

		// Start the Bubble Tea program
		p := tea.NewProgram(model)
		if err := p.Start(); err != nil {
			fmt.Println("Error starting Todo App:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
