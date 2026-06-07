package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"

	"github.com/pezops/oidc-proxy/auth"
)

// modifyRequestAuthz will modify an in-flight request in egress mode to
// insert the authorization header and JWT.
func modifyRequestAuthz(rw http.ResponseWriter, req *http.Request, manager auth.JwtManager, aud string) bool {
	token, err := manager.Token(aud)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte(fmt.Sprintf("error obtaining token: %v", err)))
		return false
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	return true
}

// validateRequestAuthz will validate an in-flight request in ingress mode by
// parsing the JWT and validating all claims.
func validateRequestAuthz(rw http.ResponseWriter, req *http.Request, manager auth.KeyManager) bool {
	az := req.Header.Get("Authorization")
	if az == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		return false
	}

	tokenSlice := strings.Split(az, " ")
	if len(tokenSlice) != 2 {
		rw.WriteHeader(http.StatusUnauthorized)
		_, _ = rw.Write([]byte("invalid authorization header format"))
		return false
	}
	tokenString := tokenSlice[1]
	if tokenString == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		_, _ = rw.Write([]byte("empty token in authorization header"))
		return false
	}

	v, err := manager.Validate(tokenString)
	if !v {
		rw.WriteHeader(http.StatusUnauthorized)
		if err != nil {
			_, _ = rw.Write([]byte(err.Error()))
		}
		return false
	}

	return true
}

// convertAudienceString converts a provided audience claim value to a slice
// of audiences strings.
func convertAudienceString(audString string) ([]string, error) {
	var aud []string

	audString = strings.TrimSpace(audString)
	if audString == "" {
		return nil, nil
	}

	err := json.Unmarshal([]byte(audString), &aud)
	if err == nil {
		log.Println("detected JSON audience list")
		return aud, nil
	}

	err = yaml.Unmarshal([]byte(audString), &aud)
	if err == nil {
		log.Println("detected YAML audience list")
		return aud, nil
	}

	return []string{audString}, nil
}

// detectValidatingKey attempts to detect the key-type for JWT validation.
func detectValidatingKey(b []byte) interface{} {
	privKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err == nil {
		log.Println("detected RSA public key")
		return privKey
	}

	decodedKey, err := base64.StdEncoding.DecodeString(string(b))
	if err == nil {
		log.Println("detected base64-encoded symmetric key")
		return decodedKey
	}

	log.Println("detected raw symmetric key")
	return b
}
