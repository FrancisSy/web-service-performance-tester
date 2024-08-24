package main

import (
	"FrancisSy/web-service-performance-tester/tester"
	"flag"
	"log"
	"net/url"
	"os"
)

/*
 * figure out how to design webtest
 */
func main() {
	// check flags
	fpath := flag.String("filepath", "", "fpath of csv")
	apiUrl := flag.String("url", "", "web service url")
	dumpToFile := flag.Bool("dump", false, "write response data to file")
	dumpfpath := flag.String("dump-filepath", "", "location of dump file")
	flag.Parse()

	if len(*fpath) == 0 {
		log.Fatal("Missing fpath")
	}

	if len(*apiUrl) == 0 {
		log.Fatal("Missing API URL")
	}

	// read the csv file
	contents := tester.ReadCsv(*fpath)

	// initialize web client
	// transport := tester.InitTransport()
	// client := &http.Client{Transport: transport}
	client := tester.InitWebClient()

	params := url.Values{}
	params.Add("hello", "hello")
	arr := []interface{}{contents[0][0], contents[1][0]}
	client.GetWithPathParams("http://pokeapi.co/api/v2/pokemon/%s/test/%s", arr)

	// call the api url with the query parameters from the csv file
	for _, c := range contents {
		res, err := client.GetWithPathParam(*apiUrl, c[0])
		if err != nil {
		}

		defer res.Body.Close()
	}

	// check to see if the file needs to be dumped to a path
	if *dumpToFile {
		if len(*dumpfpath) == 0 {
			log.Println("Dumping result file in current directory")
			_, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			// tester.DumpResponse(metadata, dir)
		} else {
			log.Printf("Dumping result file in %s\n", *dumpfpath)
			// tester.DumpResponse(metadata, *dumpfpath)
		}
	}
}
