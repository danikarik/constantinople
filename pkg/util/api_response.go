package util

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	// SUCCESS returns status 200 - OK
	SUCCESS string = "SUCCESS"
	// ERROR return status 400 - Bad request
	ERROR string = "ERROR"
)

// APIResponse represents JSON body response.
type APIResponse struct {
	Status      string      `json:"status"`
	Message     string      `json:"message,omitempty"`
	Description string      `json:"description,omitempty"`
	Result      interface{} `json:"result,omitempty"`
}

// OK returns status 200 - OK.
func OK(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, APIResponse{
		Status: SUCCESS,
	})
}

// OKWith returns status 200 - OK with additional struct.
func OKWith(w http.ResponseWriter, r *http.Request, res interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, APIResponse{
		Status: SUCCESS,
		Result: res,
	})
}

// BadRequest returns status 400 - Bad request with error message.
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, APIResponse{
		Status:  ERROR,
		Message: err.Error(),
	})
}

// BadRequestWith returns status 400 - Bad request with error message and description.
func BadRequestWith(w http.ResponseWriter, r *http.Request, err error, desc string) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, APIResponse{
		Status:      ERROR,
		Message:     err.Error(),
		Description: desc,
	})
}

// ErrorStatus returns http error status with error message.
func ErrorStatus(w http.ResponseWriter, r *http.Request, err error, code int) {
	render.Status(r, code)
	render.JSON(w, r, APIResponse{
		Status:  ERROR,
		Message: err.Error(),
	})
}
