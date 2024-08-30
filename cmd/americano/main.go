package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
  "github.com/jdkingsbury/americano/pkg/ui"
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

  p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
