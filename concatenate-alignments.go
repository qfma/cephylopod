package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Go requires active checking for error conditions which as I can tell always
// looks like this function and thus saves tons of lines.
func check(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

func ReadFastaAlignment(fname string) map[string][]string {
	f, err := os.Open(fname)
	check(err)
	defer f.Close()
	var counter int
	var species string
	aln := make(map[string][]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), ">") {
			counter += 1
			species = strings.Split(scanner.Text(), ":")[0][1:]
			aln[species] = append(aln[species], scanner.Text()[1:])
			aln[species] = append(aln[species], "")
		} else {
			seq := strings.Replace(strings.ToUpper(scanner.Text()), " ", "", -1)
			aln[species][1] = aln[species][1] + seq
		}
	}
	fmt.Println("Number of sequences in alignment: ", counter)
	return aln
}

// Takes a map containg cephalopods with the associated Genes and writes the
// sequence for the same gene for all cephalopods into a fasta file if the cephalopod
// is valid as determined by command line option
func WriteFastaAlignment(fname string, aln map[string]string) {
	file, err := os.Create(fname)
	check(err)
	defer file.Close()
	for k, v := range aln {
		file.WriteString(">" + k + "\n")
		file.WriteString(v + "\n")
	}
}

// Takes a list of os.FileInfo structs and takes the filename and splits
// the filename by "." and returns the first element (basename) of the file
func GetAllGenes(files []os.FileInfo) []string {
	var genes []string
	for _, f := range files {
		g := strings.Split(f.Name(), ".")[0]
		genes = append(genes, g)
	}
	return genes
}

// Best to implement Alignment type at some point
func GetAlignments(files *[]string) map[string]map[string][]string {
	alns := make(map[string]map[string][]string)
	for _, f := range *files {
		g := filepath.Base(f)
		g = strings.Split(g, ".")[0]
		aln := ReadFastaAlignment(f)
		alns[g] = aln
	}
	return alns
}

func GetAllSpecies(alns map[string]map[string][]string) map[string]bool {
	species := make(map[string]bool)
	for _, v := range alns {
		for k, _ := range v {
			if _, ok := species[k]; ok {
				continue
			} else {
				species[k] = true
			}
		}
	}
	return species
}

func GetAlnLength(aln map[string][]string) int {
	var AlnLength int
	for _, v := range aln {
		AlnLength = len(v[1])
		break
	}
	return AlnLength
}

func CatAln(alns map[string]map[string][]string, genes *[]string, species map[string]bool) map[string]string {
	cataln := make(map[string]string)
	for _, g := range *genes {
		AlnLength := GetAlnLength(alns[g])
		for s, _ := range species {
			if seq, ok := alns[g][s]; ok {
				cataln[s] = cataln[s] + seq[1]
			} else {
				cataln[s] = cataln[s] + strings.Repeat("-", AlnLength)
			}
		}
	}
	return cataln
}

func main() {
	//Simple command line arguments
	alnPtr := flag.String("in", "alignments", "A folder with fasta alignments")
	outfile := flag.String("out", "concatenated.fa", "The name of the concatenated output file")
	flag.Parse()

	files, _ := ioutil.ReadDir(*alnPtr)
	fmt.Println(*alnPtr)
	genes := GetAllGenes(files)
	fmt.Println(genes)

	var alnFiles []string
	for _, f := range files {
		path := *alnPtr + "/" + f.Name()
		alnFiles = append(alnFiles, path)
	}
	alns := GetAlignments(&alnFiles)
	species := GetAllSpecies(alns)
	fmt.Println(len(species))
	cataln := CatAln(alns, &genes, species)

	WriteFastaAlignment(*outfile, cataln)
	// RaxML partitioning
	CodingGenes := map[string]bool{
		"COI":   true,
		"cytb":  true,
		"H3":    true,
		"odh":   true,
		"opsin": true,
		"pax":   true,
	}

	var start, stop int
	start = 1

	for _, g := range genes {
		length := GetAlnLength(alns[g])
		stop = start + length - 1
		if _, ok := CodingGenes[g]; ok {
			fmt.Printf("DNA, %v-c1c2 = %v-%v\\3, %v-%v\\3\n", g, start, stop, start+1, stop)
			fmt.Printf("DNA, %v-c3 = %v-%v\\3\n", g, start+2, stop)
			start = stop + 1
		} else {
			fmt.Printf("DNA, %v = %v-%v\n", g, start, stop)
			start = stop + 1
		}
	}
}
