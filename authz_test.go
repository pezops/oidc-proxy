package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pezops/oidc-proxy/auth"
)

func TestModifyValidateRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://foo", strings.NewReader("test"))
	assert.Nil(t, err)
	tokenString := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2ZvbyIsInN1YiI6IjEyMzQ1Njc4OTAiLCJhdWQiOiJmb28iLCJleHAiOjQ1MTYyMzkwMjIsImlhdCI6MTUxNjIzOTAyMn0.5R_64maDqRDM5egnSgPu6CJijXo0KhqqoTAoEiMdRAM`
	retConfig := &auth.StaticTokenConfig{
		Token: tokenString,
	}
	retriever := &auth.StaticTokenRetriever{}
	err = retriever.Configure(retConfig)
	assert.Nil(t, err)
	jwtManager := auth.NewJwtManager(retriever)
	rw := httptest.NewRecorder()
	modifyRequestAuthz(rw, req, jwtManager, "foo")
	assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprintf("Bearer %v", tokenString))

	claims := &auth.ValidatableMapClaims{}
	claims.AddClaim("aud", "http://foo")
	keyManager := auth.NewStaticKeyManager(tokenString, claims)
	rw = httptest.NewRecorder()
	assert.False(t, validateRequestAuthz(rw, req, keyManager))

	claims = &auth.ValidatableMapClaims{}
	claims.AddClaim("aud", "foo")
	keyManager = auth.NewStaticKeyManager(tokenString, claims)
	rw = httptest.NewRecorder()
	assert.True(t, validateRequestAuthz(rw, req, keyManager))
	b, err := io.ReadAll(rw.Body)
	assert.Nil(t, err)
	if len(b) > 0 {
		println(string(b))
	}
}
