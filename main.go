package main

import (
	"embed"
	"os"

	"a21hc3NpZ25tZW50/app/controller"
	"a21hc3NpZ25tZW50/app/model"
	repo "a21hc3NpZ25tZW50/app/repository"
	"a21hc3NpZ25tZW50/config"

	_ "github.com/jackc/pgx/v4/stdlib"
)

//go:embed app/view/*
var Resources embed.FS

func main() {
	os.Setenv("DATABASE_URL", "postgres://postgres:12345678@localhost:5432/my_db") // Hapus jika akan melakukan deploy ke fly.io

	db := config.NewDB()
	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	conn.AutoMigrate(&model.Teacher{})
	teacherHandle := repo.NewTeacherRepo(conn)

	mainAPI := controller.NewAPI(teacherHandle, Resources)
	mainAPI.Start()
}

func FlyURL() string {
	return "https://najibfikri-be2754215.fly.dev" // TODO: replace this
}
