package work

import "net/url"

type Page struct {
	Url  *url.URL
	Html string
	Data map[string][]string
}

type ListPage struct {
	*Page
}

type ItemPage struct {
	*Page
	*ListPage
}

type DescPage struct{
	*Page
	*ItemPage
}

//type First struct {
//	HttpUrl *url.URL
//	Html string
//	Data map[string][]string
//}
//
//type Second struct{
//	HttpUrl *url.URL
//	Html string
//	Data map[string][]string
//}
//
//type Third struct{
//	HttpUrl *url.URL
//	Html string
//	Data map[string][]string
//}
//
//type Page struct{
//	*First
//	Seconds []*Second
//	Thirds []*Third
//}
