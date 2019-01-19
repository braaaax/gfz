package libgrequest

import "fmt"

// Processor : channel controlcenter
func Processor(s *State) {
	N := nrequests(s.Fuzzer.Maxes)
	urlc := make(chan string)
	errorc := make(chan error, s.Threads) // :TODO err chan

	go func() { GetURL(s, 0, s.URL, urlc) }()
	for i := 0; i < s.Threads; i++ {
		go func() {
			for {
				code, err := GoGet(s, <-urlc, s.Cookies)
				if code != nil {
					errorc <- err // res could be anything
					s.Counter.Inc()
				}
			}
		}()
	}
	for r := 0; r < N; r++ {
		<-errorc
		fmt.Printf("[+] %d/%d\r", s.Counter.v, N)
	}
}

//TODO: add --recursive
