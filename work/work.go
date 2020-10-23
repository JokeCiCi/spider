package work

import (
	"encoding/json"
	"fmt"
	"github.com/JokeCiCi/spider/download"
	"github.com/JokeCiCi/spider/process"
	"github.com/JokeCiCi/spider/store"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func StartSpider() {
	log.Println("work start")
	var FirstChain chan *ListPage = make(chan *ListPage, 10)
	var SecondChain chan *ItemPage = make(chan *ItemPage, 100)
	var ThirdChain chan *DescPage = make(chan *DescPage, 1000)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go WorkFirst(FirstChain, SecondChain, &wg)

		wg.Add(1)
		go WorkSecond(SecondChain, ThirdChain, &wg)

		wg.Add(1)
		go WorkThird(ThirdChain, &wg)
	}

	var listUrl *url.URL = &url.URL{
		Scheme: "http",
		Host:   "xx-mh.com",
		Path:   "/home/api/cate/tp/1-0-1-1-1",
	}
	listPage := &ListPage{
		Page: &Page{
			Url:  listUrl,
			Data: make(map[string][]string),
		},
	}

	FirstChain <- listPage
	wg.Wait()
	log.Println("work end")
}

func WorkFirst(firstCh chan *ListPage, secondCh chan *ItemPage, wg *sync.WaitGroup) {
	defer wg.Done()
CHCLOSED:
	for {
		select {
		case listPage, exist := <-firstCh:
			if !exist {
				break CHCLOSED
			}
			log.Printf("WorkFirst listPage:%v\n", listPage.Url.RequestURI()) // TODO

			// ① 下载：下载列表页
			// http://xx-mh.com/home/api/cate/tp/1-0-1-1-1
			var req *http.Request = &http.Request{
				Method: http.MethodGet,
				URL:    listPage.Url,
				Header: download.FakeHeader(),
			}
			html, err := download.DownloadHTML(req)
			if err != nil {
				log.Printf("WorkFirst Failed DownloadHTML err:%v\n", err)
				continue
			}
			html, err = process.UnescapeUnicode(html)
			if err != nil {
				log.Printf("WorkFirst Failed UnescapeUnicode err:%v\n", err)
				continue
			}
			listPage.Html = html

			// ② 解析：解析列表页item
			itemInfos, err := process.FindAllInfo(html, `{"id":"([0-9]+?)"[\s\S]*?"title":"([\s\S]+?)"[\s\S]*?"image":"([\s\S]+?)"[\s\S]*?"auther":"([\s\S]+?)"[\s\S]*?"desc":"([\s\S]+?)"[\s\S]*?"keyword":"([\s\S]+?)"[\s\S]*?"cover":"([\s\S]+?)"[\s\S]*?}`)
			if err != nil {
				log.Printf("WorkFirst Failed FindAllInfo err:%v\n", err)
				continue
			}
			for _, v := range itemInfos {
				var itemPage *ItemPage = &ItemPage{
					// http://xx-mh.com/home/api/cate/tp/1-0-1-1-1
					// {"id":"2137","title":"正在插入的事…会被大家发现的！","lanmu_id":null,"create_time":"2020-09-16 00:51:59","update_time":"2020-09-16 00:54:13","sort":null,"status":"1","view":"618661","image":"\/bookimages\/\/\/cover-8d0ab116baa5b212aea9380208e200cf.jpg","type":"1","auther":"Kazuma Ichihara","desc":"本来这只是我们三个臭男生的温泉旅行…却在旅馆跟班上女同学不期而遇！从松散的浴衣裙下，淫荡的汁液淌个不停！ 更夸张的是，醉醺醺的女生们要和我们玩刺激的国王游戏…！？","mark":"217","ticai_id":"9","duzhequn_id":null,"diyu_id":null,"mhstatus":"1","tjswitch":"0","isfree":"0","cjid":"1581","cjstatus":"1","cjname":"xxmh","keyword":"温泉,淫荡,多人,国王游戏","last_chapter_title":"第章","searchnums":"0","last_chapter":"第10话","isjingpin":"0","xianmian":"0","cover":"\/bookimages\/\/\/extCover-43f9b42c41bdad91293cc78dc7013f60.jpg","ishot":"0","issole":"0","isnew":"0","h":"1","vipcanread":"1","pingfen":"9.83","ticai":"都市"}
					Page: &Page{
						Url: &url.URL{
							Scheme: listPage.Url.Scheme,
							Host:   listPage.Url.Host,
							Path:   fmt.Sprintf("/home/api/chapter_list/tp/%s-%s-%s-%s", v[0], "1", "1", "10"),
						},
						Data: map[string][]string{
							"id":       []string{v[0]},
							"title":    []string{v[1]},
							"cover":    []string{strings.Replace(process.ReplaceAllString(strings.Replace(v[2],`\`,``,-1), `/+`, `/`), `/bookimages`, `http://c1-v6e9-zp1u.cangniaobbs.com/bookimages`, -1)},
							"auther":   []string{v[3]},
							"desc":     []string{v[4]},
							"keyword":  []string{v[5]},
							"extCover": []string{strings.Replace(process.ReplaceAllString(strings.Replace(v[6],`\`,``,-1), `/+`, `/`), `/bookimages`, `http://c1-v6e9-zp1u.cangniaobbs.com/bookimages`, -1)},
						},
					},
					ListPage: listPage,
				}
				log.Printf("WorkFirst itemPage:%v itemPage:%v\n", itemPage.Url, itemPage.Data) // TODO
				secondCh <- itemPage
			}

			// ③ 解析：解析构建下个列表页
			// http://xx-mh.com/home/api/cate/tp/1-0-1-1-2
			listPageNum, err := strconv.Atoi(listPage.Url.Path[strings.LastIndex(listPage.Url.Path, "-")+1:])
			if err != nil {
				log.Printf("WorkFirst Failed Atoi err:%v\n", err)
				continue
			}
			nextListPageNum := listPageNum + 1
			nextListPagePath := fmt.Sprintf("%s%d", listPage.Url.Path[:strings.LastIndex(listPage.Url.Path, "-")+1], nextListPageNum)
			nextListPage := &ListPage{
				Page: &Page{
					Url: &url.URL{
						Scheme: listPage.Url.Scheme,
						Host:   listPage.Url.Host,
						Path:   nextListPagePath,
					},
				},
			}
			log.Printf("WorkFirst nextListPage:%v\n", nextListPage.Url.RequestURI()) // TODO
			firstCh <- nextListPage
		default:
			time.Sleep(time.Second)
		}
	}
}

func WorkSecond(secondCh chan *ItemPage, thirdCh chan *DescPage, wg *sync.WaitGroup) {
	defer wg.Done()
	var stopCount uint
CHCLOSED:
	for {
		select {
		case itemPage, exist := <-secondCh:
			if !exist {
				break CHCLOSED
			}
			if len(secondCh) == 0 {
				time.Sleep(time.Second)
				stopCount++
				if stopCount == 60 {
					break CHCLOSED
				}
			} else {
				stopCount = 0
			}

			log.Printf("WorkSecond itemPage:%v\n", itemPage.Url.RequestURI()) // TODO

			// ① 下载：下载当前章节列表页
			// http://xx-mh.com/home/api/chapter_list/tp/1058-1-1-10
			var req *http.Request = &http.Request{
				Method: http.MethodGet,
				URL:    itemPage.Url,
				Header: download.FakeHeader(),
			}
			html, err := download.DownloadHTML(req)
			if err != nil {
				log.Printf("WorkSecond Failed DownloadHTML err:%v\n", err)
				continue
			}
			html, err = process.UnescapeUnicode(html)
			if err != nil {
				log.Printf("WorkSecond Failed UnescapeUnicode err:%v\n", err)
				continue
			}
			itemPage.Html = html

			// ② 解析：解析每章节信息
			// http://xx-mh.com/home/api/chapter_list/tp/1058-1-1-10
			// http://c1-v6e9-zp1u.cangniaobbs.com/bookimages/1350/224996/1.jpg
			itemInfos, err := process.FindAllInfo(itemPage.Html, `{"id":"([0-9]+?)"[\s\S]*?"title":"([\s\S]+?)"[\s\S]*?"image":"([\s\S]+?)"[\s\S]*?"imagelist":"([\s\S]+?)"[\s\S]*?}`)
			if err != nil {
				log.Printf("WorkSecond Failed FindAllInfo err:%v\n", err)
				continue
			}
			for _, v := range itemInfos {
				var descPage *DescPage = &DescPage{
					Page: &Page{
						Data: map[string][]string{
							"id":    []string{v[0]},
							"title": []string{v[1]},
							"cover": []string{strings.Replace(process.ReplaceAllString(strings.Replace(v[2],`\`,``,-1), `/+`, `/`), `/bookimages`, `http://c1-v6e9-zp1u.cangniaobbs.com/bookimages`, 1)},
							"image": strings.Split(strings.Replace(process.ReplaceAllString(v[3], `./`, `/`), `/bookimages`, `http://c1-v6e9-zp1u.cangniaobbs.com/bookimages`, -1), `,`),
						},
					},
					ItemPage: itemPage,
				}
				log.Printf("WorkSecond descPage:%v\n", descPage.Data) // TODO
				thirdCh <- descPage
			}

			// ③ 解析：判断章节列表是否有下一页
			// http://xx-mh.com/home/api/chapter_list/tp/1058-1-2-1
			lastItemInfo, err := process.FindInfo(itemPage.Html, `"lastPage":([\s\S]*?),`)
			if err != nil {
				log.Printf("WorkSecond Failed FindInfo err:%v\n", err)
				continue
			}
			if lastItemInfo[0] == "false" {
				itemPageUrlInfo, err := process.FindInfo(itemPage.Url.Path, `[\s\S]*?([0-9]+?)-([0-9]+?)-([0-9]+?)-([0-9]+?)`)
				if err != nil {
					log.Printf("WorkSecond Failed FindInfo err:%v\n", err)
					continue
				}
				itemPageNum, err := strconv.Atoi(itemPageUrlInfo[2])
				if err != nil {
					log.Printf("WorkSecond Failed Atoi err:%v\n", err)
					continue
				}
				nextItemPageNum := itemPageNum + 1
				var nextItemPage *ItemPage = &ItemPage{
					Page: &Page{
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
				secondCh <- nextItemPage
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func WorkThird(thirdCh chan *DescPage, wg *sync.WaitGroup) {
	defer wg.Done()
	var stopCount uint
CHCLOSED:
	for {
		select {
		case descPage, exist := <-thirdCh:
			if !exist {
				break CHCLOSED
			}
			stopCount = 0

			log.Printf("WorkThird descPage:%v\n", descPage.Data) // TODO
			for k, v := range descPage.Data {
				fmt.Println("descPage", k, v)
			}
			for k, v := range descPage.ItemPage.Data {
				fmt.Println("itemPage", k, v)
			}

			resourcePath := "resource"
			comicPath := path.Join(resourcePath, descPage.ItemPage.Data["title"][0])

			comicCover := descPage.ItemPage.Data["cover"][0]
			comicCoverUrl, _ := url.Parse(comicCover)
			var req *http.Request = &http.Request{
				Method: http.MethodGet,
				URL:    comicCoverUrl,
				Header: download.FakeHeader(),
			}
			req.Header.Add("Referer", "http://xx-mh.com/")
			b, err := download.DownloadFile(req)
			if err != nil {
				log.Printf("WorkThird Failed DownloadFile err:%v\n", err)
				continue
			}
			comicCoverPath := path.Join(comicPath, "cover.jpg")
			store.MkdirAll(comicCoverPath)
			store.StoreFile(comicCoverPath, b)
			log.Printf("WorkThird comicCoverPath:%v\n", comicCoverPath)

			comicExtCover := descPage.ItemPage.Data["extCover"][0]
			comicExtCoverUrl, _ := url.Parse(comicExtCover)
			req = &http.Request{
				Method: http.MethodGet,
				URL:    comicExtCoverUrl,
				Header: download.FakeHeader(),
			}
			req.Header.Add("Referer", "http://xx-mh.com/")
			b, err = download.DownloadFile(req)
			if err != nil {
				log.Printf("WorkThird Failed DownloadFile err:%v\n", err)
				continue
			}
			comicExtCoverPath := path.Join(comicPath, "extCover.jpg")
			store.StoreFile(comicCoverPath, b)
			log.Printf("WorkThird comicExtCoverPath:%v\n", comicExtCoverPath)

			b, err = json.Marshal(descPage.ItemPage.Data)
			if err != nil {
				log.Printf("WorkThird Failed Marshal err:%v\n", err)
				continue
			}
			comicReadmePath := path.Join(comicPath, "README.md")
			store.StoreFile(comicReadmePath, b)
			log.Printf("WorkThird comicReadmePath:%v\n", comicReadmePath)

			for _, v := range descPage.Data["image"] {
				var imgUrl, err = url.Parse(v)
				if err != nil {
					log.Printf("WorkThird Failed Parse err:%v\n", err)
					continue
				}
				var req *http.Request = &http.Request{
					Method: http.MethodGet,
					URL:    imgUrl,
					Header: download.FakeHeader(),
				}
				req.Header.Add("Referer", "http://xx-mh.com/")
				b, err := download.DownloadFile(req)
				if err != nil {
					log.Printf("WorkThird Failed DownloadFile err:%v\n", err)
					continue
				}
				chapterImgPath := path.Join(comicPath, descPage.Data["title"][0], path.Base(imgUrl.Path))
				err = store.MkdirAll(chapterImgPath)
				if err != nil {
					log.Printf("WorkThird Failed MkdirAll err:%v\n", err)
					continue
				}
				store.StoreFile(chapterImgPath, b)
				log.Printf("WorkThird chapterImgPath:%v\n", chapterImgPath)
			}

			chapterCover := descPage.Data["cover"][0]
			chapterCoverUrl, _ := url.Parse(chapterCover)
			req = &http.Request{
				Method: http.MethodGet,
				URL:    chapterCoverUrl,
				Header: download.FakeHeader(),
			}
			req.Header.Add("Referer", "http://xx-mh.com/")
			b, err = download.DownloadFile(req)
			if err != nil {
				log.Printf("WorkThird Failed DownloadFile err:%v\n", err)
				continue
			}
			chapterCoverPath := path.Join(comicPath, descPage.Data["title"][0], "cover.jpg")
			store.MkdirAll(chapterCoverPath)
			store.StoreFile(chapterCoverPath, b)
			log.Printf("WorkThird chapterCoverPath:%v\n", chapterCoverPath)

			b, err = json.Marshal(descPage.Data)
			if err != nil {
				log.Printf("WorkThird Failed Marshal err:%v\n", err)
				continue
			}
			chapterReadmePath := path.Join(comicPath, descPage.Data["title"][0], "README.md")
			store.StoreFile(chapterReadmePath, b)
			log.Printf("WorkThird chapterReadmePath:%v\n", chapterReadmePath)
		default:
			if len(thirdCh) == 0 {
				time.Sleep(time.Second)
				stopCount++
				if stopCount == 60 {
					break CHCLOSED
				}
			}
		}
	}
}
