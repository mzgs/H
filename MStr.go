package H

import "strings"

type Mstr string

func (m *Mstr) Remove(removePairs ...string) {
	s := RemoveFromString(string(*m), removePairs...)
	m.Set(s)
}

func (m *Mstr) Replace(replacementPairs ...string) {
	s := StringReplaceAll(m.String(), replacementPairs...)
	m.Set(s)
}

func (m *Mstr) TrimSpace() {
	s := strings.TrimSpace(m.String())
	m.Set(s)
}

func (m Mstr) Between(strStart, strEnd string) string {
	return GetTextBetween(m.String(), strStart, strEnd)
}

func (m Mstr) Contains(substr string) bool {
	return strings.Contains(m.String(), substr)
}

func (m Mstr) Split(sep string) []string {
	return strings.Split(m.String(), sep)
}

func (m Mstr) Lines() []string {
	return m.Split("\n")
}

func (m Mstr) Index(substr string) int {
	return strings.Index(m.String(), substr)
}

func (m Mstr) LastIndex(substr string) int {
	return strings.LastIndex(m.String(), substr)
}

func (m *Mstr) Set(s string) {
	*m = Mstr(s)
}

func (m *Mstr) String() string {
	return string(*m)
}

func (m Mstr) Print() {
	P(m)
}
