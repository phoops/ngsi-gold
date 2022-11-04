package ldcontext

type LdContext []any

var (
	EmptyContext   LdContext = []any{}
	DefaultContext           = EmptyContext
)
