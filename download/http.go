package download

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client *http.Client

func init() {
	rand.Seed(time.Now().UnixNano())
	client = http.DefaultClient
}

func fakeHeader(domainFilter string) http.Header {
	return http.Header{
		"User-Agent":      UserAgent(),
		"Cookie":          Cookie(domainFilter),
		"Accept":          []string{"text/html", "application/xhtml+xml", "application/xml;q=0.9", "image/avif", "image/webp", "image/apng", "*/*;q=0.8", "application/signed-exchange;v=b3;q=0.9"},
		"Accept-Encoding": []string{"deflate"},
	}
}

func DownloadHTML(httpUrl *url.URL) (string, error) {
	time.Sleep(time.Millisecond * 100 * time.Duration(rand.Intn(100)))
	var req *http.Request = &http.Request{
		Method: http.MethodGet,
		URL:    httpUrl,
		Header: fakeHeader(string([]byte(httpUrl.Host)[strings.Index(httpUrl.Host,"."):])),
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("http DefaultClient Do failed,err%v\n", err)
		return "", fmt.Errorf("DownloadHTML failed url:%v", httpUrl)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil ReadAll failed,err%v\n", err)
		return "", fmt.Errorf("DownloadHTML failed url:%v", httpUrl)
	}
	return string(b), nil
}

func DownloadFile(httpUrl *url.URL) ([]byte, error) {
	time.Sleep(time.Millisecond * 100 * time.Duration(rand.Intn(100)))
	var req *http.Request = &http.Request{
		Method: http.MethodGet,
		URL:    httpUrl,
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("http DefaultClient Do failed,err%v\n", err)
		return nil, fmt.Errorf("DownloadFile url:%v", httpUrl)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil ReadAll failed,err%v\n", err)
		return nil, fmt.Errorf("DownloadFile failed url:%v", httpUrl)
	}
	return b, nil
}
