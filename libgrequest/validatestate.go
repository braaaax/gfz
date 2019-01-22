package libgrequest

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// convPrintFilter :
func convPrintFilter(s *State, filternum string) {
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

func (s *State) readfile(fname string) {
	fn, err := os.Open(fname)
	check(err)
	defer fn.Close()
	var lines []string
	scanner := bufio.NewScanner(fn)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	s.Fuzzer.Wordlists = append(s.Fuzzer.Wordlists, lines)
}

// ParseWordlistArgs : set UrlFuzz Wordlists FuzzMap
func ParseWordlistArgs(str string, s *State) bool {
	var patzfile = "-z (file|File|FILE),[/a-zA-A0-9.-_]*"
	var patzrange = "-z (range|Range|RANGE),[0-9-]*" // put a limit
	var patzlist = "-z (list|List|LIST),[a-zA-Z0-9.-]*"
	var patwfile = "-w [/0-9a-zA-Z._-]*"
	zlistwordlist := []string{}
	zrangewordlist := []string{}
	zfile := regexp.MustCompile(patzfile)
	zrange := regexp.MustCompile(patzrange)
	wfile := regexp.MustCompile(patwfile)
	zlist := regexp.MustCompile(patzlist)
	var payloadpat = "(-z file,[/a-zA-A0-9.-_]*|-z File,[/a-zA-A0-9.-_]*|-z FILE,[/a-zA-A0-9.-_]*|-z list,[a-zA-Z0-9.-]*|-z List,[a-zA-Z0-9.-]*|-z LIST,[a-zA-Z0-9.-]*|-z range,[0-9-]*|-z Range,[0-9-]*|-z RANGE,[0-9-]*|-w [/0-9a-zA-Z._-]*)"
	payload := regexp.MustCompile(payloadpat)
	match := payload.FindAllString(str, -1)

	for N := 0; N < len(match); N++ {
		if zfile.MatchString(match[N]) {
			s.readfile(match[N][len("-z file,"):])
		}
		if zrange.MatchString(match[N]) {
			numRE := regexp.MustCompile("[0-9]+")
			numbs := numRE.FindAllString(match[N], -1)
			if len(numbs) != 0 {
				start, _ := strconv.Atoi(numbs[0])
				end, _ := strconv.Atoi(numbs[1])
				for i := start; i < end; i++ {
					zrangewordlist = append(zrangewordlist, strconv.Itoa(i))
				}
				s.Fuzzer.Wordlists = append(s.Fuzzer.Wordlists, zrangewordlist)
			}
		}
		if zlist.MatchString(match[N]) {
			zlistwordlist = strings.Split(match[N][len("-z list,"):], "-")
			if len(zlistwordlist) != 0 {
				s.Fuzzer.Wordlists = append(s.Fuzzer.Wordlists, zlistwordlist)
			}
		}
		if wfile.MatchString(match[N]) {
			s.readfile(match[N][len("-w "):])
		}
		// after payload for loop
	}
	FUZZre := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	FUZZs := FUZZre.FindAllString(str, -1)
	if len(FUZZs) != len(match) {
		return false
	}
	// setting up for GetURL
	if len(match) == 1 {
		s.Fuzzer.Indexes = append(s.Fuzzer.Indexes, 0) // just one payload/file
	} else {
		//
		s.Fuzzer.Indexes = append(s.Fuzzer.Indexes, len(s.Fuzzer.Wordlists))
		s.Fuzzer.Indexes = append(s.Fuzzer.Indexes, 0)
	}
	for _, i := range s.Fuzzer.Wordlists {
		s.Fuzzer.Maxes = append(s.Fuzzer.Maxes, len(i))
	}
	return true
}

// Validate : final input error checks before run.
func Validate(s *State, argstr, proxy string) bool {
	help := regexp.MustCompile("(-h|--help)")
	if help.MatchString(argstr) {
		return false
	}

	// parse output filter args
	f := regexp.MustCompile("--(hc|sc|hl|sl|hw|sw|hh|sh).[0-9a-zA-Z(,|)]*")
	RawFilterStringArgs := f.FindString(argstr)
	if len(RawFilterStringArgs) > 0 {
		printfilterlst := strings.Split(RawFilterStringArgs, " ")
		if len(printfilterlst) >= 2 {
			printfilterargs := printfilterlst[0]
			filternumlst := strings.Split(RawFilterStringArgs, " ")[1:]
			ParsePrintFilterArgs(s, printfilterargs)
			convPrintFilter(s, strings.Join(filternumlst, ","))
		}
	} else { // filter default
		ParsePrintFilterArgs(s, "sc")
		convPrintFilter(s, "200,301,302,403")

	}
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
	// TLS
	s.Client = &http.Client{
		Transport: &RedirectHandler{
			State: s,
			Transport: &http.Transport{
				Proxy: proxyURLFunc,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: s.InsecureSSL}},
		},
	}
	
	if len(s.Fuzzer.Wordlists) != 0 || ParseWordlistArgs(argstr, s) != false {
		return true
	}
	return false
}
