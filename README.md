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

* **Build the go binary:** `task build`
* **Clean binary directory:** `task clean`
* **Run the banner_rotation service with default config(from configs directory) locally:** `task run-locally`
* **Run the banner_rotation service with default config(from configs directory) in docker:** `task run`
* **Stop the banner_rotation service in docker:** `task stop`
* **Run test with race detector:** `task test`
* **Generate all required for project build code generated sources:** `task generate`
* **Lint project:** `task lint`
* **Install lint dependencies:** `task install-lint-deps`
* **Rollback all migrations from migrations dir:**: `task migrate_down`
* **Apply all migrations from migrations dir:** : `task migrate_up`
* **Start [evans](https://github.com/ktr0731/evans) grpc  client:** `task evans`
