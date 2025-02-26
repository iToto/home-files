## builder
FROM golang:1.23 as builder
ARG ENV=docker
WORKDIR /on-air/

# copy sources
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY go.* ./

# copy env file
RUN echo "copying env.${ENV}"
COPY configs/env.${ENV} ./config.env

# copy vendor folder
COPY vendor ./vendor

# build
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -v ./cmd/on-air

# target image
FROM alpine:3
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /on-air/on-air /on-air/on-air
COPY --from=builder /on-air/config.env /on-air/config.env

#default entry point for service
ENTRYPOINT ["/on-air/on-air"]
CMD ["-e", "/on-air/config.env", "-local"]
