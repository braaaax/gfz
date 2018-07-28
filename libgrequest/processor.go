package libgrequest

import (
	//"fmt"
	"sync"
)

//func PrintDirResult(){}

// FuzzProc : handle multi-threaded url requests created
func FuzzProc(s *State) {
	PrepareSignalHandler(s)

	urlChan := make(chan string)
	resultChan := make(chan Result)
	go func() {
		CallScanWordlist(s.Wordlists[:], s.FuzzMap, s.BaseMap, s.URLFuzz, urlChan, s.Terminate)
		close(urlChan)
	}()
	//s.Printer = PrintResult
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
	go func() {
		for r := range resultChan {
			//fmt.Printf("%+v", &r)
			
			s.Printer(s, &r)
		}
		printerGroup.Done()
	}()
	processorGroup.Wait()
	close(resultChan)
	printerGroup.Wait()
}
