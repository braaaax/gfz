package libgrequest

import (
	"regexp"
	"strings"
)

func freplace(url, wordlists string) string {
	FUZZre := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	FUZZmatch := FUZZre.FindString(url)
	return strings.Replace(url, FUZZmatch, wordlists, 1)
}

// GetURL :
func GetURL(s *State, currentloop int, u string, uchan chan string) {
	// indexes []int, maxes[]int, curloop, wordlists [][]string, u string
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
