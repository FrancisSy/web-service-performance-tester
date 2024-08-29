package webtest

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
* general idea is that creating a new FileHandler will
* create a results folder in the current directory
* where all results files will be created in
 */
type FileHandler struct {
	filename, filepath string
	settings           *FileSettings
}

func InitDefaultFileHandler(fpath string) *FileHandler {
	fname := strings.ReplaceAll("/results-"+time.Now().Format(time.RFC822)+".csv", " ", "-")
	file, err := os.Create(fpath + "/" + fname)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	return &FileHandler{
		filename: fname,
		filepath: fpath,
		settings: InitDefaultFileSettings(),
	}
}

func InitFileHandler(fpath string) *FileHandler {
	fname := strings.ReplaceAll("/results-"+time.Now().Format(time.RFC822)+".csv", " ", "-")
	file, err := os.Create(fpath + "/" + fname)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	return &FileHandler{
		filename: fname,
		filepath: fpath,
		settings: nil,
	}
}

func (f *FileHandler) WithSettings(fs *FileSettings) *FileHandler {
	f.settings = fs
	return f
}

// will refactor
func (f *FileHandler) Dump(req any, res *http.Response) {
	file, err := os.Open(f.filepath + "/" + f.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var buf bytes.Buffer
	if f.settings.writeRequest && req != nil {
		buf.Write([]byte("Request\n"))
		switch req.(type) {
		case *http.Request:
			cast, ok := req.(*http.Request)
			if !ok {
				log.Fatal("Error while trying to cast request to http.Request")
			}

			if f.settings.writeHeader {
				buf.Write([]byte("Request Headers\n"))
				// WriteHeaders(file, &cast.Header)
				for k, v := range cast.Header {
					buf.Write([]byte(fmt.Sprintf("%s : %s", k, v)))
				}
				body, err := io.ReadAll(cast.Body)
				if err != nil {
					log.Fatal("Error while trying read request body")
				}

				buf.Write([]byte("Request Body\n"))
				if err := json.Indent(&buf, body, "", "\t"); err != nil {
					log.Fatal(err)
				}
				buf.Write([]byte("\n"))
				if _, err := file.Write(append(append([]byte("Request Body\n"), body...), []byte("\n")...)); err != nil {
					log.Fatal(err)
				}
			} else {
				file.Write([]byte("No request headers were provided\n"))
			}
		case string:
			cast, ok := req.(string)
			if !ok {
				log.Fatal("Error while trying to cast request to string")
			}

			if _, err := file.Write(append([]byte("Request Parameter\n"), []byte(cast+"\n")...)); err != nil {
				log.Fatal(err)
			}
		case []string:
			cast, ok := req.([]string)
			if !ok {
				log.Fatal("Error while trying to cast request to []string")
			}

			c1 := ""
			for _, c := range cast {
				c1 += c + ","
			}

			if _, err := file.Write(append([]byte("Request Parameters\n"), []byte("["+c1[:len(c1)-1]+"]")...)); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		file.Write([]byte("No http request was provided\n"))
	}

	if f.settings.writeResponse && res != nil {
		file.Write([]byte("Response\n"))

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Error while trying to read response body contents")
		}

		var pjson bytes.Buffer
		if err := json.Indent(&pjson, body, "", "\t"); err != nil {
			log.Fatal(err)
		}

		file.Write(append(pjson.Bytes(), []byte("\n")...))
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

func Write(file *os.File, b []byte) {
	defer file.Close()

	if _, err := file.Write(b); err != nil {
		log.Fatal(err)
	}
}

func WriteToOpenFile(file *os.File, b []byte) {
	if _, err := file.Write(b); err != nil {
		log.Fatal(err)
	}
}

func WriteTo(fpath string, b []byte) {
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	if _, err := file.Write(b); err != nil {
		log.Fatal(err)
	}
}

func WriteHeaders(file *os.File, headers *http.Header) {
	for k, v := range *headers {
		if k != "Authorization" { // filter out auth headers
			v1 := ""
			for _, v2 := range v {
				v1 += v2 + ","
			}

			file.Write([]byte(string(k + " : [" + v1[:len(v1)-1] + "]\n")))
		}
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
