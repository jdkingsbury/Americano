package main

import (
	"log"
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

	// Added logging to debug file if encounter an error
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(panes.NewLayoutModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
