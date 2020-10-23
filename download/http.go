package download

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	rand.Seed(time.Now().UnixNano())
	client = http.DefaultClient
}

func FakeHeader() http.Header {
	return http.Header{
		"User-Agent":      UserAgent(),
		"Accept":          []string{"text/html", "application/xhtml+xml", "application/xml;q=0.9", "image/avif", "image/webp", "image/apng", "*/*;q=0.8", "application/signed-exchange;v=b3;q=0.9"},
		"Accept-Encoding": []string{"deflate"},
	}
}

func DownloadHTML(req *http.Request) (string, error) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("DownloadHTML Do failed err:%v\n", err)
		return "", fmt.Errorf("DownloadHTML failed url:%v", req.URL)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("DownloadHTML ReadAll failed err:%v\n", err)
		return "", fmt.Errorf("DownloadHTML failed url:%v", req.URL)
	}
	return string(b), nil
}

func DownloadFile(req *http.Request) ([]byte, error) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("DownloadHTML Do failed err:%v\n", err)
		return nil, fmt.Errorf("DownloadHTML failed url:%v", req.URL)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("DownloadHTML ReadAll failed err:%v\n", err)
		return nil, fmt.Errorf("DownloadHTML failed url:%v", req.URL)
	}
	return b, nil
}
