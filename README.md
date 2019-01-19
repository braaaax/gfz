Aims to reproduce wfuzz's functionality and versatility and combine it with go's concurrency. Based on gobuster, written by OJ (theColonial)

most basic usage:

`gofuzz -w wordlist1 -w wordlist2 http://someip/FUZZFUZ2Z`

```
[+] gofuzz: dirty fork of gobuster to reproduce the functionality of wfuzz
[+] Author: brax (https://github.com/braaaax/gofuzz)

Usage: ./gofuzz [options] -w wordlist <url>

Options:
-h/--help                     : This help.
-t N                          : Specify the number of concurrent connections (10 default).
--follow                      : Follow HTTP redirections.
-w wordlist                   : Specify a wordlist file (alias for -z file,wordlist).
--hc/hl/hw/hh N[,N]+          : Hide responses with the specified code, lines, words, or chars.
--sc/sl/sw/sh N[,N]]+         : Show responses with the specified code, lines, words, or chars.

Keyword: FUZZ, ..., FUZnZ  wherever you put these keywords wfuzz will replace them with the values of the specified payload.
Examples: gofuzz -w users.txt -w pass.txt --sc 200 http://www.site.com/log.asp?user=FUZZ&pass=FUZ2Z
          gofuzz --follow -z file,default/common.txt -z file,default/ext.txt http://somesite.com/FUZZFUZ2Z
```

foreign lib: `github.com/fatih/color`