package utils

import (
	"fmt"
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

func GenUniqueIdExt(campaignid string, ts int64, uniid string) string {
	changeid := strings.Join([]string{"wy", campaignid}, "")

	pstr := strings.Join([]string{changeid, strconv.FormatInt(ts, 10)}, "r")

	bid := []byte(uniid)
	padding := string(bid[(len(bid)/4)*3:])
	padding = strings.Replace(padding, "r", "x", -1)

	return strings.Join([]string{pstr, padding}, "r")
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value-0.005), 64)
	return value
}
