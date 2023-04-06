// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package pouchdb

import (
	"context"
	"fmt"
	"io"

	"github.com/gopherjs/gopherjs/js"

	kivik "github.com/go-kivik/kivik/v4"
	"github.com/go-kivik/kivik/v4/driver"
)

type bulkResult struct {
	*js.Object
	OK         bool   `js:"ok"`
	ID         string `js:"id"`
	Rev        string `js:"rev"`
	Error      string `js:"name"`
	StatusCode int    `js:"status"`
	Reason     string `js:"message"`
	IsError    bool   `js:"error"`
}

type bulkResults struct {
	results *js.Object
}

var _ driver.BulkResults = &bulkResults{}

func (r *bulkResults) Next(update *driver.BulkResult) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	if r.results == js.Undefined || r.results.Length() == 0 {
		return io.EOF
	}
	result := &bulkResult{}
	result.Object = r.results.Call("shift")
	update.ID = result.ID
	update.Rev = result.ID
	update.Error = nil
	if result.IsError {
		update.Error = &kivik.Error{Status: result.StatusCode, Message: result.Reason}
	}
	return nil
}

func (r *bulkResults) Close() error {
	r.results = nil // Free up memory used by any remaining rows
	return nil
}

func (d *db) BulkDocs(ctx context.Context, docs []interface{}, options map[string]interface{}) (driver.BulkResults, error) {
	result, err := d.db.BulkDocs(ctx, docs, options)
	if err != nil {
		return nil, err
	}
	return &bulkResults{results: result}, nil
}
