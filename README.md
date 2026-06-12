## GenX

GenX generates sample time-series data from one of four curve types:

- `cos` (default)
- `linear`
- `log`
- `exp`

It can run in two modes:

- CLI mode: prints `timestamp value` pairs to stdout
- Web mode: starts an HTTP preview UI

## Requirements

- Go 1.22+ (for local run)
- Docker (for container run)

## Local Run

Show available flags:

```bash
go run . --help
```

Generate CLI output (default is cosine):

```bash
go run .
```

Generate a linear series:

```bash
go run . -type linear -duration 1d -step 1h -first 10 -last 30
```

Run the web UI on port 8080:

```bash
go run . -web -port 8080
```

Then open:

```text
http://localhost:8080
```

## Flags

- `-type` curve type: `cos`, `linear`, `log`, `exp` (default `cos`)
- `-duration` total generation duration (default `1d`)
- `-step` sample interval (default `1h`)
- `-web` run HTTP web preview mode (default `false`)
- `-port` web server port (default `8080`)
- `-first` first value for linear type (default `0`)
- `-last` last value for linear type (default `1`)
- `-min` min value for cosine type (default `10`)
- `-max` max value for cosine type (default `25`)
- `-period` period for cosine type (default `1d`)

Duration values use a number plus unit: `d`, `h`, `m`, or `s`.
Examples: `2d`, `6h`, `30m`, `15s`.

## Docker

Build the image:

```bash
docker build -t genx:latest .
```

### Web mode in container

Run with an auto-published host port:

```bash
docker run -P --rm --name genx genx:latest -web
```

Find the published port:

```bash
docker port genx 8080
```

Run with a fixed host port:

```bash
docker run -p 8080:8080 --rm --name genx genx:latest -web
```

Open:

```text
http://localhost:8080
```

### CLI mode in container

```bash
docker run --rm genx:latest -type exp -duration 6h -step 30m
```

## Image layout notes

- The Dockerfile uses a multi-stage build.
- Final stage is `scratch` and copies only the compiled `/genx` binary.
- `EXPOSE 8080/tcp` documents the HTTP port used in web mode.
- `.dockerignore` limits build context to `go.mod`, `go.sum` (if present), and `*.go` files.

## Troubleshooting

- `go: command not found` on Windows PowerShell:
  Install Go and ensure the `go` binary is on your PATH, then restart the terminal.

- `docker build` fails because files are missing:
  Run the build from the project directory:

  ```bash
  cd genx-master
  docker build -t genx:latest .
  ```

- Web UI does not open with `-P`:
  Check the published host port first:

  ```bash
  docker port genx 8080
  ```

- `bind: address already in use` when using `-p 8080:8080`:
  Another process is using port 8080. Stop that process or use another host port, for example:

  ```bash
  docker run -p 8081:8080 --rm --name genx genx:latest -web
  ```

- Container exits immediately:
  This is expected if you run without `-web`; default command prints help and exits. Use `-web` for a long-running web server.
