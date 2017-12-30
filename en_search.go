package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	_ "github.com/yanyiwu/gojieba/bleve"
)

type Message struct {
	Id   string
	Body string
}

// ReadFile todo
func ReadFile(name string) (string, error) {
	content, err := ioutil.ReadFile(name)
	return string(content), err
}

func Example() {
	INDEX_DIR := "gojieba.bleve"
	var messages []Message

	filepath.Walk("/home/user/Dropbox/vimwiki/", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			content, _ := ReadFile(path)
			messages = append(messages, Message{
				Id:   path,
				Body: content,
			})
		}
		return nil
	})

	indexMapping := bleve.NewIndexMapping()
	os.RemoveAll(INDEX_DIR)
	// clean index when example finished
	defer os.RemoveAll(INDEX_DIR)

	index, err := bleve.New(INDEX_DIR, indexMapping)
	if err != nil {
		panic(err)
	}
	for _, msg := range messages {
		if err := index.Index(msg.Id, msg); err != nil {
			panic(err)
		}
	}

	querys := []string{
		"python",
	}

	fmt.Printf("%v\n", index.Stats())

	for _, q := range querys {
		req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
		req.Highlight = bleve.NewHighlight()
		res, err := index.Search(req)
		if err != nil {
			panic(err)
		}
		for _, d := range res.Hits {
			fmt.Println(d.ID)
			for k, v := range d.Fragments {
				fmt.Printf("%s\n", k)
				fmt.Printf("%s\n", v)
			}
		}
	}
}

func main() {
	Example()
}
