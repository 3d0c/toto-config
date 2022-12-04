package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/3d0c/toto-config/pkg/apiserver/models"
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
		ok             bool
		err            error
		packageName    string = chi.URLParam(r, "package")
		cm             *models.Config
		result         *models.ConfigScheme
	)

	if countryCode, ok = r.Context().Value(helpers.CountryCodeType{}).(string); countryCode == "" || !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting Country Code from context")
	}
	if percentileSeed, ok = r.Context().Value(helpers.PercentileSeedType{}).(int); percentileSeed == 0 || !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting Perentile Seed from context")
	}

	if cm, err = models.NewConfigModel(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing Config Model")
	}

	if result, err = cm.FindBy(packageName, countryCode, percentileSeed); err != nil {
		if err == models.ErrNotFound {
			return nil, http.StatusNotFound, fmt.Errorf("no SKU found for packageName '%s'", packageName)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("error finding SKU for packageName '%s' - %s", packageName, err)
	}

	return result, http.StatusOK, nil
}
