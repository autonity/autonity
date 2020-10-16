# autontiy-e2e-test-engine-chart
It setup a mini autonity network with 6 validators and a test engine which applies testcases over the network.

# usage
```
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm repo add myhelmrepo https://Jason-Zhangxin-Chen.github.io/engine-chart/
helm repo update

kubectl create namespace autonity-e2e-test
helm install myhelmrepo/autonity-e2e-test-chart --namespace autonity-e2e-test
```
