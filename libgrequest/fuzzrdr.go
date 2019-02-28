package libgrequest

import (
	"regexp"
	"strings"
	"fmt"
)

func freplace(payload, newword string) string {
	FUZZre := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	FUZZmatch := FUZZre.FindString(payload)
	if FUZZre.MatchString(payload) {
		return strings.Replace(payload, FUZZmatch, newword, 1)
	}
	return payload
}

// GetURL : Recursive function, feeds urls for GoGet into string channel.
func GetURL(s *State, currentloop int, payload string, uchan chan string) {
	if currentloop == len(s.Fuzzer.Indexes) {
		for i := 0; i < currentloop; i++ {
			payload = freplace(payload, s.Fuzzer.Wordlists[i][s.Fuzzer.Indexes[i]])
		}	
		fmt.Println("PAYLOAD", payload)
		uchan <- payload

	} else {
		for s.Fuzzer.Indexes[currentloop] = 0; s.Fuzzer.Indexes[currentloop] != s.Fuzzer.Maxes[currentloop]; s.Fuzzer.Indexes[currentloop]++ {
			GetURL(s, currentloop+1, payload, uchan)
		}
	}
}
