package apiapp

import (
	"encoding/json"
	ssov1 "github.com/PavlentiyGo/protoToDo/gen/go/sso"
	"log/slog"
	"net/http"
)

type Handlers struct {
	log *slog.Logger
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	task := TaskDTO{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		Err := NewErrDTO(err)
		http.Error(w, Err.ToString(), http.StatusInternalServerError)
	}
	ssov1.CreateTaskRequest{
		Title:       task.Title,
		Description: task.Description,
	}

}

func (h *Handlers) GetTask(w *http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) DeleteTask(w *http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) DoneTask(w *http.ResponseWriter, r *http.Request) {

}
