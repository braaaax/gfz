package libgrequest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)



// MakePostRequest :
func MakePostRequest(s *State, u, p string) (*int, *Result) {
	req, err := http.NewRequest("POST", u, strings.NewReader(p))
	if err != nil {
		return nil, nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")
	resp, err := s.Client.Do(req)
	if err != nil {
		if ue, ok := err.(*url.Error); ok {
			if strings.HasPrefix(ue.Err.Error(), "x509") {
				fmt.Println("[-] Invalid certificate")
			}
			if re, ok := ue.Err.(*RedirectError); ok {
				return &re.StatusCode, nil
			}
		}
		return nil, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var r = &Result{Url: u, Body: body}
	return &resp.StatusCode, r
}

// GoPost :
func GoPost(s *State, url, payload string) (*int, *Result) {
	return MakePostRequest(s, url, payload)
}

// PostProcessor :
func PostProcessor(s *State, u string, resultChan chan<- Result) {
	PostResp, PostRes := GoPost(s, u, s.Payload)
	if PostResp != nil {
		resultChan <- *PostRes
	}
}