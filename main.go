package main

import (
	"github.com/JokeCiCi/spider/work"
	"sync"
)

var wg sync.WaitGroup

func main() {
	//web.PrintListPage()
	//web.StartServer()
	work.StartSpider()

	// http://xx-mh.com/home/api/cate/tp/1-0-1-1-1
	// http://xx-mh.com/home/api/chapter_list/tp/2137-1-1-10
	// http://c1-v6e9-zp1u.cangniaobbs.com/bookimages/1350/224996/1.jpg
	// Referer: http://xx-mh.com/

	//	连载中
	//http://xx-mh.com/home/api/cate/tp/1-0-0-1-1
	//	已完结
	//http://xx-mh.com/home/api/cate/tp/1-0-1-1-1

	//	chaptercover
	//http://c1-v6e9-zp1u.cangniaobbs.com/bookimages/1456/cover-ed61c2fccad4e4f9803107d0daf8c668.jpg
	//Host:c1-v6e9-zp1u.cangniaobbs.com
	//Accept: image/avif,image/webp,image/apng,image/*,*/*;q=0.8
	//	Accept-Encoding: gzip, deflate
	//	Accept-Language: zh-CN,zh;q=0.9
	//Connection: keep-alive
	//Host: c1-v6e9-zp1u.cangniaobbs.com
	//Referer: http://xx-mh.com/
	//
	//	comiccover
	//http://xx-mh.com/
	//http://c1-v6e9-zp1u.cangniaobbs.com/bookimages/cover-a13ea9972552b524335daed39679d1c3.jpg
	//
	//http://xx-mh.com/
	//http://c1-v6e9-zp1u.cangniaobbs.com/bookimages/extCover-4f2de30626c6f82c5a1234197c452f2b.jpg
}
