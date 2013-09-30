package patentcluster

import (
	"time"
)

type Patent struct {
	Number     string         // patent_id number
	tags       map[string]int // hash of all the tags associated with this patent
	cluster_id string         // patent_id of the cluster to which this patent belongs
	app_date   time.Time
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
func makePatent(number, app_date string, taglist []string) *Patent {
	p := new(Patent)
	p.tags = make(map[string]int)
	for _, tag := range taglist {
		p.tags[tag] = 1
	}
	p.Number = number
	p.app_date, _ = time.Parse("Jan 02 2006", app_date)
	return p
}

/**
  Returns a string containing all tags
  separated by spaces
*/
func (p *Patent) tags_to_string() string {
	ret := ""
	for tag, _ := range p.tags {
		ret += tag + " "
	}
	return ret[0 : len(ret)-1]
}
