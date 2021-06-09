package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/pagespeedonline/v5"
)

var validInsightTypes = map[string]bool{
	"mobile":  true,
	"desktop": true,
}

func handler(c *gin.Context) {
	url := c.Param("url")[1:]
	insightType := c.Param("insightType")

	if _, ok := validInsightTypes[insightType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "insightType must be either mobile or desktop",
		})
		return
	}

	pageSpeedSvc, err := pagespeedonline.NewService(
		c.Request.Context(),
		option.WithAPIKey(os.Getenv("GOOGLE_API_KEY")),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := pageSpeedSvc.Pagespeedapi.Runpagespeed(url).Do()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"insightType": insightType,
		"url":         url,
		"pagespeed":   res.LighthouseResult.Categories.Performance.Score,
	})
	// c.Redirect(http.StatusTemporaryRedirect, "https://img.shields.io/badge/Page%20Speed%20Insights-99-green")
}

func main() {
	router := gin.Default()
	router.GET("/:insightType/*url", handler)
	router.Run()
}
