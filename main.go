/*replicate the basic functionalities of wfuzz in go*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/braaaax/gofuzz/libgrequest"
)

// ParseCmdLine : eval cmdline arguments
func ParseCmdLine(str string) *libgrequest.State {
	s := libgrequest.InitState()
	// s.Cache = libgrequest.InitSafeCache()
	s.Fuzzer = libgrequest.InitFuzz()
	s.Cookies = libgrequest.ArgString(str, "-b.[a-zA-Z0-9=/?]*")
	s.Mode = libgrequest.ArgString(str, "-m")
	s.OutputFileName = libgrequest.ArgString(str, "-o")
	s.FollowRedirect = libgrequest.ArgBool(str, "--follow")
	proxy := libgrequest.ArgString(str, "-p.htt(p|ps)://(.)*")
	s.URL = libgrequest.ArgString(str, ".htt(p|ps).+")
	s.UserAgent = libgrequest.ArgString(str, "-ua.[a-zA-Z]+")
	s.InsecureSSL = libgrequest.ArgBool(str, "-k")
	s.Quiet = libgrequest.ArgBool(str, "-q")
	s.Recursive = libgrequest.ArgBool(str, "-r")
	threads := libgrequest.ArgInt(str, "-t.[0-9]*")

	if threads > 0 {
		s.Threads = threads
	} else {
		s.Threads = 10
	}

	libgrequest.Validate(s, str, proxy)
	return s
}

func main() {
	// TESTING
	argstr := os.Args
	s := ParseCmdLine(strings.Join(argstr, " "))
	if len(s.URL) != 0 && libgrequest.IsMapFull(s.Fuzzer.Fuzzmap) {
		libgrequest.PrintTop(s)
		Code, _ := libgrequest.GoGet(s, s.URL, s.Cookies)
		if Code == nil {
			fmt.Printf("Cannot reach %s", s.URL)
			return
		}
		s.Fuzzer.Wordlists = s.SWordlists()
		if len(s.WordListFiles) != len(s.Fuzzer.Fuzzmap) {
			libgrequest.PrintHelp()
			return
		}

	} else {
		libgrequest.PrintHelp()
		return
	}
	libgrequest.Processor(s)
	return
}
