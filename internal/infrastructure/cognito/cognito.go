package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/Zeta-Manu/manu-auth/internal/config"
)

type CognitoClient struct {
	client *cognitoidentityprovider.CognitoIdentityProvider
}

func NewCognitoClient(appConfig *config.AppConfig) (*CognitoClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(appConfig.AWS.AWSRegion),
	})
	if err != nil {
		return nil, err
	}

	client := cognitoidentityprovider.New(sess)

	return &CognitoClient{client: client}, nil
}
