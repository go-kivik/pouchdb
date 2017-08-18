// Package pouchdb provides a driver for the Kivik CouchDB package. It must
// be compiled with GopherJS, and requires that the PouchDB JavaScript library
// is also loaded at runtime.
//
//  // +build js
//
//  package main
//
//  import (
//      "context"
//
//      "github.com/flimzy/kivik"
//      "github.com/go-kivik/pouchdb"
//  )
//
//  func main() {
//      client, err := kivik.New(context.TODO(), "pouch", "")
//  // ...
//  }
//
// See https://github.com/go-kivik/pouchdb#usage for details.
package pouchdb
