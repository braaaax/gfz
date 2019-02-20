package libgrequest

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func parseurl(uarg string) string {
	u, err := url.Parse(uarg)
	if err != nil {
		panic(err)
	}
	return u.Host + u.Path
}

// PrintChars :
// probably a better way to do this
func PrintChars(s *State, r *Result) {
	if s.Filter.Contains(r.Chars) == s.Show {
		if s.NoColor {
			PrintNoColorFn(s, r)
		} else {
			PrintColorFn(s, r)
		}
	}
}

// PrintWords :
func PrintWords(s *State, r *Result) {
	if s.Filter.Contains(r.Words) == s.Show {
		if s.NoColor {
			PrintNoColorFn(s, r)
		} else {
			PrintColorFn(s, r)
		}
	}
}

// PrintStatus :
func PrintStatus(s *State, r *Result) {
	if s.Filter.Contains(r.Code) == s.Show { // issue nil
		if s.NoColor {
			PrintNoColorFn(s, r)
		} else {
			PrintColorFn(s, r)
		}
	}
}

// PrintLines :
func PrintLines(s *State, r *Result) {
	if s.Filter.Contains(r.Lines) == s.Show {
		if s.NoColor {
			PrintNoColorFn(s, r)
		} else {
			PrintColorFn(s, r)
		}
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
	if !s.NoStatus {
		code := strconv.FormatInt(r.Code, 10)
		output += fmt.Sprintf(" Code %-10s ", code)
	}
	if r.Chars >= int64(0) {
		output += fmt.Sprintf(" C %-6.8s", res2string(r.Chars))
		output += fmt.Sprintf(" W %-6.8s", res2string(r.Words))
		output += fmt.Sprintf(" L %-6.8s", res2string(r.Lines))
	}
	output += fmt.Sprintf("%-30s", r.URL)
	output += "\n"
	if s.PrintBody == true {
		output += "\n"
		output += string(r.Body)
		output += "\n"
	}
	if s.OutputFile != nil {
		WriteToFile(output, s)
	}
	fmt.Printf(output)
}

// PrintColorFn : prints corized page info to stdout
func PrintColorFn(s *State, r *Result) {
	if r == nil {
		return
	}

	if s.NoColor == true {
		PrintNoColorFn(s, r)
		return
	}

	// color funcs
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()

	// write output line
	code := res2string(r.Code)
	output := ""

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
		if strings.HasPrefix(code, "5") {
			code = cyan(code)
		}
		output += fmt.Sprintf(" Code %-10s ", code)
	}

	// print result struct data
	if r.Chars >= int64(0) {
		output += fmt.Sprintf("%s %-6.8s", " C", res2string(r.Chars)) // note to self: color for result int64s add about 10 spaces
		output += fmt.Sprintf("%s %-5.5s", " W", res2string(r.Words))
		output += fmt.Sprintf("%s %-5.5s", " L", res2string(r.Lines))
	}
	output += fmt.Sprintf("%-30s", r.URL)
	output += "\n"

	if s.PrintBody == true {
		output += "\n"
		output += yellow(string(r.Body))
		output += "\n"
	}

	if s.OutputFile != nil {
		WriteToFile(output, s)
	}

	// ignore initial url
	re := regexp.MustCompile("FUZ(Z|[0-9]Z)")
	match := re.FindString(output)
	if len(match) == 0 {
		// print one line output to stdout
		fmt.Printf(output)
	}

}
