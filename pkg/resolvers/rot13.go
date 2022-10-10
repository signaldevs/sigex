package sigex

import "strings"

type Rot13Resolver struct{}

const rot13Prefix = "sigex-secret-rot13://"

// CanResolve detect if value can be resolved by this resolver
func (gr Rot13Resolver) CanResolve(value string) bool {
	return strings.HasPrefix(value, rot13Prefix)
}

func (gr Rot13Resolver) Resolve(input string) (string, error) {
	rotated := strings.ReplaceAll(input, rot13Prefix, "")
	return strings.Map(rot13, rotated), nil
}

// this is from https://www.dotnetperls.com/rot13-go
func rot13(r rune) rune {
	if r >= 'a' && r <= 'z' {
		// Rotate lowercase letters 13 places.
		if r >= 'm' {
			return r - 13
		} else {
			return r + 13
		}
	} else if r >= 'A' && r <= 'Z' {
		// Rotate uppercase letters 13 places.
		if r >= 'M' {
			return r - 13
		} else {
			return r + 13
		}
	}
	// Do nothing.
	return r
}
