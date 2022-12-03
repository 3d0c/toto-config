package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/3d0c/toto-config/pkg/helpers"
)

type sku struct{}

func skuHandler() *sku {
	return &sku{}
}

func (*sku) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		countryCode    string
		percentileSeed int
		packageName    string = chi.URLParam(r, "package")
		// cm             *models.Config
		// result         *models.ConfigScheme
	)

	if countryCode = r.Context().Value(helpers.CountryCodeType{}).(string); countryCode == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting Country Code from context")
	}
	if percentileSeed = r.Context().Value(helpers.PercentileSeedType{}).(int); percentileSeed == 0 {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting Perentile Seed from context")
	}

	return struct {
		CC      string
		Seed    int
		Package string
	}{
		CC:      countryCode,
		Seed:    percentileSeed,
		Package: packageName,
	}, http.StatusOK, nil
}
