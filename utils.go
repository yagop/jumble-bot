package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func downloadFile(url string, download_path string) (int64, error) {

	out, err := os.Create(download_path)
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	log.Println("Written bytes", n)
	return n, nil
}
