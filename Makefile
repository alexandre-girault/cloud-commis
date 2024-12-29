
all: build

.PHONY: build
build:
	mkdir -p bin
	CGO_ENABLED=0 go build \
	-ldflags=" \
	-X 'cloud-commis/config.Version=$$(cat version)-devel'\
	-X 'cloud-commis/config.BuildDate=$$(date)'" \
	-o bin/cloudcommis main.go

test:
	go test ./... -cover

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.62.2 golangci-lint run -v

.PHONY: clean
clean:
	@rm -Rf bin/*