[![Build Status](https://travis-ci.org/go-kivik/pouchdb.svg?branch=master)](https://travis-ci.org/go-kivik/pouchdb)  [![GoDoc](https://godoc.org/github.com/go-kivik/pouchdb?status.svg)](http://godoc.org/github.com/go-kivik/pouchdb)

# Kivik PouchDB

PouchDB driver for [Kivik](https://github.com/go-kivik/pouchdb).

## Usage

This package provides an implementation of the
[`github.com/go-kivik/kivik/driver`](http://godoc.org/github.com/go-kivik/kivik/driver)
interface. You must import the driver and can then use the full
[`Kivik`](http://godoc.org/github.com/go-kivik/kivik) API. Please consult the
[Kivik wiki](https://github.com/go-kivik/kivik/wiki) for complete documentation
and coding examples.

```go
// +build js

package main

import (
    "context"

    "github.com/go-kivik/kivik"
    _ "github.com/go-kivik/pouchdb" // The PouchDB driver
)

func main() {
    client, err := kivik.New(context.TODO(), "pouch", "")
    // ...
}
```

This package is intended to run in a JavaScript runtime, such as a browser or
Node.js, and must be compiled with
[GopherJS](https://github.com/gopherjs/gopherjs). At runtime, the
[PouchDB](https://pouchdb.com/download.html) JavaScript library must also be
loaded and available.

## What license is Kivik released under?

This software is released under the terms of the Apache 2.0 license. See
LICENCE.md, or read the [full license](http://www.apache.org/licenses/LICENSE-2.0).
