package utils

import (
	"crypto/sha1"
	"fmt"
	"github.com/gosexy/to"
	"io"
	"math/rand"
	"regexp"
	"time"
)

const (
	CHAR_MAP = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func FilterAllHtmlTag(html string) string {
	tagReg := regexp.MustCompile(`<(/?)(\w[^>]*?)>`)
	return tagReg.ReplaceAllString(html, "&lt;$1&gt;")
}

func FilterHarmfulHtmlTag(html string) string {
	tagReg := regexp.MustCompile(`<(/?)(script|i?frame|style|html|body|title|link|meta|form|input|textarea\?|%)([^>]*?)>`)
	html = tagReg.ReplaceAllString(html, "&lt;$1$2$3&gt;")
	attrReg := regexp.MustCompile(`<\w+\s+([^>]+)>`)

	return attrReg.ReplaceAllStringFunc(html, func(s string) string {
		eventReg := regexp.MustCompile(`\s*on\w+\s*=[^>]*?(\s|>)`)
		s = eventReg.ReplaceAllString(s, "$1")
		hrefReg := regexp.MustCompile(`(href\s*=)[^>]*?javascript:[^>]*?(\s|>)`)
		return hrefReg.ReplaceAllString(s, `$1"#"$2`)
	})
}

func GetSummary(content string) string {
	text := regexp.MustCompile(`<(/?)\w[^>]*?>`).ReplaceAllString(content, "")
	return SubstrByByte(text, 20*3)
}

func SubstrByByte(str string, length int) string {
	if len(str) < length {
		return str
	}
	bs := []byte(str)[:length]
	bl := 0
	for i := len(bs) - 1; i >= 0; i-- {
		switch {
		case bs[i] >= 0 && bs[i] <= 127:
			return string(bs[:i+1])
		case bs[i] >= 128 && bs[i] <= 191:
			bl++
		case bs[i] >= 192 && bs[i] <= 253:
			cl := 0
			switch {
			case bs[i]&252 == 252:
				cl = 6
			case bs[i]&248 == 248:
				cl = 5
			case bs[i]&240 == 240:
				cl = 4
			case bs[i]&224 == 224:
				cl = 3
			default:
				cl = 2
			}
			if bl+1 == cl {
				return string(bs[:i+cl])
			}
			return string(bs[:i])
		}
	}
	return ""
}

func Sha1(s string) string {
	t := sha1.New()
	io.WriteString(t, s)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func Num2UrlStr(n int64) string {
	l := int64(len(CHAR_MAP))
	temp := []string{}
	for n > l {
		n = n / l
		i := n % l
		temp = append(temp, CHAR_MAP[i:i+1])
	}
	s := ""
	for i := len(temp) - 1; i >= 0; i -= 1 {
		s += temp[i]
	}
	return s
}

func NewFileName() string {
	t := time.Now()
	s := to.String(t.Hour()+10+rand.Intn(10)) + to.String(t.Minute()+10+rand.Intn(10))
	s += to.String(to.String(t.Nanosecond()))
	return t.Format("2006/01/02/") + Num2UrlStr(to.Int64(s))
}
