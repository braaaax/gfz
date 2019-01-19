package libgrequest

import "fmt"

// Processor :
func Processor(s *State) {
	N := Nrequests(s.Fuzzer.Maxes)
	urlchan := make(chan string)
	errorchan := make(chan error, s.Threads) // :TODO err chan

	go func() { GetURL(s, 0, s.URL, urlchan) }()
	for i := 0; i < s.Threads; i++ {
		go func() {
			for {
				code, err := GoGet(s, <-urlchan, s.Cookies)
				if code != nil {
					errorchan <- err // res could be anything
					s.Counter.Inc()
				}
			}
		}()
	}
	for r := 0; r < N; r++ {
		<-errorchan
		fmt.Printf("%d/%d\r", s.Counter.v, N)
	}
}

//TODO: add --recursive
