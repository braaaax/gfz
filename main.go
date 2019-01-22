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
	if libgrequest.Validate(s, str, proxy) != true {
		return nil
	}
	return s
}

func main() {
	argstr := os.Args
	s := ParseCmdLine(strings.Join(argstr, " "))
	if s != nil && len(os.Args) < 1 {
		// fmt.Printf("State: %+v\n", s)
		libgrequest.PrintTop(s)
		libgrequest.Processor(s)
	} else {
		libgrequest.PrintHelp()
	}
}
