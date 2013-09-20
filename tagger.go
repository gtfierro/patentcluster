package patentcluster

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

var tagset = make(map[string]int)
var patents = [](*Patent){}
var number_of_tags float64
var sqrt_num_tags float64
var visited = make(map[string]int)

/** enumerates all tags in the taglist and inserts them into `tagset */
func extract_tags(taglist string) {
	tags := strings.Split(taglist, " ")
	for _, s := range tags {
		tagset[s] += 1
	}
}

/**
  reads buzzx.csv and accumulates all patent tags
*/
func read_file(string filename) {
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
		extract_tags(tags)
	}
	number_of_tags = float64(len(tagset))
	sqrt_num_tags = math.Sqrt(number_of_tags)
}

/**
  loops through buzzx.csv and creates a patent instance
  for each row
*/
func make_patents(string filename) {
	/* open buzzx.csv file and start counting tags */
	datafile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
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
			return
		}
		number := record[0]
		tags := record[3]
		p := makePatent(number, tags)
		patents = append(patents, p)
	}
}

/**
  for test purposes, do a pairwise comparison of all patents
*/
func pairwise_run() {
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

func main() {
	fmt.Println("Creating tag set...")
	read_file("buzzx.csv")
	fmt.Println("Accumulated", number_of_tags, "tags")
	fmt.Println("Done creating tag set!")
	fmt.Println("Making patent instances...")
	make_patents("buzzx.csv")
	fmt.Println("Finished", len(patents), "patent instances")
	p1 := patents[1]
	p2 := patents[2]
	fmt.Println(p1.JaccardDistance(p2))
	//pairwise_run()
}
