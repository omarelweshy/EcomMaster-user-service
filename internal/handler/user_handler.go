package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/omarelweshy/EcomMaster-user-service/internal/form"
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
	if err := h.UserService.RegisterUser(form.Username, form.Email, form.Password, form.FirstName, form.LastName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "registeration successful"})
}
