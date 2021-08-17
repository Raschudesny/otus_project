# Build Guide

## Checkout the source code

```
$ git clone https://github.com/Raschudesny/otus_project.git
```

## Prerequisites

[Taskfile](https://github.com/go-task/task) util must be installed.
Taskfile util installation guide could be found [here](https://taskfile.dev/#/installation).
For example if you have Golang > 1.15 installed you can do the following:
```
$ go install github.com/go-task/task/v3/cmd/task@latest
```

## Build the binaries

Run the build script from the source directory:

```
$ cd otus_project
$ go install github.com/go-task/task/v3/cmd/task@latest
$ task build
```