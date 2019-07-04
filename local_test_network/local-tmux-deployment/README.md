Please run `make-autonity`. It will launch 4 local nodes with Tendermint consensus algorithm in a tmux session called `autonity`, which wil contain 4 windows for each of the node.

# Use in Goland
- Add `BashSupport` plugin
- Add `Python Community Edition` plugin
- Configure Virtualenv Environment for python3.6 in `Build -> Execution -> Deployment -> Interpreter`
- Edit Run/Debug Configuration to create a bash configuration, where set Interpreter path to `/usr/bin/env`, Interpreter options to `bash` and Working directory to `/home/as/go/src/github.com/clearmatics/autonity/local_test_network/local-tmux-deployment`
- Run the configuration.