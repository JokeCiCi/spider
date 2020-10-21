package main

import (
	"fmt"
	"github.com/JokeCiCi/spider/download"
	"github.com/JokeCiCi/spider/process"
	"github.com/JokeCiCi/spider/web"
	"github.com/JokeCiCi/spider/work"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	//web.PrintListPage()
	web.StartServer()
	work.StartSpider()
	//fmt.Println(path.Base("/bookimages/88027/8027007/31.jpg"))

	// http://xx-mh.com/home/api/cate/tp/1-0-1-1-1
	// http://xx-mh.com/home/api/chapter_list/tp/2137-1-1-10
	// http://c1-v6e9-zp1u.cangniaobbs.com/bookimages/1350/224996/1.jpg
	// Referer: http://xx-mh.com/

	//	连载中
	//http://xx-mh.com/home/api/cate/tp/1-0-0-1-1
	//	已完结
	//http://xx-mh.com/home/api/cate/tp/1-0-1-1-1
}

func test() {
	var listUrl *url.URL = &url.URL{
		Scheme: "http",
		Host:   "xx-mh.com",
		Path:   "/home/api/cate/tp/1-0-1-1-1",
	}
	listPage := &work.ListPage{
		Page: &work.Page{
			Url:  listUrl,
			Data: make(map[string][]string),
		},
	}

	var req *http.Request = &http.Request{
		Method: http.MethodGet,
		URL:    listPage.Url,
		Header: download.FakeHeader(),
	}

	log.Printf("WorkFirst listPage:%v\n", listPage.Url.RequestURI()) // TODO

	// ① 下载：下载列表页
	// http://xx-mh.com/home/api/cate/tp/1-0-1-1-1
	html, err := download.DownloadHTML(req)
	if err != nil {
		log.Printf("WorkFirst Failed DownloadHTML err:%v\n", err)
		return
	}
	html, err = process.UnescapeUnicode(html)
	if err != nil {
		log.Printf("WorkFirst Failed UnescapeUnicode err:%v\n", err)
		return
	}
	listPage.Html = html

	// ② 解析：解析列表页item
	itemInfos, err := process.FindAllInfo(html, `{"id":"([0-9]+?)"[\s\S]*?"title":"([\s\S]+?)"[\s\S]*?"image":"([\s\S]+?)"[\s\S]*?"auther":"([\s\S]+?)"[\s\S]*?"desc":"([\s\S]+?)"[\s\S]*?"keyword":"([\s\S]+?)"[\s\S]*?"cover":"([\s\S]+?)"[\s\S]*?}`)
	if err != nil {
		log.Printf("WorkFirst Failed FindAllInfo err:%v\n", err)
		return
	}

	var itemPage *work.ItemPage
	for _, v := range itemInfos {
		itemPage = &work.ItemPage{
			Page: &work.Page{
				Url: &url.URL{
					Scheme: listPage.Url.Scheme,
					Host:   listPage.Url.Host,
					// http://xx-mh.com/home/api/chapter_list/tp/1012-1-1-10
					Path: fmt.Sprintf("/home/api/chapter_list/tp/%s-%s-%s-%s", v[0], "1", "1", "10"),
				},
				Data: map[string][]string{
					"id":      []string{v[0]},
					"title":   []string{v[1]},
					"auther":  []string{v[3]},
					"desc":    []string{v[4]},
					"keyword": []string{v[5]},
					"cover": []string{
						fmt.Sprintf("%s%s", "http://c1-v6e9-zp1u.cangniaobbs.com", strings.Replace(v[2], "\\", "", -1)),
						fmt.Sprintf("%s%s", "http://c1-v6e9-zp1u.cangniaobbs.com", strings.Replace(v[6], "\\", "", -1)),
					},
				},
			},
			ListPage: listPage,
		}
		log.Printf("WorkFirst itemPage:%v itemPage:%v\n", itemPage.Url, itemPage.Data) // TODO
	}

	// ③ 解析：解析构建下个列表页
	// http://xx-mh.com/home/api/cate/tp/1-0-1-1-2
	listPageNum, err := strconv.Atoi(listPage.Url.Path[strings.LastIndex(listPage.Url.Path, "-")+1:])
	if err != nil {
		log.Printf("WorkFirst Failed Atoi err:%v\n", err)
		return
	}
	nextListPageNum := listPageNum + 1
	nextListPagePath := fmt.Sprintf("%s%d", listPage.Url.Path[:strings.LastIndex(listPage.Url.Path, "-")+1], nextListPageNum)
	nextListPage := &work.ListPage{
		Page: &work.Page{
			Url: &url.URL{
				Scheme: listPage.Url.Scheme,
				Host:   listPage.Url.Host,
				Path:   nextListPagePath,
			},
		},
	}
	log.Printf("WorkFirst nextListPage:%v\n", nextListPage.Url.RequestURI()) // TODO

	// ------------------------------------------------------------------
	log.Printf("WorkSecond itemPage:%v\n", itemPage.Url.RequestURI()) // TODO

	// ① 下载：下载当前章节列表页
	// http://xx-mh.com/home/api/chapter_list/tp/1012-1-1-10
	req = &http.Request{
		Method: http.MethodGet,
		URL:    itemPage.Url,
		Header: download.FakeHeader(),
	}
	html, err = download.DownloadHTML(req)
	if err != nil {
		log.Printf("WorkSecond Failed DownloadHTML err:%v\n", err)
		return
	}
	html, err = process.UnescapeUnicode(html)
	if err != nil {
		log.Printf("WorkSecond Failed UnescapeUnicode err:%v\n", err)
		return
	}
	itemPage.Html = html

	// ② 解析：解析每章节信息
	itemInfos, err = process.FindAllInfo(itemPage.Html, `{"id":"([0-9]+?)"[\s\S]*?"title":"([\s\S]+?)"[\s\S]*?"imagelist":"([\s\S]+?)"[\s\S]*?}`)
	if err != nil {
		log.Printf("WorkSecond Failed FindAllInfo err:%v\n", err)
		return
	}
	for _, v := range itemInfos {
		var descPage *work.DescPage = &work.DescPage{
			Page: &work.Page{
				Data: map[string][]string{
					"id":    []string{v[0]},
					"title": []string{v[1]},
					// http://c1-v6e9-zp1u.cangniaobbs.com/bookimages/1294/222690/1.jpg
					"image": strings.Split(strings.Replace(strings.Replace(v[2], "\\", "", -1), "./", "http://c1-v6e9-zp1u.cangniaobbs.com/", -1), ","),
				},
			},
			ItemPage: itemPage,
		}
		log.Printf("WorkSecond descPage:%v\n", descPage.Data) // TODO
	}

	// ③ 解析：判断章节列表是否有下一页
	// http://xx-mh.com/home/api/chapter_list/tp/1012-1-2-10
	lastItemInfo, err := process.FindInfo(itemPage.Html, `"lastPage":([\s\S]*?),`)
	if err != nil {
		log.Printf("WorkSecond Failed FindInfo err:%v\n", err)
		return
	}
	if lastItemInfo[0] == "false" {
		itemPageUrlInfo, err := process.FindInfo(itemPage.Url.Path, `[\s\S]*?([0-9]+?)-([0-9]+?)-([0-9]+?)-([0-9]+?)`)
		if err != nil {
			log.Printf("WorkSecond Failed FindInfo err:%v\n", err)
			return
		}
		itemPageNum, err := strconv.Atoi(itemPageUrlInfo[2])
		if err != nil {
			log.Printf("WorkSecond Failed Atoi err:%v\n", err)
			return
		}
		nextItemPageNum := itemPageNum + 1
		var nextItemPage *work.ItemPage = &work.ItemPage{
			Page: &work.Page{
				Url: &url.URL{
					Scheme: itemPage.Url.Scheme,
					Host:   itemPage.Url.Host,
					Path:   fmt.Sprintf("/home/api/chapter_list/tp/%s-%s-%d-%s", itemPageUrlInfo[0], itemPageUrlInfo[1], nextItemPageNum, itemPageUrlInfo[3]),
				},
				Data: itemPage.Data,
			},
			ListPage: itemPage.ListPage,
		}
		log.Printf("WorkSecond nextItemPage:%v\n", nextItemPage.Url.RequestURI()) // TODO
	}
}
