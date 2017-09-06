```
package main

import (
    "encoding/json"
    "fmt"
    "ft"
    "io"
    "net/http"
    "strconv"
)

var db = ft.NewDB("./bolt.db")
var indexer = ft.Indexer{Db: &db} //, WordMap: make(map[string]map[string]int)}
var searcher = ft.Searcher{Db: &db}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    doc := r.URL.Query().Get("document")
    value := r.URL.Query().Get("value")
    fmt.Println("id", id, "index-document", doc, "store-value", value)

    err := indexer.AddDoc(ft.IndexDoc{Id: []byte(id), IndexValue: []byte(doc), StoreValue: []byte(value)})
    if err != nil {
        http.Error(w, "Error from server", 500)
    }
    io.WriteString(w, "indexed")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    searchTerm := r.URL.Query().Get("q")

    var diff int64 = 10
    if r.URL.Query().Get("diff") != "" {
        diff, _ = strconv.ParseInt(r.URL.Query().Get("diff"), 10, 64)
    }

    var limit = 10
    if r.URL.Query().Get("limit") != "" {
        limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
    }

    fmt.Println("searchTerm", searchTerm, "diff", diff, "limit", limit) 

    sr, err := searcher.Search(searchTerm, limit)
    if err != nil { 
        fmt.Println(err) 
        http.Error(w, "Error from server", 500)
    }

    var result = make([]map[string]string, 0)
    var highest int64 = -1

    for _, v := range sr.Items {
        if highest < v.Score {
            highest = v.Score
        }

        sDiff := highest - v.Score; if sDiff > diff {break}

        var item = make(map[string]string)
        item["id"] = fmt.Sprintf("%s", v.Id)
        item["score"] = fmt.Sprintf("%d", v.Score)
        item["value"] = fmt.Sprintf("%s", v.StoreValue)
        result = append(result, item)
    }

    j, err := json.Marshal(result)
    if err != nil {
        http.Error(w, "Error from server", 500)
    }
    io.WriteString(w, string(j))
}

func main() {
    http.HandleFunc("/index", indexHandler)
    http.HandleFunc("/search", searchHandler)
    http.ListenAndServe(":8080", nil)
}
```
