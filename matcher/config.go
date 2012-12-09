package matcher

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var whitespace = regexp.MustCompile(`.\s+`)

func regexpCompile(s string) (r *regexp.Regexp, err error) {
	normalized := whitespace.ReplaceAllStringFunc(
		strings.TrimFunc(s, unicode.IsSpace),
		func(s string) string {
			if err != nil {
				return s
			}
			if s[0] == '\\' {
				sr := strings.NewReader(s[1:])
				space, _, e := sr.ReadRune()
				if e != nil {
					err = e
				}
				return string([]rune{'\\', space})
			}
			return ""
		})
	if err != nil {
		return
	}
	return regexp.Compile(normalized)
}

func regexpMustCompile(s string) *regexp.Regexp {
	r, err := regexpCompile(s)
	if err != nil {
		panic(err)
	}
	return r
}

type Config struct {
	Tracker  string `json:"tracker"`  // Matched tracker urls.
	Basename string `json:"basename"` // Matched (root) file basenames.
	Ext      string `json:"ext"`      // Matched (nested-)file extensions.
}

func (mc Config) Matcher() *Matcher {
	m := new(Matcher)
	if mc.Tracker != "" {
		m.Tracker = regexpMustCompile(mc.Tracker)
	}
	if mc.Basename != "" {
		m.Basename = regexpMustCompile(mc.Basename)
	}
	if mc.Ext != "" {
		m.Ext = regexpMustCompile(mc.Ext)
	}
	return m
}

func (mc Config) Validate() error {
	if _, err := regexpCompile(mc.Tracker); err != nil {
		return fmt.Errorf("Matcher tracker: %v", err)
	}
	if _, err := regexpCompile(mc.Basename); err != nil {
		return fmt.Errorf("Matcher basename: %v", err)
	}
	if _, err := regexpCompile(mc.Ext); err != nil {
		return fmt.Errorf("Matcher ext: %v", err)
	}
	return nil
}
