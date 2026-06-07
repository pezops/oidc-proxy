package auth

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// A JwksKeyManager implements the KeyManager interface and supports an
// auto-refreshed JWKS URL for retrieving keys for JWT validation.
type JwksKeyManager struct {
	url            string
	jwks           keyfunc.Keyfunc
	expectedClaims *ValidatableMapClaims
}

// Validate will parse and validate a JWT token and its claims.
func (m *JwksKeyManager) Validate(tok string) (bool, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		tok, &claims, m.keyfunc,
	)
	if err != nil {
		return false, err
	}

	return m.expectedClaims.ValidateClaims(&claims)
}

// keyfunc uses the JWT kid header when present. When kid is omitted, it
// returns all verification keys and lets jwt try each key against the
// signature.
func (m *JwksKeyManager) keyfunc(token *jwt.Token) (any, error) {
	if _, ok := token.Header[jwkset.HeaderKID]; ok {
		return m.jwks.Keyfunc(token)
	}

	return m.jwks.VerificationKeySet(context.Background())
}

// NewJwksKeyManager returns a new JwksKeyManager for the specified JWKS URL.
func NewJwksKeyManager(url string, claims *ValidatableMapClaims) *JwksKeyManager {
	m := JwksKeyManager{
		url:            url,
		expectedClaims: claims,
	}

	ctx := context.Background()

	remoteJWKSets := make(map[string]jwkset.Storage)
	jwksetHTTPStorageOptions := jwkset.HTTPClientStorageOptions{
		Client:                    http.DefaultClient,
		Ctx:                       ctx,
		HTTPExpectedStatus:        http.StatusOK,
		HTTPMethod:                http.MethodGet,
		HTTPTimeout:               10 * time.Second,
		NoErrorReturnFirstHTTPReq: true, // Create storage regardless if the first HTTP request fails.
		RefreshErrorHandler: func(ctx context.Context, err error) {
			slog.Default().ErrorContext(
				ctx, "failed to refresh JWKS URL", "error", err, "url", url,
			)
		},
		RefreshInterval: 10 * time.Minute,
	}
	store, err := jwkset.NewStorageFromHTTP(url, jwksetHTTPStorageOptions)
	if err != nil {
		log.Fatalf("failed to create JWKS HTTP client for %q: %s", url, err)
	}
	_, err = store.KeyReadAll(ctx)
	if err != nil {
		log.Fatalf("failed to read JWKS keys for %q: %s", url, err)
	}
	remoteJWKSets[url] = store

	// Create the JWK Set containing HTTP clients
	jwksetHTTPClientOptions := jwkset.HTTPClientOptions{
		HTTPURLs:          remoteJWKSets,
		PrioritizeHTTP:    false,
		RefreshUnknownKID: rate.NewLimiter(rate.Every(5*time.Minute), 5),
	}
	combined, err := jwkset.NewHTTPClient(jwksetHTTPClientOptions)
	if err != nil {
		log.Fatalf("failed to create JWKS HTTP client: %s", err)
	}

	// Create the keyfunc.Keyfunc.
	keyfuncOptions := keyfunc.Options{
		Ctx:          ctx,
		Storage:      combined,
		UseWhitelist: []jwkset.USE{jwkset.UseSig},
	}
	jwks, err := keyfunc.New(keyfuncOptions)
	if err != nil {
		log.Fatalf("failed to create JWKS keyfunc: %s\n", err.Error())
	}

	m.jwks = jwks

	return &m
}
