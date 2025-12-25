package tododata

import (
	"TrueToDoList/internal/domain"
	"TrueToDoList/internal/storage"
	"context"
	"errors"
	ssov1 "github.com/PavlentiyGo/protoToDo/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type ToDoData interface {
	GetTasks(ctx context.Context, id *int64) (map[int64]domain.Task, error)
	CreateTask(ctx context.Context, title string, description string) (int64, error)
	DeleteTask(ctx context.Context, id int64) error
	DoneTask(ctx context.Context, id int64) error
}

type serverAPI struct {
	ssov1.UnimplementedToDoDataServer
	todo ToDoData
}

func RegisterServerAPI(gRPC *grpc.Server) {
	ssov1.RegisterToDoDataServer(gRPC, &serverAPI{})
}

func (s *serverAPI) GetTasks(ctx context.Context, request *ssov1.GetTasksRequest) (*ssov1.GetTasksResponse, error) {
	if err := validateGet(request.Id); err != nil {
		return nil, err
	}
	tasks, err := s.todo.GetTasks(ctx, request.Id)
	if err != nil {
		if errors.Is(err, storage.ErrNoSuchTask) {
			return nil, status.Error(codes.InvalidArgument, "no such task")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	dto := make(map[int64]*ssov1.Todo)
	for key, val := range tasks {
		dto[key] = &ssov1.Todo{
			Title:       val.Title,
			Description: val.Description,
			Done:        val.Done,
		}
	}

	return &ssov1.GetTasksResponse{
		Todos: dto,
	}, nil
}

func (s *serverAPI) CreateTask(ctx context.Context, request *ssov1.CreateTaskRequest) (*ssov1.CreateTaskResponse, error) {
	if err := validateCreate(request.GetTitle(), request.GetDescription()); err != nil {
		return nil, err
	}
	id, err := s.todo.CreateTask(ctx, request.GetTitle(), request.GetDescription())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.CreateTaskResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) DeleteTask(ctx context.Context, request *ssov1.DeleteTaskRequest) (*ssov1.DeleteTaskResponse, error) {
	if err := validateDelete(request.GetTitle()); err != nil {
		return nil, err
	}
	if err := s.todo.DeleteTask(ctx, request.GetTitle()); err != nil {
		if errors.Is(err, storage.ErrNoSuchTask) {
			return nil, status.Error(codes.InvalidArgument, "no such task")
		}
	}
	return &ssov1.DeleteTaskResponse{}, nil
}

func (s *serverAPI) DoneTask(ctx context.Context, request *ssov1.DoneTaskRequest) (*ssov1.DoneTaskResponse, error) {
	if err := validateDone(request.GetId()); err != nil {
		return nil, err
	}
	if err := s.todo.DoneTask(ctx, request.GetId()); err != nil {
		if errors.Is(err, storage.ErrNoSuchTask) {
			return nil, status.Error(codes.InvalidArgument, "no such task")
		}
	}
	return &ssov1.DoneTaskResponse{}, nil
}

func validateCreate(title string, description string) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(description) == "" {
		return status.Error(codes.InvalidArgument, "description or title is empty")
	}

	return nil
}
func validateDelete(id int64) error {
	if id <= 0 {
		return status.Error(codes.InvalidArgument, "wrong id task")
	}
	return nil
}
func validateGet(id *int64) error {
	if id == nil {
		return nil
	}
	if *id <= 0 {
		return status.Error(codes.InvalidArgument, "wrong id task")
	}
	return nil
}

func validateDone(id int64) error {
	if id <= 0 {
		return status.Error(codes.InvalidArgument, "wrong id task")
	}
	return nil
}
