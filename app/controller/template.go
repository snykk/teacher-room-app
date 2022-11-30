package controller

import (
	"net/http"
	"os"
)

func (api *API) IndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl := api.BaseViewPath()

	data := map[string]interface{}{
		"Region": os.Getenv("FLY_REGION"),
	}

	tmpl.ExecuteTemplate(w, "teacher.html.tmpl", data)
}
