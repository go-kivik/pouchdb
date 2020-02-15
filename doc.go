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
//      kivik "github.com/go-kivik/kivik/v3"
//      _ "github.com/go-kivik/pouchdb/v3" // PouchDB driver
//  )
//
//  func main() {
//      client, err := kivik.New(context.TODO(), "pouch", "")
//  // ...
//  }
//
// See https://github.com/go-kivik/pouchdb#usage for details.
package pouchdb
