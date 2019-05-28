package utils

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

var (
	// LogPingPang   = []string{"clklogl", "clklogr"}
	StatsPingPang = []string{"statsl", "statsr"}
	RsaClients    = make(map[string]*RsaClient)
)

type RsaClient struct {
	PrivateKey *rsa.PrivateKey
}

func Md5(in string) string {
	inmd5 := md5.Sum([]byte(in))
	return hex.EncodeToString(inmd5[:])
}

func Sha1(in string) string {
	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(in))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)

	return hex.EncodeToString(bs)
	//SHA1 值经常以 16 进制输出，例如在 git commit 中。使用%x 来将散列结果格式化为 16 进制字符串。
}

func RsaSign(prikeysrc string) string {
	prikeysrcs := strings.Split(prikeysrc, "^rwl^")
	if 2 != len(prikeysrcs) {
		return ""
	}

	if rsains, ok := RsaClients[prikeysrcs[0]]; ok {
		return rsains.Sign(StrToBytes(prikeysrcs[1]), crypto.SHA1)
	} else {
		return ""
	}
}

func (r *RsaClient) Sign(src []byte, hash crypto.Hash) string {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)

	rsasrc, err := rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, hash, hashed)
	if nil == err {
		bytearr := base64.StdEncoding.EncodeToString(rsasrc)
		safeurl := strings.Replace(string(bytearr), "/", "_", -1)
		safeurl = strings.Replace(safeurl, "+", "-", -1)
		safeurl = strings.Replace(safeurl, "=", "", -1)
		return safeurl
	}

	return ""
}

func StringHashCode(instr string) int32 {
	hashcode := int32(0)
	for _, cr := range instr {
		hashcode = int32(31)*hashcode + cr
	}

	return hashcode
}
