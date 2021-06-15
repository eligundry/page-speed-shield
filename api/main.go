package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	servertiming "github.com/p768lwy3/gin-server-timing"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/api/pagespeedonline/v5"
)

var (
	pageSpeedSvc *pagespeedonline.Service
	logger       *zap.Logger
	store        persistence.CacheStore
)

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

	return fmt.Sprintf("Page Speed Insights-%d-%s", score, color)
}

func runPageSpeed(ctx context.Context, url string, insightType string) (int64, error) {
	res, err := pageSpeedSvc.Pagespeedapi.Runpagespeed(url).
		Strategy(strings.ToUpper(insightType)).
		Category("PERFORMANCE").
		Context(ctx).
		Do()

	if err != nil {
		return 0, err
	}

	score := int64(res.LighthouseResult.Categories.Performance.Score.(float64) * 100)

	return score, nil
}

func handler(c *gin.Context) {
	timing := servertiming.FromContext(c)
	targetURL := c.Param("url")[1:]
	insightType := c.Param("insightType")

	if _, ok := validInsightTypes[insightType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "insightType must be either mobile or desktop",
		})
		return
	}

	pagespeedTiming := timing.NewMetric("pageSpeed").
		WithDesc("Google Page Speed Insights API").
		Start()

	score, err := runPageSpeed(
		c.Request.Context(),
		targetURL,
		insightType,
	)

	pagespeedTiming.Stop()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = "img.shields.io"
			req.Host = "img.shields.io"
			req.URL.Path = "/badge/" + constructShield(score)
		},
		ModifyResponse: func(resp *http.Response) error {
			resp.Header.Set("Cache-Control", "public, max-age=86400")
			return nil
		},
	}

	servertiming.WriteHeader(c)

	proxy.ServeHTTP(c.Writer, c.Request)
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

	logger, err = zap.NewProduction()

	if err != nil {
		log.Panicf("could not construct logger, %s", err)
	}

	store = persistence.NewInMemoryStore(time.Second)
}

func main() {
	router := gin.New()

	// middlewares
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(servertiming.Middleware())

	router.GET("/:insightType/*url", cache.CachePage(store, time.Hour*24, handler))
	router.Run()
}
