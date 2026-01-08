package apiapp

import (
	"TrueToDoList/internal/api/tododata"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type App struct {
	server *http.Server
	log    *slog.Logger
	ctx    context.Context
}

func NewApp(ctx context.Context, log *slog.Logger, port int, handlers *tododata.Handlers) *App {
	router := mux.NewRouter()
	router.Path("/tasks").Methods("POST").HandlerFunc(handlers.CreateTask)
	router.Path("/tasks").Methods("GET").HandlerFunc(handlers.GetTask)
	router.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(handlers.DeleteTask)
	router.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(handlers.DoneTask)
	serv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	return &App{
		ctx:    ctx,
		server: serv,
		log:    log,
	}
}

func (a *App) Run() {
	const op = "app.api.Run"

	log := a.log.With(slog.String("op", op))
	log.Info("starting api server")
	go a.server.ListenAndServe()
}
func (a *App) Stop() {
	const op = "app.api.Stop"

	log := a.log.With(slog.String("op", op))
	log.Info("stopping api server")

	a.server.Shutdown(a.ctx)
}
