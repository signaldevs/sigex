package sigex

type Resolver interface {
	Resolve(string) (string, error)
	CanResolve(string) bool
}
