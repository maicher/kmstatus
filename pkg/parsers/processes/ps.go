package processes

import "strings"

type PS struct {
	m map[string]int
}

func (p PS) Uniq() int {
	return len(p.m)
}

func (p PS) Find(phrase string) int {
	count, ok := p.m[phrase]
	if !ok {
		return 0
	}

	return count
}

func (p PS) FindByPrefix(phrase string) int {
	var count int
	for k, v := range p.m {
		if strings.HasPrefix(k, phrase) {
			count = count + v
		}
	}

	return count
}
