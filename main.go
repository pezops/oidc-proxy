package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/pezops/oidc-proxy/auth"
	"github.com/pezops/oidc-proxy/config"
)

var cfg = config.ProxyConfig{}

func main() {

	flagParser := flags.NewParser(nil, flags.Default)
	flagParser.NamespaceDelimiter = "-"
	_, err := flagParser.AddGroup("proxy", "", &cfg)
	if err != nil {
		return
	}
	_, err = flagParser.Parse()
	if err != nil {
		return
	}

	err = cfg.ValidateConfig()
	if err != nil {
		log.Fatalf("error validating cfg: %v\n", err.Error())
	}

	targetUrl, err := url.Parse(cfg.TargetUrl)
	if err != nil {
		log.Fatalf("error parsing target URL: %v\n", err.Error())
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: cfg.TLS.AllowInsecure},
		ExpectContinueTimeout: 1 * time.Second,
	}

	if cfg.Egress.Enabled {
		var retriever auth.JwtTokenRetriever
		var retConfig interface{}

		switch cfg.Egress.Auth.Type {
		case "static":
			retriever = new(auth.StaticTokenRetriever)
			retConfig = &cfg.Egress.Auth.Static
		case "manual":
			retriever = new(auth.ManualTokenRetriever)
			retConfig = &cfg.Egress.Auth.Manual
		case "gcp":
			retriever = new(auth.GcpTokenRetriever)
			retConfig = &cfg.Egress.Auth.Gcp
		default:
			log.Fatalln("no auth type specified")
		}

		err = retriever.Configure(retConfig)
		if err != nil {
			log.Fatalf("error configuring auth type: %v\n", err.Error())
		}

		audSlice, err := convertAudienceString(cfg.Audience)
		if err != nil {
			log.Fatalf("error parsing audience: %v\n", err.Error())
		}
		if len(audSlice) > 1 {
			log.Fatalln("error configuring audience: only one audience may be specified in egress mode")
		}

		manager := auth.NewJwtManager(retriever)
		http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
			if modifyRequestAuthz(rw, req, manager, audSlice[0]) {
				req.Host = targetUrl.Host
				body, _ := io.ReadAll(req.Body)
				req.Body = io.NopCloser(bytes.NewReader(body))
				req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
				proxy.ServeHTTP(rw, req)
			}
		})
	} else if cfg.Ingress.Enabled {
		var manager auth.KeyManager

		validClaims, err := auth.ConvertValidatableClaimString(cfg.Ingress.ValidClaims)
		if err != nil {
			log.Fatalf("error parsing expected claims: %v\n", err)
		}

		if validClaims.HasClaim("aud") {
			log.Fatal("audience claim must be specified in config, not in valid claims")
		}
		audSlice, err := convertAudienceString(cfg.Audience)
		if err != nil {
			log.Fatalf("error parsing audience: %v\n", err.Error())
		}
		validClaims.AddClaim("aud", audSlice[0])

		if cfg.Ingress.JwksUrl != "" {
			manager = auth.NewJwksKeyManager(cfg.Ingress.JwksUrl, validClaims)
		} else if cfg.Ingress.KeyData != "" {
			key := detectValidatingKey([]byte(cfg.Ingress.KeyData))
			manager = auth.NewManualKeyManager(key, validClaims)
		} else if cfg.Ingress.StaticToken != "" {
			manager = auth.NewStaticKeyManager(cfg.Ingress.StaticToken, validClaims)
		} else {
			log.Fatalln("failed to configure ingress")
		}

		http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
			if validateRequestAuthz(rw, req, manager) {
				req.Host = targetUrl.Host
				body, _ := io.ReadAll(req.Body)
				req.Body = io.NopCloser(bytes.NewReader(body))
				req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
				proxy.ServeHTTP(rw, req)
			}
		})

	}

	addr := fmt.Sprintf("%v:%v", cfg.Address, cfg.Port)

	log.Printf("listening on %v\n", addr)
	if cfg.TLS.Listen {
		err = http.ListenAndServeTLS(addr, cfg.TLS.Cert, cfg.TLS.Key, nil)
	} else {
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		log.Println(err.Error())
	}

}
