package libgrequest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// MakePostRequest :
func MakePostRequest(s *State, u, p string) (*int, error) {
	req, err := http.NewRequest("POST", u, strings.NewReader(p))
	if err != nil {
		return nil, nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=any")
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
	return &resp.StatusCode, nil
}

// GoPost :
func GoPost(s *State, url, payload string) (*int, error) {
	return MakePostRequest(s, url, payload)
}
