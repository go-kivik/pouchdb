package pouchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gopherjs/gopherjs/js"

	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kivik/driver"
	"github.com/go-kivik/kivik/errors"
	"github.com/go-kivik/pouchdb/bindings"
)

type db struct {
	db *bindings.DB

	client *client

	// compacting is set true when compaction begins, and unset when the
	// callback returns.
	compacting bool
}

var _ driver.DB = &db{}

func (d *db) AllDocs(ctx context.Context, options map[string]interface{}) (driver.Rows, error) {
	result, err := d.db.AllDocs(ctx, options)
	if err != nil {
		return nil, err
	}
	return &rows{
		Object: result,
	}, nil
}

func (d *db) Query(ctx context.Context, ddoc, view string, options map[string]interface{}) (driver.Rows, error) {
	result, err := d.db.Query(ctx, ddoc, view, options)
	if err != nil {
		return nil, err
	}
	return &rows{
		Object: result,
	}, nil
}

func (d *db) Get(ctx context.Context, docID string, options map[string]interface{}) (*driver.Document, error) {
	doc, rev, err := d.db.Get(ctx, docID, options)
	if err != nil {
		return nil, err
	}
	return &driver.Document{
		ContentLength: int64(len(doc)),
		Rev:           rev,
		Body:          ioutil.NopCloser(bytes.NewReader(doc)),
	}, nil
}

func (d *db) CreateDoc(ctx context.Context, doc interface{}, options map[string]interface{}) (docID, rev string, err error) {
	jsonDoc, err := json.Marshal(doc)
	if err != nil {
		return "", "", err
	}
	jsDoc := js.Global.Get("JSON").Call("parse", string(jsonDoc))
	return d.db.Post(ctx, jsDoc, options)
}

func (d *db) Put(ctx context.Context, docID string, doc interface{}, options map[string]interface{}) (rev string, err error) {
	jsonDoc, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	jsDoc := js.Global.Get("JSON").Call("parse", string(jsonDoc))
	if id := jsDoc.Get("_id"); id != js.Undefined {
		if id.String() != docID {
			return "", errors.Status(kivik.StatusBadAPICall, "id argument must match _id field in document")
		}
	}
	jsDoc.Set("_id", docID)
	return d.db.Put(ctx, jsDoc, options)
}

func (d *db) Delete(ctx context.Context, docID, rev string, options map[string]interface{}) (newRev string, err error) {
	return d.db.Delete(ctx, docID, rev, options)
}

func (d *db) Stats(ctx context.Context) (*driver.DBStats, error) {
	i, err := d.db.Info(ctx)
	return &driver.DBStats{
		Name:           i.Name,
		CompactRunning: d.compacting,
		DocCount:       i.DocCount,
		UpdateSeq:      i.UpdateSeq,
	}, err
}

func (d *db) Compact(_ context.Context) error {
	d.compacting = true
	go func() {
		defer func() { d.compacting = false }()
		if err := d.db.Compact(); err != nil {
			fmt.Fprintf(os.Stderr, "compaction failed: %s\n", err)
		}
	}()
	return nil
}

// CompactView  is unimplemented for PouchDB.
func (d *db) CompactView(_ context.Context, _ string) error {
	return nil
}

func (d *db) ViewCleanup(_ context.Context) error {
	d.compacting = true
	go func() { // FIXME: #14
		defer func() { d.compacting = false }()
		if err := d.db.ViewCleanup(); err != nil {
			fmt.Fprintf(os.Stderr, "view cleanup failed: %s\n", err)
		}
	}()
	return nil
}

var securityNotImplemented = errors.Status(kivik.StatusNotImplemented, "kivik: security interface not supported by PouchDB")

func (d *db) Security(ctx context.Context) (*driver.Security, error) {
	return nil, securityNotImplemented
}

func (d *db) SetSecurity(_ context.Context, _ *driver.Security) error {
	return securityNotImplemented
}
