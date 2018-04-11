package auth

import (
	"fmt"
	"net/http"

	"github.com/danikarik/constantinople/pkg/util"
	"github.com/golang/glog"
)

// Login validates signature signed by SIGN key.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request, data *UserCred) {
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
	// if resp.Status != pb.VerStatus_SUCCESS {
	// 	err := errors.New(resp.Message)
	// 	util.BadRequestWith(w, r, err, resp.Description)
	// 	return
	// }
	// username := resp.Message
	username := "920727300044"
	ok, err := a.userstate.HasUser2(username)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}
	if !ok {
		a.userstate.AddUser(username, "123456", data.EmailAddress)
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
	a.userstate.ClearCookie(w)
	a.userstate.Logout(a.userstate.Username(r))
	fmt.Println(a.userstate.IsLoggedIn("920727300044"))
	util.OK(w, r)
}

// Session returns user session data.
func (a *Auth) Session(w http.ResponseWriter, r *http.Request) {
	username := a.userstate.Username(r)
	email, err := a.userstate.Email(username)
	if err != nil {
		email = ""
	}
	util.OKWith(w, r, map[string]string{
		"serialNumber": username,
		"emailAddress": email,
	})
}
