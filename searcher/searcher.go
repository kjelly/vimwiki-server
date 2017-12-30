package searcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

// Doc todo
type Doc struct {
	Path         string
	Body         string
	ModifiedTime string
}

// New todo
func New() *Searcher {
	s := new(Searcher)
	s.Dict = make(map[string]Doc)
	return s
}

// Searcher todo
type Searcher struct {
	List         []Doc
	Dict         map[string]Doc
	index        bleve.Index
	indexMapping *mapping.IndexMappingImpl
	Version      string
}

// SearchResult todo
type SearchResult struct {
	Path      string
	Body      string
	Highlight []string
	Score     float64
}

// ReadFile todo
func ReadFile(name string) (string, error) {
	content, err := ioutil.ReadFile(name)
	return string(content), err
}

// BuildIndex todo
func (s *Searcher) BuildIndex(path string) error {
	fmt.Printf("Build index\n")
	s.Version = "SSS"
	var err error
	indexDir := "bleve"
	os.RemoveAll(indexDir)

	s.indexMapping = bleve.NewIndexMapping()
	s.index, err = bleve.New(indexDir, s.indexMapping)
	if err != nil {
		panic(err)
	}
	fmt.Printf("path=%s\n", path)
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		fmt.Printf("path=%s\n", p)
		if filepath.Ext(path) == "md" {
			content, err := ReadFile(p)
			fmt.Printf("%s\n", content)
			if err != nil {
				panic(err)
			}
			doc := Doc{
				Path:         p,
				Body:         content,
				ModifiedTime: "",
			}
			s.List = append(s.List, doc)
			s.Dict[path] = doc
			err = s.index.Index(doc.Path, doc)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
	return nil
}

// Search todo
func (s *Searcher) Search(keyword string) ([]SearchResult, error) {
	fmt.Printf("%v\n", s.index.Stats())
	ret := []SearchResult{}
	req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(keyword))
	req.Highlight = bleve.NewHighlight()
	//	index, err := bleve.Open("bleve")
	//	if err != nil {
	//		panic(err)
	//	}
	//	res, err := index.Search(req)

	res, err := s.index.Search(req)

	if err != nil {
		panic(err)
	}
	for _, d := range res.Hits {
		o := SearchResult{
			Path:      d.ID,
			Body:      s.Dict[d.ID].Body,
			Highlight: d.Fragments["Content"],
			Score:     d.Score,
		}
		ret = append(ret, o)
	}
	fmt.Printf("hits=%v\n", res.Hits)
	fmt.Printf("ret=%v\n", ret)
	fmt.Printf("%s\n", s.Version)
	return ret, err
}
