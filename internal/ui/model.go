package ui

import (
	"log"
	"time"

	"goats/internal/storage"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	state         uint
	store         *storage.Store
	notes         []storage.Note
	currNote      storage.Note
	listIndex     int
	textarea      textarea.Model
	textinput     textinput.Model
	lastSaved     time.Time
	confirmDelete bool
	width         int
	height        int
}

func NewModel(store *storage.Store) tea.Model {
	notes, err := store.GetNotes()
	if err != nil {
		log.Fatalf("Unable to get notes: %v", err)
	}

	ta := textarea.New()
	style := textarea.Style{
		Base: lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorFg),
		CursorLine: lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorFg),
		Placeholder: lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorGray),
		CursorLineNumber: lipgloss.NewStyle().
			Foreground(colorGray),
		EndOfBuffer: lipgloss.NewStyle().
			Foreground(colorGray),
		LineNumber: lipgloss.NewStyle().
			Foreground(colorGray),
		Prompt: lipgloss.NewStyle().
			Foreground(colorGray),
		Text: lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorFg),
	}
	ta.FocusedStyle = style
	ta.BlurredStyle = style

	return model{
		state:     listView,
		store:     store,
		notes:     notes,
		textarea:  ta,
		textinput: textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case listView:
			switch key {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "n":
				if !m.confirmDelete {
					m.textinput.SetValue("")
					m.textinput.Focus()
					m.currNote = storage.Note{}
					m.state = titleView
				}
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}
			case "enter":
				if m.confirmDelete && len(m.notes) > 0 {
					if err := m.store.DeleteNote(m.notes[m.listIndex]); err != nil {
						return m, tea.Quit
					}
					m.notes, _ = m.store.GetNotes()
					m.confirmDelete = false
				} else {
					m.currNote = m.notes[m.listIndex]
					m.textarea.SetValue(m.currNote.Body)
					m.textarea.Focus()
					m.textarea.CursorEnd()
					m.state = bodyView
				}
			case "d":
				if len(m.notes) > 0 && !m.confirmDelete {
					m.confirmDelete = true
				}
			case "esc":
				if m.confirmDelete {
					m.confirmDelete = false
				}
			}
		case titleView:
			switch key {
			case "enter":
				title := m.textinput.Value()
				if title != "" {
					m.currNote.Title = title
				}
				m.state = bodyView
				m.textarea.SetValue("")
				m.textarea.Focus()
				m.textarea.CursorEnd()
			case "esc":
				m.state = listView
			}
		case bodyView:
			switch key {
			case "ctrl+s":
				body := m.textarea.Value()
				m.currNote.Body = body

				var err error
				if err = m.store.SaveNote(m.currNote); err != nil {
					// TODO: Handle Error
					return m, tea.Quit
				}

				m.notes, err = m.store.GetNotes()

				if err != nil {
					// TODO: Handle Error
					return m, tea.Quit
				}

				m.lastSaved = time.Now()
			case "esc":
				m.state = listView
			}
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)
		m.textarea.SetHeight(msg.Height - 8)
	}

	return m, tea.Batch(cmds...)
}
