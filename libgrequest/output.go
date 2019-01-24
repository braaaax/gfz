package libgrequest

import (
	"fmt"
	// "io/ioutil"
	// "net/http"
	"regexp"
	"strconv"
	"strings"
	// "unicode/utf8"

	"github.com/fatih/color"
)

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
		output += fmt.Sprintf(" Status: %-10s", code)
	}
	if r.Chars >= int64(0) {
		output += fmt.Sprintf(" Chars=%-6.8s", res2string(r.Chars))
		output += fmt.Sprintf(" Words=%-6.8s", res2string(r.Words))
		output += fmt.Sprintf(" Lines=%-6.8s", res2string(r.Lines))
	}
	output += "\n"
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
	// hiblk := color.New(color.FgHiBlack).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	// white := color.New(color.FgHiWhite).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()

	// write output line
	code := res2string(r.Code)
	output := ""
	output += fmt.Sprintf("%-20s", parseurl(r.URL)) // clip the leading '/'
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
		output += fmt.Sprintf(" Status %-10s", code)
	}

	// print result struct data
	if r.Chars >= int64(0) {
		output += fmt.Sprintf("%s%s%-6.8s", " Chars", "=", res2string(r.Chars)) // note to self: color for result int64s add about 10 spaces
		output += fmt.Sprintf("%s%s%-6.8s", " Words", "=", res2string(r.Words))
		output += fmt.Sprintf("%s%s%-6.8s", " Lines", "=", res2string(r.Lines))
	}
	output += "\n"

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
