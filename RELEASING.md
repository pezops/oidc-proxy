# Release Process

Release Please maintains a release pull request from conventional commits merged into `main`. The
release pull request updates `CHANGELOG.md` and `.release-please-manifest.json`. Merging it creates a
version tag and a draft GitHub release.

The first release after adopting Release Please must be `v0.1.0` because the Go module moved from
`github.com/mbrancato/oidc-proxy` to `github.com/pezops/oidc-proxy`. Squash-merge the migration with
a breaking conventional commit title such as `feat!: move module to pezops`, or include this footer:

```text
Release-As: 0.1.0
```

## Publishing a Release

1. Merge changes into `main` using conventional commit titles.
2. Review the Release Please pull request and its proposed changelog and version.
3. Merge the Release Please pull request after its required checks pass.
4. Review the draft under GitHub Releases.
5. Publish the draft release.

Publishing runs the release workflow against the tagged source. It tests the code, builds one
multi-platform image, publishes the same digest to GHCR and Docker Hub, signs both image references,
and publishes build provenance.

For a release `v0.1.2`, the workflow publishes these tags to both registries:

- `0.1.2`
- `0.1`
- `0`
- `latest`

Floating tags are updated only when the published GitHub release is the latest stable release.

## Version Selection

Before `v1.0.0`, `fix:` and `feat:` commits produce patch releases, while breaking changes produce
minor releases. Starting with `v1.0.0`, features produce minor releases and breaking changes produce
major releases. Use a `Release-As: x.y.z` commit footer to override the calculated version.

## Weekly Rebuild

Every Tuesday at 09:00 UTC, the rebuild workflow checks out the latest published release, tests it,
and rebuilds it with the current build tooling. It intentionally moves the exact, minor, major, and
`latest` tags in both registries to the rebuilt digest. Consumers requiring immutable artifacts must
pin the image digest rather than a tag.

The rebuild workflow can also be started manually from GitHub Actions.

## Repository Configuration

The workflows require these GitHub repository settings:

- Variable `DOCKER_REPO`, for example `pezops/oidc-proxy`
- Variable `DOCKERHUB_USERNAME`
- Secret `DOCKERHUB_TOKEN`
- Secret `PEZOPS_RELEASE_DISPENSER_CLIENT_ID`
- Secret `PEZOPS_RELEASE_DISPENSER_PRIVATE_KEY`
