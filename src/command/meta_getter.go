package command

import (
	"fmt"
	"regexp"
)

type MetaGetter struct {
	regexp *regexp.Regexp
}

func NewMetaGetter(pattern string) *MetaGetter {
	return &MetaGetter{
		regexp: regexp.MustCompile(pattern),
	}
}

func (m *MetaGetter) Match(s string) (group string, name string, err error) {
	matches := m.regexp.FindStringSubmatch(s)
	if len(matches) == 0 {
		err = fmt.Errorf("'%v' did not match patttern '%v'", s, m.regexp.String())
		return
	}
	if len(matches) == 1 {
		err = fmt.Errorf("can not find group in '%v' with pattern '%v'", s, m.regexp.String())
		return
	}

	group = matches[1]
	if len(matches) < 3 {
		name = matches[0]
	} else {
		name = matches[2]
	}

	return
}

func (m *MetaGetter) String() string {
	return m.regexp.String()
}
