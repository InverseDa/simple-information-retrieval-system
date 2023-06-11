package main

import (
	"fmt"
	"information/bootstrap"
	"information/global"
	"information/search"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置，使用相对位置
	se := search.InitializeSearchEngine("/pages")
	se.Search("深圳大学粤海校区")

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
		fmt.Println("[Log] var query: ", data.Query)

		results := se.Search(data.Query)
		ret := []map[string]interface{}{}
		for _, id := range results {
			title := search.FindArticleDetails(se.Docs[id])
			ret = append(ret, map[string]interface{}{
				"content": se.Docs[id],
				"title":   title,
			})
		}

		fmt.Println("[Log] var data.Strings length: ", len(ret))

		// 返回响应，将字符串数组编码为JSON格式
		c.JSON(http.StatusOK, gin.H{"pagesString": ret})

	})

	// 启动服务器
	r.Run(":" + global.App.Config.App.Port)
}
