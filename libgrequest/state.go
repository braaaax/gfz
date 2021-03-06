/* A few of the structs below are heavilly influenced by OJ's code in gobuster */

package libgrequest

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
)

// RequestFunc :
type RequestFunc func(s *State, url, cookie, payload string) (*int, error)

// PrintResultFunc : abstraction layer to handle printing different filters
type PrintResultFunc func(s *State, r *Result)

// State :
type State struct {
	Client         *http.Client
	UserAgent      string
	FollowRedirect bool
	Username       string
	Password       string
	IncludeLength  bool
	URL            string
	Cookies        string
	StatusCodes    IntSet
	WildcardForced bool
	UseSlash       bool
	IsWildcard     bool
	ProxyURL       *url.URL
	NoStatus       bool
	NoColor        bool
	OutputFile     *os.File
	OutputFileName string
	InsecureSSL    bool
	Payload        string
	Post           bool
	PostForm       bool
	PostMulti      bool
	Request        RequestFunc
	WildcardIps    StringSet
	Show           bool
	Printer        PrintResultFunc
	Filter         IntSet
	SignalChan     chan os.Signal
	WordListFiles  []string
	Quiet          bool
	PrintBody      bool
	Terminate      bool
	Threads        int
	Counter        *SafeCounter
	FUZZs          []string
	Fuzzer         *Fuzz
	Method         string
	Commandline    string
	Headers        map[string]string
}

func (s *State) readfile(fname string) error {
	fn, err := os.Open(fname)
	if err != nil {
		fmt.Println("File not found.")
		return err
	}
	defer fn.Close()
	var lines []string
	scanner := bufio.NewScanner(fn)
	for scanner.Scan() {
		w := scanner.Text()
		lines = append(lines, w)
		//filt := regexp.MustCompile("[!@#%]")
		//if !filt.MatchString(w){
		//	lines = append(lines, w)
		//}
	}
	s.Fuzzer.Wordlists = append(s.Fuzzer.Wordlists, lines)
	return nil
}

// InitState :
func InitState() *State {
	return &State{
		StatusCodes:    IntSet{Set: map[int64]bool{}},
		WildcardIps:    StringSet{Set: map[string]bool{}},
		Filter:         IntSet{Set: map[int64]bool{}},
		NoStatus:       false,
		IsWildcard:     false,
		IncludeLength:  false,
		WildcardForced: false,
		UseSlash:       false,
		Post:           false,
		PostForm:       false,
		PostMulti:      false,
		Payload:        "",
	}
}

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	v   int
	mux sync.Mutex
}

// Inc : Increment v.
func (c *SafeCounter) Inc() {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the count
	c.v++
	c.mux.Unlock()
}

// InitSafeCounter : Return intialized SafeCounter struct pointer.
func InitSafeCounter() *SafeCounter {
	return &SafeCounter{v: 0, mux: sync.Mutex{}}
}

// Fuzz : struct to store info for GetUrl
type Fuzz struct {
	Wordlists [][]string
	Indexes   []int
	Maxes     []int
	Cmdline   []bool
}

// InitFuzz : init the Fuzz struct.
func InitFuzz() *Fuzz {
	return &Fuzz{}
}

// Result struct
type Result struct {
	URL   string
	Body  []byte
	Chars int64
	Words int64
	Lines int64
	Code  int64
}

// InitResult : process http response pointer
func InitResult(fullURL string, resp *http.Response) (*Result, error) {
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

// IntSet : Set value maps int64 to bool (from gobuster).
type IntSet struct {
	Set map[int64]bool
}

// StringSet : Set value maps string to bool (from gobuster).
type StringSet struct {
	Set map[string]bool
}

// Contains : Contains int.
func (set *IntSet) Contains(i int64) bool {
	_, found := set.Set[i]
	return found
}

// Add : Add int.
func (set *IntSet) Add(i int64) bool {
	_, found := set.Set[i]
	set.Set[i] = true
	return !found
}

// Contains : Contains string.
func (set *StringSet) Contains(s string) bool {
	_, found := set.Set[s]
	return found
}

// Add : Add string.
func (set *StringSet) Add(s string) bool {
	_, found := set.Set[s]
	set.Set[s] = true
	return !found
}
