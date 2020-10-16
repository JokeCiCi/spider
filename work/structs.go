package work

import "net/url"

type First struct {
	HttpUrl *url.URL
	Html string
	Data map[string][]string
}

type Second struct{
	HttpUrl *url.URL
	Html string
	Data map[string][]string
}

type Third struct{
	HttpUrl *url.URL
	Html string
	Data map[string][]string
}

type Page struct{
	*First
	Seconds []*Second
	Thirds []*Third
}
