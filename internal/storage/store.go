package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID        int64
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.conn, err = sql.Open("sqlite3", "./goats.db")
	if err != nil {
		return err
	}

	createTableStmt := `CREATE TABLE IF NOT EXISTS notes (
		id integer not null primary key,
		title text not null,
		body text not null,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = s.conn.Exec(createTableStmt); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetNotes() ([]Note, error) {
	rows, err := s.conn.Query("SELECT * FROM notes;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var note Note
		rows.Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt, &note.UpdatedAt)
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Store) SaveNote(note Note) error {
	now := time.Now()

	if note.ID == 0 {
		query := `INSERT INTO notes (title, body, created_at, updated_at) VALUES (?, ?, ?, ?)`
		_, err := s.conn.Exec(query, note.Title, note.Body, now, now)
		return err
	}

	query := `UPDATE notes SET title = ?, body = ?, updated_at = ? WHERE id = ?`
	_, err := s.conn.Exec(query, note.Title, note.Body, now, note.ID)
	return err
}

func (s *Store) DeleteNote(note Note) error {
	query := `DELETE FROM notes WHERE id = ?`
	_, err := s.conn.Exec(query, note.ID)
	return err
}
