package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// Read a fasta alignment from an infile
// The header is used as the key for the dictionary
// No sanity checks are performed as of now
func ReadFastaAlignment(fname string) map[string]string {

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var header string
	aln := make(map[string]string)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), ">") {
			header = scanner.Text()[1:]
			aln[header] = ""
		} else {
			// Remove potential space characters with nothing
			seq := strings.Replace(strings.ToUpper(scanner.Text()), " ", "", -1)
			aln[header] = aln[header] + seq
		}

	}

	return aln
}

// Takes a map with header/sequences and writes a sequential Phylip formated
// output file. Currently, illegal characters are replaced from the header,
// but the header is not truncated at 10 characters
func WritePhylipAlignment(fname string, aln map[string]string) {

	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sequences := len(aln)
	alnlength := GetFastaAlignmentLength(aln)

	f.WriteString(" " + strconv.Itoa(sequences) + " " + strconv.Itoa(alnlength) + "\n")

	for k, v := range aln {
		// Remove illegal characters from id
		id := ReplaceIllegalChars(k)
		f.WriteString(id + " " + v + "\n")
	}
	return
}

func ReplaceIllegalChars(id string) string {

	// These characters are not allowed and will be replaced with
	// the pipe character.
	illegal := []string{"(", ")", "[", "]", ":", ",", ";"}

	for _, char := range illegal {
		id = strings.Replace(id, char, "|", -1)
	}

	return id
}

func GetFastaAlignmentLength(aln map[string]string) int {

	var alnlength int

	for _, v := range aln {
		alnlength = len(v)
		break
	}
	return alnlength
}

func main() {

	//Simple command line arguments
	infile := flag.String("in", "infile", "An alignment in fasta format")
	outfile := flag.String("out", "outfile", "The name of the phylip output file")
	flag.Parse()

	aln := ReadFastaAlignment(*infile)
	WritePhylipAlignment(*outfile, aln)
}
