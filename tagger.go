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
func Extract_file_contents(filename string, stem bool) map[string]([]string) {
	data := make(map[string]([]string))
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
		data[record[0]] = split_tags(record[3], stem)
		for _, stem := range data[record[0]] {
			Tagset[string(stem)] += 1
		}
	}
	return data
}

/**
  loops through buzzx.csv and creates a patent instance
  for each row
  if `stem` is True, applies the Porter stemming algorithm to all tags
*/
func Make_patents(data map[string]([]string)) [](*Patent) {
	patents := [](*Patent){}
	/* open buzzx.csv file and start counting tags */
	for number, taglist := range data {
		p := makePatent(number, taglist)
		patents = append(patents, p)
	}
	return patents
}
