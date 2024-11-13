package main

import (
	"encoding/json"
	"flag"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var (
	dateQuery  = flag.String("date", "", "query for status before this date, format: "+time.DateTime)
	idQuery    = flag.String("id", "", "query for id")
	tplPath    = flag.String("tpl", "", "path to template file")
	statusPath = flag.String("path", "status.txt", "path to status output file")
	dbPath     = flag.String("db", "status.db", "path to database")
)

const defaultTpl = `{{ . }}`

type Entry struct {
	ID          string    `json:"id"`
	Posted      int       `json:"posted"`
	StatusEmoji string    `json:"status_emoji" db:"status_emoji"`
	StatusText  string    `json:"status_text" db:"status_text"`
	URL         string    `json:"url"`
	Timestamp   time.Time `json:"timestamp"`
}

func main() {
	flag.Parse()

	db, err := sqlx.Open("sqlite", *dbPath)
	checkErr(err)

	err = schema(db)
	checkErr(err)

	status := Entry{}

	if *idQuery != "" {
		status, err = findByID(db, *idQuery)
		checkErr(err)
	} else if *dateQuery != "" {
		d, err := time.Parse(time.DateTime, *dateQuery)
		checkErr(err)
		status, err = findByDate(db, d)
		checkErr(err)
	} else if fpath := os.Getenv("FNAME"); fpath != "" {
		f, err := os.ReadFile(fpath)
		checkErr(err)
		status, err = record(db, f)
		checkErr(err)
	} else {
		log.Fatal("Nothing to do")
	}

	err = writeStatus(status, *statusPath, *tplPath)
	checkErr(err)

	log.Printf("Status recorded successfully: %v", status)
}

func schema(db *sqlx.DB) error {
	q := `
	create table if not exists entries (
	    id text primary key,
		posted int,
		status_emoji text,
		status_text text,
		url text,
	    timestamp timestamp
	)
	`
	_, err := db.Exec(q)
	return err
}

func record(db *sqlx.DB, f []byte) (Entry, error) {
	entry := Entry{}
	err := json.Unmarshal(f, &entry)
	entry.Timestamp = time.Now()
	q := "insert or replace into entries values (?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(q, entry.ID, entry.Posted, entry.StatusEmoji, entry.StatusText, entry.URL, entry.Timestamp)
	return entry, err
}

func findByID(db *sqlx.DB, id string) (Entry, error) {
	entry := Entry{}
	err := db.Get(&entry, "select * from entries where id=?", id)
	return entry, err
}

func findByDate(db *sqlx.DB, date time.Time) (Entry, error) {
	entry := Entry{}
	err := db.Get(&entry, `select * from entries where timestamp < ? order by timestamp limit 1`, date)
	return entry, err
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writeStatus(status Entry, path string, tplPath string) error {
	tpl := template.Must(template.New("").Parse(defaultTpl))
	if tplPath != "" {
		tpl = template.Must(template.ParseFiles(tplPath))

	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = tpl.Execute(f, status)
	return err
}
