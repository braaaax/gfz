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

// RedirectError : redirect err struct from gobuster
type RedirectError struct {
	StatusCode int
}

func (e *RedirectError) Error() string {
	return fmt.Sprintf("%-8d", e.StatusCode)
}

// RoundTrip : roundtrip from gobuster
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

// MakeRequest : make http request
func MakeRequest(s *State, fullURL, cookie string) (*int, error) {
	req, err := http.NewRequest("GET", fullURL, nil)
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
				fmt.Println("[!] Invalid certificate, try using -k.")
			}
			if re, ok := ue.Err.(*RedirectError); ok {
				return &re.StatusCode, nil
			}
		}
		return nil, nil
	}
	defer resp.Body.Close()
	r, err := ProcessResponse(fullURL, resp)
	if err == nil {
		if s.Quiet != true {
			s.Printer(s, r)
		}
	}
	return &resp.StatusCode, nil
}

// GoGet : returs address of response statuscode and error
func GoGet(s *State, url, cookie string) (*int, error) {
	return MakeRequest(s, url, cookie)
}
