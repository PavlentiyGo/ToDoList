package app

import (
	grpcapp "TrueToDoList/internal/app/grpc"
	todos "TrueToDoList/internal/service/todo"
	"TrueToDoList/storage/postgresql"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	grpcPort int,
	apiPort int,
	pool *pgxpool.Pool,
) *App {
	storage, err := postgresql.New(ctx, pool)
	if err != nil {
		panic(err)
	}
	todo := todos.New(log, storage)
	grpcServ := grpcapp.New(log, grpcPort, todo)
	return &App{
		GRPCServer: grpcServ,
	}
}
