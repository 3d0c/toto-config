package middlewares

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/3d0c/toto-config/pkg/helpers"
)

const (
	minSeed = 1
	maxSeed = 100
)

// SetPercentile generates and sets random number between 1 and 100
func SetPercentile(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		seed int
	)

	rand.Seed(time.Now().UnixNano())

	seed = rand.Intn(maxSeed-minSeed+1) + minSeed

	ctx := r.Context()
	ctx = context.WithValue(ctx, helpers.PercentileSeedType{}, seed)

	*r = *r.WithContext(ctx)

	return nil, http.StatusOK, nil
}
