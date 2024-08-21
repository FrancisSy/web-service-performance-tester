package main

import (
	"FrancisSy/web-service-performance-tester/model"
	"FrancisSy/web-service-performance-tester/tester"
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

func main() {
	filePath := flag.String("filepath", "", "filepath of csv")
	apiUrl := flag.String("url", "", "web service url")
	dumpToFile := flag.Bool("dump", false, "write response data to file")
	dumpFilePath := flag.String("dump-filepath", "", "location of dump file")
	flag.Parse()

	if len(*filePath) == 0 {
		log.Fatal("Missing filepath")
	}

	if len(*apiUrl) == 0 {
		log.Fatal("Missing API URL")
	}

	contents := tester.ReadCsv(*filePath)
	for _, c := range contents {
		for _, s := range c {
			fmt.Println(s)
		}
	}

	entries := model.PokedexEntriesResponse{}
	data, responseTime := tester.InvokeWebService(*apiUrl)
	fmt.Printf("%.3fs\n", responseTime)

	if err := json.Unmarshal(data, &entries); err != nil {
		log.Fatal(err)
	}

	for _, e := range entries.PokemonEntries {
		fmt.Println(e)
	}

	if *dumpToFile {
		if len(*dumpFilePath) == 0 {
			log.Println("Dumping result file in current directory")
		} else {
			log.Printf("Dumping result file in %s\n", *dumpFilePath)
			tester.DumpResponse(data, *dumpFilePath)
		}
	}
}
