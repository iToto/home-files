## builder
FROM golang:1.18 as builder
ARG ENV=docker
WORKDIR /my-app/

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
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -v ./cmd/my-app

# target image
FROM alpine:3
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /my-app/my-app /my-app
COPY --from=builder /my-app/config.env /my-app/config.env

#default entry point for service
ENTRYPOINT ["/my-app/my-app"]
CMD ["-e", "/my-app/config.env", "-local"]
