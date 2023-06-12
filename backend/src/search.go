package src

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"github.com/yanyiwu/gojieba"
)

// 倒排索引
type InvertedIndex map[string][]int

// 搜索引擎
type SearchEngine struct {
	PostingList         InvertedIndex           // 倒排索引
	Docs                []string                // 文档集合
	WordSet             map[string]bool         // 词集（用map构造set，比较简单的写法）
	Terms               map[int][]string        // 每个文档ID的jieba分词的结果
	CosineSimilarityMap map[int]map[int]float64 // 两个文档之间的余弦相似度
	jieba               *gojieba.Jieba
}

// 判断是否为中文字符
func isChinese(str string) bool {
	return regexp.MustCompile(`^[\p{Han}]+$`).MatchString(str)
}

// 初始化搜索引擎
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

func (s *SearchEngine) intersect(a []int, b []int) []int {
	ret := make([]int, 0)
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			ret = append(ret, a[i])
			i++
			j++
		} else if a[i] > b[j] {
			j++
		} else {
			i++
		}
	}
	return ret
}

// 求并集，已经实现，不需要这个函数了
// func (s *SearchEngine) union(a []int, b []int) []int {
// 	hash := make([]int, 10000)
// 	ret := make([]int, 0)
// 	for _, v := range a {
// 		if hash[v] == 0 {
// 			ret = append(ret, v)
// 		}
// 		hash[v]++
// 	}
// 	for _, v := range b {
// 		if hash[v] == 0 {
// 			ret = append(ret, v)
// 		}
// 		hash[v]++
// 	}
// 	sort.Ints(ret)
// 	return ret
// }

func (s *SearchEngine) BuildInvertedIndex() {
	s.PostingList = make(InvertedIndex)
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
			s.PostingList[word] = append(s.PostingList[word], id)
			s.WordSet[word] = true
		}
	}
}

// TODO: 完善这个函数
func (s *SearchEngine) Search(query string) []int {
	ids := make(map[int]bool)
	afterIntersect := make([]int, 0)
	words := s.jieba.CutForSearch(query, true)
	for i, word := range words {
		if docIds, ok := s.PostingList[word]; ok {
			if i == 0 {
				afterIntersect = docIds
			} else {
				afterIntersect = s.intersect(afterIntersect, docIds)
			}
		}
	}
	for _, id := range afterIntersect {
		ids[id] = true
	}
	s.CosineSimlarity(query, words)

	results := make([]int, 0)
	for id := range ids {
		results = append(results, id)
	}

	TopKHeap := InitKHeap(10)
	for _, id := range results {
		TopKHeap.Push(Pair{id, s.CosineSimilarityMap[-1][id]})
	}

	ret := []int{}
	for _, val := range TopKHeap.data {
		ret = append(ret, val.first)
	}

	return ret
}

func (s *SearchEngine) ReadFile(pagesDir string) {
	s.Docs = make([]string, 0)
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
		s.Docs = append(s.Docs, string(content))

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
		for _, tmpDocId := range s.PostingList[word] {
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
		for _, tmpDocId := range s.PostingList[word] {
			if tmpDocId == docId {
				cnt++
			}
		}
		return cnt
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
	s.CosineSimilarityMap = make(map[int]map[int]float64)
	s.CosineSimilarityMap[-1] = make(map[int]float64) // docId = -1 表示query
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
			s.CosineSimilarityMap[-1][docId] = 0
		} else {
			s.CosineSimilarityMap[-1][docId] = numerator / denominator
		}
	}
}
