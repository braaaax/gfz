BUILD=go build
OUT_LINUX=gofuzz
OUT_WINDOWS=gofuzz.exe
SRC=main.go
LINUX_LDFLAGS=--ldflags "-s -w"
WIN_LDFLAGS=--ldflags "-s -w"

linux32:
	GOOS=linux GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${SRC}
