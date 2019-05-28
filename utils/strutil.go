package utils

import (
	"strconv"
	"strings"
	"unsafe"
)

func StrToBytes(str string) []byte {
	s := (*[2]uintptr)(unsafe.Pointer(&str))
	d := [3]uintptr{s[0], s[1], s[1]}
	return *(*[]byte)(unsafe.Pointer(&d))
}

func ByteToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Equal(s1, s2 string) bool {
	return 0 == strings.Compare(s1, s2)
}

func IsNotEmpty(s string) bool {
	return len(s) > 0
}

func IsEmpty(s string) bool {
	return len(s) == 0
}

func GenUniqueId(campaignid string, ts int64, uniid string) string {
	pstr := strings.Join([]string{campaignid, strconv.FormatInt(ts, 10)}, "r")
	bid := []byte(uniid)
	return strings.Join([]string{pstr, string(bid[(len(pstr) + 1):])}, "r")
}
