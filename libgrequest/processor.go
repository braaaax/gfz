package libgrequest

import "fmt"

// Processor :
func Processor(s *State) {
	N := iterations(s.Fuzzer.Maxes)
	schan := make(chan string)
	rchan := make(chan *Result, s.Threads) // :TODO err chan

	go func() { GetURLrec(s, 0, s.URL, schan) }()
	for i := 0; i < s.Threads; i++ {
		go func() {
			for {
				resp, res := GoGet(s, <-schan, s.Cookies)
				if resp != nil {
					rchan <- res // res could be anything
					s.Counter.Inc()
				}
			}
		}()
	}
	for r := 0; r < N; r++ {
		<-rchan
		fmt.Printf("%d/%d\r", s.Counter.v, N)
	}
}

//TODO: add --recursive
