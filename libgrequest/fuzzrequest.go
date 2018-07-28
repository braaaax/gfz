package libgrequest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/satori/go.uuid"
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
			body, err := ioutil.ReadAll(resp.Body)
			sbody := string(body)
			r.Body = body
			if err == nil {
				// ** address **
				r.Chars = int64(utf8.RuneCountInString(sbody))
				r.Words = int64(len(strings.Fields(sbody)))
				newlineRE := regexp.MustCompile("\n")
				r.Lines = int64(len(newlineRE.FindAllString(sbody, -1)))
			}
		} else {
			r.Chars = resp.ContentLength
			body, _ := ioutil.ReadAll(resp.Body)
			sbody := string(body)
			r.Words = int64(len(strings.Fields(sbody)))
			newlineRE := regexp.MustCompile("\n")
			r.Lines = int64(len(newlineRE.FindAllString(sbody, -1)))
		}
	}
	return &resp.StatusCode, r
}

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

// GoGet :
// Small helper to combine URL with URI then make a
// request to the generated location.
func GoGet(s *State, url, cookie string) (*int, *Result) {
	//fmt.Println(url)
	return MakeRequest(s, url, cookie)
}

// GoPost :
func GoPost(s *State, url, payload string) (*int, *Result) {
	return MakePostRequest(s, url, payload)
}

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

// GetProcessor :
func GetProcessor(s *State, u string, resultChan chan<- Result) {
	GetResp, GetResult := GoGet(s, u, s.Cookies)
	if GetResp != nil {
		resultChan <- *GetResult
	}
}

// PostProcessor :
func PostProcessor(s *State, u string, resultChan chan<- Result) {
	PostResp, PostRes := GoPost(s, u, s.Payload)
	if PostResp != nil {
		resultChan <- *PostRes
	}
}

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
