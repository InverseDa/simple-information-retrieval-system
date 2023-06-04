package main

import (
	"fmt"
	"information/bootstrap"
	"information/global"
	"information/search"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	se := search.SearchEngine{}
	se.ReadFile("/home/miaokeda/Information/backend/pages")
	se.BuildInvertedIndex()
	results := se.Search("核酸")
	fmt.Println("Search results for query", "粤海校区")
	for _, id := range results {
		fmt.Println("Document ID:", id, "Document content:", se.Docs[id])
	}
	bootstrap.InitializeConfig()

	r := gin.Default()

	// 测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 启动服务器
	r.Run(":" + global.App.Config.App.Port)
}
