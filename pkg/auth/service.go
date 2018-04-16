package auth

import (
	"net/http"

	"github.com/danikarik/constantinople/pkg/util"
	"github.com/golang/glog"
)

// Login validates signature signed by SIGN key.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request, data *UserCred) {
	// timing := servertiming.FromContext(r.Context())
	var username, email string

	if username = a.userstate.Username(r); username != "" {
		a.Session(w, r)
		return
	}

	if a.userstate.IsLoggedIn(username) {
		util.Debug("login: is already logged in: %s", username)
		util.Debug("login: set logged out: %s", username)
		a.userstate.SetLoggedOut(username)
	} else {
		util.Debug("login: is logged in: %s", a.userstate.IsLoggedIn(username))
	}

	// m := timing.NewMetric("grpc").WithDesc("PKI signature validation").Start()
	// conn, err := pb.New(r.Context(), a.addr)
	// if err != nil {
	// 	util.BadRequest(w, r, err)
	// 	return
	// }
	// defer conn.Close()
	// client := pb.NewMetropolisServiceClient(conn)
	// resp, err := client.VerifySignature(
	// 	r.Context(),
	// 	&pb.VerRequest{
	// 		SignedXml: data.SignedXML,
	// 		Flag:      pb.VerFlag_AUTH,
	// 	},
	// )
	// m.Stop()
	// if resp.Status != pb.VerStatus_SUCCESS {
	// 	err := errors.New(resp.Message)
	// 	util.BadRequestWith(w, r, err, resp.Description)
	// 	return
	// }

	// username = resp.Message
	username = "111111111"
	email = data.EmailAddress

	ok, err := a.userstate.HasUser2(username)
	if err != nil {
		util.Error("login: has user: %s", err.Error())
		util.BadRequest(w, r, err)
		return
	}
	if !ok {
		util.Debug("login: add user: %s", username)
		a.userstate.AddUser(username, "", email)
	} else {
		stateEmail, err := a.userstate.Email(username)
		if err != nil {
			util.Error("login: get email: %s", err.Error())
		} else {
			if email != stateEmail {
				util.Debug("login: remove user: %s", username)
				a.userstate.RemoveUser(username)
				util.Debug("login: re-add user: %s", username)
				a.userstate.AddUser(username, "", email)
			}
		}
	}

	if glog.V(3) {
		a.userstate.Login(w, username)
		util.OK(w, r)
		return
	}

	a.userstate.SetLoggedIn(username)
	if err := a.userstate.SetUsernameCookieOnlyHTTPS(w, username); err != nil {
		util.BadRequest(w, r, err)
		return
	}

	util.OK(w, r)
}

// Logout remove cookie and set logout out.
func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	a.userstate.Logout(a.userstate.Username(r))
	a.userstate.ClearCookie(w)
	util.OK(w, r)
}

// Session returns user session data.
func (a *Auth) Session(w http.ResponseWriter, r *http.Request) {
	username := a.userstate.Username(r)
	email, err := a.userstate.Email(username)
	if err != nil {
		util.Error("session: %s", err.Error())
	}
	util.OKWith(w, r, map[string]string{
		"serialNumber": username,
		"emailAddress": email,
	})
}
