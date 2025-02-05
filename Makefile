PROJECT_NAME := "fastburn"
PKG := "fastburn/cmd/fastburn"
PLATFORMS := linux windows
ARCHITECTURES := 386 amd64
BINARY :=fbn
LDFLAGS  :=
#LDFLAGS := "-ldflags XXX YYYY ZZZ"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean test lint

all: build

lint: ## Lint the files
	go vet ${PKG_LIST}
	staticcheck ${PKG_LIST}

test: ## Run unittests
	go test -short ${PKG_LIST}

race: dep ## Run data race detector
	go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	go test -msan -short ${PKG_LIST}

dep: ## Get the dependencies
	go get -v -d ./...

build: dep ## Build the binary file
	go build  -o ${BINARY} ${LDFLAGS} -v $(PKG) 

clean: ## Remove previous build
	rm -f $(BINARY)

cross:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell \
		export GOOS=$(GOOS); \
		export GOARCH=$(GOARCH); \
		go clean -modcache ;	\
		go build $(LDFLAGS) -v -o $(BINARY)-$(GOOS)-$(GOARCH)  $(PKG)  )))

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
