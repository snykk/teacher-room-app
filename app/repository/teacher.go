package handler

import (
	"a21hc3NpZ25tZW50/app/model"

	"gorm.io/gorm"
)

type TeacherRepo struct {
	db *gorm.DB
}

func NewTeacherRepo(db *gorm.DB) TeacherRepo {
	return TeacherRepo{db}
}

func (u *TeacherRepo) AddTeacher(teacher model.Teacher) error {
	err := u.db.Model(model.Teacher{}).Create(&teacher).Error
	return err
}

func (u *TeacherRepo) ReadTeacher() ([]model.ViewTeacher, error) {
	var data []model.ViewTeacher
	err := u.db.Table("teachers").Select("name, field_of_study, age").Where("deleted_at is NULL").Find(&data).Error
	return data, err
}

func (u *TeacherRepo) UpdateName(id uint, name string) error {
	err := u.db.Model(model.Teacher{}).Where("id = ?", id).Update("name", name).Error
	return err
}

func (u *TeacherRepo) DeleteTeacher(id uint) error {
	err := u.db.Where("id = ?", id).Delete(&model.Teacher{}).Error
	return err
}
