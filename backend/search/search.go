package search

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"github.com/yanyiwu/gojieba"
)

type InvertedIndex map[string][]int

type SearchEngine struct {
	Index   InvertedIndex
	Docs    map[int]string
	WordSet map[string]bool  // 用map构造set，比较简单的写法
	Terms   map[int][]string // jieba分词的结果
	jieba   *gojieba.Jieba
}

func isChinese(str string) bool {
	reg := regexp.MustCompile(`^[\p{Han}]+$`)
	return reg.MatchString(str)
}

func (s *SearchEngine) BuildInvertedIndex() {
	s.Index = make(InvertedIndex)
	s.WordSet = make(map[string]bool)
	s.Terms = make(map[int][]string)
	for id, doc := range s.Docs {
		words := s.jieba.CutForSearch(doc, true)
		for _, word := range words {
			if isChinese(word) {
				s.Terms[id] = append(s.Terms[id], word)
			}
		}
		for _, word := range words {
			s.Index[word] = append(s.Index[word], id)
			s.WordSet[word] = true
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

func InitializeSearchEngine(pagesDir string) *SearchEngine {
	dir, _ := os.Getwd()
	se := SearchEngine{}
	se.ReadFile(dir + pagesDir)
	se.BuildInvertedIndex()
	return &se
}

func FindArticleDetails(article string) string {
	title := ""
	for i := 0; i < len(article); {
		r, size := utf8.DecodeRuneInString(article[i:])
		if r == utf8.RuneError && size == 1 {
			// 如果解码失败，则将该字符视为乱码
			title += "?"
			i++
		} else {
			title += string(r)
			i += size
		}
		if r == '\n' {
			break
		}
	}
	return title
}

func (s *SearchEngine) TF_IDF() map[int]map[string]float64 {
	// 先定义DF函数
	df := func(word string) int {
		return len(s.Index[word])
	}
	// 再定义TF函数
	tf := func(word string, docId int) int {
		cnt := 0
		for _, tmpDocId := range s.Index[word] {
			if tmpDocId == docId {
				cnt++
			}
		}
		return cnt
	}
	// 计算IDF
	idf := make(map[string]float64)
	for word := range s.WordSet {
		idf[word] = math.Log10(float64(len(s.Docs)) / float64(df(word)))
	}
	// 计算TF-IDF
	tf_idf := make(map[int]map[string]float64)
	for docId := range s.Docs {
		tf_idf[docId] = make(map[string]float64)
		for _, word := range s.Terms[docId] {
			tf_val := tf(word, docId)
			if tf_val == 0 {
				tf_idf[docId][word] = 0
			} else {
				tf_idf[docId][word] = (1.0 + math.Log10(float64(tf_val))) * idf[word]
			}
		}
	}
	return tf_idf
}

func (s *SearchEngine) CosineSimlarity() map[int]map[int]float64 {
	// 计算TF-IDF
	cosine := make(map[int]map[int]float64)
	tf_idf := s.TF_IDF()
	for docId1 := range s.Docs {
		cosine[docId1] = make(map[int]float64)
		for docId2 := range s.Docs {
			v1 := tf_idf[docId1]
			v2 := tf_idf[docId2]
			v1Norm := 0.0
			v2Norm := 0.0

			// 计算余弦相似度
			numerator := 0.0
			for word := range v1 {
				if _, ok := v2[word]; ok {
					numerator += v1[word] * v2[word]
					v1Norm += v1[word] * v1[word]
					v2Norm += v2[word] * v2[word]
				}
			}
			denominator := math.Sqrt(v1Norm) * math.Sqrt(v2Norm)
			cosine[docId1][docId2] = numerator / denominator
		}
	}
	return cosine
}
