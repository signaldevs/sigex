package sigex

import (
	"log"
	"strings"
)

type AwsResolver struct{}

const awsPrefix = "sigex-secret-aws://"

// CanResolve detect if value can be resolved by this resolver
func (gr AwsResolver) CanResolve(value string) bool {
	return strings.HasPrefix(value, awsPrefix)
}

// Resolve gets the plaintext version of a
// secret in AWS Secrets Manager
func (gr AwsResolver) Resolve(_ string) (string, error) {
	log.Fatalln("aws secrets not yet implemented")
	// TODO: Will it get this far?
	return "", nil
}
