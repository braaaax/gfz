package libgrequest

// "fmt"

// Processor :
func Processor(s *State) {
	schan := make(chan string)
	rchan := make(chan *Result, s.Threads) // :TODO

	go func() { GetURLrec(s, 0, s.URL, schan) }()

	for i := 0; i < iterations(s.Fuzzer.Maxes); i++ {
		go func() {
			resp, res := GoGet(s, <-schan, s.Cookies)
			if resp != nil {
				rchan <- res
			}
		}()
	}

	for r := 0; r < iterations(s.Fuzzer.Maxes); r++ {
		printout(<-rchan)
	}
}

func ProcessorII(s *State) {
	schan := make(chan string)
	rchan := make(chan *Result, s.Threads) // :TODO

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
		colorize(s, <-rchan)
	}
}
