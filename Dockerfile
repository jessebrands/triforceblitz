# Common build environment shared by all commands.
#
# Each command has its own build stage to improve build times. This also
# allows for better caching, as rebuilds will only build the command that has
# changed.
FROM golang:1.23.5-bookworm AS build-environment
WORKDIR /usr/src/triforceblitz

# We cache all project dependencies here so they only need to be downloaded 
# whenever the dependencies change.
COPY go.mod go.sum /usr/src/triforceblitz/
RUN go mod download

# Copy common sources shared by all commands.
COPY ./internal /usr/src/triforceblitz/internal

# cmd/triforceblitz
FROM build-environment AS cli-build
ENV CGO_ENABLED=0

COPY ./cmd/triforceblitz /usr/src/triforceblitz/cmd/triforceblitz
RUN go build -o /usr/local/bin/triforceblitz ./cmd/triforceblitz

# cmd/server
FROM build-environment AS server-build
ENV CGO_ENABLED=0

COPY ./cmd/server /usr/src/triforceblitz/cmd/server
RUN go build -o /usr/local/bin/triforceblitz-server ./cmd/server


# cmd/updater
FROM build-environment AS updater-build 
ENV CGO_ENABLED=0

COPY ./cmd/updater /usr/src/triforceblitz/cmd/updater
RUN go build -o /usr/local/bin/triforceblitz-updater ./cmd/updater


# Runtime environment. This stage sets up the container for the server by
# installing libraries and programs that need to be present, as well as
# setting up a non-privileged user and permissions.
#
# Furthermore, it runs triforceblitz-updater to install generators.
FROM debian:bookworm-slim AS environment
ENV TRIFORCEBLITZ_GENERATORS_DIR=/usr/local/share/triforceblitz/generators
ENV TRIFORCEBLITZ_PACKAGE_CACHE_DIR=/var/cache/triforceblitz/packages
ENV TRIFORCEBLITZ_LOCK_FILE=/var/lock/triforceblitz.lock

RUN apt-get update -y && apt-get install -y \
    ca-certificates \
    python3

COPY --from=updater-build /usr/local/bin/triforceblitz-updater /usr/local/bin/

RUN triforceblitz-updater install -no-cache -b blitz
RUN useradd --system --shell /bin/bash triforceblitz


# Release stage. This is the actual image. We copy over the built server
# binary and run it.
FROM environment AS release
ENV TRIFORCEBLITZ_GENERATORS_DIR=/usr/local/share/triforceblitz/generators
ENV TRIFORCEBLITZ_PACKAGE_CACHE_DIR=/var/cache/triforceblitz/packages

COPY --from=cli-build /usr/local/bin/triforceblitz /usr/local/bin/
COPY --from=server-build /usr/local/bin/triforceblitz-server /usr/local/bin/

USER triforceblitz:triforceblitz
EXPOSE 8000
ENTRYPOINT ["triforceblitz-server"]
