package main

import (
	"os"
	"strings"

	"github.com/braaaax/grequest/libgrequest"
)

// ParseCmdLine : eval cmdline arguments
func ParseCmdLine(str string) *libgrequest.State {
	s := libgrequest.InitState()
	s.Cache = libgrequest.InitSafeCache()
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

	s.NoStatus = false
	s.IncludeLength = true
	s.WildcardForced = false
	s.UseSlash = true

	if threads > 0 {
		s.Threads = threads
	} else {
		s.Threads = 1
	}

	libgrequest.Validate(s, str, proxy)
	return s
}

func main() {
	// TESTING

	argstr := os.Args
	s := ParseCmdLine(strings.Join(argstr, " "))
	libgrequest.ProcessorII(s)

	return
}
