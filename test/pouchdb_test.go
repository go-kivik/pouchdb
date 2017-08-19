// +build js

package test

import (
	"testing"

	_ "github.com/go-kivik/pouchdb"
)

func init() {
	RegisterPouchDBSuites()
}

func TestPouchLocal(t *testing.T) {
	PouchLocalTest(t)
}

func TestPouchRemote(t *testing.T) {
	PouchRemoteTest(t)
}
