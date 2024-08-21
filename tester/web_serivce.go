package tester

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func CreateJsonRequestBody(s struct{}) string {
	json, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	return string(json)
}

func InvokeWebService(url string) ([]byte, float64) {
	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data, time.Since(start).Seconds()
}
