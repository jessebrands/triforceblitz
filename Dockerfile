FROM golang:1.23.5-bookworm AS build-stage
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /usr/src/triforceblitz
COPY go.mod /usr/src/triforceblitz/
RUN go mod download

COPY . /usr/src/triforceblitz
RUN go build -o /usr/local/bin/triforceblitz-updater ./cmd/updater

FROM build-stage AS test-stage
WORKDIR /usr/src/triforceblitz
RUN go test -v ./...
