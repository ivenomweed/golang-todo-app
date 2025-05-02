package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"
const cssVersion = "1"

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 05 * time.Second,
		WriteTimeout:      05 * time.Second,
	}
	app.infoLog.Printf(

		"Starting HTTP server in %s mode on: http://localhost:%d",
		app.config.env,
		app.config.port,
	)
	return srv.ListenAndServe()
}

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.LstdFlags)
	var cfg config
	flag.IntVar(
		&cfg.port,
		"port",
		4000,
		"Specify server port. { default: 4000 }",
	)
	flag.StringVar(
		&cfg.env,
		"env",
		"dev",
		"Specify deployment environment. { default: dev | prod }",
	)
	flag.StringVar(
		&cfg.api,
		"api",
		"http://localhost:4001",
		"Specify API url. { default: http://localhost:40001 }",
	)
	flag.Parse()
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	infoLog := log.New(
		os.Stderr,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags,
	)
	errorLog := log.New(
		os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags,
	)
	tc := make(map[string]*template.Template)
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
	}
	err := app.serve()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
