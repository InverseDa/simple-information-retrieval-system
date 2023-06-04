package search

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yanyiwu/gojieba"
)

type SearchEngine struct {
	Index InvertedIndex
	Docs  map[int]string
	jieba *gojieba.Jieba
}

type InvertedIndex map[string][]int

func (s *SearchEngine) BuildInvertedIndex() {
	s.Index = make(InvertedIndex)

	for id, doc := range s.Docs {
		words := s.jieba.CutForSearch(doc, true)
		for _, word := range words {
			s.Index[word] = append(s.Index[word], id)
		}
	}
}

func (s *SearchEngine) Search(query string) []int {
	ids := make(map[int]bool)

	words := s.jieba.CutForSearch(query, true)
	for _, word := range words {
		if docIds, ok := s.Index[word]; ok {
			for _, id := range docIds {
				ids[id] = true
			}
		}
	}

	results := make([]int, 0)
	for id := range ids {
		results = append(results, id)
	}

	return results
}

func (s *SearchEngine) ReadFile(pagesDir string) {
	s.Docs = make(map[int]string)
	s.jieba = gojieba.NewJieba()

	err := filepath.Walk(pagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("[Error] ", err)
			return err
		}
		if !info.Mode().IsRegular() || filepath.Ext(path) != ".txt" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("[Error] ", err)
			return err
		}
		id := len(s.Docs) + 1
		s.Docs[id] = string(content)

		return nil
	})
	if err != nil {
		fmt.Println("Failed to readpages directory:", err)
	}
}
