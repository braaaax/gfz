package libgrequest

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// SetFilter :
func SetFilter(s *State, filternum string) {
	for _, c := range strings.Split(filternum, ",") {
		i, err := strconv.Atoi(c)
		i64 := int64(i)
		if err != nil {
			fmt.Println(err)
		} else {
			s.Filter.Add(i64)
		}
	}
}

// SetFuzzmap : set UrlFuzz Wordlists FuzzMap
func SetFuzzmap(str string, s *State) {
	var patZ = "-z.(file|list),[/a-zA-A0-9.-]*"
	var patW = "-w.[/0-9a-zA-Z._-]*"
	var patFuzz = "FUZ(Z|[0-9]Z)"
	var patURL = "htt(p|ps)://(.)*"
	var fname []string

	zRE := regexp.MustCompile(patZ)
	wRE := regexp.MustCompile(patW)
	fuzzRE := regexp.MustCompile(patFuzz)
	urlRE := regexp.MustCompile(patURL)

	//Parse filename arguments
	if zRE.MatchString(str) {
		// if z flag is used
		zlist := ArgArray(str, patZ)
		for i := 0; i < len(zlist); i++ {
			fname = append(fname, zlist[i][len("-z file,"):])
		}
	}
	if wRE.MatchString(str) {
		// if w flag is used
		wlist := ArgArray(str, patW)
		for i := 0; i < len(wlist); i++ {
			fname = append(fname, wlist[i][len("-w "):])
		}
	}
	//parse url arguments
	URLs := urlRE.FindAllString(str, -1)
	if len(URLs) != 0 {
		FUZZs := fuzzRE.FindAllString(URLs[0], -1)
		fm := make(map[string]string)
		for index, m := range FUZZs {
			fm[fname[index]] = m
		}
		s.WordListFiles = fname
		s.Fuzzer.Fuzzmap = fm

	}
}

// Validate :
func Validate(s *State, argstr, proxy string) {
	f := regexp.MustCompile("--(hc|sc|hl|sl|hw|sw|hh|sh).[0-9a-zA-Z(,|)]*")
	filterstring := f.FindString(argstr)
	if len(filterstring) > 0 {
		filterslice := strings.Split(filterstring, " ")
		if len(filterslice) >= 2 {
			filtertag := filterslice[0]
			filternum := strings.Split(filterstring, " ")[1:]
			PrintFilter(s, filtertag)
			SetFilter(s, strings.Join(filternum, ","))
		}
	} else {
		// default
		PrintFilter(s, "sc")
		SetFilter(s, "200,301,302,403")

	}
	switch strings.ToLower(s.Mode) {
	case "POST":
		s.Processor = PostProcessor
	case "GET":
		s.Processor = GetProcessor
	default:
		s.Processor = GetProcessor
	}
	SetFuzzmap(argstr, s)
	var proxyURLFunc func(*http.Request) (*url.URL, error)

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

	s.Client = &http.Client{
		Transport: &RedirectHandler{
			State: s,
			Transport: &http.Transport{
				Proxy: proxyURLFunc,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: s.InsecureSSL}},
		},
	}
	Code, _ := GoGet(s, s.URL, s.Cookies)
	if Code == nil {
		fmt.Printf("Cannot reach %s", s.URL)
		return
	}
}
