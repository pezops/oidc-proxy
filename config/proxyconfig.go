package config

import (
	"errors"

	"github.com/pezops/oidc-proxy/auth"
)

// ProxyConfig is the base configuration for the oidc-proxy. It it used to
// generate all command line flags and configuration environment variables.
type ProxyConfig struct {
	TargetUrl string             `long:"target-url" env:"OIDC_PROXY_TARGET_URL" description:"Target URL for incoming requests" required:"true"`
	Ingress   ProxyIngressConfig `group:"ingress" namespace:"ingress" env-namespace:"OIDC_PROXY_INGRESS"`
	Egress    ProxyEgressConfig  `group:"egress" namespace:"egress" env-namespace:"OIDC_PROXY_EGRESS"`
	Audience  string             `long:"audience" env:"OIDC_PROXY_AUDIENCE" description:"Audience claim for token" required:"true"`
	Port      int                `long:"port" env:"OIDC_PROXY_PORT" description:"Port to listen for requests" default:"8080"`
	Address   string             `long:"address" env:"OIDC_PROXY_ADDRESS" description:"Address to listen for requests" default:"127.0.0.1"`
	TLS       ProxyTLSConfig     `group:"tls" namespace:"tls" env-namespace:"OIDC_PROXY_TLS"`
}

// ProxyTLSConfig contains configuration information about listening for
// requests using TLS and for making outbound requests over TLS.
type ProxyTLSConfig struct {
	Listen        bool   `long:"listen-enabled" env:"LISTEN_ENABLED" description:"Listen for requests using TLS"`
	Cert          string `long:"cert" env:"CERT" description:"Path to TLS public certificate (PEM format)"`
	Key           string `long:"key" env:"KEY" description:"Path to TLS private key (PEM format)"`
	AllowInsecure bool   `long:"allow-insecure-target" env:"ALLOW_INSECURE_TARGET" description:"Do not verify TLS for the target"`
}

// ProxyEgressConfig contains configuration data for egress mode.
type ProxyEgressConfig struct {
	Enabled bool                  `long:"enabled" env:"ENABLED" description:"Enable egress mode"`
	Auth    ProxyEgressAuthConfig `group:"egress.auth" namespace:"auth" env-namespace:"AUTH"`
}

// ProxyEgressAuthConfig contains configuration information for the selected auth method.
type ProxyEgressAuthConfig struct {
	Type   string                 `long:"type" env:"TYPE" description:"Authentication type for egress mode"`
	Static auth.StaticTokenConfig `group:"egress.auth.static" namespace:"static" env-namespace:"STATIC"`
	Manual auth.ManualTokenConfig `group:"egress.auth.manual" namespace:"manual" env-namespace:"MANUAL"`
	Gcp    auth.GcpTokenConfig    `group:"egress.auth.gcp" namespace:"gcp" env-namespace:"GCP"`
}

// ProxyIngressConfig contains configuration data for ingress mode.
type ProxyIngressConfig struct {
	Enabled     bool   `long:"enabled" env:"ENABLED" description:"Enable ingress mode"`
	JwksUrl     string `long:"jwks-url" env:"JWKS_URL" description:"JSON web key set URL for key validation"`
	KeyData     string `long:"validating-key" env:"VALIDATING_KEY" description:"Signing key for validation"`
	StaticToken string `long:"static-token" env:"STATIC_TOKEN" description:"Static identity token for validation"`
	ValidClaims string `long:"valid-claims" env:"VALID_CLAIMS" description:"Claims for validation (JSON or YAML map)"`
}

// ValidateConfig checks to make sure that the provided flags make sense and are valid.
func (p *ProxyConfig) ValidateConfig() error {
	if p.TargetUrl == "" {
		return errors.New("target URL is required")
	}

	if p.Audience == "" && p.Egress.Auth.Type != "static" {
		return errors.New("audience is required")
	}

	if p.Egress.Enabled {
		if p.Egress.Auth.Type == "" {
			return errors.New("egress mode: auth type is required")
		}

	} else if p.Ingress.Enabled {
		if p.Ingress.JwksUrl == "" && p.Ingress.KeyData == "" && p.Ingress.StaticToken == "" {
			return errors.New("ingress mode: JWKS URL, validating key, or static token is required")
		}

	} else {
		return errors.New("no direction specified, choose Ingress or Egress")
	}

	if p.TLS.Listen && (p.TLS.Cert == "" || p.TLS.Key == "") {
		return errors.New("when TLS is enabled, a certificate and key path must be specified")
	}

	return nil
}
