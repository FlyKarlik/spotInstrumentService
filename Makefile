include .env

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o ./spot-instrument-service ./cmd/spot_instrument_service/main.go

.PHONY: prepare
prepare:
	go mod download

.PHONY: spot-instrument-service
spot-instrument-service:
	./spot-instrument-service

.PHONY: clean
clean:
	rm ./spot-instrument-service

.PHONY: lint
lint:
	golangci-lint run ./... --fix

.PHONY: test
test:
	go test -v ./... --cover

.PHONE: docker-build
docker-build:
	docker build -t spot-instrument-service .