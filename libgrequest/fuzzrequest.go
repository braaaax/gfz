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
// Make a request to the given URL.
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
	var r = &Result{Url: fullUrl, Code: int64(resp.StatusCode)}
	if s.IncludeLength {
		if resp.ContentLength <= 0 {
			ProcessResults(r, resp)
		} else {
			// length of associated content is unknown
			// we're just going to process anyway
			ProcessResults(r, resp)
		}
	}
	return &resp.StatusCode, r
}

// GoGet :
func GoGet(s *State, url, cookie string) (*int, *Result) {
	//fmt.Println(url)
	return MakeRequest(s, url, cookie)
}


// GetProcessor: 
func GetProcessor(s *State, u string, resultChan chan<- Result) {
	GetResp, GetResult := GoGet(s, u, s.Cookies)
	if GetResp != nil {
		resultChan <- *GetResult
	}
}

/*
// SetupDir :
func SetupDir(s *State) bool {
	guid := uuid.Must(uuid.NewV4())
	wildcardResp, _ := GoGet(s, s.Url+fmt.Sprintf("%s", guid), s.Cookies)
	w := int64(*wildcardResp)
	if s.StatusCodes.Contains(w) {
		s.IsWildcard = true
		fmt.Println("[-] Wildcard response found:", fmt.Sprintf("%s%s", s.Url, guid), "=>", *wildcardResp)
		if !s.WildcardForced {
			fmt.Println("[-] To force processing of Wildcard responses, specify the '-fw' switch.")
		}
		return s.WildcardForced
	}

	return true
}
*/

// WriteToFile :
func WriteToFile(output string, s *State) {
	_, err := s.OutputFile.WriteString(output)
	if err != nil {
		panic("[!] Unable to write to file " + s.OutputFileName)
	}
}

// PrintFn :
func PrintFn(s *State, r *Result) {
	output := ""
	output += fmt.Sprintf("%-50s", r.Url) // to do: just the uri
	if !s.NoStatus {
		output += fmt.Sprintf(" (Status: %-4d)", r.Code)
	}
	if r.Chars >= int64(0) {
		output += fmt.Sprintf(" Chars=%-8d", r.Chars)
		output += fmt.Sprintf(" Words=%-8d", r.Words)
		output += fmt.Sprintf(" Lines=%-8d", r.Lines)
	}
	output += "\n"
	if s.OutputFile != nil {
		WriteToFile(output, s)
	}
	fmt.Printf(output)
}
