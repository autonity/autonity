# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: autonity embed-autonity-contract android ios autonity-cross evm all test clean lint lint-deps mock-gen test-fast
.PHONY: autonity-linux autonity-linux-386 autonity-linux-amd64 autonity-linux-mips64 autonity-linux-mips64le
.PHONY: autonity-linux-arm autonity-linux-arm-5 autonity-linux-arm-6 autonity-linux-arm-7 autonity-linux-arm64
.PHONY: autonity-darwin autonity-darwin-386 autonity-darwin-amd64
.PHONY: autonity-windows autonity-windows-386 autonity-windows-amd64

NPMBIN= $(shell npm bin)
GOBIN = ./build/bin
GO ?= latest
LATEST_COMMIT ?= $(shell git log -n 1 develop --pretty=format:"%H")
ifeq ($(LATEST_COMMIT),)
LATEST_COMMIT := $(shell git log -n 1 HEAD~1 --pretty=format:"%H")
endif

AUTONITY_CONTRACT_DIR = ./contracts/autonity/contract/contracts
AUTONITY_CONTRACT = Autonity.sol
GENERATED_CONTRACT_DIR = ./common/acdefault/generated
GENERATED_ABI = $(GENERATED_CONTRACT_DIR)/abi.go
GENERATED_BYTECODE = $(GENERATED_CONTRACT_DIR)/bytecode.go

# DOCKER_SUDO is set to either the empty string or "sudo" and is used to
# control whether docker is executed with sudo or not. If the user is root or
# the user is in the docker group then this will be set to the empty string,
# otherwise it will be set to "sudo".
#
# We make use of posix's short-circuit evaluation of "or" expressions
# (https://pubs.opengroup.org/onlinepubs/009695399/utilities/xcu_chap02.html#tag_02_09_03)
# where the second part of the expression is only executed if the first part
# evaluates to false. We first check to see if the user id is 0 since this
# indicates that the user is root. If not we then use 'id -nG $USER' to list
# all groups that the user is part of and then grep for the word docker in the
# output, if grep matches the word it returns the successful error code. If not
# we then echo "sudo".
DOCKER_SUDO = $(shell [ `id -u` -eq 0 ] || id -nG $(USER) | grep "\<docker\>" > /dev/null || echo sudo )

# Builds the docker image and checks that we can run the autonity binary inside it
build-docker-image:
	@$(DOCKER_SUDO) docker build -t autonity .
	@$(DOCKER_SUDO) docker run --rm autonity -h > /dev/null

autonity: embed-autonity-contract
	build/env.sh go run build/ci.go install ./cmd/autonity
	@echo "Done building."
	@echo "Run \"$(GOBIN)/autonity\" to launch autonity."

embed-autonity-contract: $(GENERATED_BYTECODE) $(GENERATED_ABI)

# NOTE previously we were using
# https://github.com/ethereum/solidity/releases/download/v0.5.1/solc-static-linux
# this binary does not produce the same bytecode as the ethereum/solc:0.5.1
# docker image. This was causing tests to fail.
$(GENERATED_BYTECODE) $(GENERATED_ABI): $(AUTONITY_CONTRACT_DIR)/$(AUTONITY_CONTRACT)
	mkdir -p $(GENERATED_CONTRACT_DIR)
	$(DOCKER_SUDO) docker run --rm -v $(CURDIR)/$(AUTONITY_CONTRACT_DIR):/contracts -v $(CURDIR)/$(GENERATED_CONTRACT_DIR):/output ethereum/solc:0.6.4 --overwrite --abi --bin -o /output /contracts/$(AUTONITY_CONTRACT)

	@echo Generating $(GENERATED_BYTECODE)
	@echo 'package generated\n' > $(GENERATED_BYTECODE)
	@echo -n 'const Bytecode = "' >> $(GENERATED_BYTECODE)
	@cat  $(GENERATED_CONTRACT_DIR)/Autonity.bin >> $(GENERATED_BYTECODE)
	@echo '"' >> $(GENERATED_BYTECODE)

	@echo Generating $(GENERATED_ABI)
	@echo 'package generated\n' > $(GENERATED_ABI)
	@echo -n 'const Abi = `' >> $(GENERATED_ABI)
	@cat  $(GENERATED_CONTRACT_DIR)/Autonity.abi | json_pp  >> $(GENERATED_ABI)
	@echo '`' >> $(GENERATED_ABI)

all: embed-autonity-contract
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
	build/env.sh go run build/ci.go test -coverage

test-fast:
	build/env.sh go run build/ci.go test

test-race-all: all
	build/env.sh go run build/ci.go test -race
	make test-race

test-race:
	go test -race -v ./consensus/tendermint/... -parallel 1
	go test -race -v ./consensus/test/... -timeout 30m

# This runs the contract tests using truffle against an autonity node instance.
test-contracts:
	@# npm list returns 0 only if the package is not installed and the shell only
	@# executes the second part of an or statment if the first fails.
	@npm list truffle > /dev/null || npm install truffle
	@npm list web3 > /dev/null || npm install web3
	@cd contracts/autonity/contract/test/autonity/ && rm -Rdf ./data && ./autonity-start.sh &
	@# Autonity can take some time to start listening on port 8545 so we allow multiple connection attempts.
	@for x in {1..10}; do \
		sleep 2 ; \
		./build/bin/autonity --exec "web3.personal.unlockAccount(eth.accounts[0], 'test', 36000)" attach http://localhost:8545 ; \
		if [ $$? -eq 0 ] ; then \
			break ; \
		fi ; \
	done
	@cd contracts/autonity/contract/ && $(NPMBIN)/truffle test && cd -

mock-gen:
	mockgen -source=consensus/tendermint/core/core_backend.go -package=core -destination=consensus/tendermint/core/backend_mock.go
	mockgen -source=consensus/protocol.go -package=consensus -destination=consensus/protocol_mock.go
	mockgen -source=consensus/consensus.go -package=consensus -destination=consensus/consensus_mock.go

lint-dead:
	@./build/bin/golangci-lint run \
		--config ./.golangci/step_dead.yml

lint:
	@echo "--> Running linter for code diff versus commit $(LATEST_COMMIT)"
	@./build/bin/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step1.yml \
	    --exclude "which can be annoying to use"

	@./build/bin/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step2.yml

	@./build/bin/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step3.yml

	@./build/bin/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step4.yml

lint-ci: lint-deps lint

test-deps:
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	cd tests/testdata && git checkout b5eb9900ee2147b40d3e681fe86efa4fd693959a

lint-deps:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ./build/bin v1.23.7

clean:
	go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*
	rm -rf $(GENERATED_CONTRACT_DIR)

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	go get -u github.com/golang/mock/mockgen
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

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
