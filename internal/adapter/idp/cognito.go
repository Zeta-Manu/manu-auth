package idp

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/Zeta-Manu/manu-auth/internal/domain/entity"
	"github.com/Zeta-Manu/manu-auth/pkg/utils"
)

type CognitoAdapter struct {
	client   *cip.Client
	poolID   string
	clientID string
}

func NewCognitoAdapter(accessKey, secretAccessKey, poolID, clientID, region string) (*CognitoAdapter, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return &CognitoAdapter{
		client:   cip.NewFromConfig(cfg),
		poolID:   poolID,
		clientID: clientID,
	}, nil
}

func (a *CognitoAdapter) Register(ctx context.Context, userRegistration entity.UserRegistration) (string, error) {
	attributes := []types.AttributeType{
		{
			Name:  aws.String("name"),
			Value: aws.String(userRegistration.Name),
		},
		{
			Name:  aws.String("email"),
			Value: aws.String(userRegistration.Email),
		},
	}

	params := &cip.SignUpInput{
		ClientId:       aws.String(a.clientID),
		Username:       aws.String(userRegistration.Email),
		Password:       aws.String(userRegistration.Password),
		UserAttributes: attributes,
	}

	result, err := a.client.SignUp(ctx, params)
	if err != nil {
		var invalidPasswordErr *types.InvalidPasswordException
		var invalidParameterErr *types.InvalidParameterException
		var usernameExistsErr *types.UsernameExistsException

		switch {
		case errors.As(err, &invalidPasswordErr):
			return "", &utils.CustomError{
				Message: "Invalid password",
				Status:  http.StatusBadRequest,
			}
		case errors.As(err, &invalidParameterErr):
			return "", &utils.CustomError{
				Message: "Invalid parameter",
				Status:  http.StatusBadRequest,
			}
		case errors.As(err, &usernameExistsErr):
			return "", &utils.CustomError{
				Message: "Username already exists",
				Status:  http.StatusConflict,
			}
		default:
			return "", &utils.CustomError{
				Message: "Failed to sign up user",
				Status:  http.StatusInternalServerError,
			}
		}
	}
	return *result.CodeDeliveryDetails.Destination, nil
}

func (a *CognitoAdapter) Login(ctx context.Context, userLogin entity.UserLogin) (*entity.LoginResult, error) {
	params := &cip.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": userLogin.Email,
			"PASSWORD": userLogin.Password,
		},
		ClientId: aws.String(a.clientID),
	}

	result, err := a.client.InitiateAuth(ctx, params)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Failed to initate auth",
			Status:  http.StatusInternalServerError,
		}
	}

	loginReturn := entity.LoginResult{
		AccessToken:  result.AuthenticationResult.AccessToken,
		ExpiresIn:    &result.AuthenticationResult.ExpiresIn,
		IdToken:      result.AuthenticationResult.IdToken,
		RefreshToken: result.AuthenticationResult.RefreshToken,
		TokenType:    result.AuthenticationResult.TokenType,
	}

	return &loginReturn, nil
}

func (a *CognitoAdapter) ConfirmRegistration(ctx context.Context, userRegistrationConfirm entity.UserRegistrationConfirm) error {
	params := &cip.ConfirmSignUpInput{
		ClientId:         aws.String(a.clientID),
		ConfirmationCode: aws.String(userRegistrationConfirm.ConfirmationCode),
		Username:         aws.String(userRegistrationConfirm.Email),
	}

	_, err := a.client.ConfirmSignUp(ctx, params)
	if err != nil {
		return &utils.CustomError{
			Message: "Failed to confirm sign up",
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (a *CognitoAdapter) ResendConfirmationCode(ctx context.Context, email entity.Email) (string, error) {
	params := &cip.ResendConfirmationCodeInput{
		ClientId: aws.String(a.clientID),
		Username: aws.String(email.Email),
	}

	result, err := a.client.ResendConfirmationCode(ctx, params)
	if err != nil {
		return "", &utils.CustomError{
			Message: "Failed to resend confirmation code",
			Status:  http.StatusInternalServerError,
		}
	}

	return *result.CodeDeliveryDetails.Destination, nil
}

func (a *CognitoAdapter) ForgotPassword(ctx context.Context, email entity.Email) (string, error) {
	params := &cip.ForgotPasswordInput{
		ClientId: aws.String(a.clientID),
		Username: aws.String(email.Email),
	}

	result, err := a.client.ForgotPassword(ctx, params)
	if err != nil {
		return "", &utils.CustomError{
			Message: "Failed to initiate forgot password",
			Status:  http.StatusInternalServerError,
		}
	}

	return *result.CodeDeliveryDetails.Destination, nil
}

func (a *CognitoAdapter) ConfirmForgotPassword(ctx context.Context, userResetPassword entity.UserResetPassword) error {
	params := &cip.ConfirmForgotPasswordInput{
		ClientId:         aws.String(a.clientID),
		Username:         aws.String(userResetPassword.Email),
		ConfirmationCode: aws.String(userResetPassword.ConfirmationCode),
		Password:         aws.String(userResetPassword.NewPassword),
	}

	_, err := a.client.ConfirmForgotPassword(ctx, params)
	if err != nil {
		return &utils.CustomError{
			Message: "Failed to confirm forgot password",
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (a *CognitoAdapter) ChangePassword(ctx context.Context, accessToken string, changePassword entity.UserChangePassword) error {
	params := &cip.ChangePasswordInput{
		AccessToken:      aws.String(accessToken),
		PreviousPassword: aws.String(changePassword.PreviousPassword),
		ProposedPassword: aws.String(changePassword.ProposedPassword),
	}

	_, err := a.client.ChangePassword(ctx, params)
	if err != nil {
		return &utils.CustomError{
			Message: "Failed to change password",
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
