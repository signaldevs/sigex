package sigex

type DefaultResolver struct{}

// CanResolve detect if value can be resolved by this resolver
func (gr DefaultResolver) CanResolve(_ string) bool {
	return true
}

// Resolve returns the input
func (gr DefaultResolver) Resolve(value string) (string, error) {
	return value, nil
}
