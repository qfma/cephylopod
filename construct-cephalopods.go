package main

import (
	"encoding/csv"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// A gene has a Name an identifier and a fasta formated sequence
type Gene struct {
	Name, ID, Fasta, Seq string
}

// A cephalopod is defined by its species Name and a number of Genes that are
// referenced in a map of Genes.
type Cephalopod struct {
	Name  string
	Genes map[string]*Gene
	Abbrv string
	Valid bool
	ID    string
}

// Set the species name abbreviation for cephalopods
// Very helpful for the following phylogenetics analysis
func (c *Cephalopod) SetAbbrv() {
	split := strings.Split(c.Name, " ")
	// Handle classic species names, such as Homo sapiens -> hsap
	// if the second part of the species is shorter than 3 chars
	// use the whole slice
	if len(split) == 2 {
		if len(split[1]) >= 3 {
			c.Abbrv = strings.ToLower(split[0][:2] + split[1][:3])
		} else {
			c.Abbrv = strings.ToLower(split[0][:2] + split[1])
		}
	}
	// Handle more complex species names, such as
	// Octopus sp. 5 MG 2004 -> ocsp.5mg2004
	if len(split) > 2 {
		c.Abbrv = strings.ToLower(split[0][:2] + strings.Join(split[1:], ""))
	}
}

// Set the species name ID for cephalopods
// Very helpful for the following phylogenetics analysis
func (c *Cephalopod) SetID(i int) {
	c.ID = "Ceph" + strconv.Itoa(i)
}

// Extract sequence from fasta, basically strip the header from the fasta file
func (g *Gene) SetSeq() {
	if g.Fasta != "" {
		g.Seq = strings.Join(strings.Split(g.Fasta, "\n")[1:], "")
	}
}

// Go requires active checking for error conditions which as I can tell always
// looks like this function and thus saves tons of lines.
func check(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

// This function takes a Genbank ID and queries the eutils API and
// retrieves the fasta entry for the corresponding ID
// At the moment a nucleotide Genbank ID is expected.
func GetGenbankFasta(GI string) string {
	fmt.Println("Adding sequence information for GI: ", GI)

	// Construct query URL
	baseurl := "http://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?"
	params := "db=nucleotide&" + "id=" + GI + "&rettype=fasta"

	// Sleep for a little bit in order to reduce server load
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(baseurl + params)
	check(err)
	defer resp.Body.Close()

	fasta, err := ioutil.ReadAll(resp.Body)
	check(err)
	if strings.HasPrefix(string(fasta), "Seq") {
		fmt.Println(string(fasta))
	}
	return strings.TrimRight(string(fasta), "\n")
}

// Takes an CSV file and reads all lines and stores them in a slice
// of slices.
func ReadAllCSV(infile string) [][]string {
	csvFile, err := os.Open(infile)
	check(err)

	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	check(err)

	return records
}

// Cephalopod constructor
// Takes a [][] as the result from reading a CSV file and constructs
// Cephalopod structs from it.
func NewCephalopods(records *[][]string) (map[string]*Cephalopod, []string) {
	cephalopods := make(map[string]*Cephalopod)
	var header []string
	for i, record := range *records {

		if i == 0 {
			// assign header line
			header = record
			continue
		}
		species := strings.TrimSpace(record[0])
		fmt.Println("Adding genes to species", species)
		for i, field := range record {
			field := strings.TrimSpace(field)
			if i == 0 {
				cephalopods[species] = &Cephalopod{Name: species,
					Genes: make(map[string]*Gene)}
				continue
			}

			if i == 1 {
				continue
			}
			gname := strings.TrimSpace(header[i])
			//Control for zeros and empty fields
			if field != "0" && field != "" {
				cephalopods[species].Genes[gname] = &Gene{Name: gname,
					ID:    field,
					Fasta: GetGenbankFasta(field)}
			} else {
				cephalopods[species].Genes[gname] = &Gene{Name: gname,
					ID: field}
			}
			g := cephalopods[species].Genes[gname]
			g.SetSeq()
		}
		c := cephalopods[species]
		c.SetAbbrv()
		c.SetID(i)
	}
	return cephalopods, header
}

func StoreCephalopods(cephalopods map[string]*Cephalopod, fname string) {
	file, err := os.Create(fname)
	check(err)
	defer file.Close()
	enc := gob.NewEncoder(file)
	if err := enc.Encode(cephalopods); err != nil {
		panic("cant encode")
	}
}

func main() {
	//Simple command line arguments
	csvPtr := flag.String("csv", "foo.csv", "A path to a CSV file")
	gobPtr := flag.String("gob", "foo.gob", "A path for the output gob file")
	flag.Parse()
	// Read the CSV and store the cephalopods as a gob file
	records := ReadAllCSV(*csvPtr)
	cephalopods, _ := NewCephalopods(&records)
	StoreCephalopods(cephalopods, *gobPtr)
}
