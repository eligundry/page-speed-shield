package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/pagespeedonline/v5"
)

var pageSpeedSvc *pagespeedonline.Service

var validInsightTypes = map[string]bool{
	"mobile":  true,
	"desktop": true,
}

func constructShield(score int64) string {
	color := "red"

	switch {
	case score >= 90:
		color = "green"
	case score >= 80 && score < 90:
		color = "orange"
	case score >= 70 && score < 80:
		color = "yellow"
	}

	return fmt.Sprintf("https://img.shields.io/badge/Page%%20Speed%%20Insights-%d-%s", score, color)
}

func runPageSpeed(ctx context.Context, url string, insightType string) (int64, error) {
	res, err := pageSpeedSvc.Pagespeedapi.Runpagespeed(url).
		Strategy(strings.ToUpper(insightType)).
		Category("PERFORMANCE").
		Do()

	if err != nil {
		return 0, err
	}

	score := int64(res.LighthouseResult.Categories.Performance.Score.(float64) * 100)

	return score, nil
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

	score, err := runPageSpeed(
		c.Request.Context(),
		url,
		insightType,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"insightType": insightType,
		"url":         url,
		"pagespeed":   score,
		"shield":      constructShield(score),
	})
}

func init() {
	var err error

	pageSpeedSvc, err = pagespeedonline.NewService(
		context.Background(),
		option.WithAPIKey(os.Getenv("GOOGLE_API_KEY")),
	)

	if err != nil {
		log.Panicf("could not construct page speed service, %s", err)
	}
}

func main() {
	router := gin.Default()
	router.GET("/:insightType/*url", handler)
	router.Run()
}
