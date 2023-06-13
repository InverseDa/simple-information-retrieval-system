package main

import (
	"information/bootstrap"
	"information/global"
	"information/src"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置，使用相对位置
	se := src.InitializeSearchEngine("/pages")
	se.Search("宣传")
	bootstrap.InitializeConfig()

	r := gin.Default()
	r.Use(cors.Default())

	// 测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// query的接口
	r.POST("/api/query", func(c *gin.Context) {

		var data struct {
			Query string `form:"query" json:"query" binding:"required"`
		}
		// 解析请求体中的JSON数据
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println("var query: ", data.Query)

		start := time.Now()
		results := se.Search(data.Query)
		end := time.Now()
		// 如果搜不到数据，进行编辑距离计算
		if len(results) == 0 {
			fuzzySearchResults := se.FuzzySearch(data.Query)
			log.Println("var fuzzySearchResults: ", fuzzySearchResults)
			c.JSON(http.StatusOK, gin.H{"status": "error", "fuzzySearchString": fuzzySearchResults})
		} else {
			ret := []map[string]interface{}{}
			// 处理结果
			for _, id := range results {
				page := se.Docs[id]
				url, title, content := src.DealDocs(page)
				ret = append(ret, map[string]interface{}{
					"url":     url,
					"title":   title,
					"content": content,
					"time":    end.Sub(start).Seconds(),
				})
			}
			log.Println("var data.Strings length: ", len(ret))

			// 返回响应，将字符串数组编码为JSON格式
			c.JSON(http.StatusOK, gin.H{"status": "success", "pagesString": ret})
		}

	})

	// 启动服务器
	r.Run(":" + global.App.Config.App.Port)
}
