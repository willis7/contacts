package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx/types"
)

// ContactFavorites is a field that contains a contact's favorites
type ContactFavorites struct {
	Colors []string `json:"colors"`
}

// Contact represents a Contact model in the database
type Contact struct {
	ID                   int
	Name, Address, Phone string

	FavoritesJSON types.JSONText    `db:"favorites"`
	Favorites     *ContactFavorites `db:"-"`

	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

var (
	connectionString = flag.String("conn", getenvWithDefault("DATABASE_URL", ""), "PostgreSQL connection string")
	listenAddr       = flag.String("addr", getenvWithDefault("LISTENADDR", ":8080"), "HTTP address to listen on")
	db               *sqlx.DB
	tmpl             = template.New("")
)

func getenvWithDefault(name, defaultValue string) string {
	val := os.Getenv(name)
	if val == "" {
		val = defaultValue
	}
	return val
}

func main() {
	flag.Parse()
	var err error

	// Templating
	tmpl.Funcs(template.FuncMap{"StringsJoin": strings.Join})
	_, err = tmpl.ParseGlob(filepath.Join(".", "templates", "*.html"))
	if err != nil {
		log.Fatalf("Unable to parse templates: %v\n", err)
	}

	// Database
	if *connectionString == "" {
		log.Fatalln("Please pass the connection string using the -conn option")
	}

	// We need to wait until postgres is available, so some connections may fail.
	// Hence, we try at least 10x with a 3 second wait inbetween
	retries := 10
	for i := 0; i <= retries; i++ {
		db, err = sqlx.Connect("pgx", *connectionString)
		if err != nil {
			log.Fatalf("Unable to establish connection: %v\n", err)
		}
		time.Sleep(time.Second * 3)
	}

	// Server
	log.Printf("listening on %s`n", *listenAddr)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(*listenAddr, nil)
}

func fetchContacts() ([]*Contact, error) {
	contacts := []*Contact{}
	err := db.Select(&contacts, "select * from contacts")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch contacts")
	}

	for _, contact := range contacts {
		err := json.Unmarshal(contact.FavoritesJSON, &contact.Favorites)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to parse JSON favourites")
		}
	}

	return contacts, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	contacts, err := fetchContacts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	tmpl.ExecuteTemplate(w, "index.html", struct{ Contacts []*Contact }{contacts})
}
