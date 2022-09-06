package sigex

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"regexp"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var secretRegex *regexp.Regexp

// IsSecretToken checks the value of a string to see if it's
// a tokenized sigex secret.
//
// Example: sigex-secret-gcp://projects/my-project/secrets/my-secret/versions/latest
func IsSecretToken(token string) bool {
	matches := secretRegex.MatchString(token)
	return matches
}

// ResolveSecretToken takes a secret token and contacts the corresponding
// secrets manager platform to return the plaintext version of the
// secret
func ResolveSecret(token string) string {
	parts := secretRegex.FindStringSubmatch(token)
	if len(parts) < 3 {
		log.Fatalln("secret token in incorrect format: ", token)
	}

	secretPlatform := parts[1]
	secretToken := parts[2]

	var secret string

	switch secretPlatform {
	case "gcp":
		secret = GetGCPSecretVersion(secretToken)
	case "aws":
		secret = GetAWSSecretVersion(secretToken)
	default:
		log.Fatalln("unsupported secret platform: " + secretPlatform)
	}

	return secret
}

// GetGCPSecretVersion gets the plaintext version of a
// secret in GCP Secrets Manager
func GetGCPSecretVersion(name string) string {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to create secretmanager client: %v", err))
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to access secret version: %v", err))
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		log.Fatalln(fmt.Errorf("data corruption detected in secret version"))
	}

	// WARNING: Do not print the secret in a production environment

	return string(result.Payload.Data)
}

// GetAWSSecretVersion gets the plaintext version of a
// secret in AWS Secrets Manager
func GetAWSSecretVersion(name string) string {
	log.Fatalln("aws secrets not yet implemented")
	return ""
}

func init() {
	secretRegex, _ = regexp.Compile(`^sigex-secret-(.*)\:\/\/(.*)$`)
}
