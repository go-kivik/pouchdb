package pouchdb

import (
	"context"
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"

	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kiviktest/kt"
)

func TestChanges(t *testing.T) {
	type tst struct {
		opts            map[string]interface{}
		status          int
		err             string
		changesErr      string
		expectedIDs     []string
		expectedLastSeq string
		expectedPending int64
	}
	tests := testy.NewTable()

	tests.Run(t, func(t *testing.T, test tst) {
		ctx := context.Background()
		client, err := kivik.New("pouch", "")
		if err != nil {
			t.Fatalf("Failed to connect to PouchDB/memdown driver: %s", err)
		}
		dbname := kt.TestDBName(t)
		defer client.DestroyDB(ctx, dbname) // nolint: errcheck
		if err := client.CreateDB(ctx, dbname); err != nil {
			t.Fatalf("Failed to create db: %s", err)
		}
		db := client.DB(ctx, dbname)
		changes, err := db.Changes(ctx, test.opts)
		testy.StatusError(t, test.err, test.status, err)
		results := []string{}
		for changes.Next() {
			results = append(results, changes.ID())
		}
		testy.Error(t, test.changesErr, changes.Err())
		if d := diff.TextSlices(test.expectedIDs, results); d != nil {
			t.Error(d)
		}
		if ls := changes.LastSeq(); ls != test.expectedLastSeq {
			t.Errorf("Unexpected last_seq: %s", ls)
		}
		if p := changes.Pending(); p != test.expectedPending {
			t.Errorf("Unexpected pending count: %d", p)
		}
	})
}
