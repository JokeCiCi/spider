package web


type Chapter struct{
	Data []string
}

type Comic struct{
	Cover []string
	Data map[string]*Chapter
}

type ListComic struct{
	Data map[string]*Comic
}

