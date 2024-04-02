# NetDiag

Network Diagnosis tooling suite for Autonity Networks. 

This toolkit is designed to evaluate and diagnose network reliability within a real-world context, aiming in the
long-run to serve as a comprehensive testing framework for sophisticated network propagation methodologies.
At present, it operates by establishing a cluster of DevP2P nodes that execute a bespoke protocol
developed specifically for network diagnostics. Those nodes are remotely controlled by a supervisory program that 
facilitates the dispatch of on-demand requests through RPC.

## Usage
### Requirements:

- You must have gloud installed along with proper adc credentials setups. 
See https://cloud.google.com/docs/authentication/external/set-up-adc for more informations.

### Deploy New GCP cluster

`./netdiag setup --gcp-project`

### Run controller

`./netdiag control --gcp-project`

This will start an interactive session where you can execute any remote control command.

### Windown deployed cluster

This command will clear

`./netdiag clean --gcp-project`