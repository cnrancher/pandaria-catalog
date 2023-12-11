#!/bin/bash

set -euo pipefail

cd $(dirname $0)/

RANCHER_VERSION="${RANCHER_VERSION:-}"
BRANCH_NAME=${BRANCH_NAME:-$(git rev-parse --abbrev-ref HEAD)}
if [[ -z "$RANCHER_VERSION" ]]; then
    case $BRANCH_NAME in
        */v2.8)
            RANCHER_VERSION="v2.8"
            ;;
        */v2.7)
            RANCHER_VERSION="v2.7"
            ;;
        */v2.6)
            RANCHER_VERSION="v2.6"
            ;;
        *)
            echo "Could not get Rancher version from git branch [$BRANCH_NAME]"
            echo "Set Rancher version to v2.7"
            RANCHER_VERSION="v2.7"
            ;;
    esac
fi

echo "Rancher version: $RANCHER_VERSION"
echo ""

go mod download && go mod verify
go build .

echo "-------- start check --------"

ERROR=""
./check --version=$RANCHER_VERSION ../ || ERROR="true"
if [[ -e "no-kube-version.txt" ]]; then
    echo ""
    echo "Some charts does not have kube-version annotations:"
    cat no-kube-version.txt
fi

if [[ -e "no-rancher-version.txt" ]]; then
    echo ""
    echo "Some charts does not have rancher-version annotations:"
    cat no-rancher-version.txt
fi

if [[ -e "image-check-failed.txt" ]]; then
    echo ""
    echo "Could not find chart images from following charts values.yaml:"
    cat image-check-failed.txt
fi

if [[ -e "system-default-registry-failed.txt" ]]; then
    echo ""
    echo "Some charts does not defined systemDefaultRegistry:"
    cat system-default-registry-failed.txt
fi

if [[ ! -z "${ERROR}"  ]]; then
    # check failed
    exit 1
fi

echo "-------- check passed --------"
