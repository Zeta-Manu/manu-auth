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

// @Summary		Sign up a new user
// @Description	Register a new user with email and password
// @Tags User
// @Accept			json
// @Produce		json
// @Param			body	body		entity.UserRegistration											true	"User registration info"
// @Success 200 {object} entity.ResponseWrapper
// @Failure 400 {object} entity.ErrorWrapper "Invalid Password or Invalid Parameter"
// @Failure 409 {object} entity.ErrorWrapper "Username Exists"
// @Failure 500 {object} entity.ErrorWrapper
// @Router			/signup [post]
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

// @Summary Confirm user registration
// @Description Confirm a user's registration using the provided confirmation information
// @Tags User
// @Accept json
// @Produce json
// @Param body body entity.UserRegistrationConfirm true "User registration confirmation info"
// @Success 200
// @Failure 400
// @Failure 408
// @Router /confirm [post]
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

// @Summary Resend confirmation code
// @Description Resend the confirmation code to the provided email address
// @Tags User
// @Accept json
// @Produce json
// @Param email body entity.Email true "Email address to resend the confirmation code to"
// @Success 200 {object} entity.ResponseWrapper
// @Failure 500 {object} entity.ErrorWrapper
// @Router /resend-confirm [post]
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

// @Summary		Log in with email and password
// @Description	Authenticate user with email and password
// @Tags User
// @Accept			json
// @Produce		json
// @Param			body	body		entity.UserLogin									true	"User login info"
// @Success 200 {object} entity.ResponseWrapper{data=entity.LoginResult}
// @Failure 400 {object} entity.ErrorWrapper "Invalid Password or Missing Parameter"
// @Failure 401 {object} entity.ErrorWrapper "Not Authorized"
// @Failure 403 {object} entity.ErrorWrapper "User Not Confirm"
// @Failure 404 {object} entity.ErrorWrapper "User Not Found"
// @Failure 500 {object} entity.ErrorWrapper
// @Router			/login [post]
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

// @Summary Forgot Password
// @Description Initiate the password reset process for a user by sending a reset link to their email
// @Tags User
// @Accept json
// @Produce json
// @Param email body entity.Email true "Email address of the user"
// @Success 200 {object} entity.ResponseWrapper
// @Failure 400 {object} entity.ErrorWrapper
// @Failure 500 {object} entity.ErrorWrapper
// @Router /forgot-password [post]
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

// @Summary Forgot Password
// @Description Initiate the password reset process for a user by sending a reset link to their email
// @Tags User
// @Accept json
// @Produce json
// @Param email body entity.UserResetPassword true "Email address of the user"
// @Success 200 {object} entity.ResponseWrapper
// @Failure 400 {object} entity.ErrorWrapper
// @Failure 500 {object} entity.ErrorWrapper
// @Router /confirm-forgot [post]
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

// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param body body entity.UserChangePassword true "User change password info"
// @Success  200
// @Failure  400 {object} entity.ErrorWrapper
// @Failure  401 {object} entity.ErrorWrapper "Not Authorized"
// @Failure  500 {object} entity.ErrorWrapper
// @Security BearerAuth
// @Router /change-password [post]
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

// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags User
// @Produce json
// @Param Authorization header string true "Bearer {token}" default(Bearer <Add access token here>)
// @Success 200 {object} entity.ResponseWrapper
// @Failure 401 {object} entity.ErrorWrapper "Not Authorized"
// @Failure 500 {object} entity.ErrorWrapper
// @Router /sub [get]
func GetSub(c *gin.Context) {
	sub, exists := c.Get("sub")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Subject not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sub})
}
