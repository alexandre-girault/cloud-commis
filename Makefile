BOOTSTRAP_VERSION=5.3.3

all: build

.PHONY: build
build:
	mkdir -p bin
	CGO_ENABLED=0 go build \
	-ldflags=" \
	-X 'cloud-commis/config.Version=$$(cat version)-devel'\
	-X 'cloud-commis/config.BuildDate=$$(date)'" \
	-o bin/cloudcommis main.go


.PHONY: frontend-libs
frontend-libs: webui/statics/css/bootstrap.min.css webui/statics/js/bootstrap.bundle.min.js webui/statics/js/htmx.min.js

tmp/bootstrap-$(BOOTSTRAP_VERSION)-dist.zip:
	mkdir -p ./tmp
	curl -Lo ./tmp/bootstrap-$(BOOTSTRAP_VERSION)-dist.zip https://github.com/twbs/bootstrap/releases/download/v$(BOOTSTRAP_VERSION)/bootstrap-$(BOOTSTRAP_VERSION)-dist.zip

webui/statics/css/bootstrap.min.css webui/statics/js/bootstrap.bundle.min.js: tmp/bootstrap-$(BOOTSTRAP_VERSION)-dist.zip
	unzip -n -d ./tmp/ ./tmp/bootstrap-$(BOOTSTRAP_VERSION)-dist.zip
	cp ./tmp/bootstrap-$(BOOTSTRAP_VERSION)-dist/css/bootstrap.min.css webui/statics/css/bootstrap.min.css
	cp ./tmp/bootstrap-$(BOOTSTRAP_VERSION)-dist/js/bootstrap.bundle.min.js webui/statics/js/bootstrap.bundle.min.js

webui/statics/js/htmx.min.js:
	curl -Lo ./tmp/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js
	cp ./tmp/htmx.min.js webui/statics/js/htmx.min.js

test:
	CC_loglevel="error" go test -v ./... -cover

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.62.2 golangci-lint run -v

.PHONY: clean
clean:
	@rm -Rf bin/*
	@rm -Rf tmp/*