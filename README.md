# Cephylopods - A phylogeny of the Kopffüßer


## Introduction

This document provides a detailed description of the methods used in the creation of a comprehensive phylogeny for cephalopods with brain volume data available.
The analysis uses data provided by [Lindgren et al.](http://www.biomedcentral.com/1471-2148/12/129).

All the code is available [on Github](https://github.com/qfma/cephylopods) under a MIT license and may be used to recreate the analysis.


### Method summary


### Sequences used

We use the table published by [Lindgren et al.](http://www.biomedcentral.com/1471-2148/12/129) and obtain the DNA sequences using the genebank IDs. The constructed Cephalopods are stored as Go structs.

```bash
./construct-cephalopods -csv lindgren-paper/Lindgren-et-al-2012-1471-2148-12-129-s2.csv -gob 2015-04-20-cephalopods.gob
```

Now we split the Cephalopods by Gene, because not every Cephalopod has sequences for all genes used in the analysis (missing data). We use only cephalopods that have at least 4 sequenced genes.

```bash
./get-species-subset -gob 2015-04-20-cephalopods.gob -minseq 4
Total number of Cephalopods:  414
Number of valid species:  181
Writing fasta file for gene:  H3
Writing fasta file for gene:  odh
Writing fasta file for gene:  COI
Writing fasta file for gene:  18S
Writing fasta file for gene:  16S
Writing fasta file for gene:  cytb
Writing fasta file for gene:  pax
Writing fasta file for gene:  12S
Writing fasta file for gene:  opsin
Writing fasta file for gene:  28S
```

### Gene alignments

Because we have brain data only from a subset of species, we use only those genes that contain those species, namely

```
12S
16S
18S
COI
```

The 12S, 16S and 18S sequences are non-coding and we can therefore just use the aligned nucleotide sequences.
For the COI sequences, we can produce a codon alignment using TranslatorX (http://translatorx.co.uk/). I modified the perl script so it would use mafft with einsi settings.

#### Excluded sequences

```
Ceph402:topac:537939951 # Excluded, some plant, no brain data
Ceph322:ropal:4003471 # Excluded, sequence degrades at end, no brain data
```

#### Changed

`Ceph269:ocvul:537938871` This is a wrong sequence GID, actually some sort of plant. Added sequence http://www.ncbi.nlm.nih.gov/nuccore/AB052253.1
`>Ceph269:ocvul:11611627`


```bash
# Use MAFFT with EINSI settings (edited perl script)
./translatorx_vLocal.pl -i ../sequences/COI/COI.valid.cleaned.fa -c 5 -o ../sequences/COI/COI.valid.cleaned.aln -t T -p F

# Use standard settings for other genes.
mafft 12S.valid.fa > 12S.valid.aln
mafft 16S.valid.fa > 16S.valid.aln
mafft 18S.valid.fa > 18S.valid.aln
```

#### Trimming

I trimmed the alignment produced by TranslatorX in order to delete gappy bits

>Ceph321:ropac:5353806/1-585 # DELETED beginning

>Deleted first and last two codons

Final result stored in COI.valid.cleaned.trimmed.codon.fa and *.valid.aln.

### Exabayes analysis

In order to analyse the species, we concatenated the alignments for all genes and added gaps for missing data.
```bash
./concatenate-alignments -in sequences/alignments -out 2015-03-28-cephalopods-12S16S18SCOI.aln
[12S 16S 18S COI]
Number of sequences in alignment:  119
Number of sequences in alignment:  171
Number of sequences in alignment:  92
Number of sequences in alignment:  125

# Exabayes partitions
DNA, 12S = 1-669
DNA, 16S = 670-1422
DNA, 18S = 1423-5674
DNA, COI-c1c2 = 5675-6334\3, 5676-6334\3
DNA, COI-c3 = 5677-6334\3

#Change to phylip format
./fasta2phylip -in cephylopods-transfer/2015-04-28-cephalopods-12S16S18SCOI.aln.fa -out cephylopods-transfer/2015-04-28-cephalopods-12S16S18SCOI.aln.phy
```

Now we run Exabayes with the following command:




### RaxML analysis



