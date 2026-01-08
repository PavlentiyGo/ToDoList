package tododata

import "encoding/json"

type TaskDTO struct {
	Title       string
	Description string
}

type ErrDTO struct {
	Err string
}
type TaskIdDTO struct {
	Id string `json:"id"`
}

func NewTaskIdDTO(id string) TaskIdDTO {
	return TaskIdDTO{Id: id}
}

func NewErrDTO(err error) ErrDTO {
	return ErrDTO{
		Err: err.Error(),
	}
}

func (e *ErrDTO) ToString() string {
	bytes, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
