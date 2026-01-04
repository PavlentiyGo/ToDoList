package grpcapp

import (
	"TrueToDoList/internal/grpc/tododata"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	port int,
	todo tododata.ToDoData,
) *App {
	gRPCServer := grpc.NewServer()
	tododata.RegisterServerAPI(gRPCServer, todo)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(nil)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("starting grpcServer")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPCServer")

	a.gRPCServer.GracefulStop()
}
