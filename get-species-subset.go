package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
)

// A gene has a Name an identifier and a fasta formated sequence
type Gene struct {
	Name, ID, Fasta, Seq string
}

// A cephalod is defined by its species Name and a number of Genes that are
// referenced in a map of Genes.
type Cephalopod struct {
	Name  string
	Genes map[string]*Gene
	Abbrv string
	Valid bool
	ID    string
}

// Go requires active checking for error conditions which as I can tell always
// looks like this function and thus saves tons of lines.
func check(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

func (c *Cephalopod) validate(minSeq *int) {
	var count int
	for _, g := range c.Genes {
		if g.Seq != "" {
			count += 1
		}
	}
	if count >= *minSeq {
		c.Valid = true
	}
}

// Takes a map containing cephalopods with the associated Genes and writes the
// sequence for the same gene for all cephalopods into a fasta file if the cephalopod
// is valid as determined by command line option
func WriteValidCephalopods(cephalopods map[string]*Cephalopod, genes []string) {
	for _, g := range genes {
		fmt.Println("Writing fasta file for gene: ", g)
		file, err := os.Create(g + ".valid.fa")
		check(err)
		defer file.Close()
		for _, c := range cephalopods {
			if c.Valid == true {
				if c.Genes[g].Seq != "" {
					file.WriteString(">" + c.ID + ":" + c.Abbrv + ":" + c.Genes[g].ID + "\n")
					file.WriteString(c.Genes[g].Seq + "\n")
				}
			}
		}
	}
}

func LoadCephalopods(fname string) (cephalopods map[string]*Cephalopod) {
	f, err := os.Open(fname)
	check(err)
	defer f.Close()

	enc := gob.NewDecoder(f)
	if err := enc.Decode(&cephalopods); err != nil {
		panic("Error: Can't decode file!")
	}
	return cephalopods
}

func GetAllGenes(c map[string]*Cephalopod) []string {
	var Genes []string
	for _, v := range c {
		for k, _ := range v.Genes {
			Genes = append(Genes, k)
		}
		break
	}
	return Genes
}

func main() {
	//Simple command line arguments
	gobPtr := flag.String("gob", "cephalopods.gob", "A path to a gob file containing cephalopod structs")
	minSeqPtr := flag.Int("minseq", 4, "Minimum number of available genes")
	flag.Parse()
	cephalopods := LoadCephalopods(*gobPtr)
	genes := GetAllGenes(cephalopods)
	var allvalid int
	for _, c := range cephalopods {
		c.validate(minSeqPtr)
		if c.Valid == true {
			// fmt.Println(c.Name)
			allvalid += 1
		}
	}
	fmt.Println("Total number of Cephalopods: ", len(cephalopods))
	fmt.Println("Number of valid species: ", allvalid)
	// for _, g := range cephalopods["Doryteuthis opalescens"].Genes {
	// 	fmt.Println(g)
	// }
	WriteValidCephalopods(cephalopods, genes)
}
