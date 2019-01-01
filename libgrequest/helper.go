package libgrequest

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	//"github.com/fatih/color"
)

// helper functions

func check(e error) {
	if e != nil {
		panic("failed")
	}
}

// IntSet :
type IntSet struct {
	Set map[int64]bool
}

// StringSet :
type StringSet struct {
	Set map[string]bool
}

// Contains :
func (set *IntSet) Contains(i int64) bool {
	_, found := set.Set[i]
	return found
}

// Add :
func (set *IntSet) Add(i int64) bool {
	_, found := set.Set[i]
	set.Set[i] = true
	return !found
}

// Contains :
func (set *StringSet) Contains(s string) bool {
	_, found := set.Set[s]
	return found
}

// Add :
func (set *StringSet) Add(s string) bool {
	_, found := set.Set[s]
	set.Set[s] = true
	return !found
}

// IsMapFull : 
func IsMapFull(fm map[string]string) bool {
	/* check whether there are keys without values */
	var r bool
	for k := range fm {
		if fm[k] != "" {
			r = true
		} else {
			r = false
		}
	}
	return r
}

// ArgBool : turn commandline pat into true or false
func ArgBool(s, p string) bool {
	re := regexp.MustCompile(p)
	match := re.FindAllString(s, -1)
	if len(match) > 0 {
		return true
	}
	return false
}

// ArgString : turn commandline pat into string
func ArgString(s, p string) string {
	re := regexp.MustCompile(p)
	match := re.FindAllString(s, -1)
	if len(match) > 0 {
		res := match[0]
		return strings.Trim(res, " ")
	}
	return ""
}

// ArgInt : turn commandline pat into int
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

// ArgArray : turn commandline pat into array
func ArgArray(s, p string) []string {
	re := regexp.MustCompile(p)
	match := re.FindAllString(s, -1)
	return match
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

func inttostring(i int) string {
	t := strconv.Itoa(i)
	return t
}

// PrintTop : beginning of output
func PrintTop(s *State) {
	//w := tabwriter.NewWriter(os.Stdout,0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Println("Target: ", s.URL)
	fmt.Println("Wordlists: ", strings.Join(s.WordListFiles, ", "))
	fmt.Println("=============================================================================================================")
	fmt.Println("URL                                                 Status        Chars          Words          Lines")
	fmt.Println("=============================================================================================================")
}

// PrintHelp :
func PrintHelp() {
	//PrintBanner()
	fmt.Printf("Usage: ./grequest [options] -w wordlist <url>\n\n")
	fmt.Printf("Options:\n")
	fmt.Println("-h/--help                     : This help")
	fmt.Println("--version                     : grequest version details")
	fmt.Println("-t N                          : Specify the number of concurrent connections (10 default)")
	fmt.Println("--follow                      : Follow HTTP redirections")
	fmt.Println("-w wordlist                   : Specify a wordlist file (alias for -z file,wordlist).")
	fmt.Println("-b cookie                     : Specify a cookie for the requests")
	fmt.Println("--hc/hl/hw/hh N[,N]+          : Hide responses with the specified code/lines/words/chars")
	fmt.Printf("--sc/sl/sw/sh N[,N]]+         : Show responses with the specified code/lines/words/chars\n")
	fmt.Println("Keyword: FUZZ, ..., FUZnZ  wherever you put these keywords wfuzz will replace them with the values of the specified payload.")
	fmt.Println("Example: - grequest -w users.txt -w pass.txt --sc 200 http://www.site.com/log.asp?user=FUZZ&pass=FUZ2Z")
}
