#!/usr/bin/env bash
set -o pipefail
set -o errexit
set -o nounset

# DEPLOY_SET=(dev-autonity-01.yml dev-autonity-02.yml dev-autonity-03.yml dev-autonity-04.yml dev-autonity-05.yml)

echo $AUTONITY_DEV_CA_CRT | base64 --decode -i > ${HOME}/ca.crt

kubectl config set-cluster dev01.autonity.io --embed-certs=true --server=${AUTONITY_DEV_CLUSTER_ENDPOINT} --certificate-authority=${HOME}/ca.crt
kubectl config set-credentials travis-default --token=$AUTONITY_DEV_USER_TOKEN
kubectl config set-context travis --cluster=dev01.autonity.io --user=travis-default --namespace=default
kubectl config use-context travis
kubectl config current-context

pwd

# kubectl apply -f dev01/

# cd ${HOME}/scripts/

# for i in "${DEPLOY_SET[@]}"
# do
# 	:
# 	kubectl apply -f $i
# done

# # function cleaner {
# #     rm -rvf "${HOME}/*"
# # }

# trap cleaner EXIT