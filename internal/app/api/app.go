package apiapp

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	handlers *Handlers
	port     int
}

func NewApp(port int, handlers *Handlers) *App {
	return &App{
		handlers: handlers,
		port:     port,
	}
}

func (a *App) Run() *http.Server {
	router := mux.NewRouter()

	serv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.port),
		Handler: router,
	}

	go serv.ListenAndServe()

	return serv
}
