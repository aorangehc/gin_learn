package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/hello_gin", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	shopGroup := r.Group("/shop", func1, func2)
	shopGroup.Use(func3)
	{
		shopGroup.GET("/index", func4, func5)
	}

	r.Run(":9999")
}
