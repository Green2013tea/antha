// antha/AnthaStandardLibrary/Packages/enzymes/Ligation.go: Part of the Antha language
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

// Package for working with enzymes; in particular restriction enzymes
package enzymes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/plasmid"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func fragmentsToDNASequences(fragments []Digestedfragment) (sequences []wtype.DNASequence, err error) {

	var errs []string

	for i, fragment := range fragments {
		seq, err := fragment.ToDNASequence("fragment" + strconv.Itoa(i))
		sequences = append(sequences, seq)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, ";"))
	}
	return
}

// fragmentsFormPlasmid checks if the two fragments can join both ends to forma plasmid
func fragmentsFormPlasmid(upfragment, downfragment Digestedfragment) bool {
	if strings.EqualFold(sequences.RevComp(upfragment.BottomStickyend_5prime), downfragment.TopStickyend_5prime) && strings.EqualFold(sequences.RevComp(downfragment.BottomStickyend_5prime), upfragment.TopStickyend_5prime) {
		return true
	}
	if strings.EqualFold(upfragment.BottomStickyend_5prime, sequences.RevComp(downfragment.BottomStickyend_5prime)) && strings.EqualFold(downfragment.TopStickyend_5prime, sequences.RevComp(upfragment.TopStickyend_5prime)) {
		return true
	}
	return false
}

// add code to check for duplicate sticky end parts to prevent simulation of assembly of all backbones
// func uniqueEnds (upFragment, downFragment Digestedfragment, endsUsedSoFar []string) bool {}
// or even better to check for presence of correct antibiotic resistance

func joinTwoParts(upstreampart []Digestedfragment, downstreampart []Digestedfragment) (assembledfragments []Digestedfragment, plasmidproducts []wtype.DNASequence, err error) {

	var sequencestojoin []string

	for _, upfragment := range upstreampart {
		for _, downfragment := range downstreampart {
			if strings.EqualFold(sequences.RevComp(upfragment.BottomStickyend_5prime), downfragment.TopStickyend_5prime) && strings.EqualFold(sequences.RevComp(downfragment.BottomStickyend_5prime), upfragment.TopStickyend_5prime) {
				sequencestojoin = append(sequencestojoin, upfragment.Topstrand, downfragment.Topstrand)
				dnastring := strings.Join(sequencestojoin, "")
				fullyassembledfragment := wtype.DNASequence{Nm: "simulatedassemblysequence", Seq: dnastring, Plasmid: true}
				plasmidproducts = append(plasmidproducts, fullyassembledfragment)
				sequencestojoin = make([]string, 0)
			}
			if strings.EqualFold(upfragment.BottomStickyend_5prime, sequences.RevComp(downfragment.BottomStickyend_5prime)) && strings.EqualFold(downfragment.TopStickyend_5prime, sequences.RevComp(upfragment.TopStickyend_5prime)) {
				sequencestojoin = append(sequencestojoin, upfragment.Topstrand, downfragment.Bottomstrand)
				dnastring := strings.Join(sequencestojoin, "")
				fullyassembledfragment := wtype.DNASequence{Nm: "simulatedassemblysequence", Seq: dnastring, Plasmid: true}
				plasmidproducts = append(plasmidproducts, fullyassembledfragment)
				sequencestojoin = make([]string, 0)
			}
			if strings.EqualFold(sequences.RevComp(upfragment.BottomStickyend_5prime), downfragment.TopStickyend_5prime) {
				sequencestojoin = append(sequencestojoin, upfragment.Topstrand, downfragment.Topstrand)
				dnastring := strings.Join(sequencestojoin, "")
				assembledfragment := Digestedfragment{dnastring, "", upfragment.TopStickyend_5prime, downfragment.TopStickyend_3prime, downfragment.BottomStickyend_5prime, upfragment.BottomStickyend_3prime}
				assembledfragments = append(assembledfragments, assembledfragment)
				sequencestojoin = make([]string, 0)
			}
			if strings.EqualFold(upfragment.BottomStickyend_5prime, sequences.RevComp(downfragment.BottomStickyend_5prime)) {
				sequencestojoin = append(sequencestojoin, upfragment.Topstrand, downfragment.Bottomstrand)
				dnastring := strings.Join(sequencestojoin, "")
				assembledfragment := Digestedfragment{dnastring, "", upfragment.TopStickyend_5prime, downfragment.BottomStickyend_3prime, downfragment.TopStickyend_5prime, upfragment.BottomStickyend_3prime}
				assembledfragments = append(assembledfragments, assembledfragment)
				sequencestojoin = make([]string, 0)
			}
		}
	}
	if len(assembledfragments) == 0 && len(plasmidproducts) == 0 {
		errstr := fmt.Sprintln("fragments aren't compatible, check ends.")
		err = fmt.Errorf(errstr)
	}
	return assembledfragments, plasmidproducts, err
}

// key function for returning arrays of partially assembled fragments and fully assembled fragments from performing typeIIS assembly on a vector and a part
func Jointwopartsfromsequence(vector wtype.DNASequence, part1 wtype.DNASequence, enzyme wtype.TypeIIs) (assembledfragments []Digestedfragment, plasmidproducts []wtype.DNASequence) {
	doublestrandedpart1 := MakedoublestrandedDNA(part1)
	digestedpart1 := DigestionPairs(doublestrandedpart1, enzyme)

	doublestrandedvector := MakedoublestrandedDNA(vector)
	digestedvector := DigestionPairs(doublestrandedvector, enzyme)

	assembledfragments, plasmidproducts, _ = joinTwoParts(digestedvector, digestedpart1)

	return assembledfragments, plasmidproducts
}

func rotateVector(vector wtype.DNASequence, enzyme wtype.TypeIIs) (wtype.DNASequence, error) {
	rotatedVector := vector.Dup()

	// the purpose of this is to ensure the RE sites go ---> xxxx <---

	if len(vector.Seq) == 0 {
		return rotatedVector, fmt.Errorf("No Sequence found for %s so cannot rotate", vector.Nm)
	}

	restrictionSites := sequences.FindAll(&rotatedVector, &wtype.DNASequence{Nm: enzyme.Name, Seq: enzyme.RecognitionSequence})

	if len(restrictionSites.Positions) > 2 {
		err := fmt.Errorf("must have 2 restriction sites to rotate vector. %d %s sites found in vector %s - cannot rotate", len(restrictionSites.Positions), enzyme.Name, vector.Nm)
		return rotatedVector, err
	}

	var fwdSites []sequences.PositionPair
	var revSites []sequences.PositionPair

	for i := range restrictionSites.Positions {
		if !restrictionSites.Positions[i].Reverse {
			fwdSites = append(fwdSites, restrictionSites.Positions[i])
		} else {
			revSites = append(revSites, restrictionSites.Positions[i])
		}
	}
	if len(revSites) > 1 {
		return rotatedVector, fmt.Errorf("%d reverse sites for %s found in vector %s", len(revSites), enzyme.Name, vector.Name())
	}

	if len(fwdSites) > 1 {
		return rotatedVector, fmt.Errorf("%d forward sites for %s found in vector %s", len(fwdSites), enzyme.Name, vector.Name())
	}

	fwdStart, _ := fwdSites[0].Coordinates(wtype.CODEFRIENDLY, wtype.IGNOREDIRECTION)

	rotatedVector = sequences.Rotate(rotatedVector, fwdStart, false)

	return rotatedVector, nil
}

func allPartOrders(parts []wtype.DNASequence) (allCombos [][]wtype.DNASequence) {

	partNumToSeq := make(map[int]wtype.DNASequence)
	var nums []int

	for i := range parts {
		partNumToSeq[i] = parts[i]
		nums = append(nums, i)
	}

	numbercombos := permutations(nums)

	allCombos = make([][]wtype.DNASequence, len(numbercombos))

	for i := range allCombos {
		var combo []wtype.DNASequence
		for _, num := range numbercombos[i] {
			combo = append(combo, partNumToSeq[num])
		}
		allCombos[i] = combo
	}

	return
}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

// FindAllAssemblyProducts will return all assembly products from a set of assembly part sequences. Unlike, JoinXnumberofparts the order of the parts is not important.
func FindAllAssemblyProducts(vector wtype.DNASequence, partsInAnyOrder []wtype.DNASequence, enzyme wtype.TypeIIs) (assembledfragments []Digestedfragment, plasmidproducts []wtype.DNASequence, err error) {

	var errs []string

	var allPartCombos [][]wtype.DNASequence = allPartOrders(partsInAnyOrder)

	for _, partOrder := range allPartCombos {
		partialassemblies, plasmids, _, err := JoinXNumberOfParts(vector, partOrder, enzyme)
		if err != nil {
			errs = append(errs, err.Error())
		}
		for i := range partialassemblies {
			assembledfragments = append(assembledfragments, partialassemblies[i])
		}
		for j := range plasmids {
			plasmidproducts = append(plasmidproducts, plasmids[j])
		}

		// if too many combinations, return as soon as a valid plasmid is found to save memory.
		if len(allPartCombos) > 24 {

			if len(plasmidproducts) > 0 {
				for _, plasmidProduct := range plasmidproducts {
					_, oris, markers, _ := plasmid.ValidPlasmid(plasmidProduct)

					if len(oris) > 0 && len(markers) > 0 {
						// if no duplicate oris and markers, likely to be a good assembly and not vector double assembly.
						if len(search.RemoveDuplicates(oris)) == len(oris) && len(search.RemoveDuplicates(markers)) == len(markers) {
							plasmidproducts = search.RemoveDuplicateSequences(plasmidproducts)
							return assembledfragments, plasmidproducts, nil
						}
					}
				}
			}
		}
	}

	plasmidproducts = search.RemoveDuplicateSequences(plasmidproducts)

	if len(errs) > 0 && len(plasmidproducts) == 0 {
		err = fmt.Errorf(strings.Join(errs, ";"))
	}

	return
}

// JoinXNumberOfParts simulates assembly of a Vector and a list of parts in order using a specified TypeIIs restriction enzyme.
// Returns an array of partially assembled fragments and fully assembled plasmid products and any error in attempting to assemble the parts.
func JoinXNumberOfParts(vector wtype.DNASequence, partsinorder []wtype.DNASequence, enzyme wtype.TypeIIs) (assembledfragments []Digestedfragment, plasmidproducts []wtype.DNASequence, inserts []wtype.DNASequence, err error) {

	var newerr error
	if vector.Seq == "" {
		err = fmt.Errorf("No Vector sequence found for vector %s", vector.Nm)
		return assembledfragments, plasmidproducts, inserts, err
	}
	// there are two cases: either the vector comes in same way parts do
	// i.e. SAPI--->xxxx<---IPAS
	// OR it comes in the other way round
	// i.e. xxxx<---IPASyyyySAPI--->zzzz
	// we have either to rotate the vector or tolerate this
	// probably best to rotate first

	rotatedvector, err := rotateVector(vector, enzyme)

	if err != nil {
		return assembledfragments, plasmidproducts, inserts, err
	}

	if len(partsinorder) == 0 {
		return nil, nil, inserts, fmt.Errorf("No parts found")
	}
	if len(partsinorder[0].Seq) == 0 {
		name := partsinorder[0].Nm
		errorstring := name + " has no sequence"
		err = fmt.Errorf(errorstring)
		return assembledfragments, plasmidproducts, inserts, err
	}
	doublestrandedpart := MakedoublestrandedDNA(partsinorder[0])
	// initialise assembledFragements with first digested part
	assembledfragments = DigestionPairs(doublestrandedpart, enzyme)

	for i := 1; i < len(partsinorder); i++ {
		if len(partsinorder[i].Seq) == 0 {
			name := partsinorder[i].Nm
			errorstring := name + " has no sequence"
			err = fmt.Errorf(errorstring)
			return assembledfragments, plasmidproducts, inserts, err
		}

		doublestrandedpart = MakedoublestrandedDNA(partsinorder[i])
		digestedpart := DigestionPairs(doublestrandedpart, enzyme)

		assembledfragments, plasmidproducts, newerr = joinTwoParts(assembledfragments, digestedpart)

		if newerr != nil {
			message := fmt.Sprint(partsinorder[i-1].Nm, " and ", partsinorder[i].Nm, ": ", newerr.Error())
			err = fmt.Errorf(message)
			return
		}
	}

	// this is the insert
	insertFragments := assembledfragments

	// now join fragment to vector
	doublestrandedvector := MakedoublestrandedDNA(rotatedvector)

	digestedvector := DigestionPairs(doublestrandedvector, enzyme)

	assembledfragments, plasmidproducts, newerr = joinTwoParts(digestedvector, insertFragments)
	if newerr != nil {
		message := fmt.Sprint(vector.Nm, " and ", partsinorder[0].Nm, ": ", newerr.Error())
		err = fmt.Errorf(message)
		return
	}

	partnames := make([]string, 0)

	for _, part := range partsinorder {
		partnames = append(partnames, part.Nm)
	}

	for _, plasmidproduct := range plasmidproducts {

		plasmidproduct.Nm = vector.Nm + "_" + strings.Join(partnames, "_")
	}
	var errs []string
	for _, vectorFragment := range digestedvector {
		for _, insertFragment := range insertFragments {
			if fragmentsFormPlasmid(vectorFragment, insertFragment) {
				insert, err := insertFragment.ToDNASequence("InsertSequence")
				if err == nil {
					inserts = append(inserts, insert)
				} else {
					errs = append(errs, err.Error())
				}
			}
		}
	}

	if len(errs) > 0 && len(inserts) == 0 {
		return assembledfragments, plasmidproducts, inserts, fmt.Errorf("no valid inserts expected to form: %s", strings.Join(errs, ";"))
	}

	return assembledfragments, plasmidproducts, inserts, nil
}

func names(seqs []wtype.DNASequence) []string {
	var nms []string
	for i := range seqs {
		nms = append(nms, seqs[i].Nm)
	}
	return nms
}

// Assemblyparameters are parameters used by the AssemblySimulator function.
type Assemblyparameters struct {
	Constructname string              `json:"construct_name"`
	Enzymename    string              `json:"enzyme_name"`
	Vector        wtype.DNASequence   `json:"vector"`
	Partsinorder  []wtype.DNASequence `json:"parts_in_order"`
}

// ToString returns a summary of the names of all components specified in the Assemblyparameters variable
func (assemblyparameters Assemblyparameters) ToString() string {
	return fmt.Sprintf("Assembly: %s, Enzyme: %s, Vector: %s, Parts: %s", assemblyparameters.Constructname, assemblyparameters.Enzymename, assemblyparameters.Vector.Nm, strings.Join(names(assemblyparameters.Partsinorder), ";"))

}

// AssemblySummary returns a summary of multiple Assemblyparameters separated by a line break for each
func AssemblySummary(params []Assemblyparameters) string {

	var summaries []string
	for _, assembly := range params {
		summaries = append(summaries, assembly.ToString())
	}

	return strings.Join(summaries, "\n")
}

// Insert will find the inserted DNA region as a linear DNA sequence from a set of assembly parameters and the assembled sequence.
func (assemblyParameters Assemblyparameters) Insert(result wtype.DNASequence) (insert wtype.DNASequence, err error) {

	// fetch enzyme properties
	enzymename := strings.ToUpper(assemblyParameters.Enzymename)

	enzyme, err := lookup.TypeIIsLookup(enzymename)

	if err != nil {
		return insert, fmt.Errorf("failure calculating insert: %s", err.Error())
	}

	_, _, inserts, err := JoinXNumberOfParts(assemblyParameters.Vector, assemblyParameters.Partsinorder, enzyme)
	if err != nil {
		return insert, fmt.Errorf("failure calculating insert: %s", err.Error())
	}

	var validInserts []wtype.DNASequence

	for i := range inserts {
		if len(sequences.FindAll(&result, &inserts[i]).Positions) == 1 {
			validInserts = append(validInserts, inserts[0])
		}
	}
	if len(validInserts) == 0 {
		return insert, fmt.Errorf("no insert sequences found which are present in assembled sequence %s. Found these: %v", result.Name(), inserts)
	}

	if len(validInserts) == 1 {
		return validInserts[0], nil
	}

	if len(validInserts) > 1 {
		return biggest(validInserts), nil
	}

	return insert, fmt.Errorf("no insert sequences found which are present in assembled sequence %s. ", result.Name())
}

// Assemblysimulator simulate assembly of Assemblyparameters: returns status, number of correct assemblies, any restriction sites found, new DNA Sequences and an error.
func Assemblysimulator(assemblyparameters Assemblyparameters) (s string, successfulassemblies int, sites []Restrictionsites, newDNASequences []wtype.DNASequence, err error) {

	// fetch enzyme properties
	enzymename := strings.ToUpper(assemblyparameters.Enzymename)

	enzyme, err := lookup.TypeIIsLookup(enzymename)

	if err != nil {
		return s, successfulassemblies, sites, newDNASequences, err
	}

	if enzyme.Class != "TypeIIs" {
		s = fmt.Sprint(enzymename, ": Incorrect Enzyme or no enzyme specified")
		err = fmt.Errorf(s)
		return s, successfulassemblies, sites, newDNASequences, err
	}
	var failedAssemblies []Digestedfragment
	var plasmidProducts []wtype.DNASequence

	if len(assemblyparameters.Partsinorder) > 6 {
		failedAssemblies, plasmidProducts, _, err = JoinXNumberOfParts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)
	} else {
		failedAssemblies, plasmidProducts, err = FindAllAssemblyProducts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)
	}
	if err != nil {
		err = fmt.Errorf("Failure Joining fragments after digestion: %s", err.Error())
		s = err.Error()
		return s, successfulassemblies, sites, plasmidProducts, err
	}

	if len(plasmidProducts) == 1 {
		sites = Restrictionsitefinder(plasmidProducts[0], []wtype.RestrictionEnzyme{bsaI, sapI, enzyme.RestrictionEnzyme})
	}

	// returns sites found in first plasmid in array! should be changed later!
	if len(plasmidProducts) > 1 {
		sites = make([]Restrictionsites, 0)
		for i := 0; i < len(plasmidProducts); i++ {
			sitesperplasmid := Restrictionsitefinder(plasmidProducts[i], []wtype.RestrictionEnzyme{bsaI, sapI, enzyme.RestrictionEnzyme})
			for _, site := range sitesperplasmid {
				sites = append(sites, site)
			}
		}
	}

	s = "hmmm I'm confused, this doesn't seem to make any sense"

	if len(plasmidProducts) == 0 && len(failedAssemblies) == 0 {
		err = fmt.Errorf("Nope! construct design %s not predicted to form any ligated parts", assemblyparameters.ToString())
		s = err.Error()
	}

	// remove invalid plasmids
	var validPlasmids []wtype.DNASequence

	for _, seq := range plasmidProducts {
		validPlasmid, _, _, _ := plasmid.ValidPlasmid(seq)
		if validPlasmid {
			validPlasmids = append(validPlasmids, seq)
		}
	}

	plasmidProducts = validPlasmids

	if len(plasmidProducts) == 1 {
		s = "Yay! this should work"
		successfulassemblies = successfulassemblies + 1
	}

	if len(plasmidProducts) > 1 {

		var errormessage string
		if err != nil {
			errormessage = err.Error()
		}
		merr := fmt.Errorf("Yay! this should work but there seems to be %d possible plasmids which could form: %+v", len(plasmidProducts), errormessage, plasmidProducts)
		s = merr.Error()
	}

	if len(plasmidProducts) == 0 && len(failedAssemblies) > 0 {

		s = fmt.Sprint("Ooh, only partial assembly expected: ", assemblyparameters.Partsinorder[(len(assemblyparameters.Partsinorder)-1)].Nm, " and ", assemblyparameters.Vector.Nm, ": ", "Not compatible, check ends")

		err = fmt.Errorf(s)

		var seqs []wtype.DNASequence

		for i, failed := range failedAssemblies {
			seq, err := failed.ToDNASequence("fragment" + strconv.Itoa(i+1))
			if err != nil {
				return s, successfulassemblies, sites, plasmidProducts, err
			}
			seqs = append(seqs, seq)
		}

		return s, successfulassemblies, sites, seqs, err

	}

	if !strings.Contains(s, "Yay! this should work") {
		err = fmt.Errorf(s)
	}
	for _, newDNASequence := range plasmidProducts {
		newDNASequence.Nm = assemblyparameters.Constructname
	}

	return s, successfulassemblies, sites, plasmidProducts, err
}

func biggest(entries []wtype.DNASequence) wtype.DNASequence {

	var value wtype.DNASequence
	var number int

	for _, str := range entries {
		if len(str.Seq) > number {
			number = len(str.Seq)
			value = str
		}
	}

	return value
}

func split(entries []wtype.DNASequence, entryPositionInSlice int) (split wtype.DNASequence, rest []wtype.DNASequence, err error) {

	if len(entries) == 0 {
		return split, rest, fmt.Errorf("no sequences to split")
	}

	if entryPositionInSlice >= len(entries) {
		return split, rest, fmt.Errorf("cannot take entry %d from slice of length %d", entryPositionInSlice, len(entries))
	}

	for i, entry := range entries {

		if i == entryPositionInSlice {
			split = entry
		} else {
			rest = append(rest, entry)
		}

	}

	return split, rest, nil
}

// MultipleAssemblies will perform simulated assemblies on multiple constructs
// and return a description of whether each was successful and how many are
// expected to work
func MultipleAssemblies(parameters []Assemblyparameters) (s string, successfulassemblies int, errors map[string]string, seqs []wtype.DNASequence) {

	seqs = make([]wtype.DNASequence, 0)
	errors = make(map[string]string) // construct to error

	for _, construct := range parameters {
		output, _, _, seq, err := Assemblysimulator(construct)
		// add first sequence only
		if len(seq) > 0 {
			seqs = append(seqs, seq[0])
		}
		if err == nil {
			successfulassemblies += 1
			continue
		} else {

			errors[construct.Constructname] = err.Error()

			if strings.Contains(err.Error(), "Failure Joining fragments after digestion") {
				sitesperpart := make([]Restrictionsites, 0)
				constructsitesstring := make([]string, 0)
				constructsitesstring = append(constructsitesstring, output)
				sitestring := ""
				enzyme, err := lookup.EnzymeLookup(construct.Enzymename)
				if err != nil {

					originalerror := errors[construct.Constructname]

					errors[construct.Constructname] = originalerror + " and " + err.Error()
				}
				sitesperpart = Restrictionsitefinder(construct.Vector, []wtype.RestrictionEnzyme{enzyme})

				if sitesperpart[0].Numberofsites != 2 {
					// need to loop through sitesperpart

					sitepositions := SitepositionString(sitesperpart[0])
					sitestring = "For " + construct.Vector.Nm + ": " + strconv.Itoa(sitesperpart[0].Numberofsites) + " sites found at positions: " + sitepositions
					constructsitesstring = append(constructsitesstring, sitestring)
				}

				for _, part := range construct.Partsinorder {
					sitesperpart = Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enzyme})
					if sitesperpart[0].Numberofsites != 2 {
						sitepositions := SitepositionString(sitesperpart[0])
						positions := ""
						if sitesperpart[0].Numberofsites != 0 {
							positions = fmt.Sprint("at positions:", sitepositions)
						}
						sitestring = fmt.Sprint("For ", part.Nm, ": ", strconv.Itoa(sitesperpart[0].Numberofsites), " sites were found ", positions)
						constructsitesstring = append(constructsitesstring, sitestring)
					}

				}
				if len(constructsitesstring) != 1 {
					message := strings.Join(constructsitesstring, "; ")
					err = fmt.Errorf(message)
				}
			}

			s = err.Error()

			if _, ok := errors[construct.Constructname]; !ok {
				errors[construct.Constructname] = s
			}

		}
	}

	if successfulassemblies == len(parameters) {
		s = "success, all assemblies seem to work"
	}
	return
}
