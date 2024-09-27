# ---------------------------------------------------------#
#                     Build web image                      #
# ---------------------------------------------------------#
FROM node:16 as web

WORKDIR /usr/src/app

COPY web/package.json ./
COPY web/yarn.lock ./

# If you are building your code for production
# RUN npm ci --omit=dev

COPY ./web .

RUN yarn && yarn build && yarn cache clean

# ---------------------------------------------------------#
#                   Build gitness image                    #
# ---------------------------------------------------------#
FROM golang:1.22-alpine3.18 as builder

RUN apk update \
    && apk add --no-cache protoc build-base git

# Setup workig dir
WORKDIR /app
RUN git config --global --add safe.directory '/app'

# Get dependancies - will also be cached if we won't change mod/sum
COPY go.mod .
COPY go.sum .

COPY Makefile .
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN make dep
RUN make tools
# COPY the source code as the last step
COPY . .

COPY --from=web /usr/src/app/dist /app/web/dist

# build
ARG GIT_COMMIT=
ARG GITNESS_VERSION_MAJOR=1
ARG GITNESS_VERSION_MINOR=0
ARG GITNESS_VERSION_PATCH=0
ARG TARGETOS TARGETARCH

RUN if [ "$TARGETARCH" = "arm64" ]; then \
    wget -P ~ https://github.com/yhsunshining/gitness/releases/download/0.0.1-alpha.0/aarch64-linux-musl-cross.tgz && \
    tar -xvf ~/aarch64-linux-musl-cross.tgz -C ~ ; \
    fi

# set required build flags
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    if [ "$TARGETARCH" = "arm64" ]; then CC=~/aarch64-linux-musl-cross/bin/aarch64-linux-musl-gcc; fi && \
    LDFLAGS="-X github.com/harness/gitness/version.GitCommit=1f7d86 -X github.com/harness/gitness/version.major=1 -X github.com/harness/gitness/version.minor=0 -X github.com/harness/gitness/version.patch=0 -extldflags '-static'" && \
    CGO_ENABLED=1 \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    CC=$CC go build -ldflags="$LDFLAGS" -o ./gitness ./cmd/gitness

### Pull CA Certs
FROM alpine:latest as cert-image

RUN apk --update add ca-certificates

# ---------------------------------------------------------#
#                   Create final image                     #
# ---------------------------------------------------------#
FROM --platform=$TARGETPLATFORM alpine/git:2.43.0 as final

# setup app dir and its content
WORKDIR /app
VOLUME /data

ENV XDG_CACHE_HOME /data
ENV GITNESS_GIT_ROOT /data
ENV GITNESS_DATABASE_DRIVER sqlite3
ENV GITNESS_DATABASE_DATASOURCE /data/database.sqlite
ENV GITNESS_METRIC_ENABLED=true
ENV GITNESS_METRIC_ENDPOINT=https://stats.drone.ci/api/v1/gitness
ENV GITNESS_TOKEN_COOKIE_NAME=token
ENV GITNESS_DOCKER_API_VERSION 1.41
ENV GITNESS_SSH_ENABLE=true
ENV GITNESS_GITSPACE_ENABLE=true

COPY --from=builder /app/gitness /app/gitness
COPY --from=cert-image /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 3000
EXPOSE 3022

ENTRYPOINT [ "/app/gitness", "server" ]
