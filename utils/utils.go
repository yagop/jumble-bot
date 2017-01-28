package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(url string, download_path string) (int64, error) {
	log.Printf("Downloading %s to %s", url, download_path)
	out, err := os.Create(download_path)
	if err != nil {
		log.Println("Error creating file")
		log.Print(err)
		return 0, err
	}
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
