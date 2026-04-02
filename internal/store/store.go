package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type ReadingItem struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	URL string `json:"url"`
	Type string `json:"type"`
	Status string `json:"status"`
	Rating int `json:"rating"`
	Notes string `json:"notes"`
	Tags string `json:"tags"`
	CompletedAt string `json:"completed_at"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"quiver.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS reading_items(id TEXT PRIMARY KEY,title TEXT NOT NULL,author TEXT DEFAULT '',url TEXT DEFAULT '',type TEXT DEFAULT 'article',status TEXT DEFAULT 'unread',rating INTEGER DEFAULT 0,notes TEXT DEFAULT '',tags TEXT DEFAULT '',completed_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *ReadingItem)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO reading_items(id,title,author,url,type,status,rating,notes,tags,completed_at,created_at)VALUES(?,?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.Author,e.URL,e.Type,e.Status,e.Rating,e.Notes,e.Tags,e.CompletedAt,e.CreatedAt);return err}
func(d *DB)Get(id string)*ReadingItem{var e ReadingItem;if d.db.QueryRow(`SELECT id,title,author,url,type,status,rating,notes,tags,completed_at,created_at FROM reading_items WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.Author,&e.URL,&e.Type,&e.Status,&e.Rating,&e.Notes,&e.Tags,&e.CompletedAt,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]ReadingItem{rows,_:=d.db.Query(`SELECT id,title,author,url,type,status,rating,notes,tags,completed_at,created_at FROM reading_items ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []ReadingItem;for rows.Next(){var e ReadingItem;rows.Scan(&e.ID,&e.Title,&e.Author,&e.URL,&e.Type,&e.Status,&e.Rating,&e.Notes,&e.Tags,&e.CompletedAt,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM reading_items WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM reading_items`).Scan(&n);return n}
