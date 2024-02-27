package controller

import (
	"net/http"

	"github.com/Zeta-Manu/manu-auth/internal/adapter/idp"
	"github.com/Zeta-Manu/manu-auth/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	idpAdapter idp.CognitoAdapter
}

func NewUserController(idpAdapter idp.CognitoAdapter) *UserController {
	return &UserController{
		idpAdapter: idpAdapter,
	}
}

func (uc *UserController) SignUp(c *gin.Context) {
	var userRegistration entity.UserRegistration
	if err := c.ShouldBindJSON(&userRegistration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := uc.idpAdapter.Register(c, userRegistration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data": result,
	}

	c.JSON(http.StatusOK, response)
}
