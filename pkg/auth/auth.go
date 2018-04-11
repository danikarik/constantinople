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
}

// Auth is a container used for userstate.
type Auth struct {
	userstate *permissions.UserState
	perm      *permissions.Permissions
	addr      string
}

// New creates a new Auth handler with the provided options.
func New(options Options) (*Auth, error) {

	if options.PKIAddress == "" {
		return nil, ErrEmptyAddress
	}

	var (
		err       error
		userstate *permissions.UserState
		perm      *permissions.Permissions
		hostname  string
	)

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

	perm = permissions.NewPermissions(userstate)
	perm.Clear()

	return &Auth{
		addr:      options.PKIAddress,
		userstate: userstate,
		perm:      perm,
	}, nil
}

// Router groups all auth handlers.
func (a *Auth) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", a.sessionHandlerfunc)
	r.Post("/", a.loginHandlerfunc)
	r.Delete("/", a.logoutHandlerfunc)
	return r
}

func (a *Auth) sessionHandlerfunc(w http.ResponseWriter, r *http.Request) {
	a.Session(w, r)
}

func (a *Auth) loginHandlerfunc(w http.ResponseWriter, r *http.Request) {
	data := &UserCred{}
	if err := render.Bind(r, data); err != nil {
		util.BadRequest(w, r, err)
		return
	}
	a.Login(w, r, data)
}

func (a *Auth) logoutHandlerfunc(w http.ResponseWriter, r *http.Request) {
	a.Logout(w, r)
}

// Middleware checks session authentication.
// func Middleware() func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			// ...
// 			next.ServeHTTP(w, r)
// 		}
// 		return http.HandlerFunc(fn)
// 	}
// }
