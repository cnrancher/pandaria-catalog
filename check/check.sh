#!/bin/bash

cd $(dirname $0)/

go mod download && go mod verify

echo "-------- start check --------"

ERROR=""
go run main.go ../ || ERROR="true"
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
