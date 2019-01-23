Aims to reproduce wfuzz's functionality and versatility. Based on gobuster.

build:

`go build`

build for windows: 

`GOOS=windows GOARCH=amd64 go build -o gofuzz.exe`

most basic usage:

`gofuzz -w wordlist1 -w wordlist2 http://someip/FUZZFUZ2Z`

```
Usage:   gofuzz [options] -w wordlist <url>
Keyword: FUZZ, ..., FUZnZ  wherever you put these keywords gofuzz will replace them with the values of the specified payload.

Options:
-h/--help                     : This help.
-w wordlist                   : Specify a wordlist file (alias for -z file,wordlist).
-z file/range/list,PAYLOAD    : Where PAYLOAD is FILENAME or 1-10 or "-" separated sequence.
--hc/hl/hw/hh N[,N]+          : Hide responses with the specified code, lines, words, or chars.
--sc/sl/sw/sh N[,N]]+         : Show responses with the specified code, lines, words, or chars.
-t N                          : Specify the number of concurrent connections (10 default).
-p URL                        : Specify proxy URL.
-b COOKIE                     : Specify cookie.
-ua USERAGENT                 : Specify user agent.
--password PASSWORD           : Specify password for basic web auth.
--username USERNAME           : Specify username.
--no-follow                   : Don't follow HTTP redirections.
--no-color                    : Monotone output.
-k                            : Strict TLS connections (skip verify = false).
-q                            : Quiet mode.

Examples: gofuzz -w users.txt -w pass.txt --sc 200 http://www.site.com/log.asp?user=FUZZ&pass=FUZ2Z
          gofuzz --follow -z file,default/common.txt -z file,default/ext.txt http://somesite.com/FUZZFUZ2Z
          gofuzz -t 32 -k --follow -w somelist.txt https://someTLSsite.com/FUZZ
```

foreign lib: 
`github.com/fatih/color`
