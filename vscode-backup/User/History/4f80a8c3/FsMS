.PHONY: run build-docker

build:
	GOARCH=amd64 CGO_ENABLED=0 go build -o bin/on-air ./cmd/on-air/main.go

run:
	TZ=UTC go run ./cmd/on-air/main.go -e ./configs/env.local -local

clean:
	rm -R bin/*
	rm -Rf vendor

docker/build:
	go mod vendor
	docker build -t {REGISTRY-URL}/on-air -f Dockerfile .
	rm -Rf vendor

docker/run:
	docker compose up -d

docker/build-and-push-image: docker/build
	docker push {REGISTRY-URL}/on-air

cloudrun/deploy:
	gcloud run deploy on-air --image {REGISTRY-URL}/on-air

