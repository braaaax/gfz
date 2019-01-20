package libgrequest

import (
	"regexp"
	"strings"
)

func freplace(url, word string) string {
	FUZZre := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	FUZZmatch := FUZZre.FindString(url)
	return strings.Replace(url, FUZZmatch, word, 1)
}

// GetURL : Recursive function, feeds urls for GoGet into string channel.
func GetURL(s *State, currentloop int, u string, uchan chan string) {
	if currentloop == len(s.Fuzzer.Indexes) {
		for i := 0; i < currentloop; i++ {
			u = freplace(u, s.Fuzzer.Wordlists[i][s.Fuzzer.Indexes[i]])
		}
		uchan <- u
	} else {
		for s.Fuzzer.Indexes[currentloop] = 0; s.Fuzzer.Indexes[currentloop] != s.Fuzzer.Maxes[currentloop]; s.Fuzzer.Indexes[currentloop]++ {
			GetURL(s, currentloop+1, u, uchan)
		}
	}
}
