package middlewares

import (
	"context"
	"net/http"

	"github.com/3d0c/toto-config/pkg/helpers"
)

const (
	AppEngineCountryHeader = "X-Appengine-Country"
)

// GeoTarget is a middleware which gets user's country.
// It tries to get this information from X-AppEngine-country HTTP Header first
// If there is no such header provided, it looks up geoip database.
func GeoTarget(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		countryCode string
	)

	if cc := r.Header.Get(AppEngineCountryHeader); len(cc) == 2 {
		countryCode = cc
	} else {
		countryCode = geoMock()
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, helpers.CountryCodeType{}, countryCode)

	*r = *r.WithContext(ctx)

	return nil, http.StatusOK, nil
}

func geoMock() string {
	return "ZZ"
}
