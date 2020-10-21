package web

type Chapter struct {
	Data map[string][]string
}

type Comic struct {
	Data     map[string][]string
	Chapters map[string]*Chapter
}

type ListPage struct {
	Data map[string]*Comic
}
