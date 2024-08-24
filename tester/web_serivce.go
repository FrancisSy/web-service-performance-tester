package tester

import (
	"FrancisSy/web-service-performance-tester/model"
	"log"
	"net/http"
	"time"
)

func InvokeWebService(url string) model.ResponseMetadata {
	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// body, err := io.ReadAll(res.Body)

	return model.ResponseMetadata{Code: res.StatusCode, Status: "SUCCESS", Url: url, Duration: time.Since(start).Seconds()}
}
