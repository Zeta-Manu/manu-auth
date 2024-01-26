package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/Zeta-Manu/manu-auth/internal/domain"
)

func (c *CognitoClient) LoginUser(userPoolID, clientID string, user domain.UserLogin) (*cognitoidentityprovider.AuthenticationResultType, error) {
	// Implement user login logic with Cognito
	input := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(clientID),
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(user.Email),
			"PASSWORD": aws.String(user.Password),
		},
	}

	result, err := c.client.InitiateAuth(input)
	if err != nil {
		return nil, err
	}

	return result.AuthenticationResult, nil
}
