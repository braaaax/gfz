package libgrequest

import (
	"net/http"
	"net/url"
	"os"
)

// MethodProc :
type MethodProc func(*State, string, string) (*int, *Result)

// ProcFunc :
type ProcFunc func(*State, string, chan<- Result)

// PrintResultFunc : abstraction layer to handle printing different filters
type PrintResultFunc func(s *State, r *Result)

// State :
type State struct {
	Client          *http.Client
	UserAgent       string
	FollowRedirect  bool
	Username        string
	Password        string
	IncludeLength   bool
	Url             string
	Cookies         string
	StatusCodes     IntSet
	WildcardForced  bool
	UseSlash        bool
	IsWildcard      bool
	ProxyURL        *url.URL
	Extensions      []string
	Verbose         bool
	Expanded        bool
	NoStatus        bool
	OutputFile      *os.File
	OutputFileName  string
	InsecureSSL     bool
	Mode            string
	Payload         string
	MethodProcessor MethodProc
	Processor       ProcFunc
	WildcardIps     StringSet
	StdIn           bool
	ShowHide        bool
	Printer         PrintResultFunc
	Filter          IntSet
	FuzzMap         map[string]string
	SignalChan      chan os.Signal
	Wordlists       []string
	Quiet           bool
	Terminate       bool
	URLFuzz         string
	BaseMap         map[string]string
	Threads         int
}

// Result :
type Result struct {
	Url   string
	Body  []byte
	Chars int64
	Words int64
	Lines int64
	Code  int64
}

// InitState :
func InitState() *State {
	return &State{
		StatusCodes: IntSet{Set: map[int64]bool{}},
		WildcardIps: StringSet{Set: map[string]bool{}},
		Filter:      IntSet{Set: map[int64]bool{}},
		IsWildcard:  false,
		StdIn:       false,
	}
}
