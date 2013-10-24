package patentcluster

import (
	"encoding/csv"
	"fmt"
	"github.com/agonopol/go-stem/stemmer"
	"io"
	"os"
	"strings"
)

var Tagset = make(map[string]int)

/** Given a space separated list of tags, returns a string slice
  of those tags. If stem is true, runs the Porter stemming algorithm
  on each tag
*/
func split_tags(taglist string, stem bool) []string {
	tags := strings.Split(taglist, " ")
	res := []string{}
	for _, s := range tags {
		stem := stemmer.Stem([]byte(s))
		res = append(res, string(stem))
	}
	return res
}

/**
  given a filename, return a map of patent number (string)
  to a slice of the tags ([]string). if stem is true,
  runs the Porter stemming on each of the tags
*/
func Extract_file_contents(filename string, stem bool) [](*Patent) {
	patents := [](*Patent){}
	datafile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer datafile.Close()
	reader := csv.NewReader(datafile)
	reader.Read() // skip first row
	/* loop through file */
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		number := record[0]
		app_date := record[1]
		tags := split_tags(record[2], stem)
		p := makePatent(number, app_date, tags)
		patents = append(patents, p)
		for _, stem := range tags {
			Tagset[string(stem)] += 1
		}
	}
	return patents
}
