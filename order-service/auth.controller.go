package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"otusHWUsers/config"
	db "otusHWUsers/db/sqlc"
	"otusHWUsers/models"
	"otusHWUsers/utils"
)

type AuthController struct {
	db  *db.Queries
}

func NewAuthController(db *db.Queries) *AuthController {
	return &AuthController{db}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var credentials *db.CreateUserParams

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword := utils.HashPassword(credentials.Password)

	args := &db.CreateUserParams{
		Name:      credentials.Name,
		Email:     credentials.Email,
		Password:  hashedPassword,
		Role:      "user",
		UpdatedAt: time.Now(),
	}

	user, err := ac.db.CreateUser(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	userResponse := models.FilteredResponse(user)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.db.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email"})
		return
	}

	if err := utils.ComparePassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid password"})
		return
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens
	time, _ := time.ParseDuration("720h")
	access_token, err := utils.CreateToken(time.Now() + time, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.SetCookie("access_token", access_token, 3600 * 24 , "/", "arch.homework", false, true)
	ctx.SetCookie("logged_in", "true", 3600 * 24 , "/", "arch.homework", false, false)

	ctx.Set("X-Token", access_token)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "arch.homework", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "arch.homework", false, true)
	ctx.Set("X-Token", "")

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
