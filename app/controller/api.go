package controller

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"text/template"

	repo "a21hc3NpZ25tZW50/app/repository"
)

type API struct {
	teacherRepo repo.TeacherRepo
	embed       embed.FS
	mux         *http.ServeMux
}

func (api *API) BaseViewPath() *template.Template {
	var tmpl = template.Must(template.ParseFS(api.embed, "app/view/*"))
	return tmpl
}

func NewAPI(teacherRepo repo.TeacherRepo, embed embed.FS) API {
	mux := http.NewServeMux()
	api := API{
		teacherRepo,
		embed,
		mux,
	}

	mux.HandleFunc("/", api.IndexPage)

	mux.Handle("/api/teacher/add", http.HandlerFunc(api.AddTeacher))
	mux.Handle("/api/teacher/read", http.HandlerFunc(api.ReadTeacher))
	mux.Handle("/api/teacher/update", http.HandlerFunc(api.UpdateTeacher))
	mux.Handle("/api/teacher/delete", http.HandlerFunc(api.DeleteTeacher))

	mux.Handle("/api/teacher/reset", http.HandlerFunc(api.ResetTeacher))

	return api
}

func (api *API) Handler() *http.ServeMux {
	return api.mux
}

func (api *API) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	fmt.Println("starting web server at http://localhost:8080")
	http.ListenAndServe(":"+port, api.Handler())
}
