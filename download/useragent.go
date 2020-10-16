package download

import "math/rand"

var uas []string

func init() {
	uas = []string{
		"Mozilla/4.0 (Macintosh) AppleWebKit/70.48 (KHTML, like Gecko) Firefox/21.0 Safari/561.92",
		"Mozilla/5.0 (Macintosh) AppleWebKit/82.91 (KHTML, like Gecko) Firefox/55.0 Safari/194.29",
		"Mozilla/4.0 (Macintosh) AppleWebKit/47.7 (KHTML, like Gecko) Firefox/49.0 Safari/565.12",
		"Mozilla/4.0 (Windows NT 5.2) AppleWebKit/95.11 (KHTML, like Gecko) Firefox/37.0 Safari/109.31",
		"Mozilla/5.0 (X11) AppleWebKit/46.81 (KHTML, like Gecko) Firefox/37.0 Safari/449.29",
		"Mozilla/4.0 (Windows NT 6.0) AppleWebKit/10.54 (KHTML, like Gecko) Edge/17.13007 Safari/160.82",
		"Mozilla/4.0 (Macintosh) AppleWebKit/54.4 (KHTML, like Gecko) Firefox/35.0 Safari/577.66",
		"Mozilla/4.0 (X11) AppleWebKit/43.39 (KHTML, like Gecko) Chrome/61.0.2730.267 Safari/378.29",
		"Mozilla/5.0 (Macintosh) AppleWebKit/60.72 (KHTML, like Gecko) Firefox/31.0 Safari/250.38",
		"Mozilla/4.0 (X11) AppleWebKit/97.45 (KHTML, like Gecko) Firefox/48.0 Safari/190.76",
		"Mozilla/5.0 (Macintosh) AppleWebKit/10.30 (KHTML, like Gecko) Chrome/46.2.2806.177 Safari/551.16",
		"Mozilla/4.0 (Windows NT 6.0) AppleWebKit/93.35 (KHTML, like Gecko) Firefox/58.0 Safari/450.61",
		"Mozilla/5.0 (Macintosh) AppleWebKit/28.81 (KHTML, like Gecko) Chrome/51.2.2116.377 Safari/466.81",
		"Mozilla/5.0 (Macintosh) AppleWebKit/21.66 (KHTML, like Gecko) Edge/13.18331 Safari/517.82",
		"Mozilla/5.0 (Macintosh) AppleWebKit/100.72 (KHTML, like Gecko) Chrome/66.2.4741.172 Safari/213.74",
		"Mozilla/5.0 (Windows NT 6.0) AppleWebKit/38.86 (KHTML, like Gecko) Firefox/23.0 Safari/129.88",
		"Mozilla/5.0 (Windows NT 5.2) AppleWebKit/82.18 (KHTML, like Gecko) Firefox/29.0 Safari/148.30",
		"Mozilla/4.0 (Windows NT 10.0) AppleWebKit/53.71 (KHTML, like Gecko) Firefox/36.0 Safari/271.29",
		"Mozilla/5.0 (X11) AppleWebKit/21.76 (KHTML, like Gecko) Firefox/45.0 Safari/303.45",
		"Mozilla/5.0 (Windows NT 6.0) AppleWebKit/51.10 (KHTML, like Gecko) Chrome/66.1.1766.488 Safari/102.27",
	}
}

func UserAgent() []string {
	return []string{uas[rand.Intn(len(uas))]}
}
