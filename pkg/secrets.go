package sigex

import (
	"log"

	sigex "github.com/signaldevs/sigex/pkg/resolvers"
)

var awsResolver, gcpResolver, rot13Resolver, defaultResolver sigex.Resolver

// ResolveSecret takes a secret token and contacts the corresponding
// secrets manager platform to return the plaintext version of the
// secret
func ResolveSecret(token string) string {
	resolvers := []sigex.Resolver{gcpResolver, awsResolver, rot13Resolver, defaultResolver}

	for _, resolver := range resolvers {
		if resolver.CanResolve(token) {
			resolved, err := resolver.Resolve(token)
			if err != nil {
				log.Fatalln(err)
			}
			return resolved
		}
	}
	log.Fatalln("unsupported secret platform in token: " + token)
	return ""
}

func init() {
	gcpResolver = sigex.GcpResolver{}
	awsResolver = sigex.AwsResolver{}
	rot13Resolver = sigex.Rot13Resolver{}
	defaultResolver = sigex.DefaultResolver{}
}
