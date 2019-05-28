/*replicate the basic functionalities of wfuzz*/
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/braaaax/gofuzz/libgrequest"
)

// ParseCmdLine : eval cmdline arguments
func ParseCmdLine(str string) *libgrequest.State {
	s := libgrequest.InitState()
	s.Counter = libgrequest.InitSafeCounter()
	s.Fuzzer = libgrequest.InitFuzz()
	s.Commandline =str
	// s.Fuzzer.Cmdline = make([]bool, 5)

	FUZZre := regexp.MustCompile("FUZ(Z|[0-9]Z)")

	s.FollowRedirect = !libgrequest.ArgBool(str, "--no-follow") // !lazydefault1: follow redirect
	s.URL = libgrequest.ArgString(str, "htt(p|ps)[^\t\n\f\r ]+$")
	// fmt.Println("URL:", s.URL)
	if len(FUZZre.FindAllString(s.URL, -1)) > 0 {
		s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,true)
	} else {
		s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,false)
	}

	s.InsecureSSL = !libgrequest.ArgBool(str, "-k") // !lazydefault2: skip verify
	threads := libgrequest.ArgInt(str, "-t.[0-9]*")
	proxy := libgrequest.ArgString(str, "-p htt(p|ps).[^\t\n\f\r ]+") // TODO
	if len(proxy) > len("-p "){
		proxy = proxy[len("-p "):]
	}
	// fmt.Println("proxy:", proxy)
	s.Quiet = libgrequest.ArgBool(str, "-q")
	s.Cookies = libgrequest.ArgString(str, "-b [^\t\n\f\r ]+")
	if len(s.Cookies) > len("-b ") {
		s.Cookies = s.Cookies[len("-b "):]
		if len(FUZZre.FindAllString(s.Cookies, -1)) > 0 {
			// fmt.Println("Cookies FUZZ")
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,true)
		} else {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,false)
		}
	}
	s.Password = libgrequest.ArgString(str, "--password [^\t\n\f\r ]+")
	if len(s.Password) > len("--password ") {
		s.Password = s.Password[len("--password "):]
		if len(FUZZre.FindAllString(s.Password, -1)) > 0 {
			// fmt.Println("password FUZZ")
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,true)
		} else {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,false)
		}
	}
	s.Username = libgrequest.ArgString(str, "--username [^\t\n\f\r ]+")
	if len(s.Username) > len("--username.") {
		s.Username = s.Username[len("--username."):]
		if len(FUZZre.FindAllString(s.Username, -1)) > 0 {
			// fmt.Println("username FUZZ")
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,true)
		} else {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,false)
		}
	}
	s.UserAgent = libgrequest.ArgString(str, "-ua [^\t\n\f\r ]+")
	if len(s.UserAgent) > len("-ua.") {
		s.UserAgent = s.UserAgent[len("-ua."):]
		if len(FUZZre.FindAllString(s.UserAgent, -1)) > 0 {
			// fmt.Println("UserAgent FUZZ")
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,true)
		} else {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline,false)
		}
	}
	// fmt.Println("CMDLINE:", s.Fuzzer.Cmdline)
	s.NoColor = libgrequest.ArgBool(str, "--no-color")
	s.PrintBody = libgrequest.ArgBool(str, "--print-body")
	// s.Recursive = libgrequest.ArgBool(str, "-r") TODO
	s.Post = libgrequest.ArgBool(str, "--post")
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
		fmt.Printf("\n[+] Time elapsed: %s\n", elapsed)
	} else {
		fmt.Println("s:", s, "args:", os.Args)
		libgrequest.PrintHelp()
	}
}
