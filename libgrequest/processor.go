package libgrequest

import (
	"fmt"

)



// Processor : channel controlcenter
func Processor(s *State) {
	N := TotalRequests(s.Fuzzer.Maxes)
	urlc := make(chan string)
	errorc := make(chan error, s.Threads)

	go func() { GetURL(s, 0, s.URL, urlc) }() // Payload is just a string with 'FUZZ'
	for i := 0; i < s.Threads; i++ {
		go func() {
			for {
				// if s.Terminate == true {break}
				code, err := GoGet(s, <-urlc, s.Cookies)
				if code != nil {
					errorc <- err
					s.Counter.Inc()
				}
			}
		}()
	}
	for r := 0; r < N; r++ {
		<-errorc
		fmt.Printf("[+] requests: %d/%d\r", s.Counter.v, N)
	}
}

//TODO: add --recursive
