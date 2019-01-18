package libgrequest

// Processor :
func Processor(s *State) {
	schan := make(chan string)
	rchan := make(chan *Result, s.Threads) // :TODO err chan

	go func() { GetURLrec(s, 0, s.URL, schan) }()

	for i := 0; i < s.Threads; i++ {
		go func() {
			for {
				resp, res := GoGet(s, <-schan, s.Cookies)
				
				if resp != nil {
					rchan <- res
				}
			}
		}()
	}

	for r := 0; r < iterations(s.Fuzzer.Maxes); r++ {
		<-rchan
	}
}

//TODO: add recursion
