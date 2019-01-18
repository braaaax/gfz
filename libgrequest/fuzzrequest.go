package libgrequest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// RedirectHandler :
type RedirectHandler struct {
	Transport http.RoundTripper
	State     *State
}

// RedirectError :
type RedirectError struct {
	StatusCode int
}

func (e *RedirectError) Error() string {
	return fmt.Sprintf("%-8d", e.StatusCode)
}

// RoundTrip :
func (rh *RedirectHandler) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if rh.State.FollowRedirect {
		return rh.Transport.RoundTrip(req)
	}
	resp, err = rh.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther,
		http.StatusNotModified, http.StatusUseProxy, http.StatusTemporaryRedirect:
		return nil, &RedirectError{StatusCode: resp.StatusCode}
	}
	return resp, err
}

// MakeRequest : 
func MakeRequest(s *State, fullUrl, cookie string) (*int, *Result) {
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, nil
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	if s.UserAgent != "" {
		req.Header.Set("User-Agent", s.UserAgent)
	}
	if s.Username != "" {
		req.SetBasicAuth(s.Username, s.Password)
	}
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
	//var r = &Result{URL: fullUrl, Code: int64(resp.StatusCode)}
	r := ProcessResults(fullUrl, resp)
	PrintFn(s, r)
	return &resp.StatusCode, r
}

// GoGet :
func GoGet(s *State, url, cookie string) (*int, *Result) {
	return MakeRequest(s, url, cookie)
}

// GetProcessor :
func GetProcessor(s *State, u string, resultChan chan<- Result) {
	GetResp, GetResult := GoGet(s, u, s.Cookies)
	if GetResp != nil {
		resultChan <- *GetResult
	}
}
