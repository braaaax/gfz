package libgrequest

import (
	"regexp"
	"strings"
	
)

func freplace(u, w string) string {
	re := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	match := re.FindString(u)
	return strings.Replace(u, match, w, 1)
}

// GetURLrec :
func GetURLrec(s *State, curloop int, u string, schan chan string) {
	// indexes []int, maxes[]int, curloop, wordlists [][]string, u string
	if curloop == len(s.Fuzzer.Indexes) {
		for i := 0; i < curloop; i++ {
			u = freplace(u, s.Fuzzer.Wordlists[i][s.Fuzzer.Indexes[i]])
		}
		
		schan <- u

	} else {
		
		for s.Fuzzer.Indexes[curloop] = 0; s.Fuzzer.Indexes[curloop] != s.Fuzzer.Maxes[curloop]; s.Fuzzer.Indexes[curloop]++ {
			GetURLrec(s, curloop+1, u, schan)
		}
	}
}
