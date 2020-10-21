package web

import (
	"html/template"
	"net/http"
)

// html/template
func StartServer() {
	var temp *template.Template = template.New("test")
	temp.Funcs(template.FuncMap{
		"inc": func(num int) int {
			num = num + 1
			return num
		},
		"isEnd": func(num, end int) bool {
			return num%end == 0
		},
	})
	tmpls, _ := temp.ParseGlob("tmpl/*")

	http.Handle("/resource/", http.StripPrefix("/resource", http.FileServer(http.Dir("resource"))))

	http.HandleFunc("/comics", func(w http.ResponseWriter, r *http.Request) {
		tmpls.ExecuteTemplate(w, "comic_list.tmpl", listPageData)
	})
	http.HandleFunc("/chapters", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		comicName := r.Form.Get("comicname")
		tmpls.ExecuteTemplate(w, "chapter_list.tmpl", listPageData.Data[comicName])
	})
	http.HandleFunc("/chapter", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		comicName := r.Form.Get("comicname")
		chapterName := r.Form.Get("chaptername")
		tmpls.ExecuteTemplate(w, "chapter.tmpl", listPageData.Data[comicName].Chapters[chapterName])
	})
	http.ListenAndServe(":80", nil)
}
