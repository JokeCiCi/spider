package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

var ListComicData *ListComic
var resourcePath string = "resource"

func Init() {
	comics, err := ioutil.ReadDir(resourcePath)
	if err != nil {
		log.Fatalf("init Failed ReadDir err:%v\n", err)
	}
	ListComicData = &ListComic{Data: make(map[string]*Comic)}
	for _, c := range comics {
		if c.IsDir() {
			comicName := c.Name()
			comicPath := path.Join(resourcePath, comicName)
			chapters, err := ioutil.ReadDir(comicPath)
			if err != nil {
				log.Fatalf("init Failed ReadDir err:%v\n", err)
			}
			comic := &Comic{Data: make(map[string]*Chapter)}
			ListComicData.Data[comicName] = comic
			for _, ch := range chapters {
				if ch.IsDir() {
					chapterName := ch.Name()
					chapter := &Chapter{}
					comic.Data[chapterName] = chapter

					chapterPath := path.Join(comicPath, chapterName)
					images, err := ioutil.ReadDir(chapterPath)
					if err != nil {
						log.Fatalf("init Failed ReadDir err:%v\n", err)
					}
					for _, img := range images {
						if !img.IsDir() {
							imageName := img.Name()
							imagePath := path.Join(chapterPath, imageName)
							chapter.Data = append(chapter.Data, imagePath)
						}
					}
				} else {
					coverName := ch.Name()
					coverPath := path.Join(comicPath, coverName)
					comic.Cover = append(comic.Cover, coverPath)
				}
			}
		}
	}

	fmt.Println(ListComicData)
}

func StartServer() {
	http.Handle("/resource/", http.StripPrefix("/resource", http.FileServer(http.Dir("resource"))))
	tmpls, err := template.ParseGlob("tmpl/*")
	if err != nil {
		log.Println("template ParseGlob err:", err)
	}

	//tmpls.Funcs(template.FuncMap{"showtime": (t time.Time, format string) string {
	//	return t.Format(format)
	//}})

	// 列出所有漫画
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		tmpls.ExecuteTemplate(w, "list_page.tmpl",ListComicData)
	})
	//// // 列出所有章节
	//http.HandleFunc("/list2", func(w http.ResponseWriter, r *http.Request) {
	//	chs := comic.ChapterObjList("漫画a")
	//	tmpls.ExecuteTemplate(w, "chapter_list.tmpl", chs)
	//})
	//
	//http.HandleFunc("/list3", func(w http.ResponseWriter, r *http.Request) {
	//	cns := comic.ChapterContents("哪有学妹这么乖", "第1话")
	//	tmpls.ExecuteTemplate(w, "chapter_list.tmpl", cns)
	//})

	http.ListenAndServe(":80", nil)
}
