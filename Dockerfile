# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/pokedex-clone/*.go ./
COPY pkg/ ./pkg/

RUN go build -o /docker-pokedex-clone

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /docker-pokedex-clone /docker-pokedex-clone

EXPOSE 5000

USER nonroot:nonroot

ENTRYPOINT ["/docker-pokedex-clone"]