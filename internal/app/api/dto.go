package apiapp

import "encoding/json"

type TaskDTO struct {
	Title       string
	Description string
}

type ErrDTO struct {
	err string
}

func NewErrDTO(err error) ErrDTO {
	return ErrDTO{
		err: err.Error(),
	}
}

func (e *ErrDTO) ToString() string {
	bytes, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
