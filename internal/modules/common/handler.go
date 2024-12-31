package common

import "github.com/gin-gonic/gin"

type Handler struct {
	useCase IUseCase
	app     *gin.Engine
}

func NewHandler(app *gin.Engine, useCase IUseCase, prefixApi string) {
	handler := &Handler{
		app:     app,
		useCase: useCase,
	}

	handler.Routes(prefixApi)
}

func (h *Handler) Routes(prefix string) {
	_ = h.app.Group(prefix)
	{
		
		// authentication.POST("/register", h.Register)
		// authentication.POST("/login", h.Login)

		// authentication.Use(middleware.AuthenticateJWT())
		// {
		// 	authentication.GET("/me", h.GetMe)
		// }
	}
}
