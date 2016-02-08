// antha/AnthaStandardLibrary/Packages/enzymes/nucleotides.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// Package for interacting with and manipulating dna sequences in extension to methods available in wtype
package sequences

import (
	"fmt"

	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	//. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"strings"

	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/entrez"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// Check for illegal nucleotides
func Illegalnucleotides(fwdsequence wtype.DNASequence) (pass bool, illegalfound []Thingfound, wobblefound []Thingfound) {
	illegal := "§1234567890-=qweryiop[]sdfhjkl;'zxvbm,./!@£$%^&*()_+" // update to include wobble nucleotides etc
	wobble := "NXBHVDMKSWRY"
	//seq := strings.ToUpper(fwdsequence.Seq)
	if strings.ContainsAny(strings.ToUpper(fwdsequence.Seq), (strings.ToUpper(illegal))) || strings.ContainsAny(fwdsequence.Seq, strings.ToLower(illegal)) == true {
		fmt.Println(pass)
		pass = false
		fmt.Println("Contains illegal characters")
		illegalarray := strings.Split(illegal, "")
		//	fmt.Println("iiiiiilllllllllegal array!!!", illegalarray)
		illegalfound = Findallthings((strings.ToUpper(fwdsequence.Seq)), illegalarray)
		//fmt.Println(len(illegalfound))

	}

	if strings.ContainsAny(strings.ToUpper(fwdsequence.Seq), (strings.ToUpper(wobble))) == true {
		pass = false
		fmt.Println("Contains wobble nucleotides")
		wobblearray := strings.Split(wobble, "")
		//fmt.Println("wobble array!!!", wobblearray)
		wobblefound = Findallthings((strings.ToUpper(fwdsequence.Seq)), wobblearray)
		//fmt.Println(len(wobblefound))

	} else {
		pass = true
		fmt.Println("illegal characters pass")
	}
	//wtype.Makeseq(Foldername, &sequence)

	return pass, illegalfound, wobblefound
}

/*
func Rev(s string) string {
	r := ""

	for i := len(s) - 1; i >= 0; i-- {
		r += string(s[i])
	}

	return r
}
func Comp(s string) string {
	r := ""

	m := map[string]string{
		"A": "T",
		"T": "A",
		"U": "A",
		"C": "G",
		"G": "C",
		"Y": "R",
		"R": "Y",
		"W": "W",
		"S": "S",
		"K": "M",
		"M": "K",
		"D": "H",
		"V": "B",
		"H": "D",
		"B": "V",
		"N": "N",
		"X": "X",
	}

	for _, c := range s {
		r += m[string(c)]
	}

	return r
}

// Reverse Complement
func RevComp(s string) string {
	return Comp(Rev(s))
}
*/
var WobbleMap = map[string][]string{
	"A": []string{"A"},
	"T": []string{"T"},
	"U": []string{"U"},
	"C": []string{"C"},
	"G": []string{"G"},
	"Y": []string{"C", "T"},
	"R": []string{"A", "G"},
	"W": []string{"A", "T"},
	"S": []string{"G", "C"},
	"K": []string{"G", "T"},
	"M": []string{"A", "C"},
	"D": []string{"A", "G", "T"},
	"V": []string{"A", "C", "G"},
	"H": []string{"A", "C", "T"},
	"B": []string{"C", "G", "T"},
	"N": []string{"A", "T", "C", "G"},
	"X": []string{"A", "T", "C", "G"},
	"-": []string{"-", "."},
	".": []string{"-", "."},
}

func Wobble(seq string) (alloptions []string) {

	arrayofarray := make([][]string, 0)
	for _, character := range seq {

		optionsforcharacterx := WobbleMap[string(character)]
		arrayofarray = append(arrayofarray, optionsforcharacterx)
		//fmt.Println("arrayofarray", arrayofarray)
	}

	alloptions = AllCombinations(arrayofarray)

	return
}

var Nucleotidegpermol = map[string]float64{
	"A":    313.2,
	"T":    304.2,
	"C":    289.2,
	"G":    329.2,
	"N":    303.7,
	"dATP": 491.2,
	"dCTP": 467.2,
	"dGTP": 507.2,
	"dTTP": 482.2,
	"dNTP": 487.0,
}

// Calculate global GC content
func GCcontent(fwdsequence string) (Percentage float64) {
	numberofAs := strings.Count(fwdsequence, "A")
	numberofTs := strings.Count(fwdsequence, "T")
	numberofCs := strings.Count(fwdsequence, "C")
	numberofGs := strings.Count(fwdsequence, "G")

	gc := (numberofCs + numberofGs)
	all := (numberofCs + numberofGs + numberofAs + numberofTs)

	percentage := float64(gc) / float64(all)
	Percentage = percentage
	return Percentage
}

// Calculate local GC content in a sliding window
func localGCContent(fwdsequence string, window int, shift int) (Pc []float64) {
	incs := len(fwdsequence) / shift
	pos := 0
	Pc = make([]float64, 0)
	for i := 0; i < (incs - 1); i++ {
		region := fwdsequence[pos : pos+window]
		gc := GCcontent(region)
		Pc = append(Pc, gc)
		pos += shift
	}
	return Pc
}

//Calculate Molecular weight of DNA
func MassDNA(fwdsequence string, phosphate5prime bool, doublestranded bool) (mw float64) {
	numberofAs := strings.Count(fwdsequence, "A")
	numberofTs := strings.Count(fwdsequence, "T")
	numberofCs := strings.Count(fwdsequence, "C")
	numberofGs := strings.Count(fwdsequence, "G")
	massofAs := (float64(numberofAs) * Nucleotidegpermol["A"])
	massofTs := (float64(numberofTs) * Nucleotidegpermol["T"])
	massofCs := (float64(numberofCs) * Nucleotidegpermol["C"])
	massofGs := (float64(numberofGs) * Nucleotidegpermol["G"])
	mw = (massofAs + massofTs + massofCs + massofGs)
	if phosphate5prime == true {
		mw = mw + 79.0 // extra for phosphate left at 5' end following digestion, not relevant for primer extension
	}
	if doublestranded == true {
		mw = 2 * mw
	}
	return mw
}

// Calclulate number of moles of a mass of DNA
func MolesDNA(mass wunit.Mass, mw float64) (moles float64) {
	massSI := mass.SIValue()
	moles = massSI / mw
	return moles
}

// calculate molar concentration of DNA sample
func GtoMolarConc(conc wunit.Concentration, mw float64) (molesperL float64) {
	concSI := conc.SIValue()
	molesperL = concSI / mw
	return molesperL
}

func MoletoGConc(molarconc float64, mw float64) (gperL wunit.Concentration) {
	gperl := molarconc * mw
	gperL = wunit.NewConcentration(gperl, "g/L")
	return gperL
}

func Moles(conc wunit.Concentration, mw float64, vol wunit.Volume) (moles float64) {
	molesperL := GtoMolarConc(conc, mw)
	moles = molesperL * vol.SIValue()
	return moles
}

func Revarrayorder(array []string) (newarray []string) {
	//for _, slice := range array {
	newarray = make([]string, 0)
	r := ""
	for i := len(array) - 1; i >= 0; i-- {
		r = string(array[i])
		newarray = append(newarray, r)
		//newarray += array()
	}
	return newarray
}

// utility function to check if item is in list. move elsewhere?
func isInList(item string, list []string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// This searches the entrez database with a query, returning the sequence.
//func retrieveVector(vector string) {
//entrez.RetrieveRecords(vector, "nucleotide", 1, "gb", "mark.gb")
//fasta = fmt.Sprintf("%s%c%s", anthapath.Dirpath(), os.PathSeparator, "mark.fasta")
//x := parser.FastaParse(fasta)
//fmt.Println(x)
//}

// This simulates the sequence assembly reaction to validate if parts will synthesise with intended manufacturer.
// Does not validate construct assembly so should be used in conjunction with enzymes.Assemblysimulator()
func ValidateSynthesis(parts []wtype.DNASequence, vector string, manufacturer string) bool {

	// type conversions
	a, _ := SynthesisStandards[manufacturer]["RepeatMax"].(int)
	b, _ := SynthesisStandards[manufacturer]["MinOrder"].(int)
	c, _ := SynthesisStandards[manufacturer]["MinLength"].(int)
	d, _ := SynthesisStandards[manufacturer]["MaxLength"].(int)
	vectorR, _ := SynthesisStandards[manufacturer]["Vector"].([]string)

	// min total order. Needs to be before local gc calc
	total := 0
	for _, part := range parts {
		total += len(part.Seq)
	}
	if total < b {
		fmt.Println("Warning: Total length of parts is less than", manufacturer, "minimum order requirement")
		return false
	}

	// check if vector is appropriate
	if isInList(vector, vectorR) == false {
		fmt.Println("Warning: Non-standard vector used for", manufacturer,
			"synthesis. Please see manufacturer instructions for the standard vector or the use of custom vectors")
	}

	for _, part := range parts {
		GCC := GCcontent(part.Seq)              // global gc
		gc := localGCContent(part.Seq, 100, 50) // local gc

		// check lengths of seq, repeat content and global gc content of each part
		if len(part.Seq) < c {
			fmt.Println("Warning:", part.Nm, "is too short to synthesise")
			//return false
		} else if len(part.Seq) > d {
			fmt.Println("Warning:", part.Nm, "is too long to sythesise")
			//return false
		} else if strings.Contains(strings.ToUpper(part.Seq), strings.Repeat("A", a)) || strings.Contains(strings.ToUpper(part.Seq), strings.Repeat("T", a)) ||
			strings.Contains(strings.ToUpper(part.Seq), strings.Repeat("C", a)) || strings.Contains(strings.ToUpper(part.Seq), strings.Repeat("A", a)) == true {
			fmt.Println("Warning:", part.Nm, "is highly repetetive and unsuitable for synthesis")
			//return false
		} else if GCC > 0.65 || GCC < 0.40 {
			fmt.Println("Warning: GC content of", part.Nm, "is too high or low for synthesis")
		}

		// check local gc content of each part in 100bp sliding window
		for _, v := range gc {
			if v < 0.25 || v > 0.80 {
				fmt.Println("Warning: Local GC content too high or low in", part.Nm)
				//return false
			}
		}
	}

	fmt.Println("Your", manufacturer, "DNA synthesis order should work")
	return true
}

var SynthesisStandards = map[string]map[string]interface{}{
	"Gen9": map[string]interface{}{
		"Vector":    []string{"pG9m-2"},
		"MaxLength": 10000,
		"MinLength": 400,
		"RepeatMax": 70,
		"MinOrder":  20000,
	},
	"DNA20": map[string]interface{}{
		"Vector":    []string{"pM265", "DSFS"},
		"MaxLength": 10000,
		"MinLength": 0,
		"RepeatMax": 70,
		"MinOrder":  0,
	},
	"GenScript": map[string]interface{}{
		"Vector":    []string{"pUC57", "pJ221", "pJ321", "pJ201"},
		"MaxLength": 10000,
		"MinLength": 0,
		"RepeatMax": 70,
		"MinOrder":  0,
	},
}
