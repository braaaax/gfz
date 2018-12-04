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
	c := SafeCounter{v: make(map[string]int)}
	f := regexp.MustCompile("--(hc|sc|hl|sl|hw|sw|hh|sh).[0-9a-zA-Z(,|)]*")
	filterstring := f.FindString(argstr)
	if len(filterstring) > 0 {
		filterslice := strings.Split(filterstring, " ")
		if len(filterslice) >= 2 {
			filtertag := filterslice[0]
			filternum := strings.Split(filterstring, " ")[1:]
			PrintFilter(s, filtertag)
			PackFilter(s, strings.Join(filternum, ","))
		}
	} else {
		// default
		PrintFilter(s, "sc")
		PackFilter(s, "200,301,302,403")
	}
	switch strings.ToLower(s.Mode) {
	case "POST":
		s.Processor = PostProcessor
	case "GET":
		s.Processor = GetProcessor
	default:
		s.Processor = GetProcessor
	}
	FuzzMapper(argstr, s)
	var proxyURLFunc func(*http.Request) (*url.URL, error)
	//
	//TODO: proxy stuff
	proxyURLFunc = http.ProxyFromEnvironment
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			panic("proxy URL is fucked")
		}
		s.ProxyURL = proxyURL
		proxyURLFunc = http.ProxyURL(s.ProxyURL)
	}
	//
	//
	s.Client = &http.Client{
		Transport: &RedirectHandler{
			State: s, 
			Transport: &http.Transport{
				Proxy: proxyURLFunc, 
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: s.InsecureSSL}},
		},
	}
	Code, _ := GoGet(s, s.Url, s.Cookies)
	if Code == nil {
		fmt.Printf("%s", s.Url)
	}
}
