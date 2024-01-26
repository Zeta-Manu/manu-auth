package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/Zeta-Manu/manu-auth/internal/domain"
)

func (c *CognitoClient) RegisterUser(userPoolID, clientID string, user domain.UserRegistration) error {
	// Implement user registration logic with Cognito
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(clientID),
		Username: aws.String(user.Email),
		Password: aws.String(user.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Name),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(user.Email),
			},
		},
	}

	_, err := c.client.SignUp(input)
	if err != nil {
		return err
	}

	return nil
}
