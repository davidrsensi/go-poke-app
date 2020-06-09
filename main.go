package main

import (
	"context"
	"go-poke-app/pkg/conf"
	"go-poke-app/pkg/database"
	"go-poke-app/pkg/handlers"
	"go-poke-app/pkg/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "poke-api : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	l.Printf("main : Started")
	defer l.Println("main : Completed")

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost"`
			Port            string        `conf:"default::8080"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			Host string `conf:"default:0.0.0.0"`
			Name string `conf:"default:mongodb"`
		}
	}

	conf.Parse(os.Args[1:], "poke-api", &cfg)

	// output configs on startup
	out, err := conf.String(&cfg)
	if err != nil {
		l.Printf("Error generating config usage, %v", err)
	}
	l.Printf("main : Config :\n%v\n", out)

	database.Init(database.Config{
		Host: cfg.DB.Host,
		Name: cfg.DB.Name,
	})

	utils.LoadTemplates("pkg/templates/*.html")

	// handlers
	ph := handlers.NewPokes(l)
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()

	getRouter.HandleFunc("/", ph.GetPokes)
	postRouter.HandleFunc("/", ph.PostPokes)

	fs := http.FileServer(http.Dir("./static/"))
	sm.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	server := &http.Server{
		Addr:         cfg.Web.Port,
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	go func() {
		l.Println("Starting server on port " + server.Addr)

		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = server.Shutdown(tc)

	if err != nil {
		l.Printf("main : Graceful shutdown did not complete in %v : %v", 30*time.Second, err)
		server.Close()
	}
}
