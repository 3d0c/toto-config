package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/3d0c/toto-config/pkg/helpers"
)

func TestGeoTarget(t *testing.T) {
	var (
		testCC = "CZ"
		status int
		err    error
	)

	cases := []struct {
		description        string
		req                *http.Request
		expectedStatusCode int
		expectedValue      string
		middleware         func(http.ResponseWriter, *http.Request) (interface{}, int, error)
		assertFn           func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			description: "Provided Country Header",
			req: withHeaders(
				httptest.NewRequest(http.MethodPut, "http:", nil),
				map[string]string{
					AppEngineCountryHeader: testCC,
				},
			),
			expectedStatusCode: http.StatusOK,
			expectedValue:      testCC,
			middleware:         GeoTarget,
			assertFn:           assert.Nil,
		},
		{
			description:        "No Country Header",
			req:                httptest.NewRequest(http.MethodPut, "http:", nil),
			expectedStatusCode: http.StatusOK,
			expectedValue:      "ZZ",
			middleware:         GeoTarget,
			assertFn:           assert.Nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.description, func(t *testing.T) {
			_, status, err = testCase.middleware(nil, testCase.req)
			testCase.assertFn(t, err)
			assert.Equal(t, testCase.expectedStatusCode, status)

			countryCode, ok := testCase.req.Context().Value(helpers.CountryCodeType{}).(string)
			assert.True(t, ok)
			assert.Equal(t, testCase.expectedValue, countryCode)
		})
	}
}

func withHeaders(req *http.Request, headers map[string]string) *http.Request {
	for h, v := range headers {
		req.Header.Set(h, v)
	}

	return req
}
