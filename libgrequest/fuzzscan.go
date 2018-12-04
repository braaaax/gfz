package libgrequest

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// FuzzReplace : helper function for ScanWordlist
func FuzzReplace(s string, fm map[string]string, c chan string) {
	/* replace FUZnZ with word and put in channel */
	re := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	match := re.FindString(s)
	res := strings.Replace(s, match, fm[match], 1)
	if len(match) > 0 {
		FuzzReplace(res, fm, c)
	} else {
		if IsMapFull(fm) == true {
			c <- res
		}
	}
}

// ScanWordlist : read words from file line by line and replace 'FUZZ' with word
func ScanWordlist(filenames []string, fuzzmap map[string]string, basefuzz map[string]string, uri string, recurse bool, c chan string, terminate bool) {
	/*scan wordlist line-by-line
	set global base_fuzz_m value according
	to variable nested loop recursive scheme */
	var scanner *bufio.Scanner
	wordlist, err := os.Open(filenames[0])
	if err != nil {
		fmt.Println(filenames[0])
		panic("failed")
	}
	defer wordlist.Close()
	scanner = bufio.NewScanner(wordlist)
	for scanner.Scan() {
		word := scanner.Text()
		if !strings.HasPrefix(word, "#") && len(word) > 0 {
			/* mark whenever there is a new FUZZ
			value and set the value of the rest to ""
			FUZZ is the "anchor" value */
			if fuzzmap[filenames[0]] == "FUZZ" && basefuzz["FUZZ"] != "" && basefuzz["FUZZ"] != word {
				for k := range basefuzz {
					if k != "FUZZ" {
						basefuzz[k] = ""
					}
				}
			}
			basefuzz[fuzzmap[filenames[0]]] = word // assign word to placeholder
			if IsMapFull(basefuzz) == true {       // if no empty values
				FuzzReplace(uri, basefuzz, c)
			}
			if recurse == true {
				CallScanWordlist(filenames[1:], fuzzmap, basefuzz, uri, c, terminate)
			}
		}
	}
}

// CallScanWordlist : variable nested loop algo for ScanWordlist
func CallScanWordlist(filenames []string, fuzzmap map[string]string, basemap map[string]string, uri string, c chan string, terminate bool) {
	if len(filenames) > 1 {
		ScanWordlist(filenames[:], fuzzmap, basemap, uri, true, c, terminate)
	} else {
		ScanWordlist(filenames[:], fuzzmap, basemap, uri, false, c, terminate)
	}
}