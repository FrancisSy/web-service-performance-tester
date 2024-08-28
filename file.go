package webtest

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	filename, filepath string
	settings           *FileSettings
}

func InitDefaultFileInfo(fpath string) *FileInfo {
	return &FileInfo{
		filename: "",
		filepath: fpath,
		settings: InitDefaultFileSettings(),
	}
}

func InitFileInfo(fpath string) *FileInfo {
	fname := strings.ReplaceAll("/results-"+time.Now().Format(time.RFC822)+".csv", " ", "-")
	os.Create(fpath + "/" + fname)
	return &FileInfo{
		filename: fname,
		filepath: fpath,
		settings: nil,
	}
}

func (f *FileInfo) WithSettings(callFlag, headerFlag, reqFlag, resFlag bool) *FileInfo {
	if f.settings == nil {
		f.settings = &FileSettings{
			writeCall:     callFlag,
			writeHeader:   headerFlag,
			writeRequest:  reqFlag,
			writeResponse: resFlag,
		}
	} else {
		f.settings.writeCall = callFlag
		f.settings.writeHeader = headerFlag
		f.settings.writeRequest = reqFlag
		f.settings.writeResponse = resFlag
	}

	return f
}

func Write(fpath string, b []byte) {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write(b); err != nil {
		log.Fatal(err)
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

func DumpResponse(res [][]byte, path string) {
	filepath := strings.ReplaceAll(path+"/results-"+time.Now().Format(time.RFC822)+".csv", " ", "-")
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
