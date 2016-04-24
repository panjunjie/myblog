package lib

import (
	"fmt"
	"log"
	"math/rand"

	"regexp"
	"strings"
	"time"
	"unicode/utf16"
)

func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

func ShowSubstr(s string, l int) string {
	s = ClearHtmlTag(s)
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}

func show_strlen(s string) int {
	sl := 0
	rs := []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			sl++
		} else {
			sl += 2
		}
	}
	return sl
}

func ClearHtmlTag(html string) string {
	if strings.TrimSpace(html) == "" {
		return ""
	}
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	html = re.ReplaceAllStringFunc(html, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	html = re.ReplaceAllString(html, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	html = re.ReplaceAllString(html, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	html = re.ReplaceAllString(html, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	html = re.ReplaceAllString(html, "\n")

	//去掉空格 &nbsp;
	html = strings.Replace(html, "&nbsp;", "", -1)
	//html = strings.Replace(html, " ", "", -1)
	html = strings.TrimSpace(html)

	return html
}

//字符串转换来unit16
func StringToUTF16(s string) []uint16 {
	return utf16.Encode([]rune(s + "\x00"))
}

func RandStr() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), r.Intn(100))
}

func CheckErr(err error, msg string) {
	if nil != err {
		log.Fatalf("发生错误：%s \n %#v", msg, err)
	}
}

func ShowErr(err error, msg string) {
	if nil != err {
		log.Printf("发生错误：%s \n %#v", msg, err)
	}
}

func ShowMsg(msg string) {
	if msg != "" {
		log.Println(msg)
	}
}

func ShowInfo(msg interface{}) {
	if nil != msg {
		log.Println(msg)
	}
}
