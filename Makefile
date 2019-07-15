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

kcn:
	build/env.sh go run build/ci.go install ./cmd/kcn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kcn\" to launch Consensus Node."

kpn:
	build/env.sh go run build/ci.go install ./cmd/kpn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kpn\" to launch Proxy Node."

ken:
	build/env.sh go run build/ci.go install ./cmd/ken
	@echo "Done building."
	@echo "Run \"$(GOBIN)/ken\" to launch Endpoint Node."

kbn:
	build/env.sh go run build/ci.go install ./cmd/kbn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kbn\" to launch bootnode."

kscn:
	build/env.sh go run build/ci.go install ./cmd/kscn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kscn\" to launch ServiceChain Consensus Node."

kspn:
	build/env.sh go run build/ci.go install ./cmd/kspn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kspn\" to launch ServiceChain Proxy Node."

ksen:
	build/env.sh go run build/ci.go install ./cmd/ksen
	@echo "Done building."
	@echo "Run \"$(GOBIN)/ksen\" to launch ServiceChain Endpoint Node."

kgen:
	build/env.sh go run build/ci.go install ./cmd/kgen
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kgen\" to launch kgen."

homi:
	build/env.sh go run build/ci.go install ./cmd/homi
	@echo "Done building."
	@echo "Run \"$(GOBIN)/homi\" to launch homi."

abigen:
	build/env.sh go run build/ci.go install ./cmd/abigen
	@echo "Done building."
	@echo "Run \"$(GOBIN)/abigen\" to launch abigen."

all:
	build/env.sh go run build/ci.go install

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
	build/env.sh go run build/ci.go test -coverage
	go tool cover -func=coverage.out -o coverage_report.txt
	go tool cover -html=coverage.out -o coverage_report.html
	@echo "Two coverage reports coverage_report.txt and coverage_report.html are generated."

fmt:
	build/env.sh go run build/ci.go fmt

lint:
	build/env.sh go run build/ci.go lint

lint-try:
	build/env.sh go run build/ci.go lint-try

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/* build/_workspace/src/

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

klay-cross: klay-linux klay-darwin klay-windows
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/klay-* $(GOBIN)/k*n-*

klay-linux: klay-linux-386 klay-linux-amd64 klay-linux-arm klay-linux-mips64 klay-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/klay-* $(GOBIN)/k*n-*

klay-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/kbn
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep 386

klay-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/kgen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/kbn
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep amd64

klay-linux-arm: klay-linux-arm-5 klay-linux-arm-6 klay-linux-arm-7 klay-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep arm

klay-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/kbn
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep arm-5

klay-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/kbn
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep arm-6

klay-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/kbn
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep arm-7

klay-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/kbn
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep arm64

klay-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/kbn
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep mips

klay-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/kbn
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep mipsle

klay-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/kbn
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep mips64

klay-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/kbn
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/k*n-linux-* | grep mips64le

klay-darwin: klay-darwin-386 klay-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/k*n-darwin-*

klay-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/kbn
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-darwin-* | grep 386

klay-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin-10.10/amd64 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin-10.10/amd64 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin-10.10/amd64 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin-10.10/amd64 -v ./cmd/kgen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin-10.10/amd64 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin-10.10/amd64 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin-10.10/amd64 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin-10.10/amd64 -v ./cmd/kbn
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-darwin-* | grep amd64

klay-windows: klay-windows-386 klay-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/k*n-windows-*

klay-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/kbn
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-windows-* | grep 386

klay-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/kcn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/kpn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/ken
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/kgen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/kscn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/kspn
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/ksen
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/kbn
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/k*n-windows-* | grep amd64
