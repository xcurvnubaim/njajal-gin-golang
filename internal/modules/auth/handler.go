package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/middleware"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/app"
	CustomValidator "github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/validator"
)

type AuthHandler struct {
	authUseCase IAuthUseCase
	app         *gin.Engine
}

func NewAuthHandler(app *gin.Engine, authUseCase IAuthUseCase, prefixApi string) {
	authHandler := &AuthHandler{
		app:         app,
		authUseCase: authUseCase,
	}

	authHandler.Routes(prefixApi)
}

func (ah *AuthHandler) Routes(prefix string) {
	authentication := ah.app.Group(prefix)
	{
		authentication.POST("/register", ah.Register)
		authentication.POST("/login", ah.Login)

		authentication.Use(middleware.AuthenticateJWT())
		{
			authentication.GET("/me", ah.GetMe)
			// authentication.GET("/all", ah.GetAllUsers)
			// authentication.GET("/username/:username", ah.GetUserByUsername)
		}
	}
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var authentication RegisterUserRequestDTO
	if err := c.ShouldBindJSON(&authentication); err != nil {
		var errMessages = CustomValidator.FormatValidationErrors(err)
		c.JSON(400, app.NewErrorResponse("Validation Error", &errMessages))
		return
	}

	// Register User
	if err := ah.authUseCase.RegisterUser(&authentication); err != nil {
		errMsg := err.Error()
		c.JSON(err.Code(), app.NewErrorResponse("Failed to register user", &errMsg))
		return
	}

	c.JSON(200, app.NewSuccessResponse("User registered successfully", &RegisterUserResponseDTO{
		Username: authentication.Username,
		PhoneNumber: authentication.PhoneNumber,
	}))
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var authentication LoginUserRequestDTO
	if err := c.ShouldBindJSON(&authentication); err != nil {
		var errMessages = CustomValidator.FormatValidationErrors(err)
		c.JSON(400, app.NewErrorResponse("Validation Error", &errMessages))
		return
	}

	// Login User
	token, err := ah.authUseCase.LoginUser(&authentication)
	if err != nil {
		errMsg := err.Error()
		c.JSON(err.Code(), app.NewErrorResponse("Failed to login user", &errMsg))
		return
	}
	
	c.JSON(200, app.NewSuccessResponse("User logged in successfully", &token))
}

func (ah *AuthHandler) GetMe(c *gin.Context) {
	// Assuming userID is extracted from the context or request
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(400, app.NewErrorResponse("User ID not found in context", nil))
		return
	}

	// Perform a type assertion to ensure userID is a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(400, app.NewErrorResponse("Invalid user ID type", nil))
		return
	}

	// Convert string to UUID
	parsedID, errUuid := uuid.Parse(userIDStr)
	if errUuid != nil {
		c.JSON(400, app.NewErrorResponse("Invalid user ID", nil))
		return
	}

	// Now call GetMe with the UUID
	user, err := ah.authUseCase.GetMe(parsedID)
	if err != nil {
		// Assuming err has a method Code() to retrieve HTTP status code
		var errMsg = err.Error()
		c.JSON(err.Code(), app.NewErrorResponse("Failed to get user data", &errMsg))
		return
	}

	c.JSON(200, app.NewSuccessResponse("User data retrieved successfully", user))
}

func (ah *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := ah.authUseCase.GetAllUser()
	if err != nil {
		errMsg := err.Error()
		c.JSON(err.Code(), app.NewErrorResponse("Failed to get all users", &errMsg))
		return
	}

	c.JSON(200, app.NewSuccessResponse("All users retrieved successfully", users))
}

func (ah *AuthHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := ah.authUseCase.GetUserByUsername(username)
	if err != nil {
		errMsg := err.Error()
		c.JSON(err.Code(), app.NewErrorResponse("Failed to get user data", &errMsg))
		return
	}

	c.JSON(200, app.NewSuccessResponse("User data retrieved successfully", user))
}
