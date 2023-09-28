package sigex

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AwsResolver struct{}

const awsPrefix = "sigex-secret-aws://"

var manager *secretsmanager.SecretsManager

// CanResolve detect if value can be resolved by this resolver
func (gr AwsResolver) CanResolve(value string) bool {
	return strings.HasPrefix(value, awsPrefix)
}

// Resolve gets the plaintext version of a
// secret in AWS Secrets Manager
func (gr AwsResolver) Resolve(input string) (string, error) {

	if manager == nil {
		sess, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})

		if err != nil {
			return "", err
		}

		manager = secretsmanager.New(sess)
	}

	awsKey := strings.ReplaceAll(input, awsPrefix, "")

	secretInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(awsKey),
	}

	secretValue, err := manager.GetSecretValue(secretInput)

	if err != nil {
		return "", err
	}

	return *secretValue.SecretString, nil
}
