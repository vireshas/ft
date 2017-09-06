package ft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

type Searcher struct {
	Db *DB
}

type SearchResultItem struct {
	Id         []byte 
	StoreValue []byte 
	Score      int64  
}

type SearchResultItems []SearchResultItem

func (s SearchResultItems) Len() int      { return len(s) }
func (s SearchResultItems) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SearchResultItems) Less(i, j int) bool {
	if s[i].Score == s[j].Score {
		return bytes.Compare(s[i].Id, s[j].Id) < 0
	}
	return s[i].Score < s[j].Score
}

type SearchResults struct {
	Items SearchResultItems
}

func NewSearcher(db *DB) *Searcher {
	return &Searcher{db}
}

func (s *Searcher) Search(search string, maxn int) (SearchResults, error) {
	sr := SearchResults{}

	searchWords := WordSplit(search)

	itemMap := make(map[string]SearchResultItem)

	for _, w := range searchWords {
		w = IndexizeWord(w)

		if CheckStopWord(w) { continue }	

		res, err := s.Db.Read([]byte(w))
		if err != nil {
			return sr, err
		}

		m := make(map[string]int)
		err = json.Unmarshal(res, &m)
		if err != nil {
			return sr, err
		}

		for docId, cnt := range m {
			sri := itemMap[docId]
			if sri.Score < 1 {
				sri.Id = []byte(docId)
			}
			sri.Score += int64(cnt)
			itemMap[docId] = sri
		}
	}

	items := make(SearchResultItems, 0, maxn)
	for _, item := range itemMap {
		items = append(items, item)
	}

	sort.Sort(sort.Reverse(items))
	fmt.Println("total matches", len(items))
	if len(items) > maxn {
		items = items[:maxn]
	}

	for i := range items {
		item := &items[i]
		v, err := s.Db.Read(item.Id)
		if err != nil {
			return sr, err
		}
		item.StoreValue = v
	}

	sr.Items = items

	return sr, nil
}
