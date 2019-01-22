Aims to reproduce wfuzz's functionality and versatility and combine it with go's concurrency. Based on gobuster, written by OJ (theColonial)

build:
`go build`

build for windows: use --no-color option
`GOOS=windows GOARCH=amd64 go build -o gofuzz.exe`

most basic usage:
`gofuzz -w wordlist1 -w wordlist2 http://someip/FUZZFUZ2Z`

```
[+] gofuzz: dirty fork of gobuster to reproduce the functionality of wfuzz
[+] Author: brax (https://github.com/braaaax/gofuzz)

Usage:   gofuzz [options] -w wordlist <url>
Keyword: FUZZ, ..., FUZnZ  wherever you put these keywords gofuzz will replace them with the values of the specified payload.

Options:
-h/--help                     : This help.
-t N                          : Specify the number of concurrent connections (10 default).
--follow                      : Follow HTTP redirections.
-k                            : Skip verify on TLS connection.
-q                            : Quiet mode.
-p URL                        : Specify proxy URL.
-b COOKIE                     : Specify cookie for web request.
-ua USERAGENT                 : Specify the user agent.
--password PASSWORD           : Specify password.
--username USERNAME           : Specify username.
-w wordlist                   : Specify a wordlist file (alias for -z file,wordlist).
--hc/hl/hw/hh N[,N]+          : Hide responses with the specified code, lines, words, or chars.
--sc/sl/sw/sh N[,N]]+         : Show responses with the specified code, lines, words, or chars.

Examples: gofuzz -w users.txt -w pass.txt --sc 200 http://www.site.com/log.asp?user=FUZZ&pass=FUZ2Z
          gofuzz --follow -z file,default/common.txt -z file,default/ext.txt http://somesite.com/FUZZFUZ2Z
          gofuzz -t 32 -k --follow -w somelist.txt https://someTLSsite.com/FUZZ
```

foreign lib: 
`github.com/fatih/color`
