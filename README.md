# calibre-rest (Go Version)

calibre-rest wraps
[calibredb](https://manual.calibre-ebook.com/generated/en/calibredb.html) to
provide a simple REST API server for your [Calibre](https://calibre-ebook.com/)
library.

**Note: This has been converted from Python to Go while maintaining the same REST API interface.**

### Disclaimer

- calibre-rest is in pre-alpha and subject to bugs and breaking changes. Please
use it at your own risk.
- This project has been tested on `amd64` Linux systems with Calibre 6.21 only.
- Contributions for testing and support on other OS platforms and Calibre versions
are greatly welcome.

## Overview

calibre-rest is a self-hosted REST API server that wraps `calibredb` to expose a
Calibre library. I wrote this project as I could not find a good
language-agnostic method to programmatically manipulate a Calibre library
(locally or remotely).

```bash
# get metadata with book id
$ curl localhost:5000/books/1

# query books with title
$ curl --get --data-urlencode "search=title:~^foo.*bar$" \
    http://localhost:5000/books

# add ebook file to library
$ curl -X POST -H "Content-Type:multipart/form-data" \
    --form "file=@foo.epub" http://localhost:5000/books

# download ebook from library (not yet implemented in Go version)
$ curl http://localhost:5000/export/1 -o bar.epub
```

See [API.md](API.md) for
full documentation of all API endpoints and examples.

calibre-rest is not meant to be a direct replacement for [Calibre Content
Server](https://manual.calibre-ebook.com/server.html). It does not have any
frontend interface and has no native access to a remote Calibre library without
an existing `calibre-server` instance. It is best used as a alternative backend
for any scripts or clients that wish to access the library remotely.

## Install

calibre-rest requires the following dependencies:

- A Calibre library on the local filesystem or served by a [Calibre content
  server](https://manual.calibre-ebook.com/generated/en/calibre-server.html)
- Go 1.19+ (for building from source)
- The `calibre` binary with the `calibredb` executable.
- Calibre's system dependencies (on Linux):

```
$ apt-get install xdg-utils, xz-utils, libopengl0, libegl1
```

The latter two are only relevant if you wish to run calibre-rest directly on
your local machine.

### Docker

Docker is the recommended method to run calibre-rest. We ship two
images:

- `kencx/calibre_rest:[version]-app` packaged without the calibre binary
- `kencx/calibre_rest:[version]-calibre` packaged with the calibre binary

The former image assumes you have an existing Calibre binary installation on
your local machine, server or Docker container (how else did you run Calibre
previously?). The binary's directory must be bind mounted to the running
container:

```yaml
version: '3.6'
services:
  calibre_rest:
    image: ghcr.io/kencx/calibre_rest:0.1.0-app
    environment:
      - "CALIBRE_REST_LIBRARY=/library"
    ports:
      - 8080:80
    volumes:
      - "/opt/calibre:/opt/calibre"
      - "./library:/library"
```

Or when paired with an existing
[linuxserver/docker-calibre](https://github.com/linuxserver/docker-calibre)
instance:

```yml
version: '3.6'
services:
  calibre:
    image: lscr.io/linuxserver/calibre
    volumes:
      - "./calibre:/opt/calibre"
      - "./library:/library"

  calibre_rest:
    image: ghcr.io/kencx/calibre_rest:0.1.0-app
    environment:
      - "CALIBRE_REST_LIBRARY=/library"
    ports:
      - 8080:80
    volumes:
      - "./calibre:/opt/calibre"
      - "./library:/library"
```

Otherwise, the larger `kencx/calibre_rest:[version]-calibre` image ships with
its own Calibre binary and only requires access to your existing Calibre
library directory.

### Build from Source

To run calibre-rest on your local machine, Calibre and its dependencies must be installed:

```console
# clone the repository
$ git clone git@github.com:kencx/calibre_rest.git
$ cd calibre_rest

# build the Go binary
$ go build -o calibre-rest .

# install calibre
$ apt-get install xdg-utils, xz-utils, libopengl0, libegl1
$ wget -nv -O- https://download.calibre-ebook.com/linux-installer.sh | sudo sh /dev/stdin
```

## Usage

Run the server with Docker:

```console
$ docker run \
    -v "/opt/calibre:/opt/calibre" \
    -v "./library:/library" \
    -p 8080:80 \
    -e "CALIBRE_REST_LIBRARY=/library" \
    ghcr.io/kencx/calibre_rest:0.1.0-app
```

or with [docker-compose](docker-compose.yml):

```console
$ docker compose up -d app
```

or directly on your local machine:

```console
$ ./calibre-rest -h

Usage: ./calibre-rest [options]

Options:
  -bind string
        Bind address HOST:PORT
  -calibre string
        Path to calibre binary directory
  -dev
        Start in dev/debug mode
  -help
        Show help
  -library string
        Path to calibre library
  -log-level string
        Log level (DEBUG, INFO, WARNING, ERROR)
  -password string
        Calibre library password
  -username string
        Calibre library username
  -version
        Print version
```

calibre-rest can access any local Calibre libraries or remote [Calibre content
server](https://manual.calibre-ebook.com/generated/en/calibre-server.html)
instances.

For the latter, authentication must be enabled and configured.
For more information, refer to the [calibredb
documentation](https://manual.calibre-ebook.com/generated/en/calibredb.html).

### Configuration

The server can be configured with the following environment variables.

| Env Variable    | Description    | Type    | Default    |
|---------------- | --------------- | --------------- | --------------- |
| `CALIBRE_REST_PATH`    | Path to `calibredb` executable    | string | `/opt/calibre/calibredb` |
| `CALIBRE_REST_LIBRARY` | Path to calibre library   | string   | `/library`   |
| `CALIBRE_REST_USERNAME` | Calibre library username  | string  |  |
| `CALIBRE_REST_PASSWORD` | Calibre library password   | string  |  |
| `CALIBRE_REST_LOG_LEVEL` | Log Level | string  | `INFO`   |
| `CALIBRE_REST_ADDR` | Server bind address | string   | `localhost:5000` |

If running directly on your local machine, we can also use flags:

```console
$ ./calibre-rest --bind localhost:5000
```

## Development

calibre-rest is built with Go 1.19+ and uses the Gorilla Mux router. Calibre should be installed to
facilitate testing.

To contribute, clone the repository and [build from source](#build-from-source).

```console
# build the application
$ go build -o calibre-rest .

# run the dev server
$ ./calibre-rest --dev

# run tests
$ go test -v
```

## Migration from Python

This version maintains the same REST API interface as the Python version, so existing clients should continue to work without changes. The main differences are:

- Built with Go instead of Python/Flask
- No longer uses Gunicorn (Go has a built-in HTTP server)
- Some endpoints are not yet fully implemented (marked as "not implemented")
- Improved performance and lower memory usage
- Single binary deployment

## Current Implementation Status

- [x] Basic HTTP server and routing
- [x] Health endpoint (`GET /health`)
- [x] Get single book (`GET /books/{id}`) 
- [x] Get books with pagination (`GET /books`)
- [x] Add book from file (`POST /books`)
- [x] Add empty book (`POST /books/empty`)
- [ ] Update book metadata (`PUT /books/{id}`) - marked as not implemented
- [ ] Delete book (`DELETE /books/{id}`) - marked as not implemented
- [ ] Export book (`GET /export/{id}`) - marked as not implemented
- [x] Configuration management
- [x] Docker support
- [x] CLI interface

## Roadmap

- [x] Support remote libraries
- [x] Pagination
- [ ] Complete all CRUD operations (update, delete)
- [ ] Export functionality
- [ ] TLS support
- [ ] Authentication
- [ ] Feature parity with `calibredb`
- [ ] S3 support