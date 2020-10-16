package main

import (
	"fmt"
	"github.com/JokeCiCi/spider/work"
	"net/url"
	"sync"
)

var wg sync.WaitGroup

func main() {
	//https://search.jd.com/Search?keyword=<KEYWORD>&page<PG NUM>&s=<ELE NUM> // SEARCH.adv_param={page:"([0-9]+?)"\s*
	// jd翻页: npage=curPage+2 s=(npage+1)/2*50 + 1
	//https://s.taobao.com/search?q=<KEYWORD>&s=<ELE NUM> // <input class="input J_Input" type="number" value="3" min="1" max="100" aria-label="页码输入框">
	// tb翻页: s=(n-1)*44
	var JDFirstChain chan *work.Page = make(chan *work.Page, 10)
	var JDSecondChain chan *work.Page = make(chan *work.Page, 100)
	var JDThirdChain chan *work.Page = make(chan *work.Page, 100)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go work.JDWorkFirst(JDFirstChain, JDSecondChain, &wg)

		wg.Add(1)
		go work.JDWorkSecond(JDSecondChain, JDThirdChain, &wg)

		wg.Add(1)
		go work.JDWorkThird(JDThirdChain, &wg)
	}

	var httpUrl *url.URL = &url.URL{
		Scheme:   "https",
		Host:     "search.jd.com",
		Path:     "/Search",
		RawQuery: url.Values{"keyword": []string{"电脑"}}.Encode(),
	}
	startPage := &work.Page{
		First: &work.First{
			HttpUrl: httpUrl,
		},
	}
	JDFirstChain <- startPage

	wg.Wait()
	fmt.Println("end")
	//httpUrl := &url.URL{
	//	Scheme: "https",
	//	Host:   "s.taobao.com",
	//	Path:   "/search",
	//	RawQuery: url.Values{
	//		"q": []string{"鞋子"},
	//	}.Encode(),
	//}
	//html,_:=download.DownloadHTML(httpUrl)
	//fmt.Println(html)

}