package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql" 
	"github.com/kasante1/paste_it_backend/internal/models"
	"github.com/go-playground/form/v4"
	"time"
	"github.com/alexedwards/scs/mysqlstore" // New import
	"github.com/alexedwards/scs/v2"
 
)

type application struct {
	errorLog *log.Logger
	infoLog	 *log.Logger
	snippets *models.SnippetModel
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager

}


func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
	return nil, err
	}
	if err = db.Ping(); err != nil {
	return nil, err
	}
	return db, nil
	}


func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
		if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}
	
	infoLog.Printf("Starting server on %s", *addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
