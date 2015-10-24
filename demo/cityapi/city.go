package city

import (
	"errors"
	"net/http"

	"golang.org/x/net/context"
)

// START OMIT
type key int

const cityKey key = 0

func NewContext(parent context.Context, city string) context.Context {
	return context.WithValue(parent, cityKey, city) // HL
}

func FromContext(parent context.Context) (string, bool) {
	city, ok := parent.Value(cityKey).(string) // HL
	return city, ok
}

func FromRequest(r *http.Request) (string, error) {
	res := r.FormValue("city")

	if res == "" {
		return "", errors.New("invalid city")
	}

	return res, nil
}

// END OMIT
