FROM golang:1.23 AS builder

COPY . /src
WORKDIR /src

# Build with caching
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache
RUN --mount=type=cache,target=/gomod-cache \
  go mod download

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
  make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    netbase \
    tzdata \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y

RUN mkdir -p /var/log/luminex
RUN touch /var/log/luminex/luminex-service.log

COPY --from=builder /src/bin /app
COPY --from=builder /src/configs /app/data/conf/

WORKDIR /app
VOLUME /data/conf

EXPOSE 8000 9000

CMD ["./server", "-conf", "/app/data/conf/"] 