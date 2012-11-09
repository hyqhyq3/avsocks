package main

import (
	"bufio"
	"container/list"
	"encoding/base64"
	"log"
	"os"
	"regexp"
	"strings"
)

type RuleManager struct {
	Regexp, Wildcard list.List
}

type Rule interface {
	Exec(string) bool
}

type RegexpRule struct {
	Exclude bool
	Expr    *regexp.Regexp
}

type WildcardRule struct {
	Exclude bool
	Expr    string
}

func (w *WildcardRule) Exec(url string) bool {
	var str, pat = url, w.Expr
	var s, p string
	var star = false

loopStart:
	for s, p = str, pat; len(s) > 0; s, p = s[1:], p[1:] {
		switch p[0] {
		case '?':
			continue
		case '*':
			star = true
			str, pat = s, p[1:]
			if len(pat) == 0 {
				return true
			}
			goto loopStart
		default:
			if (s[0] | 32) != (p[0] | 32) {
				goto starCheck
			}

		} /* endswitch */
	} /* endfor */
	if p[0] == '*' {
		p = p[1:]
	}
	return len(p) == 0

starCheck:
	if !star {
		return false
	}
	str = str[1:]
	goto loopStart
}

func (r *RegexpRule) Exec(url string) bool {
	return r.Expr.Match([]byte(url))
}

const (
	REGEXP   = 0
	WILDCARD = 1
)

func (r *RuleManager) wildcard2regexp(w string) (ret string) {
	addslash := "\\+|{}[]()^$.#"
	for rune := range addslash {
		w = strings.Replace(w, string(rune), "\\"+string(rune), -1)
	}
	ret = strings.Replace(w, "*", ".*", -1)
	ret = strings.Replace(ret, "?", ".", -1)
	return ret
}

func (r *RuleManager) Add(line string) {
	if len(line) == 0 || line[0] == ';' || line[0] == '[' || line[0] == '!' {
		return
	}
	var t int
	var exclude = false
	if line[:2] == "@@" {
		exclude = true
		line = line[2:]
	}
	if line[0] == '+' {
		line = line[1:]
	}
	if line[0] == '/' && line[len(line)-1] == '/' {
		line = line[1 : len(line)-1]
		t = REGEXP
	} else if line[0] == '|' {
		line = line[1:]
		if line[0] == '|' {
			//ip
			t = REGEXP
			line = "^[\\w\\-]+:\\/+?" + r.wildcard2regexp(line[1:])
		} else {
			//protocol
			t = WILDCARD
			if line[len(line)-1] == '|' {
				line = "@" + line[:len(line)-1]
			} else {
				line = "@" + line + "*"
			}
		}
	} else if strings.Index(line, "^") > -1 {
		t = REGEXP
		line = r.wildcard2regexp(line)
		line = strings.Replace(line, "\\^", "(?:[^\\w\\-.%\\u0080-\\uFFFF]|$)", -1)
	} else {
		t = WILDCARD
		line = "http://*" + line + "*"
	}
	var rule Rule
	switch t {
	case WILDCARD:
		rule = &WildcardRule{exclude, line}
		r.Wildcard.PushBack(rule)
	case REGEXP:
		re, e := regexp.Compile(line)
		if e != nil {
			println(line, e.Error())
			return
		}
		rule = &RegexpRule{exclude, re}
		r.Regexp.PushBack(rule)
	}

}

func NewRuleManager() *RuleManager {
	rm := &RuleManager{}
	return rm
}

func xor(a bool, b bool) bool {
	return !((a && b) || (!a && !b))
}

func (r *RuleManager) Exec(url string) bool {
	for rule := r.Wildcard.Front(); rule != nil; rule = rule.Next() {
		wildcard := rule.Value.(*WildcardRule)
		if match := wildcard.Exec(url); match {
			return xor(match, wildcard.Exclude)
		}
	}
	for rule := r.Regexp.Front(); rule != nil; rule = rule.Next() {
		regexp := rule.Value.(*RegexpRule)
		if match := regexp.Exec(url); match {
			return xor(match, regexp.Exclude)
		}
	}
	return false
}

var rm = NewRuleManager()

func main() {
	f, e := os.Open("gfwlist.txt")
	if e != nil {
		log.Fatal(e)
	}
	d := base64.NewDecoder(base64.StdEncoding, f)
	r := bufio.NewReader(d)
	l, _, _ := r.ReadLine()
	if string(l[:10]) != "[AutoProxy" {
		log.Fatal("not a autoproxy list")
	}

	for {
		l, _, e := r.ReadLine()
		if e != nil {
			break
		}
		rm.Add(string(l))
	}

	println(rm.Exec("http://www.twitter.com/23423432"))
}
