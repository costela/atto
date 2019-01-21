# atto - 10⁻¹⁸ × webserver

`atto` is a really, _really_, **really** small webserver for static files. Barely more than [go](https://golang.org/)'s standard [`http.FileServer`](https://golang.org/pkg/net/http/#example_FileServer).

## Usage

The easiest way to use it is by creating your own image based on the `costela/atto` image:

```Dockerfile
FROM costela/atto
COPY my_staticfile_dir /www
```

Integrate this into your CD pipeline and deploy it on docker swarm or kubernetes behind a reverse proxy (e.g. [traefik](https://traefik.io/)).

## Features

Basically none, but it can and will gladly serve static files. It does not (and will not) support SSL, vHosts, aliases or any of the more advanced features of full-fledged webservers.

Nevertheless, it does have the basics:

- optional directory listing (inherited from `http.FileServer`; see `--showlist` below)
- handle running under some folder under `/` (see `--prefix` below).
- graceful shutdown to avoid disrupting connection during deployment (see `--timeout.shutdown` below)

## Configuration

The following settings may be provided as command line arguments or environment variables.

| Flag | Env-Var | Description | Default |
| --- | --- | --- | --- |
| `--loglevel`, `-l` | `ATTO_LOGLEVEL` | level of logging output (any value supported by [logrus](https://github.com/sirupsen/logrus)) | `warn` |
| `--port` | `ATTO_PORT` | TCP port on which to listen | `8080` |
| `--path` | `ATTO_PATH` | path which will be served | `/www` |
| `--prefix` | `ATTO_PREFIX` | prefix under which `path` will be accessed | _none_ |
| `--showlist` | `ATTO_SHOWLIST` | whether to display folder contents | `false` |
| `--timeout.readheader` | `ATTO_TIMEOUT_READHEADER` | time to wait for request headers | `5s` |
| `--timeout.shutdown` | `ATTO_TIMEOUT_SHUTDOWN` | time to wait for ungoing requests to finish before shutting down | `30s` |

## Motivation

The main use-case for `atto` is simplifying the build/deployment process of applications with static files, so that the same mental toolset can be used both for code and assets.

It is not a replacement for proper production deployment, but can ease the cognitive load when dealing with smaller, less often touched apps.

## License

MIT - See LICENSE file.