# Changelog

## [0.1.0](https://github.com/pezops/oidc-proxy/compare/v0.0.6...v0.1.0) (2026-06-07)


### ⚠ BREAKING CHANGES

* Move to Release Please and relocate to the pezops org ([#13](https://github.com/pezops/oidc-proxy/issues/13))

### Bug Fixes

* Allow JWT validation without a `kid` header ([#16](https://github.com/pezops/oidc-proxy/issues/16)) ([1561d2d](https://github.com/pezops/oidc-proxy/commit/1561d2d356c63b2c8eb6cdf9138e5723ccf6235d))


### Documentation

* Document weekly rebuild and floating tags in releases ([#17](https://github.com/pezops/oidc-proxy/issues/17)) ([a610ee1](https://github.com/pezops/oidc-proxy/commit/a610ee1140f0d1d6fb341aac265345ab4af2a20c))


### Continuous Integration

* Add check for tidy Go modules ([#18](https://github.com/pezops/oidc-proxy/issues/18)) ([26d77c4](https://github.com/pezops/oidc-proxy/commit/26d77c4e0eb8701c246a91b049abd77c2ddce1cb))
* Move to Release Please and relocate to the pezops org ([#13](https://github.com/pezops/oidc-proxy/issues/13)) ([3ceb870](https://github.com/pezops/oidc-proxy/commit/3ceb8708e839e208f7df270ec8b50509e2719f54))

## [0.0.6](https://github.com/pezops/oidc-proxy/compare/v0.0.5...v0.0.6) (2025-12-12)

### Continuous Integration

- Update CI actions and add pull request build and test coverage
  ([#11](https://github.com/pezops/oidc-proxy/pull/11))
  ([cc3fa50](https://github.com/pezops/oidc-proxy/commit/cc3fa50))

### Miscellaneous

- Update dependencies
  ([#12](https://github.com/pezops/oidc-proxy/pull/12))
  ([2ddb0cf](https://github.com/pezops/oidc-proxy/commit/2ddb0cf))

## [0.0.5](https://github.com/pezops/oidc-proxy/compare/v0.0.4...v0.0.5) (2025-06-11)

### Bug Fixes

- Handle invalid `Authorization` header formats
  ([#9](https://github.com/pezops/oidc-proxy/pull/9))
  ([507d567](https://github.com/pezops/oidc-proxy/commit/507d567))

### Miscellaneous

- Update dependencies
  ([#10](https://github.com/pezops/oidc-proxy/pull/10))
  ([26d4eaa](https://github.com/pezops/oidc-proxy/commit/26d4eaa))

## [0.0.4](https://github.com/pezops/oidc-proxy/compare/v0.0.3...v0.0.4) (2025-02-19)

### Features

- Add validation support for arbitrary claims containing arrays of strings
  ([#6](https://github.com/pezops/oidc-proxy/pull/6))
  ([c2f3f7a](https://github.com/pezops/oidc-proxy/commit/c2f3f7a))

### Miscellaneous

- Upgrade the JWT and JWKS key-function libraries
  ([#7](https://github.com/pezops/oidc-proxy/pull/7))
  ([0fdd1d3](https://github.com/pezops/oidc-proxy/commit/0fdd1d3))

## [0.0.3](https://github.com/pezops/oidc-proxy/compare/v0.0.2...v0.0.3) (2024-11-10)

### Continuous Integration

- Add the scheduled latest-image build workflow
  ([1ce9616](https://github.com/pezops/oidc-proxy/commit/1ce9616))

### Miscellaneous

- Update and clean up dependencies
  ([#3](https://github.com/pezops/oidc-proxy/pull/3))
  ([345178c](https://github.com/pezops/oidc-proxy/commit/345178c))

## [0.0.2](https://github.com/pezops/oidc-proxy/compare/v0.0.1...v0.0.2) (2022-12-12)

### Bug Fixes

- Preserve request bodies by configuring `GetBody` on proxied client requests
  ([#2](https://github.com/pezops/oidc-proxy/pull/2))
  ([57f35e9](https://github.com/pezops/oidc-proxy/commit/57f35e9))

## [0.0.1](https://github.com/pezops/oidc-proxy/releases/tag/v0.0.1) (2022-12-09)

### Features

- Add the initial OIDC proxy implementation, authentication modes, tests, and documentation
  ([#1](https://github.com/pezops/oidc-proxy/pull/1))
  ([98a981a](https://github.com/pezops/oidc-proxy/commit/98a981a))
