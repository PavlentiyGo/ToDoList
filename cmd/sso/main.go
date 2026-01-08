package main

import (
	"TrueToDoList/internal/app"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	err := godotenv.Load(".env")
	if err != nil {
		panic(".env not found")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_DSN"))
	defer pool.Close()
	if err != nil {
		panic("error to create pgxpool")
	}
	application := app.New(context.Background(), log, 44044, 8081, pool)

	go application.GRPCServer.MustRun()
	go application.APIServer.Run()
	stop := make(chan os.Signal)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	application.APIServer.Stop()

	log.Info("app stopped")

}
