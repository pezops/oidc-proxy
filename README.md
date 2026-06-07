# oidc-proxy

A small reverse proxy to add or verify OpenID Connect (OIDC) identity tokens as authentication for clients and services
that otherwise do not support OIDC. The oidc-proxy inserts or validates
OIDC [ID tokens](https://openid.net/specs/openid-connect-core-1_0.html#IDToken) used in the `Authorization` header as
a [bearer token](https://www.rfc-editor.org/rfc/rfc6750).

<!-- TOC -->
* [Overview](#overview)
  * [Workload Authentication](#workload-authentication)
  * [Not User Authentication](#not-user-authentication)
* [Egress Mode](#egress-mode)
  * [Authentication Modes](#authentication-modes)
    * [Google Cloud Instance Identity - `gcp`](#google-cloud-instance-identity---gcp)
    * [Manual Key Signing - `manual`](#manual-key-signing---manual)
    * [Static Key - `static`](#static-key---static)
* [Ingress Mode](#ingress-mode)
  * [Token Validation Methods](#token-validation-methods)
    * [JSON Web Key Set URL](#json-web-key-set-url)
    * [Validating Key](#validating-key)
    * [Static Identity Token](#static-identity-token)
* [Releases](#releases)
* [Usage](#usage)
<!-- TOC -->

## Overview

The goal of the `oidc-proxy` is to allow services to authenticate to each other using runtime-assigned identity
information. This would normally be considered a machine account, service account or workload identity. For normal OAuth
2.0 flows with registered clients communicating on behalf of a user to an application,
see [oauth2-proxy](https://oauth2-proxy.github.io/oauth2-proxy/).

### Workload Authentication

The `oidc-proxy` uses the OIDC identity token standard to assert workload identity between two related workloads or
services, **not** on behalf of a user.

<img src="assets/oidc-service-authn.png" alt="OIDC service authentication" width="400"/>

The way a service obtains a scoped identity token is not defined by a standard flow as with OAuth 2.0, it will depend on
where the application is running. A common way that workloads are identified is via a central service that provides or
injects verifiable identities into an application at runtime. This may look like the following:

<img src="assets/oidc-service-authn-example.png" alt="OIDC service authentication" width="400"/>

Some examples of ways an application may be able to obtain an identity token include:

- Private key accessible at runtime
- [Google instance identity](https://cloud.google.com/compute/docs/instances/verifying-instance-identity#request_signature)
- [Kubernetes service accounts](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#serviceaccount-token-volume-projection)

### Not User Authentication

Normal OIDC flow is built on top of an OAuth 2.0 flow, and allows an external or third-party application to access an
application API on behalf of a user. It looks something like this (simplified):

<img src="assets/oidc-oauth-flow.png" alt="OIDC OAuth flow" width="600"/>

In this "user-driven" flow, the client is an external application (referred to as the "relying party") and uses an
access token to access a resource. The identity token is used to assert the user's identity, between the authorization
service and the client application. The identity token is signed, and validated using a public key published on the
internet.

However, there are many reasons why internal services need to communicate on behalf of themselves, not on behalf of a
user. For these, authorization can be much less complex where the services directly authorize each other.

## Egress Mode

<img src="assets/egress-use-case.png" alt="egress use case" width="600"/>


In egress mode, the oidc-proxy will add an `Authorization` header containing an identity bearer token scoped to the
target service (audience).

### Authentication Modes

In egress mode, an authentication type must be specified. Not all authentication types require additional configuration.

#### Google Cloud Instance Identity - `gcp`

The GCP authentication method works with both virtual machine instance identities, and GKE workload identities.

Optional parameters:

- Service account name

e.g.

```shell
oidc-proxy --target-url="https://foo" --audience=foo --egress-enabled --egress-auth-type=gcp
```

#### Manual Key Signing - `manual`

The manual key signing method will generate and sign a JWT identity token using a provided key (RSA or HMAC). Identity
tokens are currently generated for a validity period of one hour.

Required parameters:

- Issuer
- Subject
- Signing method
    - RS256/384/512
    - HS256/384/512
- Signing key
    - PEM formatted RSA private key
    - Base64 encoded HMAC key
    - Raw HMAC key

Optional parameters:

- Additional claims
    - JSON formatted

e.g.

```shell
oidc-proxy --target-url="https://foo" --audience=foo --egress-enabled --egress-auth-type=manual --egress-auth-manual-issuer="internal.accounts" --egress-auth-manual-subject=bar --egress-auth-manual-signing-method=hs256 --egress-auth-manual-signing-key="example"
```

#### Static Key - `static`

With the static authentication key method, a JWT identity token is provided by the environment. Static authentication
does not support renewal.

Required parameters:

- JWT identity token

e.g.

```shell
oidc-proxy --target-url="https://foo" --egress-enabled --egress-auth-type=static --egress-auth-static-token="eyJhbGciOi...."
```

## Ingress Mode

<img src="assets/ingress-use-case.png" alt="ingress use case" width="600"/>

In ingress mode, the oidc-proxy will validate the `Authorization` header of incoming HTTP requests. Unauthorized
requests will not be passed to the target service.

Additional claims may also be specified for validation on incoming identity tokens.

Optional parameters:

- Additional claims for authorization
    - JSON formatted
    - String values use [RE2 syntax](https://golang.org/s/re2syntax)

e.g.

```shell
oidc-proxy --target-url="https://foo" --ingress-enabled --ingress-jwks-url="https://www.googleapis.com/oauth2/v3/certs"
```

### Token Validation Methods

One of the following validation methods must be specified.

#### JSON Web Key Set URL

By providing a JWKS URL, incoming requests with identity tokens will be validated against keys returned the JWKS URL.
The JWKS contents will be periodically refreshed.

#### Validating Key

A validating key may be provided for identity token validation. For RSA signature methods, a PEM-formatted RSA public
key may be provided. For HMAC signature methods, a base64-encoded key or raw key may be provided.

#### Static Identity Token

While not recommended, a static identity token may be provided for validation. Keep in mind, all OIDC required claims
will be verified, including token audience and expiration claims. Signature validation will not be performed with a
static identity token.

## Releases

The Go module is `github.com/pezops/oidc-proxy`. Releases after v0.1.0 are published as signed multi-platform
container images to both registries:

- `ghcr.io/pezops/oidc-proxy`
- `docker.io/pezops/oidc-proxy`

Each release publishes its exact version tag. The latest stable release also updates the minor,
major, and `latest` tags. For example, `v0.1.2` publishes `0.1.2`, `0.1`, `0`, and `latest`. See
[RELEASING.md](RELEASING.md) for the maintainer workflow.

## Usage

```shell
oidc-proxy [OPTIONS]
```

Options and parameters:

| Option                                | Description                             | Default     | Example                                        |
|---------------------------------------|-----------------------------------------|-------------|------------------------------------------------|
| `--target-url`                        | Target URL for incoming requests        |             | `https://localhost`                            |
| `--audience`                          | Audience claim for token                |             | `https://myservice`                            |
| `--port`                              | Port to listen for requests             | `8080`      | `8080`                                         |
| `--address`                           | Address to listen for requests          | `127.0.0.1` | `localhost`                                    |
| `--ingress-enabled`                   | Enable ingress mode                     | `false`     | `true`                                         |
| `--ingress-jwks-url`                  | JSON web key set URL                    |             | `https://www.googleapis.com/oauth2/v3/certs`   |
| `--ingress-validating-key`            | Signing key                             |             | `-----BEGIN PUBLIC KEY-----\n...`              |
| `--ingress-static-token`              | Static identity token                   |             | `eyJhbG...`                                    |
| `--ingress-valid-claims`              | Claims for validation                   |             | `{"sub": "my_app", "email": "my_app@foo.com"}` |
| `--egress-enabled`                    | Enable egress mode                      | `false`     | `true`                                         |
| `--egress-auth-type`                  | Authentication type for egress mode     |             | `gcp`                                          |
| `--egress-auth-static-token`          | Static authentication identity token    |             | `eyJhbG...`                                    |
| `--egress-auth-manual-issuer`         | Manual authentication issuer claim      |             | `https://my-app`                               |
| `--egress-auth-manual-subject`        | Manual authentication subject claim     |             | `my-app`                                       |
| `--egress-auth-manual-signing-key`    | Manual authentication signing key       |             | `-----BEGIN RSA PRIVATE KEY-----\n...`         |
| `--egress-auth-manual-signing-method` | Manual authentication signing method    |             | `rs256`                                        |
| `--egress-auth-manual-claims`         | Manual authentication additional claims |             | `{"app": "my_app", "email": "my_app@foo.com"}` |
| `--egress-auth-gcp-service-account`   | GCP instance identity name              | `default`   | `other-service-account`                        |
| `--tls-listen-enabled`                | Listen for requests using TLS           | `false`     | `true`                                         |
| `--tls-cert`                          | Path to TLS public certificate          |             | `/var/opt/tls/cert.pem`                        |
| `--tls-key`                           | Path to TLS private key                 |             | `/var/opt/tls/key.pem`                         |
| `--tls-allow-insecure-target`         | Do not verify TLS for the target        | `false`     | `true`                                         |

All options may be specified using environment variables. The name of the environment variable will be prefixed
with `OIDC_PROXY` and followed by the name of the option. All dashes will become underscores in the environment variable
name.

For example, the option `--ingress-validating-key` may also be specified by setting the environment
variable `OIDC_PROXY_INGRESS_VALIDATING_KEY`.
