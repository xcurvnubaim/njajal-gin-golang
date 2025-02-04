package shortlink

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/app"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/query"
	CustomValidator "github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/validator"
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
		routes.GET("/", h.GetAllShortenerLink)
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

func (h *Handler) GetAllShortenerLink(c *gin.Context) {
	queryParams := query.NewQueryParams([]string{"original_url"})
	queryParams.Parse(c, "10")
	err := queryParams.Validate(CustomValidator.ParamValidator{
		MaxSearchLength:       100,
		AllowedOrderByColumns: []string{"created_at", "original_url"},
		MaxPageSize:           100,
	})

	if err != nil {
		errMsg := err.Error()
		c.JSON(400, app.NewErrorResponse("Validation Error", &errMsg))
		return
	}

	res, errApi := h.useCase.GetAllShortenerLink(queryParams)
	if errApi != nil {
		errMsg := errApi.Error()
		c.JSON(errApi.Code(), app.NewErrorResponse("Failed to get all shorten link", &errMsg))
		return
	}

	c.JSON(200, app.NewPaginationResponse("All shorten link retrieved successfully", res.Meta, res.Data))
}