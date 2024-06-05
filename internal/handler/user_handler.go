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
			utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", formattedErrors)
			return
		}
	}
	err := h.UserService.RegisterUser(form.Username, form.Email, form.Password, form.FirstName, form.LastName)
	if err != nil {
		switch err {
		case service.ErrUsernameTaken:
			utils.RespondWithError(c, http.StatusBadRequest, "Username already taken", nil)
			return
		case service.ErrEmailRegistered:
			utils.RespondWithError(c, http.StatusBadRequest, "Email already registered", nil)
			return
		default:
			utils.RespondWithError(c, http.StatusInternalServerError, "Registration failed", nil)
			return
		}
	}
	utils.RespondWithSuccess(c, "Registration successful", nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var form form.LoginForm

	if err := c.ShouldBind(&form); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utils.FormatValidationError(validationErrors)
			utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", formattedErrors)
			return
		}
	}
	user, err := h.UserService.LoginUser(form.Username, form.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid credentials", nil)
			return
		default:
			validationErrors := utils.FormatValidationError(err)
			utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", validationErrors)
			return
		}
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not generate token", nil)
		return
	}
	utils.RespondWithSuccess(c, "login successful", gin.H{"token": token})
}

func (h *UserHandler) Profile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "no username found", nil)
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid username format", nil)
		return
	}

	user, err := h.UserService.Repo.GetUserByUsername(usernameStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "no username found", nil)
		return
	}
	utils.RespondWithSuccess(c, "user data", gin.H{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
		"email":     user.Email,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "no username found", nil)
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid username format", nil)
		return
	}
	var form form.UpdateUserForm
	if err := c.ShouldBind(&form); err != nil {
		validationErrors := utils.FormatValidationError(err)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	updatedUser := model.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
	}

	err := h.UserService.Repo.UpdateUser(usernameStr, &updatedUser)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "could not update user profile", nil)
		return
	}
	utils.RespondWithSuccess(c, "Profile updated successfully", nil)
}
