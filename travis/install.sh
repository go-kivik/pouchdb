#!/bin/bash
set -euC
set -o xtrace

if [ "$TRAVIS_OS_NAME" == "osx" ]; then
    brew install glide
fi

glide update

function generate {
    go get -u github.com/jteeuwen/go-bindata/...
    go generate $(go list ./... | grep -v /vendor/)
}

function wait_for_server {
    printf "Waiting for $1"
    n=0
    until [ $n -gt 5 ]; do
        curl --output /dev/null --silent --head --fail $1 && break
        printf '.'
        n=$[$n+1]
        sleep 1
    done
    printf "ready!\n"
}

function setup_couch16 {
    if [ "$TRAVIS_OS_NAME" == "osx" ]; then
        return
    fi
    docker pull couchdb:1.6.1
    docker run -d -p 6000:5984 -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=abc123 --name couchdb16 couchdb:1.6.1
    wait_for_server http://localhost:6000/
    curl --silent --fail -o /dev/null -X PUT http://admin:abc123@localhost:6000/_config/replicator/connection_timeout -d '"5000"'
}

function setup_couch20 {
    if [ "$TRAVIS_OS_NAME" == "osx" ]; then
        return
    fi
    docker pull klaemo/couchdb:latest
    docker run -d -p 6001:5984 -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=abc123 --name couchdb20 klaemo/couchdb:latest
    wait_for_server http://localhost:6001/
    curl --silent --fail -o /dev/null -X PUT http://admin:abc123@localhost:6001/_users
    curl --silent --fail -o /dev/null -X PUT http://admin:abc123@localhost:6001/_replicator
    curl --silent --fail -o /dev/null -X PUT http://admin:abc123@localhost:6001/_global_changes
}

case "$1" in
    "standard")
        if [ "$TRAVIS_OS_NAME" == "linux" ]; then
            # Install nodejs and dependencies, but only for Linux
            curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
            sudo apt-get update -qq
            sudo apt-get install -y nodejs
        fi
        npm install
        # Then install GopherJS and related dependencies
        go get -u github.com/gopherjs/gopherjs

        # Source maps (mainly to make GopherJS quieter; I don't really care
        # about source maps in Travis)
        npm install source-map-support

        # Set up GopherJS for syscalls
        (
            cd $GOPATH/src/github.com/gopherjs/gopherjs/node-syscall/
            npm install --global node-gyp
            node-gyp rebuild
            mkdir -p ~/.node_libraries/
            cp build/Release/syscall.node ~/.node_libraries/syscall.node
        )

        go get -u -d -tags=js github.com/gopherjs/jsbuiltin
        # setup_couch16
        setup_couch20
        generate
    ;;
    "linter")
        go get -u gopkg.in/alecthomas/gometalinter.v1
        gometalinter.v1 --install
    ;;
esac
