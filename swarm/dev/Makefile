.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/clearmatics/autonity/cmd/swarm
	go build -o bin/autonity     github.com/clearmatics/autonity/cmd/autonity
	go build -o bin/bootnode github.com/clearmatics/autonity/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/clearmatics/autonity/swarm/...
