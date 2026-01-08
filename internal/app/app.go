package app

import (
	"TrueToDoList/internal/api/tododata"
	apiapp "TrueToDoList/internal/app/api"
	grpcapp "TrueToDoList/internal/app/grpc"
	todos "TrueToDoList/internal/service/todo"
	"TrueToDoList/storage/postgresql"
	"context"
	ssov1 "github.com/PavlentiyGo/protoToDo/gen/go/sso"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"strconv"
)

type App struct {
	GRPCServer *grpcapp.App
	APIServer  *apiapp.App
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
	cc, err := grpc.NewClient(net.JoinHostPort("localhost", strconv.Itoa(grpcPort)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	app := apiapp.NewApp(ctx, log, apiPort, tododata.NewHandlers(ctx, ssov1.NewToDoDataClient(cc)))
	return &App{
		GRPCServer: grpcServ,
		APIServer:  app,
	}
}
