package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Monokai color palette
var (
	colorBg     = lipgloss.Color("#272822")
	colorFg     = lipgloss.Color("#F8F8F2")
	colorPink   = lipgloss.Color("#F92672")
	colorGreen  = lipgloss.Color("#A6E22E")
	colorCyan   = lipgloss.Color("#66D9EF")
	colorYellow = lipgloss.Color("#F4BF75")
	colorPurple = lipgloss.Color("#AE81FF")
	colorGray   = lipgloss.Color("#75715E")
)

// Base styles
var (
	appNameStyle = lipgloss.NewStyle().
		Foreground(colorBg).
		Background(colorPink).
		Bold(true).
		Padding(0, 2)

	containerStyle = lipgloss.NewStyle().
		Background(colorBg).
		Foreground(colorFg).
		Padding(1, 2)

	headerStyle = lipgloss.NewStyle().
		Foreground(colorCyan).
		Bold(true)

	selectedStyle = lipgloss.NewStyle().
		Foreground(colorBg).
		Background(colorCyan).
		Bold(true)

	normalStyle = lipgloss.NewStyle().
		Foreground(colorFg)

	helpStyle = lipgloss.NewStyle().
		Foreground(colorGray).
		Italic(true)

	listItemStyle = lipgloss.NewStyle().
		PaddingLeft(1)

	savedStyle = lipgloss.NewStyle().
		Foreground(colorGreen).
		Bold(true)
)

func (m model) View() string {
	// Build content based on current state
	var content string

	if m.state == titleView {
		content = appNameStyle.Render("Goats") + "\n\n"
		content += headerStyle.Render("New Note") + "\n\n"
		content += m.textinput.View() + "\n\n"
		content += m.helpBar()
	}

	if m.state == listView {
		content = appNameStyle.Render("Goats") + "\n\n"
		header := headerStyle.Render(fmt.Sprintf("Notes (%d)", len(m.notes)))
		content += header + "\n\n"

		for i, n := range m.notes {
			if i == m.listIndex {
				content += selectedStyle.Render("> "+n.Title) + "\n"
			} else {
				content += listItemStyle.Render("  "+n.Title) + "\n"
			}
		}

		if m.confirmDelete {
			content += "\n" + lipgloss.NewStyle().Foreground(colorPink).Bold(true).Render("Delete this note? (enter: yes | esc: no)")
		}

		content += "\n" + m.helpBar()
	}

	if m.state == bodyView {
		content = appNameStyle.Render("Goats") + "\n\n"
		titleDisplay := m.currNote.Title
		if titleDisplay == "" {
			titleDisplay = "Untitled"
		}
		content += headerStyle.Render("Editing: "+titleDisplay) + "\n\n"
		content += m.textarea.View() + "\n\n"

		// Show save indicator if saved within last 2 seconds
		if time.Since(m.lastSaved) < 2*time.Second {
			content += savedStyle.Render("Saved!") + "\n"
		}

		content += m.helpBar()
	}

	// Wrap content in full-height container to fill terminal
	fullScreenStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(colorBg).
		Padding(1, 2)

	// Position content at top-left of full-screen container
	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, fullScreenStyle.Render(content))
}

// Help bar function
func (m model) helpBar() string {
	switch m.state {
	case listView:
		if m.confirmDelete {
			return helpStyle.Render("enter: confirm delete | esc: cancel")
		}
		return helpStyle.Render("n: new note | d: delete note | ↑↓: navigate | enter: open | q: quit")
	case titleView:
		return helpStyle.Render("enter: confirm | esc: cancel")
	case bodyView:
		return helpStyle.Render("ctrl+s: save | esc: back")
	}
	return ""
}
