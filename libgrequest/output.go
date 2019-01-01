package libgrequest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
	"strconv"

	"github.com/fatih/color"
)

// FuzzPrintChars :
func FuzzPrintChars(s *State, r *Result) {
	if s.Filter.Contains(r.Chars) == s.Show {
		PrintFn(s, r)
	}
}

// FuzzPrintWords :
func FuzzPrintWords(s *State, r *Result) {
	if s.Filter.Contains(r.Words) == s.Show {
		PrintFn(s, r)
	}
}

// FuzzPrintStatus :
func FuzzPrintStatus(s *State, r *Result) {
	if s.Filter.Contains(r.Code) == s.Show {
		PrintFn(s, r)
	}
}

// FuzzPrintLines :
func FuzzPrintLines(s *State, r *Result) {
	if s.Filter.Contains(r.Lines) == s.Show {
		PrintFn(s, r)
	}
}

// PrintFilter : switch for print filter
func PrintFilter(s *State, fs string) {
	m := regexp.MustCompile("(sl|sc|sw|sh|hc|hl|hh|hw)").FindString(fs)
	if string(m[0]) == "s" {
		s.Show = true
	} else {
		s.Show = false
	}
	switch m[1:] {
	case "c":
		s.Printer = FuzzPrintStatus
	case "l":
		s.Printer = FuzzPrintLines
	case "w":
		s.Printer = FuzzPrintWords
	case "h":
		s.Printer = FuzzPrintChars
	}
}

// ProcessResults :
func ProcessResults(r *Result, resp *http.Response) {
	//set body
	body, err := ioutil.ReadAll(resp.Body)
	r.Body = body
	sbody := string(body)
	if err == nil {
		r.Chars = int64(utf8.RuneCountInString(sbody))
		r.Words = int64(len(strings.Fields(sbody)))
		newlineRE := regexp.MustCompile("\n")
		r.Lines = int64(len(newlineRE.FindAllString(sbody, -1)))
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
	output += fmt.Sprintf("url: %-50s", r.URL) // to do: just the uri
	if !s.NoStatus {
		code := strconv.FormatInt(r.Code, 10)
		if code == "200" {color.Green(code)}
		if code == "403" {color.Red(code)}
		output += fmt.Sprintf(" (Status: %-4s)", code)
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
