package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/omarelweshy/EcomMaster-user-service/internal/form"
	"github.com/omarelweshy/EcomMaster-user-service/internal/model"
	"github.com/omarelweshy/EcomMaster-user-service/internal/service"
	"github.com/omarelweshy/EcomMaster-user-service/internal/utils"
)

type UserHandler struct {
	UserService *service.UserService
}

func (h *UserHandler) Register(c *gin.Context) {
	var form form.RegistrationForm

	if err := c.ShouldBind(&form); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utils.FormatValidationError(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{"errors": formattedErrors})
			return
		}
	}
	err := h.UserService.RegisterUser(form.Username, form.Email, form.Password, form.FirstName, form.LastName)
	if err != nil {
		switch err {
		case service.ErrUsernameTaken:
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
			return
		case service.ErrEmailRegistered:
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
			return
		default:
			validationErrors := utils.FormatValidationError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"errors": validationErrors})
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "registeration successful"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var form form.LoginForm

	if err := c.ShouldBind(&form); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utils.FormatValidationError(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{"errors": formattedErrors})
			return
		}
	}
	user, err := h.UserService.LoginUser(form.Username, form.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			validationErrors := utils.FormatValidationError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"errors": validationErrors})

		}
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "login successful", "token": token})
}

func (h *UserHandler) Profile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no username found"})
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username format"})
		return
	}

	user, err := h.UserService.Repo.GetUserByUsername(usernameStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no username found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
		"email":     user.Email,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no username found"})
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username format"})
		return
	}
	var form form.UpdateUserForm
	if err := c.ShouldBind(&form); err != nil {
		validationErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	updatedUser := model.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
	}

	err := h.UserService.Repo.UpdateUser(usernameStr, &updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user profile"})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
