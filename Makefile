BOOTSTRAP_VERSION=5.3.3
-include .env
export
APP_VERSION := $(shell git describe --tags)
LD_FLAGS="-X 'cloud-commis/config.Version=$(APP_VERSION)' -X 'cloud-commis/config.BuildDate=$(shell date)'"
CGO_ARGS=CGO_ENABLED=0 GOOS=linux


all: build

.PHONY: build

bin/cloudcommis-linux-amd64:
	mkdir -p bin
	$(CGO_ARGS) GOARCH=amd64 go build \
	-ldflags=$(LD_FLAGS) \
	-o bin/cloudcommis-linux-amd64 main.go

bin/cloudcommis-linux-arm64:
	mkdir -p bin
	$(CGO_ARGS) GOARCH=arm64 go build \
	-ldflags=$(LD_FLAGS) \
	-o bin/cloudcommis-linux-arm64 main.go

build: bin/cloudcommis-linux-amd64 bin/cloudcommis-linux-arm64
	@echo build version $(APP_VERSION)


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
	
.PHONY: docker-buildx-config
docker-buildx-config:
	if docker buildx inspect localBuilder ; then \
	echo "buildx builder OK" ; \
	else echo "creating buildx builder " ; \
	docker buildx create --driver docker-container --bootstrap --name localBuilder --use; fi
	
.PHONY: docker-build
docker-build: docker-buildx-config
	docker buildx build --no-cache --build-arg TARGET_ARCH=arm64 --provenance false --tag alexandregirault/cloud-commis:$(APP_VERSION)-devel-arm64 \
	--output type=image .
	docker buildx build --no-cache --build-arg TARGET_ARCH=amd64 --provenance false --tag alexandregirault/cloud-commis:$(APP_VERSION)-devel-amd64 \
	--output type=image .

docker-push: docker-build
	docker push alexandregirault/cloud-commis:$(APP_VERSION)-devel-arm64
	docker push alexandregirault/cloud-commis:$(APP_VERSION)-devel-amd64
	docker manifest create alexandregirault/cloud-commis:$(APP_VERSION) \
	--amend alexandregirault/cloud-commis:$(APP_VERSION)-arm64 \
	--amend alexandregirault/cloud-commis:$(APP_VERSION)-amd64
	docker manifest push alexandregirault/cloud-commis:$(APP_VERSION)

.PHONY: run
run:
	go run main.go --config=config.yaml

.PHONY: clean
clean:
	@rm -Rf bin/*
	@rm -Rf tmp/*

docker-clean:
	docker rm -vf $$(docker ps -aq)
	docker rmi -f $$(docker images -aq)
	docker buildx rm localBuilder