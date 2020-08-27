package routers

import (
	"github.com/emanuelcondo/api-todo/controllers"
	"github.com/emanuelcondo/api-todo/services"
	"github.com/emanuelcondo/api-todo/utils"
	"github.com/go-chi/chi"
)

func TodoRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(utils.InputValidator("Authorization", utils.DATATYPE_STRING, utils.SOURCE_HEADERS, true))
	router.Use(services.Authenticate())

	optionsPage := utils.InputOptions{ MinValue: 1, ParseToInt: true }
	optionsLimit := utils.InputOptions{ MinValue: 1, MaxValue: 100, ParseToInt: true }
	optionsID := utils.InputOptions{ MinValue: 1, ParseToInt: true }
	
	router.Group(func(r chi.Router) {
		r.Use(services.RequiredRole(services.ROLE_VIEWER))
		r.Use(utils.InputValidator("search", utils.DATATYPE_STRING, utils.SOURCE_QUERY, false))
		r.Use(utils.InputValidator("page", utils.DATATYPE_INTEGER, utils.SOURCE_QUERY, false, optionsPage))
		r.Use(utils.InputValidator("limit", utils.DATATYPE_INTEGER, utils.SOURCE_QUERY, false, optionsLimit))
		r.Get("/", controllers.SearchTodos)
	})

	router.Group(func(r chi.Router) {
		r.Use(services.RequiredRole(services.ROLE_EDITOR))
		r.Use(utils.InputValidator("title", utils.DATATYPE_STRING, utils.SOURCE_BODY, true))
		r.Use(utils.InputValidator("description", utils.DATATYPE_STRING, utils.SOURCE_BODY, true))
		r.Post("/", controllers.CreateTodo)
	})

	router.Group(func(r chi.Router) {
		r.Use(services.RequiredRole(services.ROLE_EDITOR))
		r.Use(utils.InputValidator("id", utils.DATATYPE_INTEGER, utils.SOURCE_PARAMS, true, optionsID))
		r.Use(utils.InputValidator("title", utils.DATATYPE_STRING, utils.SOURCE_BODY, true))
		r.Use(utils.InputValidator("description", utils.DATATYPE_STRING, utils.SOURCE_BODY, true))
		r.Put("/{id}", controllers.UpdateTodo)
	})

	return router
}
