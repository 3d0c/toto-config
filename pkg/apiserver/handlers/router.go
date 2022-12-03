package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/3d0c/toto-config/pkg/apiserver/middlewares"
	"github.com/3d0c/toto-config/pkg/config"
	"github.com/3d0c/toto-config/pkg/log"
)

// SetupRouter sets up endpoints
func SetupRouter(cfg config.Server) *chi.Mux {
	var (
		root string = "/" + cfg.APIVersion
	)

	r := chi.NewRouter()

	// Dummy endpoint. Just a stub for tests
	r.Get(
		root+"/nil",
		middlewares.Chain(
			nilHandler,
		),
	)

	// Get SKU for package
	r.Get(
		root+"/sku/{package}",
		middlewares.Chain(
			middlewares.GeoTarget,
			middlewares.SetPercentile,
			skuHandler().get,
		),
	)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.TheLogger().Debug("registered", zap.String("method", method), zap.String("route", route))
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.TheLogger().Debug("logging error", zap.Error(err))
	}

	return r
}
