package index

import (
	"testing"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/stretchr/testify/assert"
)

func makeTestDoc(id, text, source string, meta map[string]string, vector []float64) models.Document {
	return models.Document{
		ID:     id,
		Text:   text,
		Source: source,
		Meta:   meta,
		Vector: vector,
	}
}

func TestSimpleIndex_AddAndGetDocument(t *testing.T) {
	idx := NewSimpleIndex()
	doc := makeTestDoc("1", "hello world", "file1.txt", map[string]string{"author": "alice"}, []float64{1.0, 2.0})
	err := idx.AddDocument(doc)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(idx.documents))
	assert.Equal(t, doc, idx.documents[doc.ID])
}

func TestSimpleIndex_AddDocuments(t *testing.T) {
	idx := NewSimpleIndex()
	docs := []models.Document{
		makeTestDoc("1", "foo", "a.txt", nil, nil),
		makeTestDoc("2", "bar", "b.txt", nil, nil),
	}
	err := idx.AddDocuments(docs)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(idx.documents))
}

func TestSimpleIndex_DeleteDocument(t *testing.T) {
	idx := NewSimpleIndex()
	doc := makeTestDoc("1", "text", "src", nil, nil)
	_ = idx.AddDocument(doc)
	err := idx.DeleteDocument("1")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(idx.documents))
	// Try deleting non-existent
	err = idx.DeleteDocument("notfound")
	assert.Error(t, err)
}

func TestSimpleIndex_DeleteDocuments(t *testing.T) {
	idx := NewSimpleIndex()
	docs := []models.Document{
		makeTestDoc("1", "foo", "a.txt", nil, nil),
		makeTestDoc("2", "bar", "b.txt", nil, nil),
	}
	_ = idx.AddDocuments(docs)
	err := idx.DeleteDocuments([]string{"1", "2"})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(idx.documents))
}

func TestSimpleIndex_UpdateDocument(t *testing.T) {
	idx := NewSimpleIndex()
	doc := makeTestDoc("1", "old", "src", nil, nil)
	_ = idx.AddDocument(doc)
	updated := makeTestDoc("1", "new", "src", nil, nil)
	err := idx.UpdateDocument("1", updated)
	assert.NoError(t, err)
	assert.Equal(t, "new", idx.documents["1"].Text)
	// Update non-existent
	err = idx.UpdateDocument("notfound", updated)
	assert.Error(t, err)
}

func TestSimpleIndex_UpdateDocuments(t *testing.T) {
	idx := NewSimpleIndex()
	_ = idx.AddDocuments([]models.Document{
		makeTestDoc("1", "a", "a.txt", nil, nil),
		makeTestDoc("2", "b", "b.txt", nil, nil),
	})
	updates := []models.Document{
		makeTestDoc("1", "A", "a.txt", nil, nil),
		makeTestDoc("2", "B", "b.txt", nil, nil),
	}
	err := idx.UpdateDocuments(updates)
	assert.NoError(t, err)
	assert.Equal(t, "A", idx.documents["1"].Text)
	assert.Equal(t, "B", idx.documents["2"].Text)
}

func TestSimpleIndex_ConfigureAndShowConfig(t *testing.T) {
	idx := NewSimpleIndex()
	cfg := map[string]interface{}{"foo": 1, "bar": true}
	err := idx.Configure(cfg)
	assert.NoError(t, err)
	conf, err := idx.ShowConfig()
	assert.NoError(t, err)
	assert.Equal(t, cfg, conf)
	// Ensure returned config is a copy
	conf["foo"] = 999
	conf2, _ := idx.ShowConfig()
	assert.Equal(t, 1, conf2["foo"])
}

func TestSimpleIndex_CountAndSize(t *testing.T) {
	idx := NewSimpleIndex()
	doc := makeTestDoc("1", "abc", "src", map[string]string{"k": "v"}, []float64{1.1, 2.2})
	_ = idx.AddDocument(doc)
	count, err := idx.Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	size, err := idx.Size()
	assert.NoError(t, err)
	assert.Greater(t, size, 0)
}

func TestSimpleIndex_CloseFlushOptimize(t *testing.T) {
	idx := NewSimpleIndex()
	assert.NoError(t, idx.Close())
	assert.NoError(t, idx.Flush())
	assert.NoError(t, idx.Optimize())
}

func TestSimpleIndex_SearchSimple(t *testing.T) {
	idx := NewSimpleIndex()
	docs := []models.Document{
		makeTestDoc("1", "hello world", "src1", map[string]string{"author": "alice"}, nil),
		makeTestDoc("2", "foo bar", "src2", map[string]string{"author": "bob"}, nil),
	}
	_ = idx.AddDocuments(docs)
	results, err := idx.Search("hello")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "1", results[0].ID)
	results, _ = idx.Search("bob")
	assert.Len(t, results, 1)
	assert.Equal(t, "2", results[0].ID)
	results, _ = idx.Search("")
	assert.Len(t, results, 0)
}
