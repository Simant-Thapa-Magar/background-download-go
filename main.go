package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type FileInfo struct {
	FileName     string `json:"file_name"`
	DownloadLink string `json:"download_link"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GetFileInfo() []FileInfo {
	var files []FileInfo

	fileList, err := ioutil.ReadDir("./files")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileList {
		fileName := file.Name()
		files = append(files, FileInfo{
			FileName:     strings.Split(fileName, ".")[0],
			DownloadLink: fmt.Sprintf("/pdf?file=%s", fileName),
		})
	}

	return files
}

func Info(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	files := GetFileInfo()
	w.Header().Set("Content-type", "application/json")

	marshaled, err := json.Marshal(files)

	if err != nil {
		log.Fatal(err)
	}

	w.Write(marshaled)
}

func PDF(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fileQuery := r.URL.Query().Get("file")

	FILEPATH := fmt.Sprintf("./files/%s", fileQuery)

	fileStat, err := os.Stat(FILEPATH)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// get the size
	size := fileStat.Size()

	// grab the generated receipt.pdf file and stream it to browser
	streamPDFbytes, err := ioutil.ReadFile(FILEPATH)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := bytes.NewBuffer(streamPDFbytes)

	// stream straight to client(browser)
	w.Header().Set("Content-type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(int(size)))

	if _, err := b.WriteTo(w); err != nil { // <----- here!
		fmt.Fprintf(w, "%s", err)
	}

	w.Write([]byte("PDF Generated"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Info)
	mux.HandleFunc("/pdf", PDF)

	http.ListenAndServe(":8080", mux)
}
