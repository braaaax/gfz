package libgrequest

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ParseResponse : process http response pointer
func ParseResponse(fullURL string, resp *http.Response) (*Result, error) {
	//set body
	var r = &Result{URL: fullURL, Code: int64(resp.StatusCode)}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	r.Body = body
	sbody := string(body)
	if err == nil {
		r.Chars = int64(utf8.RuneCountInString(sbody))
		r.Words = int64(len(strings.Fields(sbody)))
		newlineRE := regexp.MustCompile("\n")
		r.Lines = int64(len(newlineRE.FindAllString(sbody, -1)))
	} else {
		return r, err
	}
	return r, nil
}

func parseurl(uarg string) string {
	u, err := url.Parse(uarg)
	if err != nil {
		panic(err)
	}
	return u.Host + u.Path
}

// PrintChars :
// probably a better way to do this
func PrintChars(s *State, r *Result) {
	if s.Filter.Contains(r.Chars) == s.Show {
		if s.NoColor {
			PrintNoColorFn(s, r)
		}
		PrintColorFn(s, r)
	}
}

// PrintWords :
func PrintWords(s *State, r *Result) {
	if s.Filter.Contains(r.Words) == s.Show {
		if s.NoColor {
			PrintNoColorFn(s, r)
		}
		PrintColorFn(s, r)
	}
}

// PrintStatus :
func PrintStatus(s *State, r *Result) {
	if s.Filter.Contains(r.Code) == s.Show { // issue nil
		if s.NoColor {
			PrintNoColorFn(s, r)
		}
		PrintColorFn(s, r)
	}
}

// PrintLines :
func PrintLines(s *State, r *Result) {
	if s.Filter.Contains(r.Lines) == s.Show {
		if s.NoColor {
			PrintNoColorFn(s, r)
		}
		PrintColorFn(s, r)
	}
}

// ParsePrintFilterArgs :
func ParsePrintFilterArgs(s *State, fs string) {
	m := regexp.MustCompile("(sl|sc|sw|sh|hc|hl|hh|hw)").FindString(fs)
	if string(m[0]) == "s" {
		s.Show = true
	} else {
		s.Show = false
	}
	switch m[1:] {
	case "c":
		s.Printer = PrintStatus
	case "l":
		s.Printer = PrintLines
	case "w":
		s.Printer = PrintWords
	case "h":
		s.Printer = PrintChars
	}
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
			if s.readfile(match[N][len("-z file,"):]) != nil {
				return false
			}
			s.WordListFiles = append(s.WordListFiles, match[N][len("-z file,"):])
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
				s.WordListFiles = append(s.WordListFiles, match[N][len("-z range,"):])
			}
		}
		if zlist.MatchString(match[N]) {
			zlistwordlist = strings.Split(match[N][len("-z list,"):], "-")
			if len(zlistwordlist) != 0 {
				s.Fuzzer.Wordlists = append(s.Fuzzer.Wordlists, zlistwordlist)
				s.WordListFiles = append(s.WordListFiles, match[N][len("-z list,"):])
			}
		}
		if wfile.MatchString(match[N]) {
			if s.readfile(match[N][len("-w "):]) != nil {
				return false
			}
			s.WordListFiles = append(s.WordListFiles, match[N][len("-w "):])
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
