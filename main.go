package main

import (
	//"fmt"
	"os"
	"strings"

	"github.com/braaaax/grequest/libgrequest"
)

func ParseCmdLine(str string) *libgrequest.State {
	// fmt.Println(str)
	s := libgrequest.InitState()
	s.Cookies = libgrequest.ArgString(str, "-b.[a-zA-Z0-9=/?]*")
	s.Mode = libgrequest.ArgString(str, "-m")
	s.OutputFileName = libgrequest.ArgString(str, "-o")
	s.FollowRedirect = libgrequest.ArgBool(str, "--follow")
	proxy := libgrequest.ArgString(str, "-p.htt(p|ps)://(.)*")
	s.Url = libgrequest.ArgString(str, ".htt(p|ps).+")
	//fmt.Println(s.Url)
	s.UserAgent = libgrequest.ArgString(str, "-ua.[a-zA-Z]+")
	s.InsecureSSL = libgrequest.ArgBool(str, "-k")
	s.Quiet = libgrequest.ArgBool(str, "-q")
	s.Recursive = libgrequest.ArgBool(str, "-r")
	s.Verbose = false
	s.NoStatus = false
	s.IncludeLength = true
	s.WildcardForced = false
	s.UseSlash = true

	threads := libgrequest.ArgInt(str, "-t.[0-9]*")
	if threads > 0 {
		s.Threads = threads
	} else {
		s.Threads = 10
	}
	//fmt.Printf("%+v", s)
	libgrequest.Validate(s, str, proxy)
	return s
}

func main() {
	argstr := os.Args
	if len(argstr) == 1 {
		libgrequest.PrintHelp()
		return
	}
	sp := ParseCmdLine(strings.Join(argstr, " "))
	//fmt.Printf("%+v", sp)
	libgrequest.PrintTop(sp)
	libgrequest.FuzzProc2(sp)
}
