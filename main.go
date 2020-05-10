/*replicate the basic functionalities of wfuzz*/
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/braaaax/gfuzz/libgrequest"
)

// ParseCmdLine : eval cmdline arguments
func ParseCmdLine(str string) *libgrequest.State {
	s := libgrequest.InitState()
	s.Counter = libgrequest.InitSafeCounter()
	s.Fuzzer = libgrequest.InitFuzz()
	s.Commandline = str
	FUZZre := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	s.FollowRedirect = !libgrequest.ArgBool(str, "--no-follow")   // !lazydefault1: follow redirect
	s.URL = libgrequest.ArgString(str, "htt(p|ps)[^\t\n\f\r ]+$") // htt(p|ps).[/a-zA-Z0-9:.]*
	if len(FUZZre.FindAllString(s.URL, -1)) > 0 {
		s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, true)
	} else {
		s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, false)
	}
	s.InsecureSSL = !libgrequest.ArgBool(str, "-k") // !lazydefault2: skip verify
	threads := libgrequest.ArgInt(str, "-t.[0-9]*")
	proxy := libgrequest.ArgString(str, "-p [^\t\n\f\r ]+") // TODO -p.htt(p|ps).[/a-zA-Z0-9:]*
	if len(proxy) > len("-p ") {
		proxy = proxy[len("-p "):]
		// proxy = "http://" + proxy
	}
	s.Quiet = libgrequest.ArgBool(str, "-q")
	s.Cookies = libgrequest.ArgString(str, "-b [^\t\n\f\r ]+")
	if len(s.Cookies) > len("-b ") {
		s.Cookies = s.Cookies[len("-b "):]
		if len(FUZZre.FindAllString(s.Cookies, -1)) > 0 {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, true)
		} else {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, false)
		}
	}
	s.Password = libgrequest.ArgString(str, "--password [^\t\n\f\r ]+")
	if len(s.Password) > len("--password ") {
		s.Password = s.Password[len("--password "):]
		if len(FUZZre.FindAllString(s.Password, -1)) > 0 {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, true)
		} else {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, false)
		}
	}
	s.Username = libgrequest.ArgString(str, "--username [^\t\n\f\r ]+")
	if len(s.Username) > len("--username.") {
		s.Username = s.Username[len("--username."):]
		if len(FUZZre.FindAllString(s.Username, -1)) > 0 {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, true)
		} else {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, false)
		}
	}
	s.UserAgent = libgrequest.ArgString(str, "-ua [^\t\n\f\r ]+")
	if len(s.UserAgent) > len("-ua.") {
		s.UserAgent = s.UserAgent[len("-ua."):]
		if len(FUZZre.FindAllString(s.UserAgent, -1)) > 0 {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, true)
		} else {
			s.Fuzzer.Cmdline = append(s.Fuzzer.Cmdline, false)
		}
	}
	s.Headers = make(map[string]string)
	reHeaders := regexp.MustCompile("-H [^\t\n\f\r ]+")
	match := reHeaders.FindAllString(str, -1)
	// fmt.Println(match)
	if len(match) > 0 {
		for i := range match {
			res := strings.Split(match[i], ":")
			s.Headers[strings.Replace(res[0][len("H "):], " ", "", -1)] = strings.Replace(res[1], " ", "", -1)
		}
	}

	s.NoColor = libgrequest.ArgBool(str, "--no-color")
	s.PrintBody = libgrequest.ArgBool(str, "--print-body")
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
		libgrequest.PrintHelp()
	}
}
