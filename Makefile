# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

GOBIN = $(shell pwd)/build/bin
GO ?= latest
BUILD_PARAM?=install 

OBJECTS=kcn kpn ken kscn kspn ksen kbn kgen homi
RPM_OBJECTS=$(foreach wrd,$(OBJECTS),rpm-$(wrd))
RPM_BAOBAB_OBJECTS=$(foreach wrd,$(OBJECTS),rpm-baobab-$(wrd))
TAR_LINUX_386_OBJECTS=$(foreach wrd,$(OBJECTS),tar-linux-386-$(wrd))
TAR_LINUX_amd64_OBJECTS=$(foreach wrd,$(OBJECTS),tar-linux-amd64-$(wrd))
TAR_DARWIN_amd64_OBJECTS=$(foreach wrd,$(OBJECTS),tar-darwin-amd64-$(wrd))
TAR_BAOBAB_LINUX_386_OBJECTS=$(foreach wrd,$(OBJECTS),tar-baobab-linux-386-$(wrd))
TAR_BAOBAB_LINUX_amd64_OBJECTS=$(foreach wrd,$(OBJECTS),tar-baobab-linux-amd64-$(wrd))
TAR_BAOBAB_DARWIN_amd64_OBJECTS=$(foreach wrd,$(OBJECTS),tar-baobab-darwin-amd64-$(wrd))

.PHONY: all test clean ${OBJECTS} ${RPM_OBJECTS} ${TAR_LINUX_386_OBJECTS} ${TAR_DARWIN_amd64_OBJECTS} ${TAR_LINUX_amd64_OBJECTS}

all: ${OBJECTS}
rpm-all: ${RPM_OBJECTS}
rpm-baobab-all: ${RPM_BAOBAB_OBJECTS}
tar-linux-386-all: ${TAR_LINUX_386_OBJECTS}
tar-linux-amd64-all: ${TAR_LINUX_amd64_OBJECTS}
tar-darwin-amd64-all: ${TAR_DARWIN_amd64_OBJECTS}
tar-baobab-linux-386-all: ${TAR_BAOBAB_LINUX_386_OBJECTS}
tar-baobab-linux-amd64-all: ${TAR_BAOBAB_LINUX_amd64_OBJECTS}
tar-baobab-darwin-amd64-all: ${TAR_BAOBAB_DARWIN_amd64_OBJECTS}

${OBJECTS}:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/$@

${RPM_OBJECTS}:
	./build/package-rpm.sh ${@:rpm-%=%}

${RPM_BAOBAB_OBJECTS}:
	./build/package-rpm.sh -b ${@:rpm-baobab-%=%}

${TAR_LINUX_386_OBJECTS}:
	$(eval BIN := ${@:tar-linux-386-%=%})
	./build/cross-compile.sh linux-386 ${BIN}
	./build/package-tar.sh linux-386 ${BIN}

${TAR_LINUX_amd64_OBJECTS}:
	$(eval BIN := ${@:tar-linux-amd64-%=%})
	./build/cross-compile.sh linux-amd64 ${BIN}
	./build/package-tar.sh linux-amd64 ${BIN}

${TAR_DARWIN_amd64_OBJECTS}:
	$(eval BIN := ${@:tar-darwin-amd64-%=%})
	./build/cross-compile.sh darwin-amd64 ${BIN}
	./build/package-tar.sh darwin-amd64 ${BIN}

${TAR_BAOBAB_LINUX_386_OBJECTS}:
	$(eval BIN := ${@:tar-baobab-linux-386-%=%})
	./build/cross-compile.sh linux-386 ${BIN}
	./build/package-tar.sh -b linux-386 ${BIN}

${TAR_BAOBAB_LINUX_amd64_OBJECTS}:
	$(eval BIN := ${@:tar-baobab-linux-amd64-%=%})
	./build/cross-compile.sh linux-amd64 ${BIN}
	./build/package-tar.sh -b linux-amd64 ${BIN}

${TAR_BAOBAB_DARWIN_amd64_OBJECTS}:
	$(eval BIN := ${@:tar-baobab-darwin-amd64-%=%})
	./build/cross-compile.sh darwin-amd64 ${BIN}
	./build/package-tar.sh -b darwin-amd64 ${BIN}

abigen:
	build/env.sh go run build/ci.go ${BUILD_PARAM} ./cmd/abigen
	@echo "Done building."
	@echo "Run \"$(GOBIN)/abigen\" to launch abigen."

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

# Not supported. Use lint-try instead of lint
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
	