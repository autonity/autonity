## Autonity-E2E-Test-Framework helm chart

[![Join the chat at https://gitter.im/clearmatics/autonity](https://badges.gitter.im/clearmatics/autonity.svg)](https://gitter.im/clearmatics/autonity?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Introduction
This chart deploys a **private** [Autonity](https://www.autonity.io/) network onto a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

"Ethereum based protocol enabling permissioned, decentralized and interoperable transacting member-mutual networks." - [Autonity](https://www.autonity.io)

## Examples
### Prerequisites
1. Install [Kubernetes](http://kubernetes.io) `v1.18.4`. Basic features have been tested in both Minikube and Kind with extra functionality available in EKS and GKE:
   - Local: [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
   - Local: [Kind](https://kind.sigs.k8s.io/docs/user/quick-start)
   - Cloud: [Amazon EKS](https://eksworkshop.com/prerequisites/self_paced/)
   - Cloud: [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine/docs/quickstart)
1. Install Helm [3.2.4](https://github.com/helm/helm/releases/tag/v3.2.4)
1. Initialise [the official](https://helm.sh/docs/intro/quickstart/ the official) `@stable` and the Autonity Helm charts repositories:
```bash
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm repo add charts-ose.clearmatics.com https://charts-ose.clearmatics.com/
helm repo update
```

### tl;dr
```bash
kubectl create namespace autonity-network
helm install autonity-network charts-ose.clearmatics.com/autonity-network --namespace autonity-network --version 1.8.0
```
Note: `autonity-network` versions before `1.8.0` are supported by [Helm 2](https://github.com/clearmatics/charts-ose/).

## Kubernetes objects
This chart is comprised of 4 components:
1. The 'initial jobs' which implement the bootstrapping algorithm and run once.
   1. [init-job01-ibft-keys-generator](https://github.com/clearmatics/ibft-keys-generator) will create keys set. Note: runs as a helm hook `post-install`
   1. [init-job02-ibft-genesis-configurator](https://github.com/clearmatics/ibft-genesis-configurator) will configure `genesis.json` for autonity network initialisation
1. The [autonity init](https://github.com/clearmatics/autonity-init) container for each autonity pod that downloads keys, configs and prepare chain data
1. The pods `validator-X`, are Autonity nodes that implement the [IBFT consensus algorithms](https://docs.autonity.io/IBFT/index.html)
1. The pods `observer-Y`, are Autonity nodes that connected to validators by p2p and expose JSON-RPC and WebSocket interface

## Data storage
1. The secret `account-pwd` contains the generated account password.
1. Secrets for `validators`, `observers`, `operator-governance` or `operator-treasury` contain:
   1. `0.private_key` - private key for account
1. Configmaps for `validators`, `observers`, `operator-governance` or `operator-treasury` contain:
   1. `0.address` - address
   1. `0.pub_key` - public key
1. Kubernetes [EmptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) (default) for local blockchain of `validators` and `observers`
   1. `aws_persistent_storage_enabled: true` enable AWS persistent storage for `blockchain`
   1. `gcp_persistent_storage_enabled: true` enable GCP persistent storage for `blockchain`

## Configuration
The following table lists some of the configurable parameters of the Autonity chart and their default values. This table needs extending fully.

| Parameter           | Description                                                                    | Default                                |
|---------------------|--------------------------------------------------------------------------------|----------------------------------------|
| `debug_enabled`     | Prepends log messages with call-site location                                  | `false`                                |
| `graphql_enabled`   | Enabling GraphQL query capabilities on top of HTTP RPC                         | `false`                                |
| `http_rpc.enabled   | Enabling the HTTP-RPC server                                                   | `true`                                 |
| `http_rpc.address`  | Setting the HTTP-RPC server listening interface                                | `127.0.0.1`                            |
| `http_rpc.port`     | Setting the HTTP-RPC server listening port                                     | `8545`                                 |
| `http_rpc.api`      | A list of APIs offered over the HTTP-RPC interface                             | `eth,web3,net,tendermint,txpool,debug` |
| `http_rpc.vhosts`   | A comma separated list of virtual hostnames from which to accept requests from | `\*`                                   |
| `logging_verbosity` | Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail        | `3`                                    |
| `pprof.enabled`     | HTTP server for visualization and analysis of profiling data                   | `false`                                |
| `pprof.address`     | The address the pprof server will listen on                                    | `127.0.0.1`                            |
| `pprof.port`        | The port the pprof server will start on                                        | `6060`                                 |
| `ws_rpc.enabled`    | Enabling the WS-RPC server                                                     | `true`                                 |
| `ws_rpc.address`    | Setting the WS-RPC server listening interface                                  | `127.0.0.1`                            |
| `ws_rpc.port`       | Setting the WS-RPC server listening port                                       | `8546`                                 |
| `ws_rpc.api`        | A list of APIs offered over the HTTP-RPC interface                             | `eth,web3,net,tendermint`              |
| `ws_rpc.origins`    | A comma separated list of origins from which to accept websockets requests     | `\*`                                   |

- You can change number of validators or observers using the `--set` options:

```bash
helm install autonity-network charts-ose.clearmatics.com/autonity-network \
  --namespace autonity-network
  --version 1.8.0 \
  --set validators=6,observers=2
```

- Available variables are in the [./values.yaml](values.yaml) file
- Configuration of the `autonity-network` options are available in this template [./templates/configmap_genesis_template.yaml](templates/configmap_genesis_template.yaml)
- All of the other options in the `genesis.json` file like: `validators`, `alloc`, `nodeWhiteList` will be generated automaticaly based on validators and observers list.

To get the generated genesis out of the configmap:
```bash
kubectl get configmap genesis -o yaml --export=true
```

### JSON-RPC HTTP Basic Auth
* To enable http basic auth, set:
```bash
rpc_http_basic_auth_enabled: true
```

* Generate the `htpasswd` file to ./file/htpasswd using [htpasswd](https://httpd.apache.org/docs/2.4/programs/htpasswd.html):
```bash
mkdir ./files
htpasswd -c ./files/htpasswd user1
htpasswd ./files/htpasswd user2
htpasswd ./files/htpasswd user3
```

```bash
helm install --namespace autonity-network ./ --set rpc_http_basic_auth_enabled=true
```

* Add content to values.yaml with content of the `htpasswd` file, for example in values.yaml:
```bash
htpasswd: |-
  deployer:...
```

### HTTPS for JSON-RPC
* Generate keys and certificates:
```bash
KEY_FILE=files/tls.key
CERT_FILE=files/tls.crt
DH_FILE=files/dhparam.pem

mkdir -p files
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE}

openssl dhparam -out ${DH_FILE} 4096
```

* Set variable to `rpc_https_enabled: true`
```bash
helm install --namespace autonity-network ./ --set rpc_https_enabled=true
```

* Add the `base64` encoded certificate variables to the `.yaml` file - for example - with the output of `base64 $KEY_FILE`, `base64 $CERT_FILE`, `base64 $DH_FILE`:
```bash
  rpc:
    tls_crt: ...
    tls_key: ...
    dhparam_pem: ...
```

### Node selectors
It is often desirable to bind specific Autonity services to specific Kubernetes nodes. In order to do this we have 2 possible ways.

If you have statically labled nodes (format should be validator-$i) set node_selector to 'app' and pods will be assigned to the correct physical nodes.

In certain scenarios this isn't possible (for example deploying via gce node pool or AWS ASG) in this case a strategy for node/pool assignment is as follows:

|                      | *europe-west-1a* | *europe-west-1b* | *europe-west-1c* |
|----------------------|------------------|------------------|------------------|
| *pool: validators-0* | validator-0      | validator-1      | validator-2      |
| *pool: validators-1* | validator-3      | validator-4      | validator-5      |

We must ensure that the following conditions are met:
- specify selector_method: zones
- ensure that your number of validators is divisible by your number of zones
- label nodes pool: validators-n with n being the number of pools you will need to deploy your number of validators.
- supply list of the availabiltiy zones you are using as array zones: in values.yaml (this is the value of failure-domain.beta.kubernetes.io/zone).

### Metrics by InfluxDB
You can send `autonity` metrics to InfluxDB cloud. (Disabled by default).
* Create InfluxDB2 cloud account here: [cloud2.influxdata.com](https://cloud2.influxdata.com/) It is for free for limited testing usage.
* Get your organisation name from Org Profile
* Create bucket
* Create token

Put it to [values.yaml](./values.yaml) as a `telegraf:` values

## Notes
To view available commands:
```bash
helm status autonity-network --namespace autonity-network
```

## Cleanup
```bash
helm delete autonity-network --namespace autonity-network
kubectl delete namespace autonity-network
```
