package shortlink

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/app"
)

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
	routes := h.app.Group(prefix)
	{
		routes.POST("/", h.CreateShortenerLink)
		routes.GET("/:shortenerURL", h.GetOriginalURL)
		// authentication.POST("/register", h.Register)
		// authentication.POST("/login", h.Login)

		// authentication.Use(middleware.AuthenticateJWT())
		// {
		// 	authentication.GET("/me", h.GetMe)
		// }
	}
}

func (h *Handler) CreateShortenerLink(c *gin.Context) {
	var data CreateShortenerLinkRequestDTO
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.useCase.CreateShortenerLink(&data)
	if err != nil {
		errMsg := err.Error()
		c.JSON(err.Code(), app.NewErrorResponse("Failed to create shortener link", &errMsg))
		return
	}

	c.JSON(200, app.NewSuccessResponse("Shortener link created successfully", res))
}

func (h *Handler) GetOriginalURL(c *gin.Context) {
	shortenerURL := c.Param("shortenerURL")
	log.Println(shortenerURL)
	res, err := h.useCase.GetOriginalURL(shortenerURL)
	if err != nil {
		errMsg := err.Error()
		c.JSON(err.Code(), app.NewErrorResponse("Failed to get original URL", &errMsg))
		return
	}

	c.Redirect(301, *res)
}