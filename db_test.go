package pouchdb

import (
	"context"
	"testing"

	"github.com/gopherjs/gopherjs/js"

	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kiviktest/kt"
)

func init() {
	memPouch := js.Global.Get("PouchDB").Call("defaults", map[string]interface{}{
		"db": js.Global.Call("require", "memdown"),
	})
	js.Global.Set("PouchDB", memPouch)
}

func TestPut(t *testing.T) {
	client, err := kivik.New(context.Background(), "pouch", "")
	if err != nil {
		t.Errorf("Failed to connect to PouchDB/memdown driver: %s", err)
		return
	}
	dbname := kt.TestDBName(t)
	defer client.DestroyDB(context.Background(), dbname)
	db, err := client.CreateDB(context.Background(), dbname)
	if err != nil {
		t.Fatalf("Failed to create db: %s", err)
	}
	_, err = db.Put(context.Background(), "foo", map[string]string{"_id": "bar"})
	if kivik.StatusCode(err) != kivik.StatusBadRequest {
		t.Errorf("Expected Bad Request for mismatched IDs, got %s", err)
	}
}
