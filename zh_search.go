package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/yanyiwu/gojieba"
	_ "github.com/yanyiwu/gojieba/bleve"
)

// Message todo
type Message struct {
	ID   string
	Body string
}

// ReadFile todo
func ReadFile(name string) (string, error) {
	content, err := ioutil.ReadFile(name)
	return string(content), err
}

// Example todo
func Example() {
	IndexDir := "gojieba.bleve"
	var messages []Message
	a = 1

	filepath.Walk("/home/user/Dropbox/vimwiki/", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			content, _ := ReadFile(path)
			messages = append(messages, Message{
				ID:   path,
				Body: content,
			})
		}
		return nil
	})

	indexMapping := bleve.NewIndexMapping()
	os.RemoveAll(IndexDir)
	// clean index when example finished
	defer os.RemoveAll(IndexDir)

	err := indexMapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dictpath":     gojieba.DICT_PATH,
			"hmmpath":      gojieba.HMM_PATH,
			"userdictpath": gojieba.USER_DICT_PATH,
			"idf":          gojieba.IDF_PATH,
			"stop_words":   gojieba.STOP_WORDS_PATH,
			"type":         "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	err = indexMapping.AddCustomAnalyzer("gojieba",
		map[string]interface{}{
			"type":      "gojieba",
			"tokenizer": "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	indexMapping.DefaultAnalyzer = "gojieba"

	index, err := bleve.New(IndexDir, indexMapping)
	if err != nil {
		panic(err)
	}
	for _, msg := range messages {
		if err := index.Index(msg.ID, msg); err != nil {
			panic(err)
		}
	}

	querys := []string{
		"python",
		"推薦",
	}

	for _, q := range querys {
		req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
		req.Highlight = bleve.NewHighlight()
		res, err := index.Search(req)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}

	return
}

func main() {
	Example()
}
