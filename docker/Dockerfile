# Multi-stage build for Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o calibre-rest .

# Final stage
FROM alpine:latest AS base

RUN apk add --no-cache \
    supervisor \
    nginx \
    # install calibre system dependencies (minimal for alpine)
    xdg-utils \
    xz \
    && rm -rf /var/cache/apk/*

# forward nginx log files to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log

ENV CALIBRE_CONFIG_DIRECTORY=/app/.calibre \
    CALIBRE_TEMP_DIR=/tmp \
    CALIBRE_CACHE_DIRECTORY=/tmp

RUN adduser -D -s /bin/sh calibre -u 1000 && \
    mkdir -p ${CALIBRE_CONFIG_DIRECTORY} && \
    chown -R 1000:1000 ${CALIBRE_CONFIG_DIRECTORY}

EXPOSE 80 443
CMD ["/usr/bin/supervisord"]

# App stage (without calibre binary)
FROM base as app

COPY --from=builder /app/calibre-rest /usr/local/bin/calibre-rest

COPY ./docker/supervisord-go.conf /etc/supervisor/conf.d/supervisord.conf
COPY --chmod=0755 ./docker/stop-supervisor.sh /etc/supervisor/stop-supervisor.sh

COPY ./docker/nginx.conf /etc/nginx/nginx.conf

WORKDIR /app

# Calibre builder stage
FROM base as calibre_builder
RUN apk add --no-cache wget && \
    wget -nv -O- \
    https://download.calibre-ebook.com/linux-installer.sh | sh /dev/stdin install_dir=/opt && \
    rm -rf /var/cache/apk/*

# Calibre stage (with calibre binary)
FROM calibre_builder as calibre

COPY --from=builder /app/calibre-rest /usr/local/bin/calibre-rest

COPY ./docker/supervisord-go.conf /etc/supervisor/conf.d/supervisord.conf
COPY --chmod=0755 ./docker/stop-supervisor.sh /etc/supervisor/stop-supervisor.sh

COPY ./docker/nginx.conf /etc/nginx/nginx.conf

WORKDIR /app