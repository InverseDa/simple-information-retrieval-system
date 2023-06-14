package src

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/kljensen/snowball"
	"github.com/yanyiwu/gojieba"
)

// 倒排索引
type InvertedIndex map[string][]int

// 搜索引擎
type SearchEngine struct {
	PostingList         InvertedIndex           // 倒排索引
	Docs                []string                // 文档集合
	WordSet             map[string]bool         // 词集（用map构造set，比较简单的写法）
	Terms               map[int][]string        // 每个文档ID的分词的结果
	CosineSimilarityMap map[int]map[int]float64 // 两个文档之间的余弦相似度
	tf_idf              map[int]map[string]float64
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
	se.TF_IDF_ForDocs()
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func CalculateEditDistance(s1 string, s2 string) int {
	dp := make([][]int, len(s1)+1)
	for i := 0; i <= len(s1); i++ {
		dp[i] = make([]int, len(s2)+1)
	}
	for i := 1; i <= len(s1); i++ {
		dp[i][0] = i
	}
	for j := 1; j <= len(s2); j++ {
		dp[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min(dp[i-1][j], min(dp[i][j-1], dp[i-1][j-1])) + 1
			}
		}
	}

	return dp[len(s1)][len(s2)]
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
		englishWords := regexp.MustCompile(`\b\w+\b`).FindAllString(doc, -1)
		for _, word := range words {
			if isChinese(word) {
				s.Terms[id] = append(s.Terms[id], word)
				s.PostingList[word] = append(s.PostingList[word], id)
				s.WordSet[word] = true
			}
		}
		for _, word := range englishWords {
			word, _ = snowball.Stem(word, "english", true)
			s.Terms[id] = append(s.Terms[id], word)
			s.PostingList[word] = append(s.PostingList[word], id)
			s.WordSet[word] = true
		}
	}
}

// TODO: 完善这个函数
func (s *SearchEngine) Search(query string) []int {
	ids := make(map[int]bool)
	afterIntersect := make([]int, 0)
	chineseWords := s.jieba.CutForSearch(query, true)
	englishWords := regexp.MustCompile(`\b\w+\b`).FindAllString(query, -1)
	words := make([]string, 0)
	for _, word := range chineseWords {
		if isChinese(word) {
			words = append(words, word)
		}
	}
	for _, word := range englishWords {
		word, _ = snowball.Stem(word, "english", true)
		words = append(words, word)
	}
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

// 返回1. 文档的TF-IDF向量 2. Query的TF-IDF向量
func (s *SearchEngine) TF_IDF_ForDocs() map[int]map[string]float64 {
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
	// ==============================================================================================================
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
	// ==============================================================================================================

	return tf_idf
}

func (s *SearchEngine) TF_IDF_ForQuery(query string, queryWords []string) map[string]float64 {
	// 计算Query的TF-IDF
	// query_idf := make(map[string]float64)
	query_tf := make(map[string]float64)
	for _, word := range queryWords {
		query_tf[word]++
	}
	for index := range query_tf {
		query_tf[index] /= float64(len(queryWords))
	}
	return query_tf
}

func (s *SearchEngine) CosineSimlarity(query string, queryWords []string) {
	s.CosineSimilarityMap = make(map[int]map[int]float64)
	s.CosineSimilarityMap[-1] = make(map[int]float64) // docId = -1 表示query
	query_tf := s.TF_IDF_ForQuery(query, queryWords)
	for docId := range s.Docs {
		v1 := s.tf_idf[docId]
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

func (s *SearchEngine) FuzzySearch(query string) []string {
	ret := make([]string, 0)
	for word := range s.WordSet {
		if CalculateEditDistance(query, word) <= 2 {
			ret = append(ret, word)
		}
	}
	return ret
}

func DealDocs(docs string) (string, string, string) {
	_url := regexp.MustCompile(`\[url\]:\s+(.*)`).FindStringSubmatch(docs)[1]
	parts := strings.Split(docs, "\n")
	_title := ""
	_content := ""
	isTitle := true
	for _, str := range parts {
		if str == "" || regexp.MustCompile(`^\s*$`).MatchString(str) || regexp.MustCompile(`\[url\]:\s+(.*)`).MatchString(str) {
			continue
		}
		if isTitle {
			_title = str
			isTitle = false
		} else {
			_content += str + "\n"
		}

	}
	return _url, _title, _content
}
