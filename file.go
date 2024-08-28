package webtest

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
* general idea is that creating a new FileInfo will
* create a results folder in the current directory
* where all results files will be created in
 */
type FileInfo struct {
	filename, filepath string
	settings           *FileSettings
}

func InitDefaultFileInfo(fpath string) *FileInfo {
	fname := strings.ReplaceAll("/results-"+time.Now().Format(time.RFC822)+".csv", " ", "-")
	file, err := os.Create(fpath + "/" + fname)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	return &FileInfo{
		filename: fname,
		filepath: fpath,
		settings: InitDefaultFileSettings(),
	}
}

func InitFileInfo(fpath string) *FileInfo {
	fname := strings.ReplaceAll("/results-"+time.Now().Format(time.RFC822)+".csv", " ", "-")
	file, err := os.Create(fpath + "/" + fname)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	return &FileInfo{
		filename: fname,
		filepath: fpath,
		settings: nil,
	}
}

func (f *FileInfo) WithSettings(fs *FileSettings) *FileInfo {
	f.settings = fs
	return f
}

// will refactor
func (f *FileInfo) Dump(req any, res *http.Response) {
	file, err := os.Open(f.filepath + "/" + f.filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	if f.settings.writeRequest && req != nil {
		file.Write([]byte("Request\n"))
		switch req.(type) {
		case *http.Request:
			cast, ok := req.(*http.Request)
			if !ok {
				log.Fatal("Error while trying to cast request to http.Request")
			}

			if f.settings.writeHeader {
				file.Write([]byte("Request Headers\n"))
				for k, v := range cast.Header {
					if k == "Authorization" { // filter out auth headers
						v1 := ""
						for _, v2 := range v {
							v1 += v2 + ","
						}

						file.Write([]byte(string(k + " : [" + v1[:len(v1)-1] + "]\n")))
					}
				}

				body, err := io.ReadAll(cast.Body)
				if err != nil {
					log.Fatal("Error while trying read request body")
				}

				file.Write(append(append([]byte("Request Body\n"), body...), []byte("\n")...))
			} else {
				file.Write([]byte("No request headers were provided\n"))
			}
		case string:
			cast, ok := req.(string)
			if !ok {
				log.Fatal("Error while trying to cast request to string")
			}

			file.Write(append([]byte("Request Parameter\n"), []byte(cast+"\n")...))
		case []string:
			cast, ok := req.([]string)
			if !ok {
				log.Fatal("Error while trying to cast request to []string")
			}

			c1 := ""
			for _, c := range cast {
				c1 += c + ","
			}

			file.Write(append([]byte("Request Parameters\n"), []byte("["+c1[:len(c1)-1]+"]")...))
		}
	} else {
		file.Write([]byte("No http request was provided\n"))
	}
}

type FileSettings struct {
	writeCall, writeHeader, writeRequest, writeResponse bool
}

func InitDefaultFileSettings() *FileSettings {
	return &FileSettings{
		writeCall:     true,
		writeHeader:   false,
		writeRequest:  false,
		writeResponse: false,
	}
}

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

func Write(fpath string, b []byte) {
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	if _, err := file.Write(b); err != nil {
		log.Fatal(err)
	}
}

func DumpResponse(fpath string, res [][]byte) {
	filepath := strings.ReplaceAll(fpath+"/results-"+time.Now().Format(time.RFC822)+".csv", " ", "-")
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
