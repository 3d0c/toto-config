package middlewares

import (
	"context"
	"net"
	"net/http"

	"github.com/oschwald/geoip2-golang"
	"go.uber.org/zap"

	"github.com/3d0c/toto-config/pkg/helpers"
	"github.com/3d0c/toto-config/pkg/log"
)

const (
	AppEngineCountryHeader = "X-Appengine-Country"
	DefaultCountryCode     = "ZZ"
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
		countryCode = GeoByIP(r)
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, helpers.CountryCodeType{}, countryCode)

	*r = *r.WithContext(ctx)

	return nil, http.StatusOK, nil
}

func GeoByIP(r *http.Request) string {
	var (
		db      *geoip2.Reader
		country *geoip2.Country
		err     error
	)

	if db, err = geoip2.Open("GeoIP2-City.mmdb"); err != nil {
		log.TheLogger().Error("error opening GeoIP2 database", zap.Error(err))
		return DefaultCountryCode
	}
	defer db.Close()

	if country, err = db.Country(net.ParseIP(r.RemoteAddr)); err != nil {
		log.TheLogger().Error("error getting country", zap.String("RemoteAddr", r.RemoteAddr), zap.Error(err))
		return DefaultCountryCode
	}

	return country.Country.IsoCode
}
