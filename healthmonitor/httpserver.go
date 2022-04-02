package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartHttpServer() {
	router := gin.Default()
	router.LoadHTMLGlob("resources/*.tmp")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmp", nil)
	})

	router.GET("/results", func(c *gin.Context) {
		start, _ := strconv.Atoi(c.Query("start"))
		end, _ := strconv.Atoi(c.Query("end"))

		if start < 0 || end > 99 || start > end {
			c.HTML(http.StatusOK, "error.tmp", gin.H{
				"error": "Invalid Parameters",
			})
			return
		}

		c.HTML(http.StatusOK, "results.tmp", gin.H{
			"start": start,
			"end":   end,
		})
	})

	router.Run(":5000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
