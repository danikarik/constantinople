package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	permissions "github.com/xyproto/permissions2"
)

// UserCred contains user's signed xml and email address.
type UserCred struct {
	// SignedXML is signed doc retrieved from NCALayer payload.
	SignedXML string `json:"signedXml"`
	// EmailAddress is user's email.
	EmailAddress string `json:"emailAddress,omitempty"`
}

// Bind encodes request body into user credentials.
func (uc *UserCred) Bind(r *http.Request) error {
	if uc.SignedXML == "" {
		return errors.New("auth: no bind")
	}
	return nil
}

// Options is a configuration container to setup the AUTH router.
type Options struct {
	// Hostname is redis database address.
	Hostname string
	// Password is redis database password, if "requirepass" is specified
	// in redis.conf
	Password string
}

// Auth is a container used for userstate.
type Auth struct {
	userstate *permissions.UserState
	perm      *permissions.Permissions
}

var (
	// ErrRedisConn raises if Redis conn has wrong credentials or non-reachable.
	ErrRedisConn = errors.New("auth: redis connection failed")
	defaultHost  = "127.0.0.1:6379"
)

// New creates a new Auth handler with the provided options.
func New(options Options) (*Auth, error) {
	var (
		err       error
		userstate *permissions.UserState
		perm      *permissions.Permissions
		hostname  string
	)

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
		perm = permissions.NewPermissions(userstate)
		perm.Clear()
	}

	return &Auth{
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
	fmt.Fprintf(w, "Has user bob: %v\n", a.userstate.HasUser("bob"))
	fmt.Fprintf(w, "Logged in on server: %v\n", a.userstate.IsLoggedIn("bob"))
	fmt.Fprintf(w, "Is confirmed: %v\n", a.userstate.IsConfirmed("bob"))
	fmt.Fprintf(w, "Username stored in cookies (or blank): %v\n", a.userstate.Username(r))
	fmt.Fprintf(w, "Current user is logged in, has a valid cookie and *user rights*: %v\n", a.userstate.UserRights(r))
	fmt.Fprintf(w, "Current user is logged in, has a valid cookie and *admin rights*: %v\n", a.userstate.AdminRights(r))
	fmt.Fprintf(w, "\nTry: /register, /confirm, /remove, /login, /logout, /makeadmin, /clear, /data and /admin")
}

func (a *Auth) loginHandlerfunc(w http.ResponseWriter, r *http.Request) {
	data := &UserCred{}
	if err := render.Bind(r, data); err != nil {
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"status": "error"})
		return
	}
	render.Status(r, 200)
	render.JSON(w, r, map[string]string{"status": "success"})
}

func (a *Auth) logoutHandlerfunc(w http.ResponseWriter, r *http.Request) {

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
