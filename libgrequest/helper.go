package libgrequest

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
)

// IntSet :
type IntSet struct {
	Set map[int64]bool
}

// StringSet :
type StringSet struct {
	Set map[string]bool
}

// Contains :
func (set *IntSet) Contains(i int64) bool {
	_, found := set.Set[i]
	return found
}

// Add :
func (set *IntSet) Add(i int64) bool {
	_, found := set.Set[i]
	set.Set[i] = true
	return !found
}

// Contains :
func (set *StringSet) Contains(s string) bool {
	_, found := set.Set[s]
	return found
}

// Add :
func (set *StringSet) Add(s string) bool {
	_, found := set.Set[s]
	set.Set[s] = true
	return !found
}

// PackFilter :
func PackFilter(s *State, filternum string) {
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

// BaseFuzzMap : map 'FUZZ' to "" to help replace function
func BaseFuzzMap(s *State) { //map[string]string {
	m := make(map[string]string)
	for _, val := range s.FuzzMap {
		m[val] = ""
	}
	s.BaseMap = m
}

// IsMapFull : checks whether there is a value for every key in map
func IsMapFull(fm map[string]string) bool {
	/* check whether there are keys without values */
	var r bool
	for k := range fm {
		if fm[k] != "" {
			r = true
		} else {
			r = false
		}
	}
	return r
}

// ArgBool : turn commandline pat into true or false
func ArgBool(s, p string) bool {
	re := regexp.MustCompile(p)
	match := re.FindAllString(s, -1)
	if len(match) > 0 {
		return true
	}
	return false
}

// ArgString : turn commandline pat into string
func ArgString(s, p string) string {
	re := regexp.MustCompile(p)
	match := re.FindAllString(s, -1)
	if len(match) > 0 {
		res := match[0]
		return strings.Trim(res, " ")
	}
	return ""
}

// ArgInt : turn commandline pat into int
func ArgInt(s, p string) int {
	re := regexp.MustCompile(p)
	x := re.FindAllString(s, -1)
	numRE := regexp.MustCompile("[0-9]+")
	if len(x) == 0 {
		return 0
	}
	numbs := numRE.FindString(x[0])
	res, err := strconv.Atoi(numbs)
	if err != nil {
		return 666
	}
	return res
}

// ArgArray : turn commandline pat into array
func ArgArray(s, p string) []string {
	re := regexp.MustCompile(p)
	match := re.FindAllString(s, -1)
	return match
}

// FuzzPrintChars :
func FuzzPrintChars(s *State, r *Result) {
	if s.Filter.Contains(r.Chars) == s.ShowHide {
		PrintFn(s, r)
	}
}

// FuzzPrintWords :
func FuzzPrintWords(s *State, r *Result) {
	if s.Filter.Contains(r.Words) == s.ShowHide {
		PrintFn(s, r)
	}
}

// FuzzPrintStatus :
func FuzzPrintStatus(s *State, r *Result) {
	if s.Filter.Contains(r.Code) == s.ShowHide {
		PrintFn(s, r)
	}
}

// FuzzPrintLines :
func FuzzPrintLines(s *State, r *Result) {
	if s.Filter.Contains(r.Lines) == s.ShowHide {
		PrintFn(s, r)
	}
}

// PrintFilter : switch for print filter
func PrintFilter(s *State, fs string) {
	fmt.Println("fs", fs)
	m := regexp.MustCompile("(sl|sc|sw|sh|hc|hl|hh|hw)").FindString(fs)
	if string(m[0]) == "s" {
		s.ShowHide = true
	} else {
		s.ShowHide = false
	}
	switch m[1:] {
	case "c":
		s.Printer = FuzzPrintStatus
	case "l":
		s.Printer = FuzzPrintLines
	case "w":
		s.Printer = FuzzPrintWords
	case "h":
		s.Printer = FuzzPrintChars
	}
}

// FuzzMapper : set Fuzzy.UrlFuzz Fuzzy.Wordlists Fuzzy.FuzzMap
func FuzzMapper(str string, s *State) {
	var patZ = "-z.(file|list),[/a-zA-A0-9.-]*"
	var patW = "-w.[/0-9a-zA-Z._-]*"
	var patFuzz = "FUZ(Z|[0-9]Z)"
	var patURL = "htt(p|ps)://(.)*"
	var wordlists []string

	zRE := regexp.MustCompile(patZ)
	wRE := regexp.MustCompile(patW)
	fuzzRE := regexp.MustCompile(patFuzz)
	urlRE := regexp.MustCompile(patURL)

	if zRE.MatchString(str) {
		zlist := ArgArray(str, patZ)
		for i := 0; i < len(zlist); i++ {
			wordlists = append(wordlists, zlist[i][len("-z file,"):])
		}
	}
	if wRE.MatchString(str) {
		wlist := ArgArray(str, patW)
		for i := 0; i < len(wlist); i++ {
			wordlists = append(wordlists, wlist[i][len("-w "):])
		}
	}

	URLs := urlRE.FindAllString(str, -1)
	FUZZs := fuzzRE.FindAllString(URLs[0], -1)

	fm := make(map[string]string)
	for index, m := range FUZZs {
		fm[wordlists[index]] = m
	}
	s.URLFuzz = URLs[0]
	s.Wordlists = wordlists
	s.FuzzMap = fm
	BaseFuzzMap(s)
}

// PrepareSignalHandler : Signal handler straight from gobuster to catch CTRL+C
func PrepareSignalHandler(s *State) {
	s.SignalChan = make(chan os.Signal, 1)
	signal.Notify(s.SignalChan, os.Interrupt)
	go func() {
		for range s.SignalChan { // for _ := range
			// caught CTRL+C
			if !s.Quiet {
				fmt.Println("[!] Keyboard interrupt detected, terminating.")
				s.Terminate = true
			}
		}
	}()
}

// PrintTop : beginning of output
func PrintTop(s *State) {
	fmt.Println("Target: ", s.URLFuzz)
	fmt.Println("Wordlists: ", strings.Join(s.Wordlists, ", "))
	fmt.Println("=============================================================================================================")
	fmt.Println("URL                                                 Status        Chars          Words          Lines")
	fmt.Println("=============================================================================================================")
}

// PrintHelp :
func PrintHelp() {
	//PrintBanner()
	fmt.Printf("Usage: ./gfuzz.py [options] -w wordlist <url>\n\n")
	fmt.Printf("Options:\n")
	fmt.Println("-h/--help                     : This help")
	fmt.Println("--version                     : Wfuzz version details")
	fmt.Println("-p addr                       : Use Proxy in format http//ip:port")
	fmt.Println("-t N                          : Specify the number of concurrent connections (10 default)")
	fmt.Println("--follow                      : Follow HTTP redirections")
	fmt.Println("-w wordlist                   : Specify a wordlist file (alias for -z file,wordlist).")
	fmt.Println("-b cookie                     : Specify a cookie for the requests")
	fmt.Println("--hc/hl/hw/hh N[,N]+          : Hide responses with the specified code/lines/words/chars")
	fmt.Printf("--sc/sl/sw/sh N[,N]]+         : Show responses with the specified code/lines/words/chars\n")
	fmt.Println("Keyword: FUZZ, ..., FUZnZ  wherever you put these keywords wfuzz will replace them with the values of the specified payload.")
	fmt.Println("Example: - gfuzz -w users.txt -w pass.txt --sc 200 http://www.site.com/log.asp?user=FUZZ&pass=FUZ2Z")
}
