package tester

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"
)

func ReadCsv(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	res, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func DumpResponse(res [][]byte, path string) {
	filepath := strings.ReplaceAll(path+"/results-"+time.Now().Format(time.RFC3339)+".csv", " ", "-")
	f, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		if _, err := f.Write(r); err != nil {
			panic(err)
		}
	}
}
