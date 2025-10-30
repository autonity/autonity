.PHONY: build-docker-image build-docker-image-alltools
.PHONY: all autonity upcheck release android ios
.PHONY: contracts compile-contracts compile-contracts-upgrades compile-contracts-standard 4byte bindings
.PHONY: test test-fast test-race-all test-race
.PHONY: test-contracts test-contracts-fast test-contracts-pre start-autonity start-ganache
.PHONY: test-contracts-truffle test-contracts-truffle-fast docker-e2e-test
.PHONY: mock-gen generate lint-dead lint lint-ci lint-deps clean devtools

# |-----------|
# |	VARIABLES |
# |-----------|

BINDIR = ./build/bin
GO ?= latest
LATEST_COMMIT ?= $(shell git log -n 1 develop --pretty=format:"%H")
# useful for CI
ifeq ($(LATEST_COMMIT),)
	LATEST_COMMIT := $(shell git log -n 1 HEAD~1 --pretty=format:"%H")
endif
SOLC_VERSION = 0.8.30
SOLC_BINARY = $(BINDIR)/solc_static_linux_v$(SOLC_VERSION)
GOBINDATA_VERSION = 3.23.0
GOBINDATA_BINARY = $(BINDIR)/go-bindata
ABIGEN_BINARY = $(BINDIR)/abigen

CONTRACTS_BASE_DIR = ./autonity/solidity
CONTRACTS_DIR = $(CONTRACTS_BASE_DIR)/contracts
CONTRACTS_UPGRADES_DIR = $(CONTRACTS_BASE_DIR)/contracts/upgrades
CONTRACTS_TEST_DIR = $(CONTRACTS_BASE_DIR)/test
GENERATED_CONTRACTS_DIR = ./params/generated

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

# |--------|
# |	DOCKER |
# |--------|

build-docker-image:
	@$(DOCKER_SUDO) docker build -t autonity .

build-docker-image-alltools:
	@$(DOCKER_SUDO) docker build -f Dockerfile.alltools -t autonity-alltools .

# |----------|
# |	BUILDING |
# |----------|

all: $(ABIGEN_BINARY) contracts
	go run build/ci.go install

autonity:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/autonity ./cmd/autonity
	@echo "Done building."
	@echo "Run \"$(BINDIR)/autonity\" to launch autonity."

$(ABIGEN_BINARY):
	go build -o $(ABIGEN_BINARY) ./cmd/abigen

upcheck: $(SOLC_BINARY)
	go build -o $(BINDIR)/upcheck ./cmd/upcheck

release: autonity contracts
	mkdir -p ./build/release/$(VERSION)
	cd ./build/bin && tar -czvf ../release/$(VERSION)/autonity-linux-amd64-$(VERSION).tar.gz autonity
	cd ./params/generated && tar -czvf ../../build/release/$(VERSION)/protocol-contracts-abi-$(VERSION).tar.gz *.abi

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

# |---------------------|
# |	COMPILING CONTRACTS |
# |---------------------|

contracts: compile-contracts 4byte bindings

compile-contracts: compile-contracts-standard compile-contracts-upgrades

compile-contracts-standard: $(SOLC_BINARY) $(GOBINDATA_BINARY) $(CONTRACTS_DIR)/*.sol $(ABIGEN_BINARY)
	@echo "compiling protocol contracts"
	@$(call gen-contract,,Autonity)
	@$(call gen-contract,,Oracle)
	@$(call gen-contract,,Accountability)
	@$(call gen-contract,,OmissionAccountability)
	@$(call gen-contract,,UpgradeManager)
	@$(call gen-contract,,InflationController)
	@$(call gen-contract,asm/,ACU)
	@$(call gen-contract,asm/,SupplyControl)
	@$(call gen-contract,asm/,Stabilization)
	@$(call gen-contract,asm/,Auctioneer)
	@echo "compiling test protocol contracts"
	@$(call gen-contract,test-contract/,AccountabilityTest)
	@$(call gen-contract,test-contract/,AutonityTest)
	@$(call gen-contract,test-contract/,AutonityUpgradeTest)
	@$(call gen-contract,test-contract/,OmissionAccountabilityTest)

compile-contracts-upgrades: $(SOLC_BINARY) $(GOBINDATA_BINARY) $(CONTRACTS_DIR)/*.sol $(ABIGEN_BINARY)
	@echo "compiling protocol contract upgrades"
	@$(call gen-contract,upgrades/0/,Oracle0)
	@$(call gen-contract,upgrades/1/,UpgradeManager1)
	@echo "compiling protocol contract test upgrades"
	@$(call gen-contract,upgrades/tests/,ACUTestUpgrade)
	@$(call gen-contract,upgrades/tests/,AuctioneerTestUpgrade)
	@$(call gen-contract,upgrades/tests/,InflationControllerTestUpgrade)
	@$(call gen-contract,upgrades/tests/,StabilizationTestUpgrade)
	@$(call gen-contract,upgrades/tests/,SupplyControlTestUpgrade)
	@$(call gen-contract,upgrades/tests/,UpgradeManagerTestUpgrade)

4byte:
	@echo "update 4byte selector for clef"
	./build/generate_4bytedb.sh $(SOLC_BINARY)
	cd signer/fourbyte && go generate

bindings: $(ABIGEN_BINARY)
	@$(call gen-bindings,bindings/,bindings)
	@# NOTE: the sed substitution is required to generate the same runtime bytecode that got deployed on mainnet.
	@# it simply substitutes the metadata hash of the oracle upgraded bytecode with the metadata hash used on mainnet.
	@# the hashes are different because the mainnet upgrade had been built in the `autonity/solidity/contracts/` folder.
	@# for more details see https://docs.soliditylang.org/en/latest/metadata.html
	sed -i s/9332f2eb7a6ca90dad8644e835fd01558f41287c53fb4f634a9a437d9ac6c96a/397c7e11019699c95916f85b08a3150696523f8159ab4f3bc820cc82275c34bc/ ./autonity/bindings/bindings.go ./autonity/tests/bindings.go

# params:
# $(1) --> contract folder (can be empty if the contract is in $(CONTRACTS_DIR) )
# $(2) --> contract filename (which should be == with the contract's name)
define gen-contract
	$(SOLC_BINARY) --overwrite --optimize --optimize-runs 10000 --evm-version london --abi --bin --bin-runtime --metadata --userdoc --devdoc -o $(GENERATED_CONTRACTS_DIR) $(CONTRACTS_DIR)/$(1)$(2).sol

	@# NOTE: the sed substitution is required to generate the same runtime bytecode that got deployed on mainnet.
	@# it simply substitutes the metadata hash of the oracle upgraded bytecode with the metadata hash used on mainnet.
	@# the hashes are different because the mainnet upgrade had been built in the `autonity/solidity/contracts/` folder.
	@# for more details see https://docs.soliditylang.org/en/latest/metadata.html
	@if [ "$(2)" = "Oracle0" ]; then \
		sed -i s/9332f2eb7a6ca90dad8644e835fd01558f41287c53fb4f634a9a437d9ac6c96a/397c7e11019699c95916f85b08a3150696523f8159ab4f3bc820cc82275c34bc/ $(GENERATED_CONTRACTS_DIR)/$(2).bin $(GENERATED_CONTRACTS_DIR)/$(2).bin-runtime; \
	fi

	@echo Generating bytecode for $(2)
	@echo 'package generated' > $(GENERATED_CONTRACTS_DIR)/$(2).go
	@echo 'import (' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@echo '	"strings"' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@echo '' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@echo '	"github.com/autonity/autonity/accounts/abi"' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@echo '	"github.com/autonity/autonity/common"' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@echo '	"github.com/autonity/autonity/crypto"' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@echo ')' >> $(GENERATED_CONTRACTS_DIR)/$(2).go

	@echo -n 'var $(2)Bytecode = common.Hex2Bytes("' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@cat $(GENERATED_CONTRACTS_DIR)/$(2).bin >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@printf '")\n\n' >> $(GENERATED_CONTRACTS_DIR)/$(2).go

	@echo -n 'var $(2)RuntimeBytecode = common.Hex2Bytes("' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@cat $(GENERATED_CONTRACTS_DIR)/$(2).bin-runtime >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@printf '")\n\n' >> $(GENERATED_CONTRACTS_DIR)/$(2).go

	@echo "var $(2)CodeHash = crypto.Keccak256Hash($(2)RuntimeBytecode)\n" >> $(GENERATED_CONTRACTS_DIR)/$(2).go

	@echo Generating Abi for $(2)
	@echo -n 'var $(2)Abi,_ = abi.JSON(strings.NewReader(`' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@cat  $(GENERATED_CONTRACTS_DIR)/$(2).abi | json_pp  >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@echo '`))' >> $(GENERATED_CONTRACTS_DIR)/$(2).go
	@gofmt -s -w $(GENERATED_CONTRACTS_DIR)/$(2).go
endef

# params:
# $(1) --> contract folder (can be empty if the contract is in $(CONTRACTS_DIR))
# $(2) --> contract filename
define gen-bindings
	@echo "Generating bindings for $(2).sol"
	$(ABIGEN_BINARY) --pkg bindings --solc $(SOLC_BINARY) \
	--sol $(CONTRACTS_DIR)/$(1)$(2).sol \
	--out ./autonity/bindings/$(2).go

	@echo "Generating testing bindings for $(2).sol"
	$(ABIGEN_BINARY) --test \
	--pkg tests --solc $(SOLC_BINARY) \
	--sol $(CONTRACTS_DIR)/$(1)$(2).sol \
	--out ./autonity/tests/$(2).go
endef

# |---------|
# |	TESTING |
# |---------|

test: all
	go run build/ci.go test -coverage

test-fast:
	go run build/ci.go test

test-race-all: all
	go run build/ci.go test -race
	make test-race

test-race:
	go test -race -v ./consensus/tendermint/... -parallel 1

test-contracts: test-contracts-truffle

test-contracts-fast: test-contracts-truffle-fast

# prerequisites for testing contracts
test-contracts-pre:
	@# npm list returns 0 only if the package is not installed and the shell only
	@# executes the second part of an or statement if the first fails.
	@# Nov, 2022, the latest release of Truffle, v5.6.6 does not works for the tests.
	@echo "check and install truffle.js"
	@npm list truffle > /dev/null || npm install truffle
	@echo "check and install web3.js"
	@npm list web3 > /dev/null || npm install web3
	@echo "check and install keccak256"
	@npm list keccak256 > /dev/null || npm install keccak256
	@echo "check and install truffle-assertions.js"
	@npm list truffle-assertions > /dev/null || npm install truffle-assertions
	@echo "check and install ganache"
	@npm list ganache > /dev/null || npm install ganache
	@npx truffle version

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
	@#validator_management.js is ran only against Autonity, since ganache does not implement the POP  logic
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test validator_management.js && cd -
	@echo "killing test autonity network and cleaning chaindata"
	@-pkill autonity
	@cd $(CONTRACTS_TEST_DIR)/autonity/ && rm -Rdf ./data

# This runs the contract tests using truffle against a Ganache instance
test-contracts-truffle-fast: contracts test-contracts-pre start-ganache
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test autonity.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test oracle.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test accountability.js && cd -
	@cd $(CONTRACTS_TEST_DIR) && npx truffle test protocol.js && cd -
	@echo "killing ganache"
	@-pkill -f "ganache"

docker-e2e-test: contracts
	build/env.sh go run build/ci.go install
	cd docker_e2e_test && sudo python3 test_via_docker.py ..

# |-----------|
# |	UTILITIES |
# |-----------|

mock-gen:
	mockgen -source=consensus/tendermint/core/interfaces/core_backend.go -package=interfaces -destination=consensus/tendermint/core/interfaces/core_backend_mock.go
	mockgen -source=consensus/tendermint/accountability/fault_detector.go -package=accountability -destination=consensus/tendermint/accountability/fault_detector_mock.go
	mockgen -source=consensus/consensus.go -package=consensus -destination=consensus/consensus_mock.go
	mockgen -source=accounts/abi/bind/backend.go -package=bind -destination=accounts/abi/bind/backend_mock.go
	mockgen -source=consensus/tendermint/core/interfaces/gossiper.go -package=interfaces -destination=consensus/tendermint/core/interfaces/gossiper_mock.go
	mockgen -source=consensus/tendermint/core/interfaces/broadcaster.go -package=interfaces -destination=consensus/tendermint/core/interfaces/broadcaster_mock.go
	mockgen -source=consensus/tendermint/router/interfaces/interfaces.go -package=mocks -destination=consensus/tendermint/router/mocks/interfaces_mock.go
	mockgen -source=consensus/tendermint/router/ping/pinger.go -package=mocks -destination=consensus/tendermint/router/mocks/pinger_mock.go
	mockgen -source=consensus/tendermint/router/cache/cache.go -package=mocks -destination=consensus/tendermint/router/mocks/cache_mock.go
	mockgen -source=event/subscription.go -package=mocks -destination=consensus/tendermint/router/mocks/subscription_mock.go

generate:
	cd core/types/ && go generate

lint-dead:
	@./.github/tools/golangci-lint run \
		--config ./.golangci/step_dead.yml

lint:
	@echo "--> Running linter for code diff versus commit $(LATEST_COMMIT)"
	@./.github/tools/golangci-lint run \
	    --new-from-rev=$(LATEST_COMMIT) \
	    --config ./.golangci/step1.yml \

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

lint-deps:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.github/tools v1.64.2

clean:
	go clean -cache
	rm -f $(BINDIR)/*
	rm -f $(GENERATED_CONTRACTS_DIR)/*.abi $(GENERATED_CONTRACTS_DIR)/*.bin* $(GENERATED_CONTRACTS_DIR)/*.doc*

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

$(SOLC_BINARY):
	mkdir -p $(BINDIR)
	wget -O $(SOLC_BINARY) https://github.com/ethereum/solidity/releases/download/v$(SOLC_VERSION)/solc-static-linux
	chmod +x $(SOLC_BINARY)

$(GOBINDATA_BINARY):
	mkdir -p $(BINDIR)
	wget -O $(GOBINDATA_BINARY) https://github.com/kevinburke/go-bindata/releases/download/v$(GOBINDATA_VERSION)/go-bindata-linux-amd64
	chmod +x $(GOBINDATA_BINARY)

