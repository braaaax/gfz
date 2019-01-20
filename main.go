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
	s.Counter = libgrequest.InitSafeCounter()
	s.Fuzzer = libgrequest.InitFuzz()
	s.FollowRedirect = !libgrequest.ArgBool(str, "--no-follow") // !lazydefault1: follow redirect
	s.URL = libgrequest.ArgString(str, " htt(p|ps).+")
	s.InsecureSSL = !libgrequest.ArgBool(str, "-k") // !lazydefault2: skip verify
	threads := libgrequest.ArgInt(str, "-t.[0-9]*")
	proxy := libgrequest.ArgString(str, "-p.htt(p|ps)://(.)*")
	s.Quiet = libgrequest.ArgBool(str, "-q")
	s.Cookies = libgrequest.ArgString(str, "-b.[a-zA-Z0-9=/?]*")
	s.Password = libgrequest.ArgString(str, "--password.[a-zA-Z0-9=/?]*")
	s.Username = libgrequest.ArgString(str, "--username.[a-zA-Z0-9=/?]*")
	s.UserAgent = libgrequest.ArgString(str, "-ua.[a-zA-Z]+")
	s.NoColor = libgrequest.ArgBool(str, "--no-color")
	s.OutputFileName = libgrequest.ArgString(str, "-o")

	// s.Recursive = libgrequest.ArgBool(str, "-r")
	// s.Mode = libgrequest.ArgString(str, "-m")

	if threads > 0 {
		s.Threads = threads
	} else {
		s.Threads = 10
	}

	libgrequest.Validate(s, str, proxy)
	return s
}

func main() {
	argstr := os.Args
	s := ParseCmdLine(strings.Join(argstr, " "))
	if len(s.URL) != 0 && libgrequest.IsMapFull(s.Fuzzer.Fuzzmap) {
		libgrequest.PrintTop(s)
		Code, _ := libgrequest.GoGet(s, s.URL, s.Cookies)
		if Code == nil {
			fmt.Printf("Cannot reach %s", s.URL)
			return
		}
		s.Fuzzer.Wordlists = s.SetWordlists()
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
