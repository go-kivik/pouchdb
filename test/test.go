package test

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/flimzy/kivik"
	"github.com/go-kivik/kiviktest"
	"github.com/go-kivik/kiviktest/kt"
	"github.com/gopherjs/gopherjs/js"
)

func init() {
	if pouchDB := js.Global.Get("PouchDB"); pouchDB != js.Undefined {
		memPouch := js.Global.Get("PouchDB").Call("defaults", map[string]interface{}{
			"db": js.Global.Call("require", "memdown"),
		})
		js.Global.Set("PouchDB", memPouch)
	}
}

// RegisterPouchDBSuites registers the PouchDB test suites.
func RegisterPouchDBSuites() {
	kiviktest.RegisterSuite(kiviktest.SuitePouchLocal, kt.SuiteConfig{
		"PreCleanup.skip": true,

		// Features which are not supported by PouchDB
		"Log.skip":         true,
		"Flush.skip":       true,
		"Security.skip":    true, // FIXME: Perhaps implement later with a plugin?
		"SetSecurity.skip": true, // FIXME: Perhaps implement later with a plugin?
		"DBUpdates.skip":   true,

		"AllDBs.skip":   true, // FIXME: Find a way to test with the plugin
		"CreateDB.skip": true, // FIXME: No way to validate if this works unless/until allDbs works
		"DBExists.skip": true, // FIXME: Maybe fix this if/when allDBs works?

		"AllDocs/Admin.databases":                        []string{},
		"AllDocs/RW/group/Admin/WithDocs/UpdateSeq.skip": true,

		"Find/Admin.databases":                []string{},
		"Find/RW/group/Admin/Warning.warning": "no matching index found, create an index to optimize query time",

		"Explain.databases": []string{},
		"Explain.plan": &kivik.QueryPlan{
			Index: map[string]interface{}{
				"ddoc": nil,
				"name": "_all_docs",
				"type": "special",
				"def":  map[string]interface{}{"fields": []interface{}{map[string]string{"_id": "asc"}}},
			},
			Selector: map[string]interface{}{"_id": map[string]interface{}{"$gt": nil}},
			Options: map[string]interface{}{
				"bookmark":  "nil",
				"conflicts": false,
				"r":         []int{49},
				"sort":      map[string]interface{}{},
				"use_index": []interface{}{},
			},
			Fields: []interface{}{},
			Range: map[string]interface{}{
				"start_key": nil,
			},
		},

		"Query/RW/group/Admin/WithDocs/UpdateSeq.skip": true,

		"Version.version":        `^6\.\d\.\d$`,
		"Version.vendor":         `^PouchDB$`,
		"Version.vendor_version": `^6\.\d\.\d$`,

		"Get/RW/group/Admin/bogus.status":  kivik.StatusNotFound,
		"Get/RW/group/NoAuth/bogus.status": kivik.StatusNotFound,

		"Rev/RW/group/Admin/bogus.status":  kivik.StatusNotFound,
		"Rev/RW/group/NoAuth/bogus.status": kivik.StatusNotFound,

		"Delete/RW/Admin/group/MissingDoc.status":       kivik.StatusNotFound,
		"Delete/RW/Admin/group/InvalidRevFormat.status": kivik.StatusBadRequest,
		"Delete/RW/Admin/group/WrongRev.status":         kivik.StatusConflict,

		"Stats/Admin.skip": true, // No predefined DBs for Local PouchDB

		"BulkDocs/RW/Admin/group/Mix/Conflict.status": kivik.StatusConflict,

		"GetAttachment/RW/group/Admin/foo/NotFound.status": kivik.StatusNotFound,

		"GetAttachmentMeta/RW/group/Admin/foo/NotFound.status": kivik.StatusNotFound,

		"PutAttachment/RW/group/Admin/Conflict.status": kivik.StatusConflict,

		// "DeleteAttachment/RW/group/Admin/NotFound.status": kivik.StatusNotFound, // https://github.com/pouchdb/pouchdb/issues/6409
		"DeleteAttachment/RW/group/Admin/NoDoc.status": kivik.StatusNotFound,

		"Put/RW/Admin/group/LeadingUnderscoreInID.status": kivik.StatusBadRequest,
		"Put/RW/Admin/group/Conflict.status":              kivik.StatusConflict,

		"CreateIndex/RW/Admin/group/EmptyIndex.status":   kivik.StatusInternalServerError,
		"CreateIndex/RW/Admin/group/BlankIndex.status":   kivik.StatusBadRequest,
		"CreateIndex/RW/Admin/group/InvalidIndex.status": kivik.StatusInternalServerError,
		"CreateIndex/RW/Admin/group/NilIndex.status":     kivik.StatusInternalServerError,
		"CreateIndex/RW/Admin/group/InvalidJSON.status":  kivik.StatusBadRequest,

		"GetIndexes.databases": []string{},

		"DeleteIndex/RW/Admin/group/NotFoundDdoc.status": kivik.StatusNotFound,
		"DeleteIndex/RW/Admin/group/NotFoundName.status": kivik.StatusNotFound,

		"Replicate.skip": true, // No need to do this for both Local and Remote

		"Query/RW/group/Admin/WithoutDocs/ScanDoc.status": kivik.StatusBadRequest,
	})
	kiviktest.RegisterSuite(kiviktest.SuitePouchRemote, kt.SuiteConfig{
		// Features which are not supported by PouchDB
		"Log.skip":         true,
		"Flush.skip":       true,
		"Session.skip":     true,
		"Security.skip":    true, // FIXME: Perhaps implement later with a plugin?
		"SetSecurity.skip": true, // FIXME: Perhaps implement later with a plugin?
		"DBUpdates.skip":   true,

		"PreCleanup.skip": true,

		"AllDBs.skip": true, // FIXME: Perhaps a workaround can be found?

		"CreateDB/RW/NoAuth.status":         kivik.StatusUnauthorized,
		"CreateDB/RW/Admin/Recreate.status": kivik.StatusPreconditionFailed,

		"DBExists.databases":              []string{"_users", "chicken"},
		"DBExists/Admin/_users.exists":    true,
		"DBExists/Admin/chicken.exists":   false,
		"DBExists/NoAuth/_users.exists":   true,
		"DBExists/NoAuth/chicken.exists":  false,
		"DBExists/RW/group/Admin.exists":  true,
		"DBExists/RW/group/NoAuth.exists": true,

		"DestroyDB/RW/NoAuth/NonExistantDB.status": kivik.StatusNotFound,
		"DestroyDB/RW/Admin/NonExistantDB.status":  kivik.StatusNotFound,
		"DestroyDB/RW/NoAuth/ExistingDB.status":    kivik.StatusUnauthorized,

		"AllDocs.databases":                                  []string{"_replicator", "_users", "chicken"},
		"AllDocs/Admin/_replicator.expected":                 []string{"_design/_replicator"},
		"AllDocs/Admin/_replicator.offset":                   0,
		"AllDocs/Admin/_users.expected":                      []string{"_design/_auth"},
		"AllDocs/Admin/chicken.status":                       kivik.StatusNotFound,
		"AllDocs/NoAuth/_replicator.status":                  kivik.StatusUnauthorized,
		"AllDocs/NoAuth/_users.status":                       kivik.StatusUnauthorized,
		"AllDocs/NoAuth/chicken.status":                      kivik.StatusNotFound,
		"AllDocs/Admin/_replicator/WithDocs/UpdateSeq.skip":  true,
		"AllDocs/Admin/_users/WithDocs/UpdateSeq.skip":       true,
		"AllDocs/RW/group/Admin/WithDocs/UpdateSeq.skip":     true,
		"AllDocs/RW/group/Admin/WithoutDocs/UpdateSeq.skip":  true,
		"AllDocs/RW/group/NoAuth/WithDocs/UpdateSeq.skip":    true,
		"AllDocs/RW/group/NoAuth/WithoutDocs/UpdateSeq.skip": true,

		"Find.databases":                       []string{"chicken", "_duck"},
		"Find/Admin/chicken.status":            kivik.StatusNotFound,
		"Find/Admin/_duck.status":              kivik.StatusNotFound,
		"Find/NoAuth/chicken.status":           kivik.StatusNotFound,
		"Find/NoAuth/_duck.status":             kivik.StatusUnauthorized,
		"Find/RW/group/Admin/Warning.warning":  "no matching index found, create an index to optimize query time",
		"Find/RW/group/NoAuth/Warning.warning": "no matching index found, create an index to optimize query time",

		"Explain.databases":             []string{"chicken", "_duck"},
		"Explain/Admin/chicken.status":  kivik.StatusNotFound,
		"Explain/Admin/_duck.status":    kivik.StatusNotFound,
		"Explain/NoAuth/chicken.status": kivik.StatusNotFound,
		"Explain/NoAuth/_duck.status":   kivik.StatusUnauthorized,
		"Explain.plan": &kivik.QueryPlan{
			Index: map[string]interface{}{
				"ddoc": nil,
				"name": "_all_docs",
				"type": "special",
				"def":  map[string]interface{}{"fields": []interface{}{map[string]string{"_id": "asc"}}},
			},
			Selector: map[string]interface{}{"_id": map[string]interface{}{"$gt": nil}},
			Options: map[string]interface{}{
				"bookmark":        "nil",
				"conflicts":       false,
				"execution_stats": false,
				"r":               []int{49},
				"sort":            map[string]interface{}{},
				"use_index":       []interface{}{},
				"stable":          false,
				"stale":           false,
				"update":          true,
				"skip":            0,
				"limit":           25,
				"fields":          "all_fields",
			},
			Fields: []interface{}{},
			Range:  nil,
			Limit:  25,
		},

		"CreateIndex/RW/Admin/group/EmptyIndex.status":    kivik.StatusBadRequest,
		"CreateIndex/RW/Admin/group/BlankIndex.status":    kivik.StatusBadRequest,
		"CreateIndex/RW/Admin/group/InvalidIndex.status":  kivik.StatusBadRequest,
		"CreateIndex/RW/Admin/group/NilIndex.status":      kivik.StatusBadRequest,
		"CreateIndex/RW/Admin/group/InvalidJSON.status":   kivik.StatusBadRequest,
		"CreateIndex/RW/NoAuth/group/EmptyIndex.status":   kivik.StatusBadRequest,
		"CreateIndex/RW/NoAuth/group/BlankIndex.status":   kivik.StatusBadRequest,
		"CreateIndex/RW/NoAuth/group/InvalidIndex.status": kivik.StatusBadRequest,
		"CreateIndex/RW/NoAuth/group/NilIndex.status":     kivik.StatusBadRequest,
		"CreateIndex/RW/NoAuth/group/InvalidJSON.status":  kivik.StatusBadRequest,
		"CreateIndex/RW/NoAuth/group/Valid.status":        kivik.StatusInternalServerError, // COUCHDB-3374

		"GetIndexes.databases":                     []string{"_replicator", "_users", "_global_changes"},
		"GetIndexes/Admin/_replicator.indexes":     []kivik.Index{kt.AllDocsIndex},
		"GetIndexes/Admin/_users.indexes":          []kivik.Index{kt.AllDocsIndex},
		"GetIndexes/Admin/_global_changes.indexes": []kivik.Index{kt.AllDocsIndex},
		"GetIndexes/NoAuth/_replicator.indexes":    []kivik.Index{kt.AllDocsIndex},
		"GetIndexes/NoAuth/_users.indexes":         []kivik.Index{kt.AllDocsIndex},
		"GetIndexes/NoAuth/_global_changes.skip":   true, // Pouch connects to the DB before searching the Index, so this test fails
		"GetIndexes/NoAuth/_global_changes.status": kivik.StatusUnauthorized,
		"GetIndexes/RW.indexes": []kivik.Index{kt.AllDocsIndex,
			{
				DesignDoc: "_design/foo",
				Name:      "bar",
				Type:      "json",
				Definition: map[string]interface{}{
					"fields": []map[string]string{
						{"foo": "asc"},
					},
					"partial_filter_selector": map[string]string{},
				},
			},
		},

		"DeleteIndex/RW/Admin/group/NotFoundDdoc.status":  kivik.StatusNotFound,
		"DeleteIndex/RW/Admin/group/NotFoundName.status":  kivik.StatusNotFound,
		"DeleteIndex/RW/NoAuth/group/NotFoundDdoc.status": kivik.StatusNotFound,
		"DeleteIndex/RW/NoAuth/group/NotFoundName.status": kivik.StatusNotFound,

		"Query/RW/group/Admin/WithDocs/UpdateSeq.skip":  true,
		"Query/RW/group/NoAuth/WithDocs/UpdateSeq.skip": true,

		"Version.version":        `^6\.\d\.\d$`,
		"Version.vendor":         `^PouchDB$`,
		"Version.vendor_version": `^6\.\d\.\d$`,

		"Get/RW/group/Admin/bogus.status":  kivik.StatusNotFound,
		"Get/RW/group/NoAuth/bogus.status": kivik.StatusNotFound,

		"Rev/RW/group/Admin/bogus.status":  kivik.StatusNotFound,
		"Rev/RW/group/NoAuth/bogus.status": kivik.StatusNotFound,

		"Delete/RW/Admin/group/MissingDoc.status":        kivik.StatusNotFound,
		"Delete/RW/Admin/group/InvalidRevFormat.status":  kivik.StatusBadRequest,
		"Delete/RW/Admin/group/WrongRev.status":          kivik.StatusConflict,
		"Delete/RW/NoAuth/group/MissingDoc.status":       kivik.StatusNotFound,
		"Delete/RW/NoAuth/group/InvalidRevFormat.status": kivik.StatusBadRequest,
		"Delete/RW/NoAuth/group/WrongRev.status":         kivik.StatusConflict,
		"Delete/RW/NoAuth/group/DesignDoc.status":        kivik.StatusUnauthorized,

		"Stats.databases":             []string{"_users", "chicken"},
		"Stats/Admin/chicken.status":  kivik.StatusNotFound,
		"Stats/NoAuth/chicken.status": kivik.StatusNotFound,

		"BulkDocs/RW/NoAuth/group/Mix/Conflict.status": kivik.StatusConflict,
		"BulkDocs/RW/Admin/group/Mix/Conflict.status":  kivik.StatusConflict,

		"GetAttachment/RW/group/Admin/foo/NotFound.status":  kivik.StatusNotFound,
		"GetAttachment/RW/group/NoAuth/foo/NotFound.status": kivik.StatusNotFound,

		"GetAttachmentMeta/RW/group/Admin/foo/NotFound.status":  kivik.StatusNotFound,
		"GetAttachmentMeta/RW/group/NoAuth/foo/NotFound.status": kivik.StatusNotFound,

		"PutAttachment/RW/group/Admin/Conflict.status":         kivik.StatusConflict,
		"PutAttachment/RW/group/NoAuth/Conflict.status":        kivik.StatusConflict,
		"PutAttachment/RW/group/NoAuth/UpdateDesignDoc.status": kivik.StatusUnauthorized,
		"PutAttachment/RW/group/NoAuth/CreateDesignDoc.status": kivik.StatusUnauthorized,

		// "DeleteAttachment/RW/group/Admin/NotFound.status":  kivik.StatusNotFound, // COUCHDB-3362
		// "DeleteAttachment/RW/group/NoAuth/NotFound.status": kivik.StatusNotFound, // COUCHDB-3362
		"DeleteAttachment/RW/group/Admin/NoDoc.status":      kivik.StatusConflict,
		"DeleteAttachment/RW/group/NoAuth/NoDoc.status":     kivik.StatusConflict,
		"DeleteAttachment/RW/group/NoAuth/DesignDoc.status": kivik.StatusUnauthorized,

		"Put/RW/Admin/group/LeadingUnderscoreInID.status":  kivik.StatusBadRequest,
		"Put/RW/Admin/group/Conflict.status":               kivik.StatusConflict,
		"Put/RW/NoAuth/group/DesignDoc.status":             kivik.StatusUnauthorized,
		"Put/RW/NoAuth/group/LeadingUnderscoreInID.status": kivik.StatusBadRequest,
		"Put/RW/NoAuth/group/Conflict.status":              kivik.StatusConflict,

		"Replicate.NotFoundDB": func() string {
			var dsn string
			for _, env := range []string{"KIVIK_TEST_DSN_COUCH21", "KIVIK_TEST_DSN_COUCH20", "KIVIK_TEST_DSN_COUCH16", "KIVIK_TEST_DSN_CLOUDANT"} {
				dsn = os.Getenv(env)
				if dsn != "" {
					break
				}
			}
			parsed, _ := url.Parse(dsn)
			parsed.User = nil
			return strings.TrimSuffix(parsed.String(), "/") + "/doesntexist"

		}(),
		"Replicate.prefix":                                       "none",
		"Replicate.timeoutSeconds":                               5,
		"Replicate.mode":                                         "pouchdb",
		"Replicate/RW/Admin/group/MissingSource/Results.status":  kivik.StatusUnauthorized,
		"Replicate/RW/Admin/group/MissingTarget/Results.status":  kivik.StatusUnauthorized,
		"Replicate/RW/NoAuth/group/MissingSource/Results.status": kivik.StatusUnauthorized,
		"Replicate/RW/NoAuth/group/MissingTarget/Results.status": kivik.StatusUnauthorized,

		"Query/RW/group/Admin/WithoutDocs/ScanDoc.status":  kivik.StatusBadRequest,
		"Query/RW/group/NoAuth/WithoutDocs/ScanDoc.status": kivik.StatusBadRequest,

		// "ViewCleanup/RW/NoAuth.status": kivik.StatusUnauthorized, # FIXME: #14
	})
}

// PouchLocalTest runs the PouchDB tests against a local database.
func PouchLocalTest(t *testing.T) {
	client, err := kivik.New(context.Background(), "pouch", "")
	if err != nil {
		t.Errorf("Failed to connect to PouchDB driver: %s", err)
		return
	}
	clients := &kt.Context{
		RW:    true,
		Admin: client,
	}
	kiviktest.RunTestsInternal(clients, kiviktest.SuitePouchLocal, t)
}

// PouchRemoteTest runs the PouchDB tests against a remote CouchDB database.
func PouchRemoteTest(t *testing.T) {
	kiviktest.DoTest(kiviktest.SuitePouchRemote, "KIVIK_TEST_DSN_COUCH21", t)
}
