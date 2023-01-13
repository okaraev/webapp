package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type httpresponse struct {
	Status  bool
	Message string
}

var buckets = []float64{.05, .1, .25, .5, .9, .99, .999, 1, 2.5}
var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "webapp_processed_requests_total",
		Help: "The total number of requests",
	})
)

var (
	hist = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "webapp",
			Name:      "http_server_request_duration_seconds",
			Help:      "Histogram of response time for handler in seconds",
			Buckets:   buckets,
		},
		[]string{"route", "method", "status_code"},
	)
)

func GetSample(c *gin.Context) {
	rand := rand.Intn(100)
	sleepDuration := time.Duration(rand * int(time.Millisecond))
	time.Sleep(sleepDuration)
	opsProcessed.Inc()
	obs := float64(rand) / 1000
	hist.WithLabelValues("/get", "GET", "200").Observe(obs)
	c.IndentedJSON(200, httpresponse{Status: true, Message: "Success"})
}

func main() {

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9101", nil)
	}()
	prometheus.MustRegister(hist)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8086"
	}
	router := gin.Default()
	router.GET("/get", GetSample)
	router.Run(":" + port)
}
