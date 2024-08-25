package main

import (
	tester "FrancisSy/web-service-performance-tester/webtest"
	"flag"
	"fmt"
	"log"
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

	/*
	* initialize web client with headers
	*
	* can also initialize custom client
	* transport := tester.InitTransport()
	* client := &http.Client{Transport: transport}
	 */
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	client := tester.InitWebClient().Headers(&headers)

	// GET request with the query parameters from the csv file
	for _, c := range contents {
		res, err := client.Get(*apiUrl + c[0])
		if err != nil {
		}

		defer res.Body.Close()
	}

	// exmaple GET request with path parameters
	// request with query parameters would be similar
	arr := []interface{}{contents[0][0], contents[1][0]}
	client.Get(fmt.Sprintf("http://pokeapi.co/api/v2/pokemon/%s/test/%s", arr...))

	// example POST request
	for _, c := range contents {
		res, err := client.Post(*apiUrl, []byte(fmt.Sprintf(`
		{
			"name": "%s",
			"region": "%s"
		}
		`, c[0], c[1])))

		if err != nil {
		}

		defer res.Body.Close()
	}

	// still working on this
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
