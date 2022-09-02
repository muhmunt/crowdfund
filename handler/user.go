package handler

import (
	"fmt"
	"go_crowdfund/auth"
	"go_crowdfund/helper"
	"go_crowdfund/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(http.StatusUnprocessableEntity, "Register Failed!", "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	createUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse(http.StatusUnprocessableEntity, "Register Failed!", "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(createUser.ID)

	if err != nil {
		response := helper.APIResponse(http.StatusUnprocessableEntity, "Register Failed!", "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(createUser, token)
	response := helper.APIResponse(http.StatusOK, "Account successfully registered", "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(http.StatusUnprocessableEntity, "Login Failed!", "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"error": err.Error()}

		response := helper.APIResponse(http.StatusUnprocessableEntity, "Login Failed!", "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)

	if err != nil {
		response := helper.APIResponse(http.StatusUnprocessableEntity, "Login Failed!", "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token)
	response := helper.APIResponse(http.StatusOK, "Login successfully", "success", formatter)

	c.JSON(http.StatusOK, response)
}

