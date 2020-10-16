package download

import (
	"fmt"
	"github.com/zellyn/kooky"
	"log"
	"net/http"
	"os/user"
	"runtime"
	"time"
)

type Browser uint8

const (
	Chrome Browser = iota
	Firefox
	Safari
)

func ReadCookies(browser Browser, domainFilter, nameFilter string, expireAfter time.Time) ([]*http.Cookie, error) {
	usr, _ := user.Current()
	switch {
	case runtime.GOOS == "windows":
		switch browser {
		case Chrome:
			cookiesFile := fmt.Sprintf("%s/AppData/Local/Google/Chrome/User Data/Default/Cookies", usr.HomeDir)
			kookyCookies, err := kooky.ReadChromeCookies(cookiesFile, domainFilter, nameFilter, expireAfter)
			var cookies []*http.Cookie
			for _, v := range kookyCookies {
				cookies = append(cookies, &http.Cookie{Name: v.Name, Value: v.Value})
			}
			return cookies, err
		case Firefox:
			return nil, fmt.Errorf("ReadCookies failed unsopport browser:%v", browser)
		case Safari:
			return nil, fmt.Errorf("ReadCookies failed unsopport browser:%v", browser)
		default:
			return nil, fmt.Errorf("ReadCookies failed unsopport browser:%v", browser)
		}
	case runtime.GOOS == "linux":
		switch browser {
		case Chrome:
			cookiesFile := fmt.Sprintf("%s/Library/Application Support/Google/Chrome/Default/Cookies", usr.HomeDir)
			kookyCookies, err := kooky.ReadChromeCookies(cookiesFile, domainFilter, nameFilter, expireAfter)
			var cookies []*http.Cookie
			for _, v := range kookyCookies {
				cookies = append(cookies, &http.Cookie{Name: v.Name, Value: v.Value})
			}
			return cookies, err
		case Firefox:
			return nil, fmt.Errorf("ReadCookies failed unsopport browser:%v", browser)
		case Safari:
			return nil, fmt.Errorf("ReadCookies failed unsopport browser:%v", browser)
		default:
			return nil, fmt.Errorf("ReadCookies failed unsopport browser:%v", browser)
		}
	case runtime.GOOS == "darwin":
		return nil, fmt.Errorf("ReadCookies failed unsopport os:%v", runtime.GOOS)
	default:
		return nil, fmt.Errorf("ReadCookies failed unsopport os:%v", runtime.GOOS)
	}
}

func ConvertCookie(cookies []*http.Cookie) ([]string, error) {
	if cookies == nil || len(cookies) == 0 {
		return nil, fmt.Errorf("Cookie failed cookies is invalid cookies:%v", cookies)
	}
	var cookie []string
	for _, c := range cookies {
		s := fmt.Sprintf("%s=%s", c.Name, c.Value)
		cookie = append(cookie, s)
	}
	return cookie, nil
}

func Cookie(domainFilter string) []string {
	cookies, err := ReadCookies(Chrome, domainFilter, "", time.Time{})
	if err != nil {
		log.Printf("Cookie failed err:%v\n", err)
		return []string{}
	}
	cookie, err := ConvertCookie(cookies)
	if err != nil {
		log.Printf("Cookie failed err:%v\n", err)
		return []string{}
	}
	return cookie
}
