<p align="center">
    <h1 align="center">Banner rotation </h1>
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/Raschudesny/otus_project)](https://goreportcard.com/report/github.com/Raschudesny/otus_project)
![ci](https://github.com/Raschudesny/otus_project/actions/workflows/ci.yaml/badge.svg)

---
_Simple service for banner rotation purposes._
> This is OTUS "Golang Developer Professional" final project repository. Service technical specification can be [found here](docs/tz.md).
---

## Global description

This service purpose is to provide the most popular (in "clicks" metric) banner.

## Documentation

Some cool documentation here :)

## Build

See the [build guide](docs/build-guide.md).

## Commands
Available tasks for this project:

* **Alias for local:build command:** `task build`
* **Clean binary directory:** `task clean`
* **Build docker image of service** `task docker:build`
* **Runs the banner rotation service with all required environment(db, redis) in docker via docker-compose** `task docker:run`
* **Stop the bannerrotation service in docker** `task docker:stop`
* **Start [evans](https://github.com/ktr0731/evans) grpc  client:** `task evans`
* **Generate all required for project build code generated sources:** `task generate`
* **Install lint dependencies:** `task install-lint-deps`
* **Lint project:** `task lint`
* **Build the go binary locally:** `task local:build`
* **Run the banner_rotation service with default config(from configs directory) locally:** `task local:run`
* **Apply all migrations from migrations dir:** : `task migrations:up`
* **Rollback all migrations from migrations dir:**: `task migrations:down`
* **Alias for docker:run command:** `task run`
* **Alias for docker:stop command:** `task stop`
* **Run test with race detector:** `task test`
