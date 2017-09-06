package ft

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type IndexDoc struct {
	Id         []byte
	IndexValue []byte
	StoreValue []byte
}

func IndexizeWord(w string) string {
	return strings.TrimSpace(strings.ToLower(w))
}

type WordSplitter func(string) []string

type WordCleaner func(string) string

type StopWordChecker func(string) bool

var wordizeRe *regexp.Regexp

func Wordize(t string) []string {
	return wordizeRe.Split(t, -1)
}

func WordSplit(t string) []string {
    t = strings.Replace(t, "-", " ", -1)
    t = strings.Replace(t, "_", " ", -1)
	return strings.Fields(t)
}

func WordTrim(t string) string {
	return strings.Trim(t, "")
}

type Indexer struct {
	WordMap       map[string]map[string]int
	WordSplit     WordSplitter
	WordClean     WordCleaner
	StopWordCheck StopWordChecker
	Db *DB
}

func (idx *Indexer) AddDoc(idoc IndexDoc) error {
	docId := string(idoc.Id)

	idx.WordSplit = WordSplit
	idx.WordClean = IndexizeWord
	idx.StopWordCheck = CheckStopWord
	idx.WordMap = make(map[string]map[string]int)

	idx.Db.Write([]byte(docId), idoc.StoreValue)

	words := append(idx.WordSplit(string(idoc.IndexValue)), idx.WordSplit(string(idoc.StoreValue))...)
	for _, word := range words {
		word = idx.WordClean(word)

		if idx.StopWordCheck != nil {
			if idx.StopWordCheck(word) {
				continue
			}
		}

		if idx.WordMap[word] == nil {
			idx.WordMap[word] = make(map[string]int)
		}

		c := idx.WordMap[word][docId] + 1
		idx.WordMap[word][docId] = c
	}

	idx.Commit()

	return nil
}

func (idx *Indexer) ReadPreviousValue(word string, docCountMap map[string]int) (map[string]int, error) {

	res, err := idx.Db.Read([]byte(word))
	if err != nil {
		return docCountMap, err
	}

	m := make(map[string]int)
	err = json.Unmarshal(res, &m)
	if err != nil {
		return docCountMap, err
	}

	for docId, cnt := range m {
		docCountMap[docId] = cnt
	}

	return docCountMap, nil
}

func (idx *Indexer) Commit() error {
	for word, m := range idx.WordMap {
        updatedWord, err := idx.ReadPreviousValue(word, m) 
		if err != nil {
			//fmt.Println("Indexer: error reading previous value", err)
		}	

		b, err := json.Marshal(updatedWord)
		if err != nil {
			fmt.Println("Indexer: json marshal", err)
			return err
		}	

		err = idx.Db.Write([]byte(word), b)
	}

	idx.WordMap = make(map[string]map[string]int)
	return nil
}
