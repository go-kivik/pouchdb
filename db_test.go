package pouchdb

import (
	"context"
	"net/http"
	"testing"

	"github.com/gopherjs/gopherjs/js"
	"gitlab.com/flimzy/testy"

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
	client, err := kivik.New("pouch", "")
	if err != nil {
		t.Errorf("Failed to connect to PouchDB/memdown driver: %s", err)
		return
	}
	dbname := kt.TestDBName(t)
	ctx := context.Background()
	defer client.DestroyDB(ctx, dbname) // nolint: errcheck
	if e := client.CreateDB(ctx, dbname); e != nil {
		t.Fatalf("Failed to create db: %s", e)
	}
	_, err = client.DB(ctx, dbname).Put(ctx, "foo", map[string]string{"_id": "bar"})
	testy.StatusError(t, "id argument must match _id field in document", http.StatusBadRequest, err)
}
