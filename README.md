# Goats

Goats is a simple, terminal-based note-taking application written in Go. It utilizes the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework to provide a clean and interactive Terminal User Interface (TUI).

## Features

*   **Create Notes:** Easily add new notes with a title and body.
*   **Edit Notes:** Update existing notes directly from the terminal.
*   **Delete Notes:** Remove unwanted notes with a confirmation step.
*   **Persistent Storage:** All notes are saved locally using SQLite (`goats.db`).
*   **Keyboard Navigation:** Fully navigable using keyboard shortcuts (Vim-style navigation supported).

## Installation

### Prerequisites

*   Go 1.25 or higher

### Build from Source

1.  Clone the repository:
    ```bash
    git clone <repository-url>
    cd goats
    ```

2.  Download dependencies:
    ```bash
    go mod download
    ```

3.  Run the application:
    ```bash
    go run cmd/goats/main.go
    ```

    Or build the binary:
    ```bash
    go build -o goats cmd/goats/main.go
    ./goats
    ```

## Usage

Once the application is running, you can use the following keybindings to interact with it:

### List View (Main Screen)
*   `j` / `Down Arrow`: Move selection down.
*   `k` / `Up Arrow`: Move selection up.
*   `n`: Create a new note.
*   `d`: Delete the selected note (requires confirmation).
    *   Press `Enter` to confirm deletion.
    *   Press `Esc` to cancel deletion.
*   `Enter`: Open the selected note for editing.
*   `q` / `Ctrl+c`: Quit the application.

### Title View (Creating a Note)
*   Type your title.
*   `Enter`: Confirm title and move to the body editor.
*   `Esc`: Cancel and return to the list view.

### Body View (Editing a Note)
*   Type to edit the note content.
*   `Ctrl+s`: Save the note.
*   `Esc`: Return to the list view (unsaved changes will be lost unless saved).

## Tech Stack

*   **Language:** Go
*   **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea)
*   **Components:** [Bubbles](https://github.com/charmbracelet/bubbles) (Textarea, Textinput)
*   **Styling:** [Lip Gloss](https://github.com/charmbracelet/lipgloss)
*   **Database:** [go-sqlite3](https://github.com/mattn/go-sqlite3)
