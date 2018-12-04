package libgrequest

import (
	"sync"
)

//func PrintDirResult(){}

// FuzzProc : handle multi-threaded url requests created
func FuzzProc2(s *State) {
	PrepareSignalHandler(s)

	urlChan := make(chan string)
	resultChan := make(chan Result)
	go func() {
		CallScanWordlist(s.Wordlists[:], s.FuzzMap, s.BaseMap, s.URLFuzz, urlChan, s.Terminate)
		close(urlChan)
	}()
	processorGroup := new(sync.WaitGroup)
	processorGroup.Add(s.Threads)
	printerGroup := new(sync.WaitGroup)
	printerGroup.Add(1)
	for i := 0; i < s.Threads; i++ {
		go func() {
			for u := range urlChan {
				if s.Terminate {
					break
				}
				s.Processor(s, u, resultChan)
			}
			processorGroup.Done()
		}()
	}
	// reads from resultChan
	go func() {
		for r := range resultChan {
			s.Printer(s, &r)
		}
		printerGroup.Done()
	}()
	processorGroup.Wait()
	close(resultChan)
	printerGroup.Wait()
}
