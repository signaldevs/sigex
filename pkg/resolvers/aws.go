package sigex

import (
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"strings"
)

type AwsResolver struct{}

const awsPrefix = "sigex-secret-aws://"

// CanResolve detect if value can be resolved by this resolver
func (gr AwsResolver) CanResolve(value string) bool {
	return strings.HasPrefix(value, awsPrefix)
}

var (
	secretCache, _ = secretcache.New()
)

// Resolve gets the plaintext version of a
// secret in AWS Secrets Manager
func (gr AwsResolver) Resolve(input string) (string, error) {
	awsKey := strings.ReplaceAll(input, awsPrefix, "")
	return secretCache.GetSecretString(awsKey)
}
