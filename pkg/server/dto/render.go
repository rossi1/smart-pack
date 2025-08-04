package dto

import (
	"net/http"

	"github.com/go-chi/render"
)

// Read is for reading request body to dto struct.
func Read(r *http.Request, v interface{}) error {
	return render.Decode(r, v)
}

// Write is for writing response body from a dto struct.
func Write(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.Respond(w, r, v)
}
