package model

type Content struct {
	StudentId    string `json:"student_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	CurrentClass string `json:"current_class" binding:"required"`
}
