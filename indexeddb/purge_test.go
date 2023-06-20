package indexeddb

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gopherjs/gopherjs/js"
	"gitlab.com/flimzy/testy"

	kivik "github.com/go-kivik/kivik/v4"
	"github.com/go-kivik/kiviktest/v4/kt"
	_ "github.com/go-kivik/pouchdb/v4" // PouchDB driver we're testing
)

func init() {
	js.Global.Call("require", "fake-indexeddb/auto")
	indexedDBPlugin := js.Global.Call("require", "pouchdb-adapter-indexeddb")
	pouchDB := js.Global.Get("PouchDB")
	pouchDB.Call("plugin", indexedDBPlugin)
	idbPouch := pouchDB.Call("defaults", map[string]interface{}{
		"adapter": "indexeddb",
	})
	js.Global.Set("PouchDB", idbPouch)
}

func TestPurge(t *testing.T) {
	client, err := kivik.New("pouch", "")
	if err != nil {
		t.Errorf("Failed to connect to PouchDB/memdown driver: %s", err)
		return
	}
	v, _ := client.Version(context.Background())
	pouchVer := v.Version

	t.Run("not found", func(t *testing.T) {
		if strings.HasPrefix(pouchVer, "7.") {
			t.Skipf("Skipping PouchDB 8 test for PouchDB %v", pouchVer)
		}
		const wantErr = "not_found: missing"
		client, err := kivik.New("pouch", "")
		if err != nil {
			t.Errorf("Failed to connect to PouchDB/memdown driver: %s", err)
			return
		}
		dbname := kt.TestDBName(t)
		ctx := context.Background()
		t.Cleanup(func() {
			_ = client.DestroyDB(ctx, dbname)
		})
		if e := client.CreateDB(ctx, dbname); e != nil {
			t.Fatalf("Failed to create db: %s", e)
		}
		_, err = client.DB(dbname).Purge(ctx, map[string][]string{"foo": {"1-xxx"}})
		if !testy.ErrorMatches(wantErr, err) {
			t.Errorf("Unexpected error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		if strings.HasPrefix(pouchVer, "7.") {
			t.Skipf("Skipping PouchDB 8 test for PouchDB %v", pouchVer)
		}
		client, err := kivik.New("pouch", "")
		if err != nil {
			t.Errorf("Failed to connect to PouchDB/memdown driver: %s", err)
			return
		}
		const docID = "test"
		dbname := kt.TestDBName(t)
		ctx := context.Background()
		t.Cleanup(func() {
			_ = client.DestroyDB(ctx, dbname)
		})
		if e := client.CreateDB(ctx, dbname); e != nil {
			t.Fatalf("Failed to create db: %s", e)
		}
		db := client.DB(dbname)
		rev, err := db.Put(ctx, docID, map[string]string{"foo": "bar"})
		if err != nil {
			t.Fatal(err)
		}
		result, err := db.Purge(ctx, map[string][]string{docID: {rev}})
		if err != nil {
			t.Fatal(err)
		}
		want := &kivik.PurgeResult{
			Seq: 0,
			Purged: map[string][]string{
				docID: {rev},
			},
		}
		if d := cmp.Diff(want, result); d != "" {
			t.Error(d)
		}
	})
}
