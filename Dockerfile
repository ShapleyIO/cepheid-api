FROM golang:1.22.1 as base

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

RUN mkdir ~/.ssh
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts \
RUN git config --global url."git@github.com".insteadOf "https://github.com"

ENV PROJECT shapley-cepheid
WORKDIR $GOPATH/src/github.com/ShapleyIO/$PROJECT
RUN git config --global --add safe.directory $GOPATH/src/github.com/ShapleyIO/cepheid

FROM base as builder

RUN apt-get update && apt-get install --yes --quiet \
    netcat-openbsd unzip \
    && rm -rf /var/lib/apt/lists/*

ENV OAPI_CODEGEN_VERSION 1.13.2

VOLUME $GOPATH/src/github.com/ShapleyIO/$PROJECT

# Install Tools
COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY vendor vendor
COPY tools.go ./tools.go
# COPY bin/configure bin/configure
# RUN chmod 775 bin/configure
# RUN bin/configure
RUN GOFLAGS='' go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v${OAPI_CODEGEN_VERSION}

FROM base as golang-builder
ARG BUILD_DATE
ARG REVISION
ARG OAPI_CODEGEN_VERSION
COPY . .

FROM golang-builder as api-builder
RUN --mount=type=cache,target=/root/.cache/go-build bin/build

FROM alpine:3.18 as api
COPY --from=api-builder /go/src/github.com/ShapleyIO/shapley-cepheid/dist/cepheid-api /cepheid-api
ENV INTERFACE="[::]"
EXPOSE 80
ENTRYPOINT [ "/cepheid-api" ]
