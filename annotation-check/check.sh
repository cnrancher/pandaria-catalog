#!/bin/bash

cd $(dirname $0)/

go mod download && go mod verify

ERROR=""
go run main.go ../charts/ || ERROR="true"
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

if [[ ! -z "${ERROR}"  ]]; then
    # check failed
    exit 1
fi
