# azactl

A simple cli to launch container dev environment based on profile using [docker compose](https://docs.docker.com/compose/) behind the scene.
- Jekyll
- Golang
- Helm

## Prerequisites

- Docker compose

## Usage

```sh
azactl run -p [PROFILE]
```

### Help

```sh
azactl run -h

Run a container based environment of the desired profile.
Available profile:
- go
- helm
- jekyll

Usage:
  azactl run [flags]
Flags:
  -h, --help             help for run
  -p, --profile string   Profile tu use [jekyll/go/helm] (required)
  -v, --version string   Image version to use (default latest)
```
