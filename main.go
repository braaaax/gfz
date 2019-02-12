/*replicate the basic functionalities of wfuzz in go*/
package main

import (
	"os"
	"strings"
	"time"
	"fmt"

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
	s.PrintBody = libgrequest.ArgBool(str, "--print-body")
	// s.Recursive = libgrequest.ArgBool(str, "-r")
	s.PostForm = libgrequest.ArgBool(str, "--post-form")
	s.PostMulti = libgrequest.ArgBool(str, "--post-multipart")

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
	// start := time.Now()
	argstr := os.Args
	s := ParseCmdLine(strings.Join(argstr, " "))
	if s != nil && len(os.Args) > 1 {
		if s.NoColor {
			libgrequest.PrintTopNoColor(s)
		} else {
			libgrequest.PrintTopColor(s)
		}
		start := time.Now()
		libgrequest.Processor(s)
		elapsed := time.Since(start)
		fmt.Printf("\nRequests took: %s\n", elapsed)
	} else {
		libgrequest.PrintHelp()
	}
}
