package patentcluster

import (
	"strings"
)

type Patent struct {
	number string
	tags   map[string]int
}

func (p *Patent) JaccardDistance(target *Patent) float64 {
	var count, union float64
	count = 0
	union = float64(len(p.tags))
	for tag, _ := range target.tags {
		if p.tags[tag] > 0 {
			count += 1
		} else {
			union += 1
		}
	}
	return 1 - count/union
}

/**
  given a string representing a patent number and
  a string representing the space-delimited list of
  tags for a patent, returns a reference to a Patent
  object
*/
func makePatent(number, tagstring string) *Patent {
	p := new(Patent)
	p.tags = make(map[string]int)
	for _, tag := range strings.Split(tagstring, " ") {
		p.tags[tag] = 1
	}
	p.number = number
	return p
}
