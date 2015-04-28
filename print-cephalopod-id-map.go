package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"sort"
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

func WriteSortedMap(fname string, cephalopods map[string]*Cephalopod) {
	var IDs []string

	for _, c := range cephalopods {
		IDs = append(IDs, c.ID)
	}

	sort.Strings(IDs)

	file, err := os.Create(fname)
	check(err)
	defer file.Close()

	for _, id := range IDs {
		for _, c := range cephalopods {
			if c.ID == id {
				file.WriteString(c.ID + "\t" + c.Abbrv + "\t" + c.Name + "\n")
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

func main() {
	//Simple command line arguments
	gobPtr := flag.String("gob", "cephalopods.gob", "A path to a gob file containg cephalopod structs")
	minSeqPtr := flag.String("map", "ceph.map", "Cephalopod mapping")
	flag.Parse()
	cephalopods := LoadCephalopods(*gobPtr)
	WriteSortedMap(*minSeqPtr, cephalopods)
}
