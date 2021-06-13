package server

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessages_405WhenMethodIsNotSupported(t *testing.T) {
	t.Parallel()
	h := noDbHandler(t)
	for _, method := range allMethodsExcept("GET", "POST", "OPTIONS") {
		t.Run("method "+method, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, requestEmpty(t, method, "/messages"))
			require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
			require.Equal(t, "Allow: GET, POST", rr.Body.String())
		})
	}
}

func TestMessages_CanListOptions(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	noDbServe(t, rr, requestEmpty(t, "OPTIONS", "/messages"))
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "Allow: GET, POST", rr.Body.String())
}
