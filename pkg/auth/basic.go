package auth

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/danikarik/constantinople/pkg/util"
)

const (
	basic        = "basic"
	defaultRealm = "Restricted"
)

var (
	// DefaultBasicAuthOptions is the default BasicAuth middleware config.
	DefaultBasicAuthOptions = BasicAuthOptions{
		Realm: defaultRealm,
	}
)

// BasicAuthOptions defines the config for BasicAuth middleware.
type BasicAuthOptions struct {
	Validator BasicAuthValidator
	Realm     string
}

// BasicAuthValidator defines a function to validate BasicAuth credentials.
type BasicAuthValidator func(string, string) bool

// BasicAuth is a basic auth middleware.
func BasicAuth(fn BasicAuthValidator) func(next http.Handler) http.Handler {
	opts := DefaultBasicAuthOptions
	opts.Validator = fn
	return BasicAuthWithOptions(opts)
}

// BasicAuthWithOptions returns an BasicAuth middleware with config.
func BasicAuthWithOptions(opts BasicAuthOptions) func(next http.Handler) http.Handler {
	if opts.Realm == "" {
		opts.Realm = defaultRealm
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				util.ErrorStatus(w, r, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}
			if ok = opts.Validator(u, p); !ok {
				util.ErrorStatus(w, r, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}
			realm := defaultRealm
			if opts.Realm != defaultRealm {
				realm = strconv.Quote(opts.Realm)
			}
			w.Header().Set("WWW-Authenticate", basic+" realm="+realm)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
