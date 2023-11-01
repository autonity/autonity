# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: autonity contracts android ios autonity-cross evm all test clean lint mock-gen test-fast test-contracts test-contracts-truffle-fast test-contracts-truffle start-autonity start-ganache test-contracts-pre test-contracts-fast

BINDIR = ./build/bin
GO ?= latest
LATEST_COMMIT ?= $(shell git log -n 1 develop --pretty=format:"%H")
ifeq ($(LATEST_COMMIT),)
LATEST_COMMIT := $(shell git log -n 1 HEAD~1 --pretty=format:"%H")
endif
SOLC_VERSION = 0.8.21
SOLC_BINARY = $(BINDIR)/solc_static_linux_v$(SOLC_VERSION)
GOBINDATA_VERSION = 3.23.0
GOBINDATA_BINARY = $(BINDIR)/go-bindata
ABIGEN_BINARY = $(BINDIR)/abigen

CONTRACTS_BASE_DIR = ./autonity/solidity
CONTRACTS_DIR = $(CONTRACTS_BASE_DIR)/contracts
CONTRACTS_TEST_DIR = $(CONTRACTS_BASE_DIR)/test
GENERATED_CONTRACT_DIR = ./params/generated

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

# Builds the docker image and checks that we can run the autonity binary inside
# it.
build-docker-image:
	@$(DOCKER_SUDO) docker build -t autonity .
	@$(DOCKER_SUDO) docker run --rm autonity -h > /dev/null

autonity:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/autonity ./cmd/autonity
	@echo "Done building."
	@echo "Run \"$(BINDIR)/autonity\" to launch autonity."

bindings:
	@echo Generating protocol contracts bindings
	$(ABIGEN_BINARY)  --pkg autonity --solc $(SOLC_BINARY) --sol "$(CONTRACTS_DIR)/bindings.sol" --out ./autonity/bindings.go

# Builds Autonity without contract compilation, useful with alpine containers not supporting
# glibc for solc.
autonity-docker:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/autonity ./cmd/autonity
	@echo "Done building."
	@echo "Run \"$(BINDIR)/autonity\" to launch autonity."

define gen-contract
	$(SOLC_BINARY) --overwrite --optimize --optimize-runs 10000 --evm-version london --abi --bin -o $(GENERATED_CONTRACT_DIR) $(CONTRACTS_DIR)/$(1)$(2).sol

	@echo Generating bytecode for $(2)
	@echo 'package generated' > $(GENERATED_CONTRACT_DIR)/$(2).go
	@echo 'import "strings"' >> $(GENERATED_CONTRACT_DIR)/$(2).go
	@echo 'import "github.com/autonity/autonity/accounts/abi"' >> $(GENERATED_CONTRACT_DIR)/$(2).go
	@echo 'import "github.com/autonity/autonity/common"' >> $(GENERATED_CONTRACT_DIR)/$(2).go

	@echo -n 'var $(2)Bytecode = common.Hex2Bytes("' >> $(GENERATED_CONTRACT_DIR)/$(2).go
	@cat  $(GENERATED_CONTRACT_DIR)/$(2).bin >> $(GENERATED_CONTRACT_DIR)/$(2).go
	@echo '")\n' >> $(GENERATED_CONTRACT_DIR)/$(2).go

	@echo Generating Abi for $(2)
	@echo -n 'var $(2)Abi,_ = abi.JSON(strings.NewReader(`' >> $(GENERATED_CONTRACT_DIR)/$(2).go
	@cat  $(GENERATED_CONTRACT_DIR)/$(2).abi | json_pp  >> $(GENERATED_CONTRACT_DIR)/$(2).go
	@echo '`))' >> $(GENERATED_CONTRACT_DIR)/$(2).go
	@gofmt -s -w $(GENERATED_CONTRACT_DIR)/$(2).go

endef

contracts: $(SOLC_BINARY) $(GOBINDATA_BINARY) $(CONTRACTS_DIR)/*.sol $(ABIGEN_BINARY)
	@$(call gen-contract,,Autonity)
	@$(call gen-contract,,Oracle)
	@$(call gen-contract,,AutonityUpgradeTest)
	@$(call gen-contract,,Accountability)
	@$(call gen-contract,asm/,ACU)
	@$(call gen-contract,asm/,SupplyControl)
	@$(call gen-contract,asm/,Stabilization)
	# update 4byte selector for clef
	build/generate_4bytedb.sh $(SOLC_BINARY)
	cd signer/fourbyte && go generate
	# Generate go bindings
	@echo Generating protocol contracts bindings
	$(ABIGEN_BINARY)  --pkg autonity --solc $(SOLC_BINARY) --sol $(CONTRACTS_DIR)/bindings.sol --out ./autonity/bindings.go



$(SOLC_BINARY):
	mkdir -p $(BINDIR)
	wget -O $(SOLC_BINARY) https://github.com/ethereum/solidity/releases/download/v$(SOLC_VERSION)/solc-static-linux
	chmod +x $(SOLC_BINARY)

$(GOBINDATA_BINARY):
	mkdir -p $(BINDIR)
	wget -O $(GOBINDATA_BINARY) https://github.com/kevinburke/go-bindata/releases/download/v$(GOBINDATA_VERSION)/go-bindata-linux-amd64
	chmod +x $(GOBINDATA_BINARY)

$(ABIGEN_BINARY):
	go build -o $(ABIGEN_BINARY) ./cmd/abigen

all: contracts
	make bindings
	go run build/ci.go install

android:
	go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/geth.aar\" to use the library."
	@echo "Import \"$(GOBIN)/geth-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"

ios:
	go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(BINDIR)/autonity.framework\" to use the library."

test: all
	go run build/ci.go test -coverage

test-fast:
	go run build/ci.go test

test-race-all: all
	go run build/ci.go test -race
	make test-race

test-race:
	go test -race -v ./consensus/tendermint/... -parallel 1
	go test -race -v ./consensus/test/... -timeout 30m

test-contracts: test-contracts-asm test-contracts-truffle

test-contracts-fast: test-contracts-asm test-contracts-truffle-fast

# prerequisites for testing contracts
test-contracts-pre:
	@# npm list returns 0 only if the package is not installed and the shell only
	@# executes the second part of an or statement if the first fails.
	@# Nov, 2022, the latest release of Truffle, v5.6.6 does not works for the tests.
	@echo "check and install truffle.js"
	@npm list truffle > /dev/null || npm install truffle
	@echo "check and install web3.js"
	@npm list web3 > /dev/null || npm install web3
	@echo "check and install truffle-assertions.js"
	@npm list truffle-assertions > /dev/null || npm install truffle-assertions
	@echo "check and install ganache"
	@npm list ganache > /dev/null || npm install ganache
	@npx truffle version

test-contracts-asm:
	@echo "check and install ape framework"
	@ape > /dev/null || pipx install eth-ape || { pipx uninstall eth-ape; exit 1; }
	@echo "check and install hardhat"
	@cd $(CONTRACTS_BASE_DIR) && npm list hardhat > /dev/null || npm install hardhat
	@echo "install ape framework plugins"
	@cd $(CONTRACTS_BASE_DIR) && ape plugins install -y --verbosity ERROR .
	@echo "run tests for the asm contracts"
	@cd $(CONTRACTS_BASE_DIR) && ape --verbosity WARNING test --network ::hardhat ./test/asm

# start an autonity network for contract tests
start-autonity:
	@echo "starting autonity test network"
	@cd $(CONTRACTS_TEST_DIR)/autonity/ && nohup ./autonity-start.sh >/dev/null 2>&1 &
	@# give some time for autonity to start
	@sleep 10
	@# check that autonity started correctly and is listening on the correct port
	@pgrep autonity
	@lsof -i :8545 | grep autonity

# start a ganache network for fast contract tests
start-ganache:
	@echo "starting ganache"
	@nohup npx ganache --chain.allowUnlimitedContractSize --chain.allowUnlimitedInitCodeSize --gasLimit 0x1fffffffffffff >/dev/null 2>&1 &
	@sleep 2
	@pgrep -f ganache
	@lsof -i :8545 | grep node

# This runs the contract tests using truffle against an Autonity node instance.
test-contracts-truffle: autonity contracts test-contracts-pre start-autonity
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test autonity.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test oracle.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test liquid.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test accountability.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test protocol.js && cd -
	@#refund.js is ran only against Autonity, since ganache does not implement the oracle vote refund logic
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test refund.js && cd -
	@echo "killing test autonity network and cleaning chaindata"
	@-pkill autonity
	@cd $(CONTRACTS_TEST_DIR)/autonity/ && rm -Rdf ./data

# This runs the contract tests using truffle against a Ganache instance
test-contracts-truffle-fast: contracts test-contracts-pre start-ganache
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test autonity.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test oracle.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test liquid.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test accountability.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test protocol.js && cd -
	@echo "killing ganache"
	@-pkill -f "ganache"

docker-e2e-test: contracts
	build/env.sh go run build/ci.go install
	cd docker_e2e_test && sudo python3 test_via_docker.py ..

mock-gen:
	mockgen -source=consensus/tendermint/core/interfaces/core_backend.go -package=interfaces -destination=consensus/tendermint/core/interfaces/core_backend_mock.go
	mockgen -source=consensus/tendermint/accountability/fault_detector.go -package=accountability -destination=consensus/tendermint/accountability/fault_detector_mock.go
	mockgen -source=consensus/protocol.go -package=consensus -destination=consensus/protocol_mock.go
	mockgen -source=interfaces.go -package=ethereum -destination=interfaces_mock.go
	mockgen -source=consensus/consensus.go -package=consensus -destination=consensus/consensus_mock.go
	mockgen -source=consensus/tendermint/core/interfaces/tendermint.go -package=interfaces -destination=consensus/tendermint/core/interfaces/tendermint_mock.go
	mockgen -source=accounts/abi/bind/backend.go -package=bind -destination=accounts/abi/bind/backend_mock.go


lint-dead:
	@./.github/tools/golangci-lint run \
		--config ./.golangci/step_dead.yml

lint:
	@echo "--> Running linter for code diff versus commit $(LATEST_COMMIT)"
	@./.github/tools/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step1.yml \
	    --exclude "which can be annoying to use"

	@./.github/tools/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step2.yml

	@./.github/tools/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step3.yml

	@./.github/tools/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step4.yml

lint-ci: lint

test-deps:
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	cd tests/testdata && git checkout b5eb9900ee2147b40d3e681fe86efa4fd693959a

lint-deps:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ./build/bin v1.53.3

clean:
	go clean -cache
	rm -fr build/_workspace/pkg/ $(BINDIR)/*
	rm -rf $(GENERATED_CONTRACT_DIR)/*.abi $(GENERATED_CONTRACT_DIR)/*.bin

# The devtools target installs tools required for 'go generate'.
# You need to put $BINDIR (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	go get -u go.uber.org/mock/mockgen
	env BINDIR= go get -u golang.org/x/tools/cmd/stringer
	env BINDIR= go install github.com/kevinburke/go-bindata/v4/...@latest
	env BINDIR= go install github.com/fjl/gencodec@latest
	env BINDIR= go get -u github.com/golang/protobuf/protoc-gen-go
	go build -o $BINDIR/abigen ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
