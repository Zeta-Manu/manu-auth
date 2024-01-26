package service

import (
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/Zeta-Manu/manu-auth/internal/config"
	"github.com/Zeta-Manu/manu-auth/internal/domain"
	"github.com/Zeta-Manu/manu-auth/internal/infrastructure/cognito"
)

// UserService represents the user service
type UserService struct {
	cognitoClient *cognito.CognitoClient
	appConfig     *config.AppConfig
}

// NewUserService creates a new user service
func NewUserService(appConfig *config.AppConfig) (*UserService, error) {
	cognitoClient, err := cognito.NewCognitoClient(appConfig)
	if err != nil {
		return nil, err
	}

	return &UserService{cognitoClient: cognitoClient, appConfig: appConfig}, nil
}

// RegisterUserHandler handles user registration
func (u *UserService) RegisterUserHandler(user domain.UserRegistration) error {
	// Use u.cognitoClient to register a user
	return u.cognitoClient.RegisterUser(u.appConfig.AWS.CognitoUserPoolID, u.appConfig.AWS.CognitoClientID, user)
}

// LoginUserHandler handles user login
func (u *UserService) LoginUserHandler(user domain.UserLogin) (*cognitoidentityprovider.AuthenticationResultType, error) {
	// Use u.cognitoClient to log in a user
	return u.cognitoClient.LoginUser(u.appConfig.AWS.CognitoUserPoolID, u.appConfig.AWS.CognitoClientID, user)
}
