package patentcluster

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
    "github.com/agonopol/go-stem/stemmer"
)

var Tagset = make(map[string]int)
var number_of_tags float64
var sqrt_num_tags float64
var visited = make(map[string]int)

/** enumerates all tags in the taglist and inserts them into `Tagset 
    if stem is True, applies the Porter stemming algorithm to all tags
*/
func extract_tags(taglist string, stem bool) {
	tags := strings.Split(taglist, " ")
	for _, s := range tags {
        stem := stemmer.Stem([]byte(s))
		Tagset[string(stem)] += 1
	}
}

/**
  reads buzzx.csv and accumulates all patent tags
  if `stem` is True, applies the Porter stemming algorithm to all tags
*/
func Read_file(filename string, stem bool) {
	/* open buzzx.csv file and start counting tags */
	datafile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer datafile.Close()
	reader := csv.NewReader(datafile)
	number_of_records := 0
	/* loop through file */
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		number_of_records += 1
		tags := record[3]
		extract_tags(tags, stem)
	}
	number_of_tags = float64(len(Tagset))
	sqrt_num_tags = math.Sqrt(number_of_tags)
}

/**
  loops through buzzx.csv and creates a patent instance
  for each row
  if `stem` is True, applies the Porter stemming algorithm to all tags
*/
func Make_patents(filename string, stem bool) [](*Patent) {
    patents := [](*Patent){}
	/* open buzzx.csv file and start counting tags */
	datafile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer datafile.Close()
	reader := csv.NewReader(datafile)
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
		tags := record[3]
		p := makePatent(number, tags)
		patents = append(patents, p)
	}
    return patents
}

/**
  for test purposes, do a pairwise comparison of all patents
*/
func pairwise_run(patents [](*Patent)) {
	vectorfile, _ := os.Create("pairwise.txt")
	w := bufio.NewWriter(vectorfile)
	for _, p1 := range patents {
		for _, p := range patents {
			key := p1.number + p.number
			key2 := p.number + p1.number
			if (visited[key] == 0 || visited[key2] == 0) && p1.number != p.number {
				res := p1.JaccardDistance(p)
				fmt.Fprintln(w, res, p.number, p1.number)
				visited[key] = 1
			}
		}
	}
}

//func main() {
//	fmt.Println("Creating tag set...")
//	Read_file("buzzx.csv", true)
//	fmt.Println("Accumulated", number_of_tags, "tags")
//	fmt.Println("Done creating tag set!")
//	fmt.Println("Making patent instances...")
//    patents := Make_patents("buzzx.csv", true)
//	fmt.Println("Finished", len(patents), "patent instances")
//	p1 := patents[1]
//	p2 := patents[2]
//	fmt.Println(p1.JaccardDistance(p2))
//	//pairwise_run()
//}
