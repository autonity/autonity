# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: autonity android ios autonity-cross swarm evm all test clean
.PHONY: autonity-linux autonity-linux-386 autonity-linux-amd64 autonity-linux-mips64 autonity-linux-mips64le
.PHONY: autonity-linux-arm autonity-linux-arm-5 autonity-linux-arm-6 autonity-linux-arm-7 autonity-linux-arm64
.PHONY: autonity-darwin autonity-darwin-386 autonity-darwin-amd64
.PHONY: autonity-windows autonity-windows-386 autonity-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

autonity:
	build/env.sh go run build/ci.go install ./cmd/autonity
	@echo "Done building."
	@echo "Run \"$(GOBIN)/autonity\" to launch autonity."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/autonity.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/autonity.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

test-race: all
	build/env.sh go run build/ci.go test -race

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

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

swarm-devtools:
	env GOBIN= go install ./cmd/swarm/mimegen

# Cross Compilation Targets (xgo)

autonity-cross: autonity-linux autonity-darwin autonity-windows autonity-android autonity-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/autonity-*

autonity-linux: autonity-linux-386 autonity-linux-amd64 autonity-linux-arm autonity-linux-mips64 autonity-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-*

autonity-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/autonity
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep 386

autonity-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/autonity
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep amd64

autonity-linux-arm: autonity-linux-arm-5 autonity-linux-arm-6 autonity-linux-arm-7 autonity-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep arm

autonity-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/autonity
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep arm-5

autonity-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/autonity
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep arm-6

autonity-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/autonity
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep arm-7

autonity-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/autonity
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep arm64

autonity-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/autonity
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep mips

autonity-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/autonity
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep mipsle

autonity-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/autonity
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep mips64

autonity-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/autonity
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/autonity-linux-* | grep mips64le

autonity-darwin: autonity-darwin-386 autonity-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/autonity-darwin-*

autonity-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/autonity
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-darwin-* | grep 386

autonity-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/autonity
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-darwin-* | grep amd64

autonity-windows: autonity-windows-386 autonity-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/autonity-windows-*

autonity-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/autonity
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-windows-* | grep 386

autonity-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/autonity
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/autonity-windows-* | grep amd64
