package libgrequest

import (
	"bufio"
	"net/http"
	"net/url"
	"os"
	"sync"
)

// MethodProc :
type MethodProc func(*State, string, string) (*int, *Result)

// ProcFunc :
type ProcFunc func(*State, string, chan<- Result)

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
	OutputFile     *os.File
	OutputFileName string
	InsecureSSL    bool
	Mode           string
	Payload        string
	Processor      ProcFunc
	WildcardIps    StringSet
	Show           bool
	Printer        PrintResultFunc
	Filter         IntSet
	SignalChan     chan os.Signal
	WordListFiles  []string
	Quiet          bool
	Recursive      bool
	Terminate      bool
	Threads        int
	Counter          *SafeCounter
	FUZZs          []string
	Fuzzer         *Fuzz
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
	}
}

type Fuzz struct {
	Wordlists [][]string
	Indexes   []int
	Maxes     []int
	Fuzzmap   map[string]string
}

func InitFuzz() *Fuzz {
	return &Fuzz{}
}

/*
func (f *Fuzz) SetWordlist() [][]string {
	wordlists := [][]string{}
	var scanner *bufio.Scanner
	for fn := range f.Fuzzmap {
		wordlist, err := os.Open(fn)
		check(err)
		defer wordlist.Close()
		scanner = bufio.NewScanner(wordlist)
		scanner.Split(bufio.ScanWords)
		var words []string
		for scanner.Scan() {
			words = append(words, scanner.Text())

		}
		wordlists = append(wordlists, words)
	}
	// setting up for rloop
	f.Indexes = append(f.Indexes, len(wordlists))
	f.Indexes = append(f.Indexes, 0)
	for _, i := range wordlists {
		f.Maxes = append(f.Maxes, len(i))
	}

	return wordlists
}
*/

func (s *State) SWordlists() [][]string {
	wordlists := [][]string{}
	for _, filename := range s.WordListFiles {
		fn, err := os.Open(filename)
		check(err)
		defer fn.Close()
		var lines []string
		scanner := bufio.NewScanner(fn)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		wordlists = append(wordlists, lines)

	}
	// setting up for rloop
	if len(wordlists) == 1 {
		s.Fuzzer.Indexes = append(s.Fuzzer.Indexes, 0)
	} else {
		s.Fuzzer.Indexes = append(s.Fuzzer.Indexes, len(wordlists))
		s.Fuzzer.Indexes = append(s.Fuzzer.Indexes, 0)
	}
	for _, i := range wordlists {
		s.Fuzzer.Maxes = append(s.Fuzzer.Maxes, len(i))
	}
	return wordlists
}

// Result :
type Result struct {
	URL   string
	Body  []byte
	Chars int64
	Words int64
	Lines int64
	Code  int64
}

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	v   int
	mux sync.Mutex
}

// Inc :
func (c *SafeCounter) Inc() {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the count
	c.v++
	c.mux.Unlock()
}

// InitSafeCounter :
func InitSafeCounter() *SafeCounter {
	return &SafeCounter{v: 0, mux: sync.Mutex{}}
}
