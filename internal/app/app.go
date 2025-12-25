package app

import (
	grpcapp "TrueToDoList/internal/app/grpc"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	apiPort int,
	storagePath string,
) *App {
	grpcServ := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcServ,
	}
}
