package search

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"unicode/utf8"

	"github.com/yanyiwu/gojieba"
)

type InvertedIndex map[string][]int

type SearchEngine struct {
	Index              InvertedIndex
	Docs               map[int]string
	WordSet            map[string]bool         // 用map构造set，比较简单的写法
	Terms              map[int][]string        // jieba分词的结果
	CosineSimlarityMap map[int]map[int]float64 //两个文档之间的余弦相似度
	jieba              *gojieba.Jieba
}

func isChinese(str string) bool {
	return regexp.MustCompile(`^[\p{Han}]+$`).MatchString(str)
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
		for _, word := range s.Terms[id] {
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
	s.CosineSimlarity(query, words)

	results := make([]int, 0)
	for id := range ids {
		results = append(results, id)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return s.CosineSimlarityMap[-1][results[i]] > s.CosineSimlarityMap[-1][results[j]]
	})

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

func (s *SearchEngine) TF_IDF(query string, queryWords []string) (map[int]map[string]float64, map[string]float64) {
	// 先定义DF函数
	df := func(word string) int {
		cnt := 0
		hash := make(map[int]bool)
		for _, tmpDocId := range s.Index[word] {
			if !hash[tmpDocId] {
				cnt++
				hash[tmpDocId] = true
			}
		}
		return cnt
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
	//////////////////////////////////////////////////////////////////////////////
	// 计算文档的IDF
	idf := make(map[string]float64)
	for word := range s.WordSet {
		idf[word] = math.Log10(float64(len(s.Docs)) / float64(df(word)))
	}
	// 计算文档的TF-IDF
	tf_idf := make(map[int]map[string]float64)
	for docId := range s.Docs {
		tf_idf[docId] = make(map[string]float64)
		for _, word := range s.Terms[docId] {
			tf_val := tf(word, docId)
			if tf_val == 0 {
				tf_idf[docId][word] = 0
			} else {
				abc := float64(tf_val) * idf[word]
				tf_idf[docId][word] = abc
			}
		}
	}
	//////////////////////////////////////////////////////////////////////////////
	// 计算Query的TF-IDF
	// query_idf := make(map[string]float64)
	query_tf := make(map[string]float64)
	for _, word := range queryWords {
		query_tf[word]++
	}
	for index := range query_tf {
		query_tf[index] /= float64(len(queryWords))
	}
	return tf_idf, query_tf
}

func (s *SearchEngine) CosineSimlarity(query string, queryWords []string) {
	s.CosineSimlarityMap = make(map[int]map[int]float64)
	s.CosineSimlarityMap[-1] = make(map[int]float64) // docId = -1 表示query
	tf_idf, query_tf := s.TF_IDF(query, queryWords)
	for docId := range s.Docs {
		v1 := tf_idf[docId]
		v2 := query_tf
		v1Norm := 0.0
		v2Norm := 0.0
		// 计算余弦相似度
		numerator := 0.0
		for word := range v2 {
			if _, ok := v1[word]; ok {
				numerator += v1[word] * v2[word]
				v1Norm += v1[word] * v1[word]
				v2Norm += v2[word] * v2[word]
			}
		}
		denominator := math.Sqrt(v1Norm) * math.Sqrt(v2Norm)
		if denominator == 0 {
			s.CosineSimlarityMap[-1][docId] = 0
		} else {
			s.CosineSimlarityMap[-1][docId] = numerator / denominator
		}
	}
}
