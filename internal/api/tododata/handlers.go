package tododata

import (
	"TrueToDoList/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ssov1 "github.com/PavlentiyGo/protoToDo/gen/go/sso"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handlers struct {
	todo ssov1.ToDoDataClient
	ctx  context.Context
}

func NewHandlers(ctx context.Context, todo ssov1.ToDoDataClient) *Handlers {
	return &Handlers{
		todo: todo,
		ctx:  ctx,
	}
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	task := TaskDTO{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		Err := NewErrDTO(err)
		http.Error(w, Err.ToString(), http.StatusInternalServerError)
		return
	}
	resp, err := h.todo.CreateTask(r.Context(), &ssov1.CreateTaskRequest{Title: task.Title, Description: task.Description})
	if err != nil {
		Err := NewErrDTO(err)
		http.Error(w, Err.ToString(), http.StatusInternalServerError)
		return
	}

	response := NewTaskIdDTO(strconv.Itoa(int(resp.GetId())))

	bytes, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		Err := NewErrDTO(err)
		http.Error(w, Err.ToString(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(bytes); err != nil {
		panic("error with write json")
	}
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		resp, err := h.todo.GetTasks(h.ctx, &ssov1.GetTasksRequest{Id: nil})
		if err != nil {
			Err := NewErrDTO(err)
			http.Error(w, Err.ToString(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			Err := NewErrDTO(err)
			http.Error(w, Err.ToString(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(bytes); err != nil {
			panic(err)
		}
	} else {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		idPointer := int64(idInt)
		resp, err := h.todo.GetTasks(h.ctx, &ssov1.GetTasksRequest{Id: &idPointer})
		if err != nil {
			if errors.Is(err, storage.ErrNoSuchTask) {
				http.Error(w, "no such task", http.StatusBadRequest)
				return
			}
			Err := NewErrDTO(err)
			fmt.Println(Err.ToString())
			http.Error(w, Err.ToString(), http.StatusBadRequest)
			return
		}
		bytes, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			Err := NewErrDTO(err)
			http.Error(w, Err.ToString(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(bytes); err != nil {
			panic(err)
		}
	}
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	_, err = h.todo.DeleteTask(h.ctx, &ssov1.DeleteTaskRequest{Id: int64(id)})
	if err != nil {
		if errors.Is(err, storage.ErrNoSuchTask) {
			http.Error(w, "no such task", http.StatusBadRequest)
			return
		}
		Err := NewErrDTO(err)
		http.Error(w, Err.ToString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *Handlers) DoneTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	_, err = h.todo.DoneTask(h.ctx, &ssov1.DoneTaskRequest{Id: int64(id)})
	if err != nil {
		if errors.Is(err, storage.ErrNoSuchTask) {
			http.Error(w, "no such task", http.StatusBadRequest)
			return
		}
		Err := NewErrDTO(err)
		http.Error(w, Err.ToString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
