package libgrequest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

func parseurl(arg string) string {
	u, err := url.Parse(arg)
	if err != nil {
		panic(err)
	}
	return u.Path
}

func resulttostring(arg int64) string {
	return strconv.FormatInt(arg, 10)
}

// PrintChars :
// probably a better way to do this
func PrintChars(s *State, r *Result) {
	if s.Filter.Contains(r.Chars) == s.Show {
		PrintFn(s, r)
	}
}

// PrintWords :
func PrintWords(s *State, r *Result) {
	if s.Filter.Contains(r.Words) == s.Show {
		PrintFn(s, r)
	}
}

// PrintStatus :
func PrintStatus(s *State, r *Result) {
	if s.Filter.Contains(r.Code) == s.Show { // issue nil
		PrintFn(s, r)
	}
}

// PrintLines :
func PrintLines(s *State, r *Result) {
	if s.Filter.Contains(r.Lines) == s.Show {
		PrintFn(s, r)
	}
}

// ParsePrintFilterArgs :
func ParsePrintFilterArgs(s *State, fs string) {
	m := regexp.MustCompile("(sl|sc|sw|sh|hc|hl|hh|hw)").FindString(fs)
	if string(m[0]) == "s" {
		s.Show = true
	} else {
		s.Show = false
	}
	switch m[1:] {
	case "c":
		s.Printer = PrintStatus
	case "l":
		s.Printer = PrintLines
	case "w":
		s.Printer = PrintWords
	case "h":
		s.Printer = PrintChars
	}
}

// ProcessResponse :
func ProcessResponse(fullUrl string, resp *http.Response) (*Result, error) {
	//set body
	var r = &Result{URL: fullUrl, Code: int64(resp.StatusCode)}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	r.Body = body
	sbody := string(body)
	if err == nil {
		r.Chars = int64(utf8.RuneCountInString(sbody))
		r.Words = int64(len(strings.Fields(sbody)))
		newlineRE := regexp.MustCompile("\n")
		r.Lines = int64(len(newlineRE.FindAllString(sbody, -1)))
	} else {
		return r, err
	}
	return r, nil
}

// WriteToFile :
func WriteToFile(output string, s *State) {
	_, err := s.OutputFile.WriteString(output)
	if err != nil {
		panic("[!] Unable to write to file " + s.OutputFileName)
	}
}

// PrintNoColorFn : print page into to stdout
func PrintNoColorFn(s *State, r *Result) {
	output := ""
	output += fmt.Sprintf("%-20s", parseurl(r.URL)) // to do: just the uri
	if !s.NoStatus {
		code := strconv.FormatInt(r.Code, 10)
		if code == "200" {
			color.Green(code)
		}
		if code == "403" {
			color.Red(code)
		}
		output += fmt.Sprintf(" Status: %-8s", code)
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

// PrintFn : prints corized page info to stdout
func PrintFn(s *State, r *Result) {
	if r == nil {
		return
	}
	// blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	code := resulttostring(r.Code)
	output := ""

	output += fmt.Sprintf("%-20s", parseurl(r.URL)) // to do: just the uri
	if !s.NoStatus {
		if strings.HasPrefix(code, "2") {
			code = green(code)
		}
		if strings.HasPrefix(code, "3") {
			code = yellow(code)
		}
		if strings.HasPrefix(code, "4") {
			code = red(code)
		}
		output += fmt.Sprintf(" Status=%-8s", code)
	}
	if r.Chars >= int64(0) {
		output += fmt.Sprintf(" Chars=%-8s", yellow(r.Chars))
		output += fmt.Sprintf(" Words=%-8s", yellow(r.Words))
		output += fmt.Sprintf(" Lines=%-8s", yellow(r.Lines))
	}
	output += "\n"

	if s.OutputFile != nil {
		WriteToFile(output, s)
	}
	re := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	match := re.FindString(output)
	if len(match) == 0 {
		fmt.Printf(output)
	}

}
