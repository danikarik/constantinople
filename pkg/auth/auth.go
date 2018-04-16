package auth

import (
	"errors"
	"net/http"

	"github.com/danikarik/constantinople/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	permissions "github.com/xyproto/permissions2"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	// ErrRedisConn raises if Redis conn has wrong credentials or non-reachable.
	ErrRedisConn = errors.New("auth: failed redis connection")
	// ErrEmptyAddress raises if PKI service address is not specified or empty.
	ErrEmptyAddress = errors.New("auth: empty pki address")
	defaultHost     = "127.0.0.1:6379"
	validate        *validator.Validate
)

// UserCred contains user's signed xml and email address.
type UserCred struct {
	// SignedXML is signed doc retrieved from NCALayer payload.
	SignedXML string `json:"signedXml" validate:"required"`
	// EmailAddress is user's email.
	EmailAddress string `json:"emailAddress" validate:"email"`
}

// Bind encodes request body into user credentials.
func (uc *UserCred) Bind(r *http.Request) error {
	return util.Validate(uc, validate)
}

// Options is a configuration container to setup the AUTH router.
type Options struct {
	// PKIAddress is a address of signature validation service.
	PKIAddress string
	// Hostname is redis database address.
	Hostname string
	// Password is redis database password, if "requirepass" is specified
	// in redis.conf
	Password string
	// Validate is used to struct schema validation.
	Validate *validator.Validate
	// Debug mode.
	Debug bool
}

// Auth is a container used for userstate.
type Auth struct {
	userstate *permissions.UserState
	perm      *permissions.Permissions
	addr      string
	debug     bool
}

// New creates a new Auth handler with the provided options.
func New(options Options) (*Auth, error) {
	var (
		err       error
		userstate *permissions.UserState
		hostname  string
	)
	if options.PKIAddress == "" {
		return nil, ErrEmptyAddress
	}
	if options.Validate == nil {
		validate = validator.New()
	}
	if options.Hostname != "" {
		hostname = options.Hostname
	} else {
		hostname = defaultHost
	}
	if options.Password != "" {
		if userstate, err = permissions.NewUserStateWithPassword2(
			hostname,
			options.Password,
		); err != nil {
			return nil, ErrRedisConn
		}
	} else {
		if userstate, err = permissions.NewUserState2(
			0,
			true,
			hostname,
		); err != nil {
			return nil, ErrRedisConn
		}
	}
	return &Auth{
		addr:      options.PKIAddress,
		userstate: userstate,
		perm:      permissions.NewPermissions(userstate),
		debug:     options.Debug,
	}, nil
}

// Router groups all auth handlers.
func (a *Auth) Router(pattern string) (string, http.Handler) {
	a.perm.Clear()
	a.perm.SetPublicPath([]string{"/", pattern + "/login"})
	a.perm.SetUserPath([]string{pattern + "/session", pattern + "/logout"})
	a.perm.SetDenyFunction(denyHandler)
	r := chi.NewRouter()
	r.Use(a.Middleware())
	r.Get("/session", a.sessionHandler)
	r.Post("/login", a.loginHandler)
	r.Delete("/logout", a.logoutHandler)
	return pattern, r
}

func (a *Auth) sessionHandler(w http.ResponseWriter, r *http.Request) {
	a.Session(w, r)
}

func (a *Auth) loginHandler(w http.ResponseWriter, r *http.Request) {
	data := &UserCred{}
	if err := render.Bind(r, data); err != nil {
		util.BadRequest(w, r, err)
		return
	}
	a.Login(w, r, data)
}

func (a *Auth) logoutHandler(w http.ResponseWriter, r *http.Request) {
	a.Logout(w, r)
}

func denyHandler(w http.ResponseWriter, r *http.Request) {
	util.ErrorStatus(w, r, errors.New("permission denied"), http.StatusForbidden)
}

// Middleware checks session authentication.
func (a *Auth) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the user has the right admin/user rights
			if a.perm.Rejected(w, r) {
				// Deny the request
				a.perm.DenyFunction()(w, r)
				return
			}
			// Serve the requested page
			next.ServeHTTP(w, r)
		})
	}
}
