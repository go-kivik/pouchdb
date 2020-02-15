#!/bin/bash
set -euC

function join_list {
    local IFS=","
    echo "$*"
}

case "$1" in
    "standard")
        gopherjs test $(go list ./... | grep -v /vendor/ | grep -Ev 'kivik/(serve|auth|proxy)')
    ;;
    "linter")
        go install # to make gotype (run by gometalinter) happy
        golangci-lint run ./...
    ;;
esac
