package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/3d0c/toto-config/pkg/helpers"
)

func TestPercentile(t *testing.T) {
	var (
		status int
		err    error
	)

	t.Run("Percentile Range", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "http:", nil)
		_, status, err = SetPercentile(nil, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, status)

		obtained, ok := req.Context().Value(helpers.PercentileSeedType{}).(int)

		assert.True(t, ok)

		if obtained <= 0 || obtained > 100 {
			t.Fatalf("Expected range between 0 and 100, obtained - %d\n", obtained)
		}
	})
}
