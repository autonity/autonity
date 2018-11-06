#!/usr/bin/env bash
set -o pipefail
set -o errexit
set -o nounset

echo "executed deploy script"

# TIMER_VAR=1
# DEPLOY_SET=(dev-autonity-01.yml dev-autonity-02.yml dev-autonity-03.yml dev-autonity-04.yml dev-autonity-05.yml)

# echo $AUTONITY_DEV_CA_CRT | base64 --decode -i > ${HOME}/ca.crt

# kubectl config set-cluster our-k8s-cluster --embed-certs=true --server=${AUTONITY_DEV_CLUSTER_ENDPOINT} --certificate-authority=${HOME}/ca.crt
# kubectl config set-credentials travis-default --token=$SA_TOKEN
# kubectl config set-context travis --cluster=$CLUSTER_NAME --user=travis-default --namespace=default
# kubectl config use-context travis
# kubectl config current-context

# for i in "${DEPLOY_SET[@]}"
# do
# 	:
# 	kubectl apply -f $i
# done

# # function cleaner {
# #     rm -rvf "${HOME}/*"
# # }

# trap cleaner EXIT