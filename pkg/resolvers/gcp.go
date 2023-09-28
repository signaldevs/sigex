package sigex

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type GcpResolver struct{}

const gcpPrefix = "sigex-secret-gcp://"

// CanResolve detect if value can be resolved by this resolver
func (gr GcpResolver) CanResolve(value string) bool {
	return strings.HasPrefix(value, gcpPrefix)
}

// Resolve gets the plaintext version of a
// secret in GCP Secrets Manager
func (gr GcpResolver) Resolve(input string) (string, error) {
	// remove prefix:
	name := strings.ReplaceAll(input, gcpPrefix, "")

	// examples:
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to create secretmanager client: %v", err))
	}
	defer func(client *secretmanager.Client) {
		err := client.Close()
		if err != nil {
			log.Println(fmt.Errorf("error closing secret manager client: %v", err))
		}
	}(client)

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
	return string(result.Payload.Data), nil
}
