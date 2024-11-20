package utils_http

import (
	"encoding/json"
	"net/http"
)

func ParseBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}
