.PHONY: build install snapshot dist test vet lint fmt run clean docker
OUT := keyvault-certdeploy
PKG := github.com/emgag/keyvault-certdeploy
PKG_LIST := $(shell go list ${PKG}/... )
GO_FILES := $(shell find . -name '*.go' )

all: build

build:
	CGO_ENABLED=0 GOOS=linux go build -a -v -o ${OUT} ${PKG}

install:
	CGO_ENABLED=0 GOOS=linux go install -a -v ${PKG}

snapshot:
	goreleaser --snapshot --rm-dist

dist:
	goreleaser --rm-dist

test:
	@go test -v ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

fmt:
	@gofmt -l -w -s ${GO_FILES}

clean:
	-@rm -vf ${OUT}
	-@rm -vrf dist


