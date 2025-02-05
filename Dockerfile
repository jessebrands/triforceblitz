FROM golang:1.23.5-bookworm AS build
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /usr/src/triforceblitz
COPY go.mod go.sum /usr/src/triforceblitz/
RUN go mod download

COPY . /usr/src/triforceblitz
RUN go build -o /usr/local/bin/triforceblitz-server ./cmd/server
RUN go build -o /usr/local/bin/triforceblitz-updater ./cmd/updater

FROM build AS test
WORKDIR /usr/src/triforceblitz
RUN go test -v ./...

FROM debian:bookworm-slim AS environment
RUN apt-get update -y && apt-get install -y \
    ca-certificates \
    python3

RUN mkdir -p /usr/local/share/triforceblitz/generators

RUN useradd --system --shell /bin/bash triforceblitz
RUN chown triforceblitz:triforceblitz -R /usr/local/share/triforceblitz

FROM environment AS release
ENV TRIFORCEBLITZ_GENERATORS_DIR=/usr/local/share/triforceblitz/generators

COPY --from=build /usr/local/bin/triforceblitz-* /usr/local/bin/

USER triforceblitz:triforceblitz

EXPOSE 8000
ENTRYPOINT ["triforceblitz-server"]
