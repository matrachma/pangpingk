OS      := $(shell uname -a | cut -f 1 -d ' ' | tr [:upper:] [:lower:])
ARCH    := $(shell uname -m)
TAG     := $(shell git describe --tags --always)
TIMESTAMP := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

all: build

build:
	@go install
	@cd cmd/pangpingk && go build -ldflags="-X main.appBuildTime=$(TIMESTAMP) -X main.appVersion=$(TAG)"

release: build
	@echo "Packaging pangpingk ${TAG} for ${OS}"
	@cd cmd/pangpingk && tar -czf pangpingk-${TAG}-${OS}-${ARCH}.tar.gz pangpingk

clean:
	@rm -f cmd/pangpingk/pangpingk-*.tar.gz cmd/pangpingk/pangpingk

buildall: clean build
