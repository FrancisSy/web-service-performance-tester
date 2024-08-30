package webtest

import (
	"bytes"
	"encoding/csv"
	"log"
	"os"
)

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

func FormatStringArray(buf *bytes.Buffer, arr []string, delimiter string) *bytes.Buffer {
	buf.WriteString("[")
	for i := range arr {
		if i != len(arr)-1 {
			buf.WriteString("\"" + arr[i] + "\"" + delimiter + " ")
		} else {
			buf.WriteString("\"" + arr[i] + "\"")
		}
	}
	buf.WriteString("]")
	return buf
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

func WriteToFile(file *os.File, b []byte) {
	defer file.Close()
	if _, err := file.Write(b); err != nil {
		log.Fatal(err)
	}
}
