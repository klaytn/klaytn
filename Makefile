# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: klay-cross all test clean
.PHONY: klay-linux klay-linux-386 klay-linux-amd64 klay-linux-mips64 klay-linux-mips64le
.PHONY: klay-linux-arm klay-linux-arm-5 klay-linux-arm-6 klay-linux-arm-7 klay-linux-arm64
.PHONY: klay-darwin klay-darwin-386 klay-darwin-amd64
.PHONY: klay-windows klay-windows-386 klay-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest
BUILD_PARAM?=install

kcn:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/kcn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kcn\" to launch Consensus Node."

kpn:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/kpn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kpn\" to launch Proxy Node."

ken:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/ken
	@echo "Done building."
	@echo "Run \"$(GOBIN)/ken\" to launch Endpoint Node."

kbn:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/kbn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kbn\" to launch bootnode."

kscn:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/kscn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kscn\" to launch ServiceChain Consensus Node."

kspn:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/kspn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kspn\" to launch ServiceChain Proxy Node."

ksen:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/ksen
	@echo "Done building."
	@echo "Run \"$(GOBIN)/ksen\" to launch ServiceChain Endpoint Node."

kgen:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/kgen
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kgen\" to launch kgen."

homi:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/homi
	@echo "Done building."
	@echo "Run \"$(GOBIN)/homi\" to launch homi."

abigen:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/abigen
	@echo "Done building."
	@echo "Run \"$(GOBIN)/abigen\" to launch abigen."

all:
	build/env.sh go run build/ci.go ${BUILD_PARAM}

test:
	build/env.sh go run build/ci.go test

test-seq:
	build/env.sh go run build/ci.go test -p 1

test-datasync:
	build/env.sh go run build/ci.go test -p 1 ./datasync/...

test-networks:
	build/env.sh go run build/ci.go test -p 1 ./networks/...

test-tests:
	build/env.sh go run build/ci.go test -p 1 ./tests/...

test-others:
	build/env.sh go run build/ci.go test -p 1 -exclude datasync,networks,tests

cover:
	build/env.sh go run build/ci.go cover -coverprofile=coverage.out
	go tool cover -func=coverage.out -o coverage_report.txt
	go tool cover -html=coverage.out -o coverage_report.html
	@echo "Two coverage reports coverage_report.txt and coverage_report.html are generated."

fmt:
	GOFLAGS= GO111MODULE=off build/env.sh go run build/ci.go fmt

# Not supported. Use lint-try intead of lint
#lint:
#	build/env.sh env GOFLAGS= GO111MODULE=off go run build/ci.go lint

lint-try:
	GOFLAGS= GO111MODULE=off build/env.sh go run build/ci.go lint-try

clean:
	./build/clean_go_build_cache.sh
	chmod -R +w ./build/_workspace/pkg/
	rm -fr build/_workspace/pkg/ $(GOBIN)/* build/_workspace/src/

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOFLAGS= GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOFLAGS= GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOFLAGS= GOBIN= go get -u github.com/fjl/gencodec
	env GOFLAGS= GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOFLAGS= GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
