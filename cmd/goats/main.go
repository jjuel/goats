package main

import (
	"log"

	"goats/internal/storage"
	"goats/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	store := &storage.Store{}
	if err := store.Init(); err != nil {
		log.Fatalf("Unable to init store: %v", err)
	}

	m := ui.NewModel(store)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Unable to run tui: %v", err)
	}
}
