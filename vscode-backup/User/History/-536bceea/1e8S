.PHONY: run build-docker

build:
	GOARCH=amd64 CGO_ENABLED=0 go build -o bin/my-app ./cmd/my-app/main.go

run:
	TZ=UTC go run ./cmd/my-app/main.go -e ./configs/env.local -local

clean:
	rm -R bin/*
	rm -Rf vendor

docker/build:
	go mod vendor
	docker build -t {REGISTRY-URL}/my-app -f Dockerfile .
	rm -Rf vendor

docker/run:
	docker compose up -d

docker/build-and-push-image: docker/build
	docker push {REGISTRY-URL}/my-app

cloudrun/deploy:
	gcloud run deploy my-app --image {REGISTRY-URL}/my-app

