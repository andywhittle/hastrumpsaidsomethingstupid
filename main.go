package main

import (
	"net/http"

	"github.com/andywhittle/hastrumpsaidsomethingstupid/search"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/images", "./images")
	router.StaticFile("/yup", "templates/yup.html")
	router.StaticFile("/na", "templates/na.html")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		s := search.BBCNews{Keyword: "trump"}
		c.HTML(
			http.StatusOK,
			"index.tmpl",
			struct {
				Headlines []string
			}{
				s.Headlines(),
			})
	})

	router.Run(":8898")
}
