# How to run

## build
```bash
$ go build -o main
```

## Run
```bash
# Listening to ICMP packets requires root privileges
$ sudo ./main
```

## Usage
```bash
$ curl http://localhost:8080/trace/{ip}
```
```bash
$ curl http://localhost:8080/trace/{ip}?maxHops={max_hops}
```

## Example
```bash
$ curl http://localhost:8080/trace/8.8.8.8
```
```bash
$ curl http://localhost:8080/trace/google.com
```
```bash
$ curl http://localhost:8080/trace/8.8.8.8?maxHops=10
```