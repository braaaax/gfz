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

// ConvertPrintFilter :
func ConvertPrintFilter(s *State, filternum string) {
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

// ParseWordlistArgs : set UrlFuzz Wordlists FuzzMap
func ParseWordlistArgs(str string, s *State) {
	var patZ = "-z.(file|list),[/a-zA-A0-9.-]*"
	var patW = "-w.[/0-9a-zA-Z._-]*"
	var patFuzz = "FUZ(Z|[0-9]Z)"
	var patURL = "htt(p|ps)://(.)*"
	var fname []string

	wordlistargsz := regexp.MustCompile(patZ)
	wordlistargsw := regexp.MustCompile(patW)
	fuzzRE := regexp.MustCompile(patFuzz)
	urlRE := regexp.MustCompile(patURL)

	//Parse filename arguments
	if wordlistargsz.MatchString(str) {
		// if z flag is used
		wordlistz := ArgArray(str, patZ)
		for i := 0; i < len(wordlistz); i++ {
			fname = append(fname, wordlistz[i][len("-z file,"):])
		}
	}
	if wordlistargsw.MatchString(str) {
		// if w flag is used
		wordlistw := ArgArray(str, patW)
		for i := 0; i < len(wordlistw); i++ {
			fname = append(fname, wordlistw[i][len("-w "):])
		}
	}
	//parse url arguments
	URLs := urlRE.FindAllString(str, -1)
	if len(URLs) != 0 {
		FUZZs := fuzzRE.FindAllString(URLs[0], -1)
		mapfuzzs2wordlists := make(map[string]string)
		for index, fzN := range FUZZs {
			mapfuzzs2wordlists[fname[index]] = fzN
		}
		s.WordListFiles = fname
		s.Fuzzer.Fuzzmap = mapfuzzs2wordlists

	}
}

// Validate :
func Validate(s *State, argstr, proxy string) {
	f := regexp.MustCompile("--(hc|sc|hl|sl|hw|sw|hh|sh).[0-9a-zA-Z(,|)]*")
	RawFilterStringArgs := f.FindString(argstr)
	if len(RawFilterStringArgs) > 0 {
		printfilterlst := strings.Split(RawFilterStringArgs, " ")
		if len(printfilterlst) >= 2 {
			printfilterargs := printfilterlst[0]
			filternumlst := strings.Split(RawFilterStringArgs, " ")[1:]
			ParsePrintFilterArgs(s, printfilterargs)
			ConvertPrintFilter(s, strings.Join(filternumlst, ","))
		}
	} else {
		// default
		ParsePrintFilterArgs(s, "sc")
		ConvertPrintFilter(s, "200,301,302,403")

	}
	/*switch strings.ToLower(s.Mode) {
	case "POST":
		s.Processor = PostProcessor
	case "GET":
		s.Processor = GoGet
	default:
		s.Processor = GoGet
	}*/
	ParseWordlistArgs(argstr, s)
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
		fmt.Printf("Cannot reach %s\n", s.URL)
		return
	}
}
