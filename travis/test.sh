#!/bin/bash
set -euC

if [ "${TRAVIS_OS_NAME:-}" == "osx" ]; then
    # We don't have docker in OSX, so skip these tests
    unset KIVIK_TEST_DSN_COUCH16
    unset KIVIK_TEST_DSN_COUCH20
fi

function join_list {
    local IFS=","
    echo "$*"
}

case "$1" in
    "standard")
        unset KIVIK_TEST_DSN_COUCH16
        gopherjs test $(go list ./... | grep -v /vendor/ | grep -Ev 'kivik/(serve|auth|proxy)')
    ;;
    "linter")
        diff -u <(echo -n) <(gofmt -e -d $(find . -type f -name '*.go' -not -path "./vendor/*"))
        go install # to make gotype (run by gometalinter) happy
        gometalinter.v1 --config .linter_test.json
        gometalinter.v1 --config .linter.json
    ;;
esac
