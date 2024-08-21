package tester

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func ReadCsv(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var result [][]string
	for _, r := range records {
		result = append(result, r)
	}

	return result
}

func DumpResponse(res []byte, path string) {
	filepath := strings.ReplaceAll(path+"/results-"+time.DateTime, " ", "-")
	f, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	w, err := f.Write(res)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote %d bytes\n", w)
}
