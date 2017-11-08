package main

import (
	"net/http"
	"time"

	"github.com/andywhittle/hastrumpsaidsomethingstupid/search"
	"github.com/gin-gonic/gin"
)

const timeout = 10 * time.Second

// Router builds the gin engine routing for the app
func Router(client search.Gettable) *gin.Engine {
	router := gin.Default()
	router.Static("/images", "./images")
	router.StaticFile("/yup", "templates/yup.html")
	router.StaticFile("/na", "templates/na.html")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		s := search.BBCNews{Client: client, Keyword: "trump"}
		hl, err := s.Headlines()
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.HTML(
				http.StatusOK,
				"index.tmpl",
				struct{ Headlines []string }{hl},
			)
		}
	})

	return router
}

func main() {
	c := http.Client{Timeout: timeout}
	Router(&c).Run(":8898")
}
