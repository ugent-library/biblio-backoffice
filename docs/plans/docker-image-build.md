# Plan: rebuild the container image with current Go and CA certificates

**Status:** proposed 2026-10-26. Copied from ingest's plan (`ingest/docs/archive/docker-image-build.md`), where part A shipped in v1.1.1 on 2026-10-26. Not started here.

## Reusing this plan in another project

This plan is written to be copied to other Go projects that build an Alpine image on Docker Hub and run on Nomad. Everything outside *This project* and *Part B* applies to such a project as it stands. To reuse it:

1. Copy this file to the other project's `docs/plans/`.
2. Fill in *This project* with that project's values. The steps in part A refer to them by name (the binary path, the build and test commands).
3. Replace *Part B* with that project's own cleanup, or delete it.
4. Work through part A there and record the answers to the open questions in that copy.

There is no need to repeat the test in *What we tested*: it shows how Go's `crypto/x509` and Alpine's CA bundle behave, which is the same for every Go binary on Alpine.

## This project

| | biblio-backoffice |
| --- | --- |
| Outbound HTTPS calls | Crossref, DataCite, arXiv, Europe PMC, Plato, ORCID, the Handle server, the people and projects APIs, citeproc, the OIDC provider, S3, Elasticsearch |
| Build stage today | `FROM golang:alpine`, runs `go get -v ./...`, builds with `go build -o app -v`, then `go install github.com/jackc/tern/v2@latest` |
| Final stage today | `FROM alpine:latest`, no package upgrade, copies the whole `/build` directory to `/dist` |
| Binaries in the image | `/dist/app`, and `/dist/tern` for database migrations. Both are built in the build stage, so both must be built with the fixed Go |
| `go` directive in go.mod | `go 1.24.0`, with `toolchain go1.24.1` |
| `hooks/build` | exists. Passes `SOURCE_BRANCH`, `SOURCE_COMMIT` and `IMAGE_NAME`; no `--pull`, no `SNAPSHOT_EPOCH` |
| Build and test commands | `go build -o app`, `go test ./...` (the Cypress tests in `.github/workflows/cypress-tests.yml` run on GitHub) |
| Running containers | the server (`server start`) and the API (`api start`), plus any periodic jobs (`update-oai`, `update-embargoes`, ...) that run from this image. *The full list of Nomad jobs is looked up at deploy time (A5, A6)* |
| Branches | `main` is released; `dev` collects work before a release and runs on the development environment |
| Merge commit message | git's default, `Merge branch '<branch>'` |
| Releases | annotated tags `vX.Y.Z` plus an `update CHANGELOG.md for new version X.Y.Z` commit; no GitHub releases. Latest is `v1.0.78` |
| Docker Hub build rules | branch `dev` builds `ugentlib/biblio-backoffice:dev`, branch `main` builds `:latest`, git tag `vX.Y.Z` builds `:release-X.Y.Z` (no `v`; latest is `:release-1.0.78`). Found from the published images and tags on 2026-10-27. Production runs a release tag (today `:release-1.0.78`) |

The Dockerfile declares `ARG SOURCE_BRANCH`, `SOURCE_COMMIT` and `IMAGE_NAME` and copies them into `ENV`. Unlike ingest, the existing `hooks/build` already fills them in. A3 extends that hook instead of adding one.

The upgrade is a large step: from Go 1.24 to 1.27.3 (or 1.26.10, see *Open questions*). Go 1.24 has no release that carries the fix. Expect dependency and vet changes, and run the full test suite, including Cypress, before merging.

## Why now

On 17 October 2026 a researcher reported that the private keys of eight eMudhra intermediate CA certificates could be recovered from their public certificates. Anyone holding those keys can issue a TLS certificate that a client trusts for any domain. Browsers blocked the intermediates between 17 and 20 October, and the root programs distrusted `emSign ECC Root CA - G3` on 21 October. A Go service is not a browser: every outbound HTTPS call trusts whatever Go's `crypto/x509` and the image's CA bundle trust. Go does not check revocation, so the revocation of the intermediates on 18 and 19 October does not protect us.

The [disclosure site](https://emudhra.we-have-a-problem.com/) publishes the affected certificates in `impacted-intermediates.json`: sixteen certificates for the eight keys. Each intermediate exists twice, once issued by an emSign root and once cross-signed by Sectigo's `AAA Certificate Services`, and both copies carry the same key.

Two upstream fixes exist, and each one closes the hole in an Alpine image on its own:

- **Go 1.27.3** (released 2026-10-19, [CL 851047](https://go-review.googlesource.com/c/go/+/851047), CVE-2026-111943) makes `x509.Certificate.Verify` drop chains through the compromised intermediates, identified by the SHA-256 hash of their public key. That covers both copies of every intermediate, whatever root the chain ends in. The fix is compiled into the binary. Go 1.26.10 carries the same fix for projects still on Go 1.26.
- **Alpine `ca-certificates-bundle 20261020-r0`** (v3.24, v3.23, v3.22 and edge) removes `emSign ECC Root CA - G3`, following Mozilla's NSS 3.131.1.

### What we tested (2026-10-26, in ingest)

A small Go HTTPS client, built once with Go 1.27.2 and once with Go 1.27.3, run in Alpine images with the old (`20260909-r0`) and the new (`20261020-r0`) bundle, against the four test hosts. The bundle contents were checked by SHA-256 fingerprint:

| Bundle | `emSign ECC Root CA - G3` | `emSign Root TLS CA - G3` | `AAA Certificate Services` |
| --- | --- | --- | --- |
| `20260909-r0` | present | absent | absent |
| `20261020-r0` | absent | absent | absent |

| Go | Bundle | emSign chain | Sectigo chain |
| --- | --- | --- | --- |
| 1.27.2 | old | **accepted** | rejected |
| 1.27.3 | old | rejected (`no valid chains built`) | rejected |
| 1.27.2 | new | rejected (`certificate signed by unknown authority`) | rejected |
| 1.27.3 | new | rejected | rejected |

The hosts that serve certificates never logged to Certificate Transparency gave the same results, since Go does not enforce CT. So:

- An image is exposed only with an old Go and an old bundle together. Any image built before 20 October most likely has both.
- The Sectigo cross-signature does not reach an Alpine image: Alpine's bundle, built from Mozilla's, does not contain `AAA Certificate Services`. It matters for trust stores that do contain it, such as those of macOS and Windows, which is one reason the Go fix is worth having alongside the bundle.
- The error text shows which fix rejected a chain: `no valid chains built` is Go's filter, `certificate signed by unknown authority` is the bundle.

We take both fixes, so that a later mistake in one (a stale bundle, a toolchain left behind) does not reopen the hole. Because each fix was shown to work, the checks below only confirm the versions: the fixed Go in every binary and bundle `20261020-r0` or later in the image.

## Why a rebuild may pick up neither

- `FROM golang:alpine` and `FROM alpine:latest` are floating tags. Docker uses whatever copy of the tag it already has unless the build runs with `--pull`, so a rebuild can produce a binary built with an older Go, on an old bundle.
- Without a package upgrade in the final stage, the image has whatever bundle its base image shipped with. And a `RUN apk upgrade` alone does not help for long: Docker reuses a `RUN` layer as long as its parent layer and its command text are unchanged; it never looks at the Alpine package index. The layer would run once and stay cached.
- Docker Hub runs `hooks/build` instead of its own `docker build` when that file exists. Its default build does not pass `--build-arg`, so any `ARG` in the Dockerfile stays empty.

**Model:** `ugent-library/depot-management` solves the caching problem. Its `hooks/build` passes `--build-arg SNAPSHOT_EPOCH="$(date +%s)"`, and its Dockerfile declares `ARG SNAPSHOT_EPOCH` right before the package-install `RUN`. A declared `ARG` is part of the environment of every `RUN` after it in the same stage, so a new value on every build makes that `RUN` miss the cache and run again.

## Open questions

- **Go 1.27.3 or 1.26.10?** Both carry the fix. 1.27.3 matches ingest and stays supported longer; 1.26.10 is a smaller jump from 1.24. Either way this is a jump of several minor versions. Decide before A1. *biblio-backoffice: 1.27.3 (2026-10-26). Needs thorough testing; the main risks are the JSON v2 backend, the Green Tea GC, and (if go.mod is raised) container-aware GOMAXPROCS and the refusal of SHA-1 signatures in TLS 1.2.*
- **Pin tern.** `go install github.com/jackc/tern/v2@latest` builds whatever tern is newest on the day of the build. Pin it to a version that builds with the chosen Go, so a rebuild does not change the migration tool unnoticed. *biblio-backoffice: pin it, no preference for a version (v2.4.3 is the latest and needs Go 1.25 or later). The pin must keep working with whatever Go the project is on.*
- **Periodic rebuilds.** Docker Hub builds only on a push or a manual trigger. `SNAPSHOT_EPOCH` makes a build that runs pick up new packages; it does not make a build happen. If a project goes months without a push, its image keeps an old Go and an old bundle. Decide how rebuilds happen between code changes: by hand after each Go point release and Alpine security update, or on a schedule from outside Docker Hub. Check what Docker Hub currently offers for triggering builds before choosing. (Ingest chose by hand.) *biblio-backoffice: by hand (2026-10-27).*
- **Raise the `go` directive in go.mod?** Raising it to the fixed version makes a build on an older toolchain fail outright, because the official golang images set `GOTOOLCHAIN=local`. It also forces every developer onto that version (Go downloads it automatically unless `GOTOOLCHAIN=local` is set locally). Pinning the image tag already fixes the Docker build; this would add a check that stops a stale toolchain anywhere. Drop or update the `toolchain go1.24.1` line in the same change. (Ingest raised it.) *biblio-backoffice: raised to `go 1.27.3` and the `toolchain` line dropped (2026-10-27). `go mod tidy` regrouped the `require` blocks (new 1.27 behaviour); every module version is unchanged. Test results are the same as with the old go.mod: the build succeeds, the unit tests in `mutate` and `views` pass, and two failures predate the upgrade. 141 `go vet` "non-constant format string" findings on `gotext.Locale.Get` stop ten packages from building under `go test` (none of them has tests), and the `client/cmd` integration tests need a running API server on `:3999`.*
- **Docker Hub "Build caching" setting.** Look up whether it is on for the repository (Builds, Configure automated builds). If it is off, every build already starts from scratch and the EPOCH argument only matters if caching is turned on later. Keep the argument either way, so the Dockerfile does not depend on that setting. *biblio-backoffice (2026-10-27): not looked up; there's no need to turn it off. With caching on, `SNAPSHOT_EPOCH` still reruns `apk upgrade` on every build, `--pull` makes a republished base image a cache miss, and with tern pinned no layer downloads something the Dockerfile doesn't name.*

## Part A: pick up both fixes

The fix branches from `main`, not `dev`, so that releasing it does not also release whatever is waiting on `dev`. It is merged into both, so `dev` does not lose it at the next release.

### A0. Branch

- [x] Switch to `main` and bring it up to date: `git switch main && git pull --ff-only`.
- [x] Create the feature branch: `git switch -c docker-image-build`. Steps A1 to A3 are commits on this branch. *biblio-backoffice: branched from `main` at `bf9ef72a` (2026-10-26). `main` is ahead of `v1.0.78` with #1780 (`set_doi` and `remove_doi` batch operations) and an `npm audit fix` that only touches dev dependencies in `package-lock.json`. Kept `main` as the base, so the release in A6 ships #1780 too.*

### A1. Go

- [x] Change the build stage to `FROM golang:1.27.3-alpine` (or `golang:1.26.10-alpine`, see *Open questions*).
- [x] Decide the go.mod question above; if yes, change the `go` directive to the same version in the same commit, and run `go mod tidy`.
- [x] Pin the tern version in the `go install` line. *biblio-backoffice: `v2.4.3`.*
- [x] Build locally with `docker build --pull` and check which toolchain built each binary: copy it out with `docker create` and `docker cp <container>:<binary> ./app`, for every binary in *This project* (`/dist/app` and `/dist/tern`), then run `go version` on each copy on the host. It must print the fixed Go version. *biblio-backoffice (2026-10-27): both print `go1.27.3`; `tern version` prints `v2.4.3`.*

### A2. CA bundle, kept current

- [x] Change the final stage to `FROM alpine:3.24`, the release `alpine:latest` points to today. Pinning it makes moving to a new Alpine release a deliberate change rather than something a rebuild does unnoticed.
- [x] Right after that `FROM`, add the argument and the upgrade. The `ARG` must sit inside the final stage; one declared before the first `FROM` is not visible there. Upgrade all packages rather than naming one: the base image ships `ca-certificates-bundle`, not `ca-certificates`, so `apk upgrade ca-certificates` would change nothing, and the other packages get their security fixes too.

  ```dockerfile
  FROM alpine:3.24
  # hooks/build passes a new value on every build, so this layer is never served from cache.
  ARG SNAPSHOT_EPOCH
  RUN apk upgrade --no-cache
  ```

- [x] Check the result: `docker run --rm --entrypoint apk <image> list -I ca-certificates-bundle` shows `ca-certificates-bundle-20261020-r0` or later. Grepping the bundle for "emSign" proves nothing: `/etc/ssl/certs/ca-certificates.crt` holds PEM blocks only, without the root names, so the count is always 0. *biblio-backoffice (2026-10-27): `ca-certificates-bundle-20261020-r0` on Alpine 3.24.3; `alpine:latest` and `alpine:3.24` have the same digest. The first build showed the `apk upgrade` step as `CACHED`, reusing ingest's identical layer from the same machine, and values for `SNAPSHOT_EPOCH` that had been used before were cached too. A new value reruns the step and a repeated value is cached, so `hooks/build` must pass a value that is unique for each build, which `date +%s` does.*

### A3. Docker Hub build hook

- [x] Extend the existing `hooks/build` with `--pull` and the `SNAPSHOT_EPOCH` argument, so it reads:

  ```bash
  #!/bin/bash
  if [[ -z "$SOURCE_BRANCH" ]]; then
  	export SOURCE_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
  fi
  if [[ -z "$SOURCE_COMMIT" ]]; then
  	export SOURCE_COMMIT="$(git rev-parse HEAD)"
  fi

  echo "SOURCE_BRANCH: $SOURCE_BRANCH"
  echo "SOURCE_COMMIT: $SOURCE_COMMIT"
  echo "IMAGE_NAME: $IMAGE_NAME"

  docker build \
      --pull \
      --build-arg SOURCE_BRANCH=$SOURCE_BRANCH \
      --build-arg SOURCE_COMMIT=$SOURCE_COMMIT \
      --build-arg IMAGE_NAME=$IMAGE_NAME \
      --build-arg SNAPSHOT_EPOCH="$(date +%s)" \
      -f $DOCKERFILE_PATH \
      -t $IMAGE_NAME .
  ```

- [x] Check that the file is still executable (`git ls-files -s hooks/build` shows mode `100755`). *biblio-backoffice: it was `100644` since it was added in 2023 (`574d223b`); now `100755` via `chmod +x` and `git update-index --chmod=+x`. Docker Hub ran it anyway: the published `:dev` and `:latest` images have `SOURCE_BRANCH` and `SOURCE_COMMIT` filled in.*
- [x] Run the hook locally with `DOCKERFILE_PATH=Dockerfile IMAGE_NAME=biblio-backoffice:local hooks/build` to check it works outside Docker Hub. *biblio-backoffice (2026-10-27): it prints the three variables, the `apk upgrade` step runs instead of being cached, the image has `SOURCE_BRANCH` and `SOURCE_COMMIT` in its environment, and the bundle is `20261020-r0`.*
- [x] Look up the Docker Hub "Build caching" setting and record the answer in this plan. *biblio-backoffice: not needed, see* Open questions.

### A4. Check locally

- [x] Run the build and test commands from *This project* with the fixed Go toolchain. *biblio-backoffice (2026-10-27): go.mod is raised to `go 1.27.3`. Compared against a baseline of `main` with Go 1.24.1, run in the devcontainer against a separate `biblio_test` database and `biblio_test_*` indexes, both reset before each run:*
  - *Unit tests (`mutate`, `views`): pass on both.*
  - *`client/cmd` integration tests against `api start`: 31 pass and 24 fail on both. The output is identical apart from paths, IDs and timings; the failures are outdated tests (expected error formats, `doi`/`url` fields, time zones, a missing person in Mongo, index alias names).*
  - *API server logs: identical apart from one closure name in stack traces (1.27 names inlined closures differently).*
  - *`server start`: boots; `/`, `/status` and `/static/manifest.json` return 200 and `/login` redirects to OIDC.*
  - *Cypress: not run locally (decided 2026-10-27). The nightly CI run against the development environment covers it after A5.*

  *The devcontainer needed changes to start: the MinIO images are no longer published (so the stack ran without `minio` and `mc`), its Go image was bumped locally to `2-1.27-bookworm` (`apt-get install` needs `-y` on bookworm; not committed with this change), and its startup script expects the repo at `/src`, which only holds under VS Code.*
- [x] Look up the Docker Hub build rules and which image tag production runs, and record them in *This project*. Steps A5 and A6 depend on them.

### A5. Merge into dev and test there

- [ ] Merge the feature branch into `dev` with a merge commit: `git switch dev && git pull --ff-only && git merge --no-ff docker-image-build`, with the merge commit message from *This project*.
- [ ] Push `dev`. Docker Hub builds from what is pushed, so the development environment can only be tested after this push.
- [ ] In the Docker Hub build log, check that the hook ran (it prints `SOURCE_BRANCH`, `SOURCE_COMMIT` and `IMAGE_NAME`) and that the `apk upgrade` step ran rather than showing `CACHED`.
- [ ] Pull the image Docker Hub built and repeat the checks from A1 and A2: `go version` on both binaries and the bundle version. Also check that `env` shows non-empty `SOURCE_BRANCH` and `SOURCE_COMMIT`.
- [ ] Deploy to the development environment. Open a shell in every running allocation and run the command below. Type the paths themselves: a placeholder such as `<binary>` makes the shell read `<` as a redirect and fail with a syntax error.

  ```sh
  for b in /dist/app /dist/tern; do grep -a -o -E -m1 "go1\.[0-9]+\.[0-9]+" $b | head -1; done; apk list -I ca-certificates-bundle
  ```

  It must print the fixed Go version twice and `ca-certificates-bundle-20261020-r0` or later. The image has no Go toolchain, but each binary holds its Go version as text, which BusyBox `grep` can read.
- [ ] Run a migration with tern on the development database (or `tern status`) to check the new tern build works.
- [ ] Check that the services work as before on the development environment.
- [ ] Trigger a second Docker Hub build of `dev` without code changes and check that the `apk upgrade` step runs again. Skip this if build caching is off.

### A6. Merge into main and release

- [ ] Merge the feature branch into `main` with a merge commit: `git switch main && git pull --ff-only && git merge --no-ff docker-image-build`, with the merge commit message from *This project*. Merge the feature branch, not `dev`, so nothing else from `dev` comes along.
- [ ] Push `main`. Until production runs the new image, its containers keep the old Go and the old bundle, the one combination that accepts the compromised chain.
- [ ] Decide whether a release is needed. It is if production runs an image built from a tag: a new tag is then the only way to get the new Dockerfile into production, because rebuilding the old tag rebuilds the old Dockerfile. If production runs the image built from `main`, the push is enough.
- [ ] If a release is needed: add a `v1.0.79` entry to `CHANGELOG.md` that says the image is rebuilt with the fixed Go and the updated CA bundle because of the eMudhra key compromise and lists #1780 under "Added", commit it as `update CHANGELOG.md for new version 1.0.79`, create `v1.0.79` as an annotated tag on `main` and push the tag. No GitHub release; this project only uses tags.
- [ ] Deploy to production and run the same check in every running allocation as in A5.
- [ ] Merge `main` back into `dev` so `dev` has the CHANGELOG entry.
- [ ] Delete the feature branch, locally and on GitHub if it was pushed.
- [ ] Decide the periodic rebuild question and record the answer here.

## Part B: Dockerfile cleanup (biblio-backoffice only)

None of this is needed for the incident. It can ship after part A, in its own commits.

- [ ] Drop `RUN go get -v ./...`. `go build` downloads the modules itself; if a separate layer for downloads is wanted, `COPY go.mod go.sum ./` then `RUN go mod download` before `COPY . .` makes it cacheable.
- [ ] Copy only what the image needs into the final stage. `COPY --from=build /build .` puts the whole source tree in `/dist`. The image needs `app`, `tern`, `tern.docker.conf` and `db/migrations`; check what else the binary reads from disk at runtime (for example `static/` or `locales/`, unless they are embedded) before narrowing it.
- [ ] Use the `ENV key=value` form. The `ENV key value` form is deprecated and makes BuildKit print warnings.

## Retiring the plan

- [ ] Describe how the image is built (Docker Hub, `hooks/build`, the `SNAPSHOT_EPOCH` argument, pinned base images, how rebuilds are triggered) in the project's README, so the next person to bump Go or Alpine knows which lines to change.
- [ ] Move this plan to `docs/archive/` or delete it.
