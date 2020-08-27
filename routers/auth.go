package routers

import (
	"github.com/emanuelcondo/api-todo/controllers"
	"github.com/emanuelcondo/api-todo/utils"
	"github.com/go-chi/chi"
)

func AuthRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(utils.InputValidator("email", utils.DATATYPE_STRING, utils.SOURCE_BODY, true))
		r.Use(utils.InputValidator("password", utils.DATATYPE_STRING, utils.SOURCE_BODY, true))
		r.Post("/login", controllers.Login)
	})

	return router
}
