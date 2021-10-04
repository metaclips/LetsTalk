package v1controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/maps"
)

func GetFaculties(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	faculties := make([]string, len(maps.Faculties))

	i := 0
	for faculty := range maps.Faculties {
		faculties[i] = faculty
		i++
	}

	apihelpers.Respond(w, http.StatusOK, faculties, "successfully retrieved faculties", "")
}

func GetDepartments(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	departments := make([]string, len(maps.Departments))

	i := 0
	for department := range maps.Departments {
		departments[i] = department
		i++
	}

	apihelpers.Respond(w, http.StatusOK, departments, "successfully retrieved departments", "")
}
