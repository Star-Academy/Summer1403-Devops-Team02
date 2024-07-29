# <Build stage>
FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./traceroute-api
# </Build stage>

# <Final stage>
FROM alpine:3.20.0 AS final

WORKDIR /app

COPY --from=build /app/traceroute-api ./traceroute-api

EXPOSE 8080

CMD ["./traceroute-api"]
# </Final stage>
