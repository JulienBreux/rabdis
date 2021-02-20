# BUILD
FROM golang:1.16-alpine3.13 as build

ARG VERSION=dev
ARG COMMIT=n/a
ARG RAW_DATE=n/a

RUN apk --update upgrade \
    && apk --no-cache --no-progress add git ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /tmp/app

# Prepare modules.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy application.
COPY . .

# Build the Go app
RUN go build \
        -ldflags="-s -w -X 'github.com/julienbreux/rabdis/pkg/version.Version=${VERSION}' -X 'github.com/julienbreux/rabdis/pkg/version.Commit=${COMMIT}' -X 'github.com/julienbreux/rabdis/pkg/version.RawDate=${RAW_DATE}'" \
        -o ./bin/rabdis ./cmd/rabdis

## IMAGE
FROM alpine:3.13

RUN apk --no-cache --no-progress add ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

COPY --from=build /tmp/app/bin/rabdis /usr/local/bin/rabdis

RUN addgroup -g 1000 -S rabdis && \
    adduser -u 1000 -S rabdis -G rabdis

# Use an unprivileged user.
USER rabdis:rabdis

EXPOSE 9090

ENTRYPOINT ["/usr/local/bin/rabdis"]
