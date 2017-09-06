package ft

import (
	"fmt"
	"ft"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFt(t *testing.T) {
	fmt.Printf("TestFt\n")

    db := ft.NewDB("./bolt_test.db")
	idx := ft.Indexer{Db: &db}

	idx.AddDoc(ft.IndexDoc{Id: []byte(`blah1`), StoreValue: []byte(`store this`), IndexValue: []byte(`test of the emergency broadcast system`)})
	idx.AddDoc(ft.IndexDoc{Id: []byte(`blah2`), StoreValue: []byte(`store this stuff too, yeah store it`), IndexValue: []byte(`every good boy does fine`)})
	idx.AddDoc(ft.IndexDoc{Id: []byte(`blah3`), StoreValue: []byte(`more storage here`), IndexValue: []byte(`a taco in the hand is worth two in the truck`)})

	var searcher = ft.Searcher{Db: &db}
	sr, _ := searcher.Search("store", 10)

	assert.Equal(t, len(sr.Items), 2, "Search results should be of length 2")
	assert.Equal(t, string(sr.Items[0].StoreValue), "store this stuff too, yeah store it", "First result")
	assert.Equal(t, string(sr.Items[1].StoreValue), "store this", "Second result")
	assert.Equal(t, int(sr.Items[1].Score), int(1), "Second score")
	assert.Equal(t, int(sr.Items[0].Score), int(2), "First score")
}
