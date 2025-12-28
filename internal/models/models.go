package models

type TaskDTO struct {
	Header      string `json:"header" validate:"required"`
	Description string `json:"description"`
	Finished    bool   `json:"finished"`
}

type TaskDomain struct {
	ID          uint   `json:"id" validate:"required"`
	Header      string `json:"header" validate:"required"`
	Description string `json:"description" validate:"required"`
	Finished    bool   `json:"finished" validate:"required"`
}
