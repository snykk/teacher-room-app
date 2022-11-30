package model

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Name         string `json:"name"`
	FieldOfStudy string `json:"field_of_study"`
	Age          int    `json:"age"`
}

type ViewTeacher struct {
	Name         string `json:"name"`
	FieldOfStudy string `json:"field_of_study"`
	Age          int    `json:"age"`
}

type UpdateTeacher struct {
	Id      int    `json:"id"`
	NewName string `json:"new_name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
