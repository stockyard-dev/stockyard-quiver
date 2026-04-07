package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ db *sql.DB }
type ReadingItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Rating      int    `json:"rating"`
	Notes       string `json:"notes"`
	Tags        string `json:"tags"`
	CompletedAt string `json:"completed_at"`
	CreatedAt   string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "quiver.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS reading_items(id TEXT PRIMARY KEY,title TEXT NOT NULL,author TEXT DEFAULT '',url TEXT DEFAULT '',type TEXT DEFAULT 'article',status TEXT DEFAULT 'unread',rating INTEGER DEFAULT 0,notes TEXT DEFAULT '',tags TEXT DEFAULT '',completed_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
	resource TEXT NOT NULL,
	record_id TEXT NOT NULL,
	data TEXT NOT NULL DEFAULT '{}',
	PRIMARY KEY(resource, record_id)
)`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string        { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string          { return time.Now().UTC().Format(time.RFC3339) }
func (d *DB) Create(e *ReadingItem) error {
	e.ID = genID()
	e.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO reading_items(id,title,author,url,type,status,rating,notes,tags,completed_at,created_at)VALUES(?,?,?,?,?,?,?,?,?,?,?)`, e.ID, e.Title, e.Author, e.URL, e.Type, e.Status, e.Rating, e.Notes, e.Tags, e.CompletedAt, e.CreatedAt)
	return err
}
func (d *DB) Get(id string) *ReadingItem {
	var e ReadingItem
	if d.db.QueryRow(`SELECT id,title,author,url,type,status,rating,notes,tags,completed_at,created_at FROM reading_items WHERE id=?`, id).Scan(&e.ID, &e.Title, &e.Author, &e.URL, &e.Type, &e.Status, &e.Rating, &e.Notes, &e.Tags, &e.CompletedAt, &e.CreatedAt) != nil {
		return nil
	}
	return &e
}
func (d *DB) List() []ReadingItem {
	rows, _ := d.db.Query(`SELECT id,title,author,url,type,status,rating,notes,tags,completed_at,created_at FROM reading_items ORDER BY created_at DESC`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []ReadingItem
	for rows.Next() {
		var e ReadingItem
		rows.Scan(&e.ID, &e.Title, &e.Author, &e.URL, &e.Type, &e.Status, &e.Rating, &e.Notes, &e.Tags, &e.CompletedAt, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}
func (d *DB) Update(e *ReadingItem) error {
	_, err := d.db.Exec(`UPDATE reading_items SET title=?,author=?,url=?,type=?,status=?,rating=?,notes=?,tags=?,completed_at=? WHERE id=?`, e.Title, e.Author, e.URL, e.Type, e.Status, e.Rating, e.Notes, e.Tags, e.CompletedAt, e.ID)
	return err
}
func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM reading_items WHERE id=?`, id)
	return err
}
func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM reading_items`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []ReadingItem {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (title LIKE ?)"
		args = append(args, "%"+q+"%")
	}
	if v, ok := filters["type"]; ok && v != "" {
		where += " AND type=?"
		args = append(args, v)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(`SELECT id,title,author,url,type,status,rating,notes,tags,completed_at,created_at FROM reading_items WHERE `+where+` ORDER BY created_at DESC`, args...)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []ReadingItem
	for rows.Next() {
		var e ReadingItem
		rows.Scan(&e.ID, &e.Title, &e.Author, &e.URL, &e.Type, &e.Status, &e.Rating, &e.Notes, &e.Tags, &e.CompletedAt, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Stats() map[string]any {
	m := map[string]any{"total": d.Count()}
	rows, _ := d.db.Query(`SELECT status,COUNT(*) FROM reading_items GROUP BY status`)
	if rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}
	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
