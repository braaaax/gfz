package libgrequest

import (
	"net/http"
	"net/url"
	"os"
	"sync"
)

// MethodProc :
//type MethodProc func(*State, string, string) (*int, *Result)

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
	NoColor        bool
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
	Counter        *SafeCounter
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

// Fuzz : struct to store info for GetUrl
type Fuzz struct {
	Wordlists [][]string
	Indexes   []int
	Maxes     []int
	Fuzzmap   map[string]string
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
