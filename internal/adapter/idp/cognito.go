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

func (a *CognitoAdapter) Register(ctx context.Context, userRegistration entity.UserRegistration) (*types.CodeDeliveryDetailsType, error) {
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

	input := &cip.SignUpInput{
		ClientId:       aws.String(a.clientID),
		Username:       aws.String(userRegistration.Email),
		Password:       aws.String(userRegistration.Password),
		UserAttributes: attributes,
	}

	result, err := a.client.SignUp(ctx, input)
	if err != nil {
		var invalidPasswordErr *types.InvalidPasswordException
		var invalidParameterErr *types.InvalidParameterException
		var usernameExistsErr *types.UsernameExistsException

		switch {
		case errors.As(err, &invalidPasswordErr):
			return nil, &utils.CustomError{
				Message: "Invalid password",
				Status:  http.StatusBadRequest,
			}
		case errors.As(err, &invalidParameterErr):
			return nil, &utils.CustomError{
				Message: "Invalid parameter",
				Status:  http.StatusBadRequest,
			}
		case errors.As(err, &usernameExistsErr):
			return nil, &utils.CustomError{
				Message: "Username already exists",
				Status:  http.StatusConflict,
			}
		default:
			return nil, &utils.CustomError{
				Message: "Failed to sign up user",
				Status:  http.StatusInternalServerError,
			}
		}
	}
	return result.CodeDeliveryDetails, nil
}
