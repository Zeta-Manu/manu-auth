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
		uc.logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		uc.logger.Info("User registered successfully", zap.String("Email", userRegistration.Email))
	}

	response := gin.H{
		"data": result,
	}

	c.JSON(http.StatusOK, response)
}

func (uc *UserController) ConfirmSignUp(c *gin.Context) {
	var userRegistrationConfirm entity.UserRegistrationConfirm
	if err := c.ShouldBindJSON(&userRegistrationConfirm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.idpAdapter.ConfirmRegistration(c, userRegistrationConfirm)
	if err != nil {
		var customErr *utils.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.Status, gin.H{"error": customErr.Message})
			uc.logger.Error("User confirm registration failed", zap.String("error", customErr.Message))
			return
		}
		uc.logger.Error("Failed to confirm user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		uc.logger.Info("User confirm successfully", zap.String("Email", userRegistrationConfirm.Email))
	}

	c.Status(http.StatusOK)
}

func (uc *UserController) ResendConfirmationCode(c *gin.Context) {
	var email entity.Email
	result, err := uc.idpAdapter.ResendConfirmationCode(c, email.Email)
	if err != nil {
		var customErr *utils.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.Status, gin.H{"error": customErr.Message})
			uc.logger.Error("User resend confirm registration code failed", zap.String("error", customErr.Message))
			return
		}
		uc.logger.Error("Failed to resend confirm registraion code", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		uc.logger.Info("User resend confirm successfully", zap.String("Email", email.Email))
	}

	response := gin.H{
		"data": result,
	}

	c.JSON(http.StatusOK, response)
}

func (uc *UserController) LogIn(c *gin.Context) {
	var userLogin entity.UserLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := uc.idpAdapter.Login(c, userLogin)
	if err != nil {
		var customErr *utils.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.Status, gin.H{"error": customErr.Message})
			uc.logger.Error("User login failed", zap.String("error", customErr.Message))
			return
		}
		uc.logger.Error("Failed to login user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		uc.logger.Info("User login successfully", zap.String("Email", userLogin.Email))
	}

	response := gin.H{
		"data": result,
	}

	c.JSON(http.StatusOK, response)
}

func (uc *UserController) ForgotPassword(c *gin.Context) {
	var email entity.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := uc.idpAdapter.ForgotPassword(c, email.Email)
	if err != nil {
		var customErr *utils.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.Status, gin.H{"error": customErr.Message})
			uc.logger.Error("User forgot password", zap.String("error", customErr.Message))
			return
		}
		uc.logger.Error("Failed to forgot password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		uc.logger.Info("User forgot password successfully", zap.String("Email", email.Email))
	}

	response := gin.H{
		"data": result,
	}

	c.JSON(http.StatusOK, response)
}

func (uc *UserController) ConfirmForgotPassword(c *gin.Context) {
	var userResetPassword entity.UserResetPassword
	if err := c.ShouldBindJSON(&userResetPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.idpAdapter.ConfirmForgotPassword(c, userResetPassword)
	if err != nil {
		var customErr *utils.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.Status, gin.H{"error": customErr.Message})
			uc.logger.Error("User confirm forgot password failed", zap.String("error", customErr.Message))
			return
		}
		uc.logger.Error("Failed to confirm forgot password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		uc.logger.Info("User confirm forgot password successfully", zap.String("Email", userResetPassword.Email))
	}

	c.Status(http.StatusOK)
}

func (uc *UserController) ChangePassword(c *gin.Context) {
	var userChangePassword entity.UserChangePassword
	if err := c.ShouldBindJSON(&userChangePassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := uc.idpAdapter.ChangePassword(c, token.(string), userChangePassword)
	if err != nil {
		var customErr *utils.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.Status, gin.H{"error": customErr.Message})
			uc.logger.Error("User change password failed", zap.String("error", customErr.Message))
			return
		}
		uc.logger.Error("Failed change password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		uc.logger.Info("User change password successfully")
	}

	c.Status(http.StatusOK)
}

func GetSub(c *gin.Context) {
	sub, exists := c.Get("sub")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Subject not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sub})
}
