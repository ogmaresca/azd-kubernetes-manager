# Download trusted root certificates
FROM alpine:3.15 AS ca-certificates

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

# Download modules
FROM golang:1.17-alpine3.15 AS base

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

RUN go get

RUN go build -ldflags="-w -s" -o /go/bin/azd-kubernetes-manager

# Use a distroless final base
FROM scratch AS final

EXPOSE 10102
EXPOSE 10902

USER 1000

WORKDIR /home/azd-kubernetes-manager

ENV HOME /home/azd-kubernetes-manager

COPY --from=build --chown=1000 /go/bin/azd-kubernetes-manager /bin/azd-kubernetes-manager

ENTRYPOINT ["/bin/azd-kubernetes-manager"]
CMD []
