package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

var listPageData *ListPage
var resourcePath string = "resource"

func init() {
	comicDirs, err := ioutil.ReadDir(resourcePath)
	if err != nil {
		log.Fatalf("init Failed ReadDir err:%v\n", err)
	}
	listPageData = &ListPage{Data: make(map[string]*Comic)}
	for _, cd := range comicDirs {
		if cd.IsDir() {
			comicName := cd.Name()
			comic := &Comic{
				Data:     make(map[string][]string),
				Chapters: make(map[string]*Chapter),
			}
			listPageData.Data[comicName] = comic
			comic.Data["name"] = []string{comicName}

			comicPath := path.Join(resourcePath, comicName) // ① 解析comic目录
			comicFiles, err := ioutil.ReadDir(comicPath)
			if err != nil {
				log.Fatalf("init Failed ReadDir err:%v\n", err)
			}
			for _, cf := range comicFiles {
				if !cf.IsDir() {
					comicFileName := cf.Name()
					if strings.HasSuffix(comicFileName, "jpg") {
						comicCoverPath := path.Join(comicPath, comicFileName) // ② 解析comic目录文件
						comic.Data["cover"] = append(comic.Data["cover"], comicCoverPath)
					}
				} else {
					chapterName := cf.Name()
					chapter := &Chapter{
						Data: make(map[string][]string),
					}
					comic.Chapters[chapterName] = chapter

					chapterPath := path.Join(comicPath, chapterName) // ③ 解析chapter目录
					chapterFiles, err := ioutil.ReadDir(chapterPath)
					if err != nil {
						log.Fatalf("init Failed ReadDir err:%v\n", err)
					}
					for _, chf := range chapterFiles {
						if !chf.IsDir() {
							chapterFileName := chf.Name()
							if strings.HasPrefix(chapterFileName, "cover") {
								chapterCoverPath := path.Join(chapterPath, chapterFileName) // ④ 解析chapter目录文件
								chapter.Data["cover"] = []string{chapterCoverPath}
							} else {
								chapterImgPath := path.Join(chapterPath, chapterFileName)
								chapter.Data["image"] = append(chapter.Data["image"], chapterImgPath)
							}
						}
					}
				}
			}
		}
	}
}

func PrintListPage() {
	for comicName, comic := range listPageData.Data {
		for comicKey, comicValue := range comic.Data {
			fmt.Println(comicName, comicKey, comicValue)
		}
		for chapterName, chapter := range comic.Chapters {
			for chapterKey, chapterValue := range chapter.Data {
				fmt.Println(chapterName, chapterKey, chapterValue)
			}
		}
	}
}
