package pouchdb

import (
	"context"
	"testing"

	"github.com/flimzy/kivik"
	"github.com/flimzy/kivik/errors"
	"github.com/flimzy/kivik/test/kt"
	"github.com/gopherjs/gopherjs/js"
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
	if err = client.CreateDB(context.Background(), dbname); err != nil {
		t.Fatalf("Failed to create db: %s", err)
	}
	db, err := client.DB(context.Background(), dbname)
	if err != nil {
		t.Fatalf("Failed to connect to db: %s", err)
	}
	_, err = db.Put(context.Background(), "foo", map[string]string{"_id": "bar"})
	if errors.StatusCode(err) != kivik.StatusBadRequest {
		t.Errorf("Expected Bad Request for mismatched IDs, got %s", err)
	}
}
