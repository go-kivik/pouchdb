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
	"net/http"
	"testing"

	"github.com/gopherjs/gopherjs/js"
	"gitlab.com/flimzy/testy"

	kivik "github.com/go-kivik/kivik/v4"
	"github.com/go-kivik/kiviktest/v4/kt"
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
