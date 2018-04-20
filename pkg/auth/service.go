package auth

import (
	"errors"
	"net/http"

	pb "github.com/danikarik/constantinople/pkg/proto"
	"github.com/danikarik/constantinople/pkg/util"
	servertiming "github.com/mitchellh/go-server-timing"
	"github.com/xyproto/cookie"
)

const (
	sessionIDKey = "sessionid"
)

// Login validates signature signed by SIGN key.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request, data *UserCred) {
	var user, email, code string

	if user = a.userstate.Username(r); user != "" {
		if sid, ok := cookie.SecureCookie(
			r,
			sessionIDKey,
			a.userstate.CookieSecret(),
		); ok {
			if a.userstate.CorrectPassword(user, sid) {
				a.Session(w, r)
				return
			}
		}
	}

	timing := servertiming.FromContext(r.Context())
	m := timing.NewMetric("grpc").WithDesc("PKI signature validation").Start()
	conn, err := pb.New(r.Context(), a.addr)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}
	defer conn.Close()
	client := pb.NewMetropolisServiceClient(conn)
	resp, err := client.VerifySignature(
		r.Context(),
		&pb.VerRequest{
			SignedXml: data.SignedXML,
			Flag:      pb.VerFlag_AUTH,
		},
	)
	m.Stop()

	if resp.Status != pb.VerStatus_SUCCESS {
		err := errors.New(resp.Message)
		util.BadRequestWith(w, r, err, resp.Description)
		return
	}

	user = resp.Message
	email = data.EmailAddress

	ok, err := a.userstate.HasUser2(user)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}
	if ok {
		a.userstate.RemoveUser(user)
	}

	code, err = a.userstate.GenerateUniqueConfirmationCode()
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}

	a.userstate.AddUser(user, code, email)
	cookie.SetSecureCookiePathWithFlags(
		w,
		sessionIDKey,
		code,
		a.userstate.CookieTimeout(user),
		"/",
		a.userstate.CookieSecret(),
		false,
		true,
	)

	a.userstate.Login(w, user)
	util.OK(w, r)
}

// Logout remove cookie and set logout out.
func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	u := a.userstate.Username(r)
	a.userstate.Logout(u)
	a.userstate.ClearCookie(w)
	a.userstate.RemoveUser(u)
	util.OK(w, r)
}

// Session returns user session data.
func (a *Auth) Session(w http.ResponseWriter, r *http.Request) {
	user := a.userstate.Username(r)
	email, err := a.userstate.Email(user)
	if err != nil {
		util.Error("session: %s", err.Error())
	}
	util.OKWith(w, r, map[string]string{
		"serialNumber": user,
		"emailAddress": email,
	})
}
