package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdkingsbury/americano/internal/tui/panes"
)

func main() {
	saveState := exec.Command("tput", "smcup")
	saveState.Stdout = os.Stdout
	saveState.Run()

	defer func() {
		restoreState := exec.Command("tput", "rmcup")
		restoreState.Stdout = os.Stdout
		restoreState.Run()
	}()

	p := tea.NewProgram(panes.NewLayoutModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
