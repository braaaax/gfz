package libgrequest

import (
	// "fmt"
	"regexp"
	"strings"
)

func freplace(payload, newword string) string {
	// fmt.Println("payload:", payload, "newword", newword)
	FUZZre := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	FUZZmatch := FUZZre.FindString(payload)
	if FUZZre.MatchString(payload) {
		return strings.Replace(payload, FUZZmatch, newword, 1)
	}
	return payload
}

// GetURL : Recursive function, feeds urls for GoGet into string channel.
func GetURL(s *State, currentloop int, cli string, pchan chan string) {
	if currentloop == len(s.Fuzzer.Indexes) {
		for i := 0; i < currentloop; i++ {
			cli = freplace(cli, s.Fuzzer.Wordlists[i][s.Fuzzer.Indexes[i]])
			// fmt.Println("payload out:", cli)
		}
		pchan <- cli

	} else {
		for s.Fuzzer.Indexes[currentloop] = 0; s.Fuzzer.Indexes[currentloop] != s.Fuzzer.Maxes[currentloop]; s.Fuzzer.Indexes[currentloop]++ {
			GetURL(s, currentloop+1, cli, pchan)
		}
	}
}
