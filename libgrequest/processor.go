package libgrequest

import (
	//"fmt"
)

// Processor : 
func Processor(s *State) {
	schan := make(chan string)
	rchan := make(chan *Result, 2) // :TODO

	go func() { GetURLrec(s, 0, s.URL, schan) }()

	for i := 0; i < 8; i++ {
		go func() {
			resp, res := GoGet(s, <-schan, s.Cookies)
			if resp != nil {
				rchan <- res
			}
		}()
	}
	
	for r := 0; r < 8; r++ {
		PrintFn(s, <-rchan)
	}
}
