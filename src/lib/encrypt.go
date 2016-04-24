package lib

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

//对字符串进行MD5哈希
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func MD5Ext(s string, upper bool) string {
	s = MD5(s)
	if upper {
		return strings.ToUpper(s)
	}
	return s
}

//对字符串进行SHA1哈希
func SHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
