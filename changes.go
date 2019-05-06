package pouchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gopherjs/gopherjs/js"

	"github.com/go-kivik/kivik/driver"
	"github.com/go-kivik/pouchdb/bindings"
)

type changesFeed struct {
	changes *js.Object
	feed    chan *driver.Change
	err     error
}

var _ driver.Changes = &changesFeed{}

type changeRow struct {
	*js.Object
	ID      string     `js:"id"`
	Seq     string     `js:"seq"`
	Changes *js.Object `js:"changes"`
	Doc     *js.Object `js:"doc"`
	Deleted bool       `js:"deleted"`
}

func (c *changesFeed) Next(row *driver.Change) error {
	if c.err != nil {
		return c.err
	}
	newRow, ok := <-c.feed
	if !ok {
		return io.EOF
	}
	*row = *newRow
	return nil
}

func (c *changesFeed) Close() error {
	c.changes.Call("cancel")
	return nil
}

// LastSeq returns an empty string.
func (c *changesFeed) LastSeq() string {
	return ""
}

// Pending returns 0
func (c *changesFeed) Pending() int64 {
	return 0
}

func (c *changesFeed) change(change *changeRow) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				_ = c.Close()
				if e, ok := r.(error); ok {
					c.err = e
				} else {
					c.err = fmt.Errorf("%v", r)
				}
			}
		}()
		changedRevs := make([]string, 0, change.Changes.Length())
		for i := 0; i < change.Changes.Length(); i++ {
			changedRevs = append(changedRevs, change.Changes.Index(i).Get("rev").String())
		}
		var doc json.RawMessage
		if change.Doc != js.Undefined {
			doc = json.RawMessage(js.Global.Get("JSON").Call("stringify", change.Doc).String())
		}
		row := &driver.Change{
			ID:      change.ID,
			Seq:     change.Seq,
			Deleted: change.Deleted,
			Doc:     doc,
			Changes: changedRevs,
		}
		c.feed <- row
	}()
}

func (c *changesFeed) complete(info *js.Object) {
	close(c.feed)
}

func (c *changesFeed) error(e *js.Object) {
	c.err = bindings.NewPouchError(e)
}

func (d *db) Changes(ctx context.Context, options map[string]interface{}) (driver.Changes, error) {
	changes, err := d.db.Changes(ctx, options)
	if err != nil {
		return nil, err
	}

	feed := make(chan *driver.Change, 32)
	c := &changesFeed{
		changes: changes,
		feed:    feed,
	}

	changes.Call("on", "change", c.change)
	changes.Call("on", "complete", c.complete)
	changes.Call("on", "error", c.error)
	return c, nil
}
