# Download modules
FROM golang:1.13-alpine3.11 AS base

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /go/src

ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

COPY go.mod /go/src/go.mod

RUN go mod download

# Compile
FROM base AS build

COPY main.go /go/src/main.go
COPY pkg /go/src/pkg

RUN go build -ldflags="-w -s" -o /go/bin/azd-kubernetes-manager

# Use alpine as final base stage
FROM alpine:3.11 AS final

EXPOSE 10102
EXPOSE 10902

RUN adduser -D -g '' -u 1000 azd-kubernetes-manager

USER 1000

WORKDIR /home/azd-kubernetes-manager

COPY --from=build --chown=azd-kubernetes-manager /go/bin/azd-kubernetes-manager /bin/azd-kubernetes-manager

ENTRYPOINT ["/bin/azd-kubernetes-manager"]
CMD []
