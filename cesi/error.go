package cesi

const (
	ErrMissingCredentials = Error("missing credentials")
	ErrInitConnection     = Error("error while initializing connection")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
