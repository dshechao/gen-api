package main

import (
	"net/http"

	"github.com/dshechao/gen-api/gen"
	yaag_gin "github.com/dshechao/gen-api/gin/v1"
)

func main() {
	gen.Init(&gen.Config{On: true, DocTitle: "Gin", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
	r := gin.Default()
	r.Use(yaag_gin.Document())
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "Hello World!"})
	})
	r.GET("/plain", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"result": "Hello World!"})
	})
	r.GET("/complex", func(c *gin.Context) {
		value := c.Query("key")
		c.JSON(http.StatusOK, gin.H{"value": value})
	})
	r.Run(":8080")
}
