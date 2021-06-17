package handler

import (
	"bytes"
	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func requestEmpty(t *testing.T, method, url string) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(nil))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return req
}

func errorJson(t *testing.T, err error) string {
	v, err := apperrors.ToJSON(err)
	require.Nil(t, err)
	return string(v)
}

func TestGetQueryParams_errorOnInvalidPageSize(t *testing.T) {
	cases := []string{
		"badSize",
		"0",
		"-1",
		"true",
	}
	for _, value := range cases {
		t.Run("invalid pageSize="+value, func(t *testing.T) {
			_, _, _, err := GetQueryParams("", requestEmpty(t, "GET", "/messages?pageSize="+value))
			require.Equal(t, "{\"errors\":[\"invalid pageSize value\"]}", errorJson(t, err))
		})
	}
}

func TestGetQueryParams_errorOnInvalidPageStartIndex(t *testing.T) {
	cases := []string{
		"badSize",
		"0",
		"-1",
		"true",
	}
	for _, value := range cases {
		t.Run("invalid pageSize="+value, func(t *testing.T) {
			_, _, _, err := GetQueryParams("", requestEmpty(t, "GET", "/messages?pageStartIndex="+value))
			require.Equal(t, "{\"errors\":[\"invalid pageStartIndex value\"]}", errorJson(t, err))
		})
	}
}
