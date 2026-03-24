package ui

import (
	"log"
	"time"
	"unicode"

	"goats/internal/storage"

	"github.com/charmbracelet/bubbles/cursor"
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

const (
	insertMode uint = iota
	normalMode
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
	editorMode    uint
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
		state:      listView,
		store:      store,
		notes:      notes,
		textarea:   ta,
		textinput:  textinput.New(),
		editorMode: normalMode,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case listView:
			return m.updateList(msg)
		case titleView:
			return m.updateTitle(msg)
		case bodyView:
			return m.updateBody(msg)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)
		m.textarea.SetHeight(msg.Height - 8)
	}

	return m, nil
}

func (m model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "n":
		if !m.confirmDelete {
			m.currNote = storage.Note{}
			m.textinput.SetValue("")
			m.textinput.CursorEnd()
			m.state = titleView
			m.enterInsertMode()
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
			if m.listIndex >= len(m.notes) && m.listIndex > 0 {
				m.listIndex--
			}
			m.confirmDelete = false
		} else if len(m.notes) > 0 {
			m.currNote = m.notes[m.listIndex]
			m.textarea.SetValue(m.currNote.Body)
			m.textarea.CursorStart()
			m.state = bodyView
			return m, m.enterNormalMode()
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

	return m, nil
}

func (m model) updateTitle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editorMode == insertMode {
		switch msg.String() {
		case "esc":
			return m, m.enterNormalMode()
		case "enter":
			if m.textinput.Value() != "" {
				m.currNote.Title = m.textinput.Value()
			}
			m.textarea.SetValue("")
			m.textarea.CursorStart()
			m.state = bodyView
			return m, m.enterInsertMode()
		}

		var cmd tea.Cmd
		m.textinput, cmd = m.textinput.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "esc":
		m.state = listView
		return m, m.enterNormalMode()
	case "enter":
		if m.textinput.Value() != "" {
			m.currNote.Title = m.textinput.Value()
		}
		m.textarea.SetValue("")
		m.textarea.CursorStart()
		m.state = bodyView
		return m, m.enterInsertMode()
	case "i":
		return m, m.enterInsertMode()
	case "a":
		pos := m.textinput.Position()
		if pos < len([]rune(m.textinput.Value())) {
			m.textinput.SetCursor(pos + 1)
		}
		return m, m.enterInsertMode()
	case "I":
		m.textinput.CursorStart()
		return m, m.enterInsertMode()
	case "A":
		m.textinput.CursorEnd()
		return m, m.enterInsertMode()
	case "h":
		if pos := m.textinput.Position(); pos > 0 {
			m.textinput.SetCursor(pos - 1)
		}
	case "l":
		pos := m.textinput.Position()
		if pos < len([]rune(m.textinput.Value())) {
			m.textinput.SetCursor(pos + 1)
		}
	case "0":
		m.textinput.CursorStart()
	case "$":
		m.textinput.CursorEnd()
	case "w":
		m.textinput.SetCursor(nextWordStart([]rune(m.textinput.Value()), m.textinput.Position()))
	case "b":
		m.textinput.SetCursor(prevWordStart([]rune(m.textinput.Value()), m.textinput.Position()))
	case "x":
		m.deleteTitleRuneAtCursor()
	}

	return m, nil
}

func (m model) updateBody(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+s":
		return m.saveCurrentNote()
	}

	if m.editorMode == insertMode {
		if msg.String() == "esc" {
			return m, m.enterNormalMode()
		}

		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "esc":
		m.state = listView
		return m, m.enterNormalMode()
	case "i":
		return m, m.enterInsertMode()
	case "a":
		m.textarea, _ = m.textarea.Update(keyMsg(tea.KeyRight))
		return m, m.enterInsertMode()
	case "I":
		m.textarea.CursorStart()
		return m, m.enterInsertMode()
	case "A":
		m.textarea.CursorEnd()
		return m, m.enterInsertMode()
	case "h":
		m.textarea, _ = m.textarea.Update(keyMsg(tea.KeyLeft))
	case "j":
		m.textarea.CursorDown()
	case "k":
		m.textarea.CursorUp()
	case "l":
		m.textarea, _ = m.textarea.Update(keyMsg(tea.KeyRight))
	case "0":
		m.textarea.CursorStart()
	case "$":
		m.textarea.CursorEnd()
	case "w":
		m.textarea, _ = m.textarea.Update(altRuneKeyMsg('f'))
	case "b":
		m.textarea, _ = m.textarea.Update(altRuneKeyMsg('b'))
	case "x":
		m.textarea, _ = m.textarea.Update(keyMsg(tea.KeyDelete))
	}

	return m, nil
}

func (m model) saveCurrentNote() (tea.Model, tea.Cmd) {
	m.currNote.Body = m.textarea.Value()

	if err := m.store.SaveNote(m.currNote); err != nil {
		return m, tea.Quit
	}

	notes, err := m.store.GetNotes()
	if err != nil {
		return m, tea.Quit
	}

	m.notes = notes
	m.lastSaved = time.Now()

	if m.currNote.ID == 0 {
		for i := range m.notes {
			if m.notes[i].Title == m.currNote.Title && m.notes[i].Body == m.currNote.Body {
				m.currNote = m.notes[i]
				m.listIndex = i
				break
			}
		}
		return m, nil
	}

	for i := range m.notes {
		if m.notes[i].ID == m.currNote.ID {
			m.currNote = m.notes[i]
			m.listIndex = i
			break
		}
	}

	return m, nil
}

func (m *model) enterInsertMode() tea.Cmd {
	m.editorMode = insertMode
	switch m.state {
	case titleView:
		m.textarea.Blur()
		return tea.Batch(
			m.textinput.Cursor.SetMode(cursor.CursorBlink),
			m.textinput.Focus(),
		)
	case bodyView:
		m.textarea.Blur()
		m.textinput.Blur()
		return tea.Batch(
			m.textarea.Cursor.SetMode(cursor.CursorBlink),
			m.textarea.Focus(),
		)
	default:
		m.textinput.Blur()
		m.textarea.Blur()
	}

	return nil
}

func (m *model) enterNormalMode() tea.Cmd {
	m.editorMode = normalMode
	switch m.state {
	case titleView:
		m.textarea.Blur()
		return tea.Batch(
			m.textinput.Cursor.SetMode(cursor.CursorStatic),
			m.textinput.Focus(),
		)
	case bodyView:
		m.textinput.Blur()
		return tea.Batch(
			m.textarea.Cursor.SetMode(cursor.CursorStatic),
			m.textarea.Focus(),
		)
	default:
		m.textinput.Blur()
		m.textarea.Blur()
	}

	return nil
}

func (m *model) deleteTitleRuneAtCursor() {
	value := []rune(m.textinput.Value())
	pos := m.textinput.Position()
	if pos >= len(value) {
		return
	}

	value = append(value[:pos], value[pos+1:]...)
	m.textinput.SetValue(string(value))
	if pos > len(value) {
		m.textinput.SetCursor(len(value))
	} else {
		m.textinput.SetCursor(pos)
	}
}

func keyMsg(keyType tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg(tea.Key{Type: keyType})
}

func altRuneKeyMsg(r rune) tea.KeyMsg {
	return tea.KeyMsg(tea.Key{Type: tea.KeyRunes, Runes: []rune{r}, Alt: true})
}

func nextWordStart(value []rune, pos int) int {
	if pos >= len(value) {
		return len(value)
	}

	idx := pos
	for idx < len(value) && isWordRune(value[idx]) {
		idx++
	}
	for idx < len(value) && !isWordRune(value[idx]) {
		idx++
	}
	return idx
}

func prevWordStart(value []rune, pos int) int {
	if pos <= 0 {
		return 0
	}

	idx := pos - 1
	for idx > 0 && !isWordRune(value[idx]) {
		idx--
	}
	for idx > 0 && isWordRune(value[idx-1]) {
		idx--
	}
	return idx
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
