package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartHttpServer() {
	router := gin.Default()
	router.LoadHTMLGlob("resources/*.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/", func(c *gin.Context) {
		startDateString := c.PostForm("start")
		endDateString := c.PostForm("end")

		fmt.Println(startDateString, endDateString)
	})

	router.Run(":5000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
