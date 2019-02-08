package libgrequest

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"

	// foreign
	"github.com/fatih/color"
)

// ArgBool : Turn commandline pat into true or false
func ArgBool(s, p string) bool {
	re := regexp.MustCompile(p)
	match := re.FindAllString(s, -1)
	if len(match) > 0 {
		return true
	}
	return false
}

// ArgString : Turn commandline pat into string
func ArgString(s, p string) string {
	re := regexp.MustCompile(p)
	match := re.FindAllString(s, -1)
	if len(match) > 0 {
		res := match[0]
		return strings.Trim(res, " ")
	}
	return ""
}

// ArgInt : Turn commandline pat into int
func ArgInt(s, p string) int {
	re := regexp.MustCompile(p)
	x := re.FindAllString(s, -1)
	numRE := regexp.MustCompile("[0-9]+")
	if len(x) == 0 {
		return 0
	}
	numbs := numRE.FindString(x[0])
	res, err := strconv.Atoi(numbs)
	if err != nil {
		return 666
	}
	return res
}

// PrepareSignalHandler : Signal handler straight from gobuster to catch CTRL+C
func PrepareSignalHandler(s *State) {
	s.SignalChan = make(chan os.Signal, 1)
	signal.Notify(s.SignalChan, os.Interrupt)
	go func() {
		for range s.SignalChan {
			// caught CTRL+C
			if !s.Quiet {
				fmt.Println("[!] Keyboard interrupt detected, terminating.")
				s.Terminate = true
			}
		}
	}()
}

func int2string(i int) string {
	t := strconv.Itoa(i)
	return t
}

func res2string(arg int64) string {
	return strconv.FormatInt(arg, 10)
}

// convPrintFilter :
func convPrintFilter(s *State, filternum string) {
	for _, c := range strings.Split(filternum, ",") {
		i, err := strconv.Atoi(c)
		i64 := int64(i)
		if err != nil {
			fmt.Println(err)
		} else {
			s.Filter.Add(i64)
		}
	}
}
func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func TotalRequests(maxes []int) int {
	c := 1
	for i := range maxes {
		c = maxes[i] * c
	}
	return c
}

// PrintTopColor : beginning of output
func PrintTopColor(s *State) {
	ye := color.New(color.FgYellow).SprintFunc()
	wordlists := strings.Join(s.WordListFiles, ", ")
	fmt.Printf("\n")
	fmt.Println("[+] Target: ", ye(s.URL))
	fmt.Println("[+] Wordlists: ", ye(wordlists))
	fmt.Printf("\n")
}

func PrintTopNoColor(s *State) {
	wordlists := strings.Join(s.WordListFiles, ", ")
	fmt.Printf("\n")
	fmt.Println("[+] Target: ", s.URL)
	fmt.Println("[+] Wordlists: ", wordlists)
	fmt.Printf("\n")
}

// PrintHelp : cli help info
func PrintHelp() {
	fmt.Printf("\n[+] gofuzz: dirty fork of gobuster to reproduce the functionality of wfuzz\n")
	fmt.Printf("[+] Author: brax (https://github.com/braaaax/gofuzz)\n")
	fmt.Printf("\nUsage:   gofuzz [options] <url>\n")
	fmt.Printf("Keyword: FUZZ, ..., FUZnZ  wherever you put these keywords gofuzz will replace them with the values of the specified payload.\n\n")
	fmt.Printf("Options:\n")
	fmt.Println("-h/--help                     : This help.")
	fmt.Println("-w wordlist                   : Specify a wordlist file (alias for -z file,wordlist).")
	fmt.Println("-z file/range/list,PAYLOAD    : Where PAYLOAD is FILENAME or 1-10 or \"-\" separated sequence.")
	fmt.Println("--hc/hl/hw/hh N[,N]+          : Hide responses with the specified code, lines, words, or chars.")
	fmt.Println("--sc/sl/sw/sh N[,N]]+         : Show responses with the specified code, lines, words, or chars.")
	fmt.Println("-t N                          : Specify the number of concurrent connections (10 default).")
	fmt.Println("-p URL                        : Specify proxy URL.")
	fmt.Println("-b COOKIE                     : Specify cookie.")
	fmt.Println("-ua USERAGENT                 : Specify user agent.")
	fmt.Println("--password PASSWORD           : Specify password for basic web auth.")
	fmt.Println("--username USERNAME           : Specify username.")
	fmt.Println("--no-follow                   : Don't follow HTTP redirections.")
	fmt.Println("--no-color                    : Monotone output.")
	fmt.Println("--print-body                  : Print response body to stdout.")
	fmt.Println("-k                            : Strict TLS connections (skip verify = false).")
	fmt.Println("-q                            : Quiet mode.")
	fmt.Printf("\n")
	fmt.Println("Examples: gofuzz -w users.txt -w pass.txt --sc 200 http://www.site.com/log.asp?user=FUZZ&pass=FUZ2Z")
	fmt.Println("          gofuzz -z file,default/common.txt -z list,-.php http://somesite.com/FUZZFUZ2Z")
	fmt.Println("          gofuzz -t 32 -w somelist.txt https://someTLSsite.com/FUZZ")
}
