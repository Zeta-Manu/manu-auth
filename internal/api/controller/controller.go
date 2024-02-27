package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Zeta-Manu/manu-auth/internal/adapter/idp"
	"github.com/Zeta-Manu/manu-auth/internal/domain/entity"
	"github.com/Zeta-Manu/manu-auth/pkg/utils"
)

type UserController struct {
	logger     *zap.Logger
	idpAdapter idp.CognitoAdapter
}

func NewUserController(idpAdapter idp.CognitoAdapter, logger *zap.Logger) *UserController {
	return &UserController{
		idpAdapter: idpAdapter,
		logger:     logger,
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
		var customErr *utils.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.Status, gin.H{"error": customErr.Message})
			uc.logger.Error("User registration failed", zap.String("error", customErr.Message))
			return
		}
		uc.logger.Error("Faile to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		uc.logger.Info("User registered successfully", zap.String("Email", userRegistration.Email))
	}

	response := gin.H{
		"data": result,
	}

	c.JSON(http.StatusOK, response)
}
