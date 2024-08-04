# <Build stage>
FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /app

RUN adduser -D -g '' -u 10001 builder
RUN chown -R builder:builder /app
USER builder

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./traceroute-api
# </Build stage>

# <Final stage>
FROM alpine:3.20.0 AS final

WORKDIR /app

LABEL org.opencontainers.image.source https://github.com/Star-Academy/Summer1403-Devops-Team12

COPY --from=build --chown=root:root /app/traceroute-api ./traceroute-api

EXPOSE 8080

CMD ["./traceroute-api"]
# </Final stage>
