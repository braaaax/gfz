package libgrequest

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	//"fmt"
)

func (s *State) setPayload(str string) {
	var patpostform = "--post-form [^\t\n\f\r ]+"
	var patmultpart = "--post-multipart [^\t\n\f\r ]+"
	FUZZre := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	postform := regexp.MustCompile(patpostform)
	multpartform := regexp.MustCompile(patmultpart)

	if postform.MatchString(str) {
		postformstr := postform.FindString(str)[len("--post-form "):]
		if len(postformstr) > 0 {
			if FUZZre.MatchString(postformstr) {
				s.Post = true
				s.PostForm = true
				s.Payload = postformstr
			}
		}
		return
	}
	if multpartform.MatchString(str) {
		// f := mustOpen(multpartform.FindString(str)[len("--post-multipart "):])
		// b := []byte{}
		// somebytes, _ := f.Read(b)
		// s.Payload = string(somebytes)
		// if len(s.Payload) > 0 {
		s.Payload = multpartform.FindString(str)[len("--post-multipart "):]
		if FUZZre.MatchString(s.Payload) {
			s.Post = true
			s.PostMulti = true
		}
		// }
		// defer f.Close()
		return
	}
}

// ParseWordlistArgs : set UrlFuzz Wordlists FuzzMap
func ParseWordlistArgs(str string, s *State) bool {
	var patzfile = "-z (file|File|FILE),[^\t\n\f\r ]+"
	var patzrange = "-z (range|Range|RANGE),[0-9-]*" // put a limit
	var patzlist = "-z (list|List|LIST),[^\t\n\f\r ]+"
	var patwfile = "-w [^\t\n\f\r ]+"
	zlistwordlist := []string{}
	zrangewordlist := []string{}
	zfile := regexp.MustCompile(patzfile)
	zrange := regexp.MustCompile(patzrange)
	wfile := regexp.MustCompile(patwfile)
	zlist := regexp.MustCompile(patzlist)
	var payloadpat = "(-z file,[^\t\n\f\r ]+|-z File,[^\t\n\f\r ]+|-z FILE,[^\t\n\f\r ]+|-z list,[^\t\n\f\r ]+|-z List,[^\t\n\f\r ]+|-z LIST,[^\t\n\f\r ]+|-z range,[0-9-]*|-z Range,[0-9-]*|-z RANGE,[0-9-]*|-w [^\t\n\f\r ]+)"
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
		// if more than one FUZZ
		s.Fuzzer.Indexes = append(s.Fuzzer.Indexes, len(s.Fuzzer.Wordlists))
		for i := 1; i < len(s.Fuzzer.Wordlists); i++ {
			s.Fuzzer.Indexes = append(s.Fuzzer.Indexes, 0)
		}
	}
	for _, i := range s.Fuzzer.Wordlists {
		s.Fuzzer.Maxes = append(s.Fuzzer.Maxes, len(i))
	}
	/*fmt.Println("s.Fuzzer.Indexes", s.Fuzzer.Indexes)
	fmt.Println("s.Fuzzer.Wordlists", s.Fuzzer.Maxes)
	fmt.Println("s.Fuzzer.Maxes", s.Fuzzer.Maxes)*/
	if len(s.Fuzzer.Wordlists) != 0 {
		return true
	}
	return false
}

// Validate : final input error checks before run.
func Validate(s *State, argstr, proxy string) bool {
	help := regexp.MustCompile("(--help)")
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
	// set proxy info TODO parse urls at the cmdline
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

	s.setPayload(argstr)
	if s.PostForm {
		s.Request = GoPostForm
	} else if s.PostMulti {
		s.Request = GoPostMultiPart
	} else {
		if s.Post {s.Method = "POST"} else {s.Method = "GET"}
		s.Request = GoGet
	}

	if len(s.Fuzzer.Wordlists) != 0 || ParseWordlistArgs(argstr, s) != false {
		return true
	}
	return false
}
