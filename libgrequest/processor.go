package libgrequest

import (
	"fmt"
)

// Processor : channel controlcenter
func Processor(s *State) {
	N := TotalRequests(s.Fuzzer.Maxes)
	urlc := make(chan string)
	codec := make(chan *int, s.Threads)

	go func() { GetURL(s, 0, s.URL, urlc) }() // Payload is just a string with 'FUZZ'

	for i := 0; i < s.Threads; i++ {
		go func() {
			for {
				// if s.Terminate == true {break}
				code, _ := GoGet(s, <-urlc, s.Cookies)
				codec<-code
			}
		}()
	}
	for r := 0; r < N; r++ {
		<-codec
		fmt.Printf("[+] requests: %d/%d\r", s.Counter.v, N)
	}
}

//TODO: add --recursive
