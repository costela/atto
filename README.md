[![Go Report Card](https://goreportcard.com/badge/github.com/costela/atto)](https://goreportcard.com/report/github.com/costela/atto)
[![Docker Hub Version](https://img.shields.io/badge/dynamic/json.svg?label=docker%20hub&url=https%3A%2F%2Findex.docker.io%2Fv1%2Frepositories%2Fcostela%2Fatto%2Ftags&query=%24[-1:].name&colorB=green)](https://hub.docker.com/r/costela/atto)

# atto - 10⁻¹⁸ × webserver

`atto` is a really, _really_, **really** small webserver for static files. Barely more than [go](https://golang.org/)'s standard [`http.FileServer`](https://golang.org/pkg/net/http/#example_FileServer).

## Usage

The easiest way to use it is by creating your own image based on the `costela/atto` image:

```Dockerfile
FROM costela/atto
COPY /path_to_my_staticfile_dir/ /www/
```

Integrate this into your CD pipeline and deploy it on e.g. [kubernetes](https://kubernetes.io/), behind a reverse proxy (e.g. [traefik](https://traefik.io/)).

## Features

Basically none, but it can and will gladly serve static files. It does not (and will not) support SSL, vHosts, aliases or any of the more advanced features of full-fledged webservers.

Nevertheless, it does have the basics:

- optional directory listing (inherited from `http.FileServer`; see `--showlist` option)
- optional transparent compression of content (see `--compress` option)
- handle running under some folder below `/` (see `--prefix` option)
- graceful shutdown to avoid disrupting long-running connections during deployment (see `--timeout.shutdown` option)

## Configuration

The following settings may be provided as command line arguments or environment variables.

| Flag | Env-Var | Description | Default |
| --- | --- | --- | --- |
| `--compress` | `ATTO_COMPRESS` | whether to transparently compress served files | `true` |
| `--loglevel`, `-l` | `ATTO_LOGLEVEL` | level of logging output (any value supported by [logrus](https://github.com/sirupsen/logrus)) | `warn` |
| `--port` | `ATTO_PORT` | TCP port on which to listen | `8080` |
| `--path` | `ATTO_PATH` | path which will be served | `.` |
| `--path404` | `ATTO_PATH` | path to a file returned when the requested content cannot be found | `404.html` |
| `--prefix` | `ATTO_PREFIX` | prefix under which `path` will be accessed | _none_ |
| `--showlist` | `ATTO_SHOWLIST` | whether to display folder contents | `false` |
| `--timeout.readheader` | `ATTO_TIMEOUT_READHEADER` | time to wait for request headers | `5s` |
| `--timeout.shutdown` | `ATTO_TIMEOUT_SHUTDOWN` | time to wait for ungoing requests to finish before shutting down | `30s` |

## Motivation

The main use-case for `atto` is simplifying the build/deployment process of applications with static files, so that the same mental toolset can be used both for code and assets.

It is not a replacement for proper production deployment, but can ease the cognitive load when dealing with smaller, less often touched apps.

## License

MIT - See LICENSE file.
