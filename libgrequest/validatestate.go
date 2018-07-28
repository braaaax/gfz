package libgrequest

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Validate :
func Validate(s *State, argstr, proxy string) {
	// var proxy string
	f := regexp.MustCompile("--(hc|sc|hl|sl|hw|sw|hh|sh).[0-9a-zA-Z(,|)]*")
	filterstring := f.FindString(argstr)
	fmt.Println("filterstring", filterstring, len(filterstring))
	if len(filterstring) > 0 {
		filterslice := strings.Split(filterstring, " ")
		fmt.Println("filterslice", filterslice)
		if len(filterslice) >= 2 {
			filtertag := filterslice[0]
			fmt.Println("filtertag", filtertag)
			filternum := strings.Split(filterstring, " ")[1:]
			//fmt.Println("XXXXXXX")
			PrintFilter(s, filtertag)
			//fmt.Println("s.Printer", s.Printer)
			PackFilter(s, strings.Join(filternum, ","))
		}

	} else {
		// fmt.Println("else condition [filterstring is greater than 0")
		PrintFilter(s, "sc")
		PackFilter(s, "200, 301, 302, 403")
	}
	switch strings.ToLower(s.Mode) {
	case "POST":
		// POST
		s.Processor = PostProcessor
		// fmt.Println("POST")
	case "GET":
		// GET
		s.Processor = GetProcessor
		// fmt.Println("GET")
	default:
		// GET
		s.Processor = GetProcessor
		// fmt.Println("GET")
	}
	// fmt.Println("argstr", argstr)
	FuzzMapper(argstr, s)
	var proxyURLFunc func(*http.Request) (*url.URL, error)
	proxyURLFunc = http.ProxyFromEnvironment
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			panic("proxy URL is fucked")
		}
		s.ProxyURL = proxyURL
		proxyURLFunc = http.ProxyURL(s.ProxyURL)
	}
	s.Client = &http.Client{
		Transport: &RedirectHandler{
			State: s, Transport: &http.Transport{Proxy: proxyURLFunc, TLSClientConfig: &tls.Config{InsecureSkipVerify: s.InsecureSSL}},
		},
	}
	Code, _ := GoGet(s, s.Url, s.Cookies)
	if Code == nil {
		fmt.Printf("%s", s.Url)
	}
}
