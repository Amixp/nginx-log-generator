package main

import (
	"fmt"
	gofakeit "github.com/brianvoe/gofakeit"
	"github.com/caarlos0/env/v9"
	"github.com/dustinkirkland/golang-petname"
	go_uuid "github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

type config struct {
	Rate             float32 `env:"RATE" envDefault:"1"`
	IPv4Percent      int     `env:"IPV4_PERCENT" envDefault:"100"`
	StatusOkPercent  int     `env:"STATUS_OK_PERCENT" envDefault:"80"`
	PathMinLength    int     `env:"PATH_MIN" envDefault:"1"`
	PathMaxLength    int     `env:"PATH_MAX" envDefault:"5"`
	PercentageGet    int     `env:"GET_PERCENT" envDefault:"60"`
	PercentagePost   int     `env:"POST_PERCENT" envDefault:"30"`
	PercentagePut    int     `env:"PUT_PERCENT" envDefault:"0"`
	PercentagePatch  int     `env:"PATCH_PERCENT" envDefault:"0"`
	PercentageDelete int     `env:"DELETE_PERCENT" envDefault:"0"`
}

func initRand() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func randFloat(min, max float32) float32 {
	res := min + float32(rand.Int31n(int32(max))) + float32(rand.Int31n(99)+1)/100

	return res
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func main() {
	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)
	checkMinMax(&cfg.PathMinLength, &cfg.PathMaxLength)

	ticker := time.NewTicker(time.Second / time.Duration(cfg.Rate))

	gofakeit.Seed(time.Now().UnixNano())

	var ip, httpMethod, path, referrer, userAgent string
	var statusCode, bodyBytesSent, upstream_status int
	var timeLocal time.Time

	//httpVersion = "HTTP/1.1"
	referrer = "-"

	initRand()

	for range ticker.C {
		timeLocal = time.Now()

		ip = weightedIPVersion(cfg.IPv4Percent)
		httpMethod = weightedHTTPMethod(cfg.PercentageGet, cfg.PercentagePost, cfg.PercentagePut, cfg.PercentagePatch, cfg.PercentageDelete)
		path = randomPath(cfg.PathMinLength, cfg.PathMaxLength)
		statusCode = weightedStatusCode(cfg.StatusOkPercent)
		bodyBytesSent = realisticBytesSent(statusCode)
		userAgent = gofakeit.UserAgent()
		uuid := go_uuid.NewV4().String()

		httpHost := petname.Generate(3, ".")
		serverName := petname.Generate(1, ".")
		requestTime := randFloat(0.01, 301.98)
		upstreamConnectTime := randFloat(0.01, 301.98)
		upstreamHeaderTime := randFloat(0.01, 301.98)
		upstreamResponseTime := randFloat(0.01, 301.98)
		pid := rand.Intn(1000)
		httpReferer := referrer
		upstreamCacheStatus := "-"
		upstreamAddr := "-"
		requestUri := path
		proxiedUri := path
		serverProtocol := "HTTP/1.1"
		requestLength := bodyBytesSent
		httpXRequestedWith := "-"
		scheme := "http"

		fmt.Printf("\"%s\" \"%s\" \"-\" \"%s\" \"%s\" \"80\" \"11\" \"22\" \"33\" \"44\" \"55\" \"66\" \"77\" \"%s\" \"%v\" \"%v\" \"%v\" \"%v\" \"%v\" \"%v\" \"%v\" \"%s\" \"%s\" \"%s\" \"%s\" \"%s\" \"%s\" \"%v\" \"%v\" \"%s\" \"%s\" \"%v\" \"%v\" \"%v\"\n", uuid, ip, httpHost, serverName, timeLocal.Format("02/Jan/2006:15:04:05 -0700"), statusCode, pid, requestTime, upstream_status, upstreamConnectTime, upstreamHeaderTime, upstreamResponseTime, upstreamCacheStatus, upstreamAddr, httpMethod, requestUri, proxiedUri, serverProtocol, requestLength, bodyBytesSent, httpReferer, userAgent, httpXRequestedWith, scheme, bodyBytesSent)
	}
}

func realisticBytesSent(statusCode int) int {
	if statusCode != 200 {
		return gofakeit.Number(30, 120)
	}

	return gofakeit.Number(800, 3100)
}

func weightedStatusCode(percentageOk int) int {
	roll := gofakeit.Number(0, 100)
	if roll <= percentageOk {
		return 200
	}

	return gofakeit.SimpleStatusCode()
}

func weightedHTTPMethod(percentageGet, percentagePost, percentagePut, percentagePatch, percentageDelete int) string {
	if percentageGet+percentagePost >= 100 {
		panic("HTTP method percentages add up to more than 100%")
	}

	roll := gofakeit.Number(0, 100)
	if roll <= percentageGet {
		return "GET"
	} else if roll <= percentagePost {
		return "POST"
	} else if roll <= percentagePut {
		return "PUT"
	} else if roll <= percentagePatch {
		return "PATCH"
	} else if roll <= percentageDelete {
		return "DELETE"
	}

	return gofakeit.HTTPMethod()
}

func weightedIPVersion(percentageIPv4 int) string {
	roll := gofakeit.Number(0, 100)
	if roll <= percentageIPv4 {
		return gofakeit.IPv4Address()
	} else {
		return gofakeit.IPv6Address()
	}
}

func randomPath(min, max int) string {
	var path strings.Builder
	length := gofakeit.Number(min, max)

	path.WriteString("/")

	for i := 0; i < length; i++ {
		if i > 0 {
			path.WriteString(gofakeit.RandString([]string{"-", "-", "_", "%20", "/", "/", "/"}))
		}
		path.WriteString(gofakeit.BuzzWord())
	}

	path.WriteString(gofakeit.RandString([]string{".hmtl", ".php", ".htm", ".jpg", ".png", ".gif", ".svg", ".css", ".js"}))

	result := path.String()
	return strings.Replace(result, " ", "%20", -1)
}

func checkMinMax(min, max *int) {
	if *min < 1 {
		*min = 1
	}
	if *max < 1 {
		*max = 1
	}
	if *min > *max {
		*min, *max = *max, *min
	}
}
