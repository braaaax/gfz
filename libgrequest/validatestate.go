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
	var patZ = "-z (file|File|FILE),[/a-zA-A0-9.-_]*"
	var patZ2 = "-z (range|Range|RANGE),[0-9-]*" // put a limit
	var patW = "-w [/0-9a-zA-Z._-]*"
	var patFuzz = "FUZ(Z|[0-9]Z)"
	var patURL = "htt(p|ps)://(.)*"
	var fname []string

	wordlistargszfile := regexp.MustCompile(patZ)
	wordlistargszrange := regexp.MustCompile(patZ2)
	wordlistargswfile := regexp.MustCompile(patW)
	fuzzRE := regexp.MustCompile(patFuzz)
	urlRE := regexp.MustCompile(patURL)

	//Parse filename arguments
	if wordlistargszfile.MatchString(str) && !wordlistargszrange.MatchString(str) {
		// if z flag is used
		wordlistz := ArgArray(str, patZ)
		for i := 0; i < len(wordlistz); i++ {
			fname = append(fname, wordlistz[i][len("-z file,"):])
		}
	}
	if wordlistargszrange.MatchString(str) {
		rangewordlists := [][]string{}
		numberstrings := []string{}
		// get string from argparse paterrn

		patresult := wordlistargszrange.FindAllString(str, -1)
		if len(patresult) == 0 {
			fmt.Println("patresult: ", patresult)

		}
		// return list of number therein
		numRE := regexp.MustCompile("[0-9]+")
		numbs := numRE.FindAllString(patresult[0], -1)
		if len(numbs) != 0 {
			start, _ := strconv.Atoi(numbs[0])
			end, _ := strconv.Atoi(numbs[1])
			for i := start; i < end; i++ {
				numberstrings = append(numberstrings, strconv.Itoa(i))
			}
			rangewordlists = append(rangewordlists, numberstrings)
			s.Fuzzer.Wordlists = rangewordlists
		}
	}
	if wordlistargswfile.MatchString(str) {
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
		if len(fname) == len(FUZZs) {
			mapfuzzs2wordlists := make(map[string]string)
			for index, fzN := range FUZZs {
				mapfuzzs2wordlists[fname[index]] = fzN
			}
			s.WordListFiles = fname
			s.Fuzzer.Fuzzmap = mapfuzzs2wordlists
		}
	}
}

// Validate : final input error checks before run.
func Validate(s *State, argstr, proxy string) {
	// parse output filter args
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

	// set mode for when GoPost is added
	/*switch strings.ToLower(s.Mode) {
	case "POST":
		s.Processor = GoPost
	case "GET":
		s.Processor = GoGet
	default:
		s.Processor = GoGet
	}*/

	ParseWordlistArgs(argstr, s)

	// set proxy info
	var proxyURLFunc func(*http.Request) (*url.URL, error)
	proxyURLFunc = http.ProxyFromEnvironment
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			panic("[!] proxy URL is fucked")
		}
		s.ProxyURL = proxyURL
		proxyURLFunc = http.ProxyURL(s.ProxyURL)
	}

	// Client struct initialization
	s.Client = &http.Client{
		Transport: &RedirectHandler{
			State: s,
			Transport: &http.Transport{
				Proxy: proxyURLFunc,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: s.InsecureSSL}},
		},
	}

	// Pass/Fail
	if len(s.URL) == 0 {
		return
	}
	// fmt.Printf("\n[+] target URL: " + s.URL)
	Code, _ := GoGet(s, s.URL, s.Cookies)
	if Code == nil {
		fmt.Printf("{!] Cannot reach %s\n", s.URL)
		return
	}
}
