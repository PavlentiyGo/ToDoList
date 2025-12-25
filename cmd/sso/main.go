package main

import (
	"TrueToDoList/internal/app"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	application := app.New(log, 44044, 8081, "")

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()

	log.Info("app stopped")

}
