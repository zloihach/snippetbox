package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func dbConnect(dsn string) (*sql.DB, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("postgres", "postgres://web:web@localhost:5432/snippetbox/sslmode=disable", "Postgres connection string")
	flag.Parse()

	db, err := dbConnect(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
