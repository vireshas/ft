package ft

import (
	"fmt"
	"ft"
	"testing"
	"github.com/stretchr/testify/assert"
)

var db = ft.NewDB("/tmp/bolt_test.db")
func TestFt(t *testing.T) {
	fmt.Printf("TestFt\n")

	idx := ft.Indexer{Db: &db}

	idx.AddDoc(ft.IndexDoc{Id: []byte(`bla1`), StoreValue: []byte(`do store this`), IndexValue: []byte(`test of the emergency broadcast system`)})
	idx.AddDoc(ft.IndexDoc{Id: []byte(`bla2`), StoreValue: []byte(`do store this stuff too, yeah just store it`), IndexValue: []byte(`every good boy does fine`)})
	idx.AddDoc(ft.IndexDoc{Id: []byte(`bla3`), StoreValue: []byte(`more storage here`), IndexValue: []byte(`a taco in the hand is worth two in the truck`)})
	idx.AddDoc(ft.IndexDoc{Id: []byte(`blah1`), StoreValue: []byte(`do store this`), IndexValue: []byte(`test of the emergency broadcast system`)})
	idx.AddDoc(ft.IndexDoc{Id: []byte(`blah2`), StoreValue: []byte(`do store this stuff too, yeah just store it`), IndexValue: []byte(`every good boy does fine`)})
	idx.AddDoc(ft.IndexDoc{Id: []byte(`blah3`), StoreValue: []byte(`more storage here`), IndexValue: []byte(`a taco in the hand is worth two in the truck`)})

	var searcher = ft.Searcher{Db: &db}
	sr, _ := searcher.Search("store", 10)

	assert.Equal(t, len(sr.Items), 4, "Search results should be of length 4")
}
