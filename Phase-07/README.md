# How to run

# Login to ghcr.io

```bash
$ docker login ghcr.io
```

## Pull

```bash
$ docker compose pull
```

## Run

```bash
$ docker compose up
```

## Test the API

```bash
$ curl -X GET http://localhost:9123/trace/www.google.com?maxHops=15
```