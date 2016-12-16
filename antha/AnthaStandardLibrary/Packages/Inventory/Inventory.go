// antha/AnthaStandardLibrary/Packages/Inventory/Inventory.go: Part of the Antha language
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

// Some example part structures for reference and testing purposes
package Inventory

import (
	"fmt"
	"os"
	"path/filepath"

	parser "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"log"
	//"os"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
)

// what about cutting 3prime to recognition site?
type TypeIIs struct {
	wtype.RestrictionEnzyme
	Name                              string
	Isoschizomers                     []string
	Topstrand3primedistancefromend    int
	Bottomstrand5primedistancefromend int
}

type LogicalRestrictionEnzyme struct {
	// other fields required but for now the main things are...
	RecognitionSequence               string
	EndLength                         int
	Name                              string
	Prototype                         string
	Topstrand3primedistancefromend    int
	Bottomstrand5primedistancefromend int
	MethylationSite                   string   //"attr, <4>"
	CommercialSource                  []string //string "attr, <5>"
	References                        []int
	RemoteCutter                      bool
}

var features []wtype.Feature

var SapI = wtype.RestrictionEnzyme{"GCTCTTC", 3, "SapI", "", 1, 4, "", []string{"N"}, []int{91, 1109, 1919, 1920}, "TypeIIs"}
var isoschizomers = []string{"BspQI", "LguI", "PciSI", "VpaK32I"}
var SapIenz = TypeIIs{SapI, "SapI", isoschizomers, 1, 4}
var BsaI = wtype.RestrictionEnzyme{"GGTCTC", 4, "BsaI", "Eco31I", 1, 5, "?(5)", []string{"N"}, []int{814, 1109, 1912, 1995, 1996}, "TypeIIs"}
var BsaIenz = TypeIIs{BsaI, "BsaI", []string{"none"}, 1, 5}

var TypeIIsEnzymeproperties = map[string]TypeIIs{
	"SAPI": SapIenz,
	"BSAI": BsaIenz,
}

// Parts for SapI assembly
var Vector1 = wtype.DNASequence{"dv_ts_CmR_1mod", strings.ToUpper("CACTATAGGGCGAATTGGCGGAAGGCCGTCAAGGCCGCATGCTCTTCTCGCTCGGTTGCCGCCGGGCGTTTTTTATTGGTGAGAATACGGTTTCCCTCTAGAAATAATTTTGTTTACCCGATACAGTGATTATGAAAGGTTTGCGGGGCATAGCTACGACTTGCTTAGCTACGTGCGAGGGGAGAAACTTTTGCGTATTTGTATGTTCACCCGTCTACTACCCATGCCCGGAGATTATGTAGGTTGTGAGATGCGGGAGAAGTTCTCGACCTTCCCGTGGGACGTCAACCTATCCCTTAATAGAGCATTCCGCTCGACGCGAGAATTCGAGCTCGGTACCCGGGGATCCTCTAGAGTCGACCTGCAGGCATGCAAGCTTAGGTACGGTTTCCCTCTAGAAATAATTTTGTTTAACGATGACTACTGGTATTGGCACAAACCTGATTCCAATTTGAGCAAGGCTATGTGCCATCTCGATACTCGTTCTTAACTCAACAGAAGATGCTTTGTGCATACAGCCCCTCGTTTATTATTTATCTCCTCAGCCAGCCGCTGTGCTTTCAGTGGATTTCGGATAACAGAAAGGCCGGGAAATACCCAGCCTCGCTTTGTAACGGAGTAGACGAAAGTGATTGCGCCTACCCGGATATTATCGTGAGGATGCGAAACAGACGAAGAATCCATGGGTATGGACAGTTTTCCCTTTGATATGTAACGGTGAACAGTTGTTCTACTTTTGTTTGTTAGTCTTGATGCTTCACTGATAGATACAAGAGCCATAAGAACCTTCAGGTTTGCCGGCTGAAAGCGCTATTTCTTCCAGAATTGCCATGATTTTTTCCCCACGGGAGGCGTCACTGGCTCCCGTGTTGTCGGCAGCTTTGATTCGATAAGCAGCATCGCCTGTTTCAGGCTGTCTATGTGTGACTGTTGAGCTGTAACAAGTTGTCTCAGGTGTTCAATTTCATGTTCTAGTTGCTTTGTTTTACTGGTTTCACCTGTTCTATTAGGTGTTACATGCTGTTCATCTGTTACATTGTCGATCTGTTCATGGTGAACAGCTTTAAATGCACCAAAAACTCGTAAAAGCTCTGATGTATCTATCTTTTTTACACCGTTTTCATCTGTGCATATGGACAGTTTTCCCTTTGATATCTAACGGTGAACAGTTGTTCTACTTTTGTTTGTTAGTCTTGATGCTTCACTGATAGATACAAGAGCCATAGAAGAGCCTGGGCCTCATGGGCCTTCCGCTCACTGCCCGCTTTCCAG"), true, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}
var Part1 = wtype.DNASequence{"dv_ts_CmR_2mod", strings.ToUpper("GCTCTTCTcataagaaccTTAgatccttccgtatttagccagtatgttctctagtgtggttcgttgtttttgcgtgagccatgagaacgaaccattgagatcatgcttactttgcatgtcactcaaaaattttgcctcaaaactggtgagctgaatttttgcagttaaagcatcgtgtagtgtttttcttagtccgttacgtaggtaggaatctgatgtaatggttgttggtattttgtcaccattcatttttatctggttgttctcaagttcggttacgagatccatttgtctatctagttcaacttggaaaatcaacgtatcagtcgggcggcctcgcttatcaaccaccaatttcatattgctgtaagtgtttaaatctttacttattggtttcaaaacccattggttaagccttttaaactcatggtagttattttcaagcattaacatgaacttaaattcatcaaggctaatctctatatttgccttgtgagttttcttttgtgttagttcttttaataaccactcataaatcctcatagagtatttgttttcaaaagacttaacatgttccagattatattttatgaatttttttaactggaaaagataaggcaatatctcttcactaaaaactaattctaatttttcgcttgagaacttggcatagtttgtccactggaaaatctcaaagcctttaaccaaaggattcctgatttccacagttctcgtcatcagctctctggttgctttagctaatacaccataagcattttccctactgatgttcatcatctgaacgtattggttataagtgaacgataccgtccgttctttccttgtagggttttcaatcgtggggttgagtagtgccacacagcataaaattagcttggtttcatgctccgttaagtcatagcgactaatcgctagttcatttgctttgaaaacaactaattcagacatacatctcaattggtctaggtgattttaatcactataccaattgagatgggctagtcaatgataattactagtccttttcctttgagttgtgggtatctgtaaattctgctagacctttgctggaaaacttgtaaattctgctagaccctctgtaaattccgctagacctttgtgtgttttttttgtttatattcaagtggttataatttatagaataaagaaagaataaaaaaagataaaaagaatagatcccagccctgtgtataactcactactttagtcagttccgcagtattacaaaaggatgtcgcaaacgctgtttgctcctctacaaaacagaccttaaaaccctaaaggcttaagtagcaccctcgcaagctcgggcaaatcgctgaatattccttttgtctccgaccatcaggcacctgagtcgctgtctttttcgtgacattcagttcgctgcgctcacggctctggcagtgaatgggggtAGAAGAGC"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}
var Part2 = wtype.DNASequence{"dv_ts_CmR_3mod", strings.ToUpper("GCTCTTCTggtaaatggcactacaggcgccttttatggattcatgcaaggaaactacccataatacaagaaaagcccgtcacgggcttctcagggcgttttatggcgggtctgctatgtggtgctatctgactttttgctgttcagcagttcctgccctctgattttccagtctgaccacttcggattatcccgtgacaggtcattcagactggctaatgcacccagtaaggcagcggtatcatcaacaggcttacccgtcttactgtcaaccagacccgccaggataagcaatccggcagactggtacagagcatggtcacgggctttacgggcggctctggcttcggctcgcttttctgcctgtatcaggttcatgactttcgctctggcttcctgcgctttctctgctgcctgttcagccttcatcagcagggacagcgttttgatatcctgtactactttagtcagttccgcagtattacaaaaggatgtcgcaaacgctgtttgctcctctacaaaacagaccttaaaaccctaaaggcttaagtagcaccctcgcaagctcggttgcggccgcaatcgggcaaatcgctgaatattccttttgtctccgaccatcaggcacctgagtcgctgtctttttcgtgacattcagttcgctgcgctcacggctctggcagtgaatgggggtaaatggcactacaggcgccttttatggattcatgcaaggaaactacccataatacaagaaaagcccgtcacgggcttctcagggcgttttatggcgggtctgctatgtggtgctatctgactttttgctgttcagcagttcctgccctctgattttccagtctgaccacttcggattatcccgtgacaggtcattcagactggctaatgcacccagtaaggcagcggtatcatcaaTCCGCTAGCGCTGATGTCCGGCGGTGCTTTTGCCGTTACGCACCACCCCGTCAGTAGCTGAACAGGAGGGACAGCTGATAGAAACAGAAGCCACTGGAGCACCTCAAAAACACCATCATACACTAAATCAGTAAGTTGGCAGCATCACCCGACGCACTTTGCGCCGAATAAATACCTGTGACGGAAGATCACTTCGCAGAATAAATAAATCCTGGTGTCCCTGTTGATACCGGGAAGCCCTGGGCCAACTTTTGGCGAAAATGAGACGTTGATCGGCACGTAAGAGGTTCCAACTTTCACCATAATGAAATAAGATCACTACCGGGCGTATTTTTTGAGTTATCGAGATTTTCAGGAGCTAAGGAAGCTAAAATGGAGAAAAAAATCACTGGATATACCACCGTTGATATATCCCAATGGCATCGTAAAGAACATTTTGAGGCATTTCAGTCAGTTGCTCAATGTACCTATAACCAGACCGTTCAGCTGGATATTACGGCCTTTTTAAAGACCGTAAAGAAAAATAAGCACAAGTTTTATCCGGCCTTTATTCACATTCTTGCCCGCCTGATGAATGCTCATCCGGAATTCCGTATGGCAATGAAAGACGGTGAGCTGGTGATATGGGATAGTGTTCACCCTTGTTACACCGTTTTCCATGAGCAAACTGAAACGTTTTCATCGCTCTGGAGTGAATACCACGACGATTTCCGGCAGTTTCTACACATATATTCGCAAGATGTGGCGTGTTACGGTGAAAACCTGGCCTATTTCCCTAAAGGGTTTATTGAGAATATGTTTTTCGTCTCAGCCAATCCCTGGGTGAGTTTCACCAGTTTTGATTTAAACGTGGCCAATATGGACAACTTCTTCGCCCCCGTTTTCACCATGGGCAAATATTATACGCAAGGCGACAAGGTGCTGATGCCGCTGGCGATTCAGGTTCATCATGCCGTCTGTGATGGCTTCCATGTCGGCAGAATGCTTAATGAATTACAACAGTACTGCGATGAGTGGCAGGGCGGGGCGTAATTTTTTTAAGGCAGTTATTGGTGCCCTTAAACGCCTGGTGCTACGCCTGAATAAGTGATAATAAGCGGATGAATGGCAGAAATTCGAAAGCAAATTCGACCCGGTCGTCGGTTCAGGGCAGGGTCGTTAAATAGCCGCTTATGTCTATTGCTGGTTTACCGGTTTATTGACTACCGGAAGCAGTGTGACCGTGTGCTTCTCAAATGCCTGAGGCCAGTTTGCTCAGGCTCTCCCCGTGGAGGTAATAATTGACGATATGATCATTTATTCTGCCTCCCAGAGCCTGATAAAAACGGTTAGCGCTTCGTTAATACAGATGTAGGTGTTCCACAGGGTAGCAgcgggactctggggttcgacagcacgcttggactcctgttgatagatccagtaatgacctcagaactccatctggatttgttcagaacgcAGAAGAGC"), false, false, wtype.Overhang{0, -1, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}

// Parts to add overhangs to
var pAGM1311_R = wtype.DNASequence{"pAGM1311_R", strings.ToUpper("gtgaagGAGGagagaccgcagctggcacgacaggtttcccgactggaaagcgggcagtgagcgcaacgcaattaatgtgagttagctcactcattaggcaccccaggctttacactttatgcttccggctcgtatgttgtgtggaattgtgagcggataacaatttcacacaggaaacagctatgaccatgattacgccaagcttgcatgcctgcaggtcgactctagaggatccccgggtaccgagctcgaattcactggccgtcgttttacaacgtcgtgactgggaaaaccctggcgttacccaacttaatcgccttgcagcacatccccctttcgccagctggcgtaatagcgaagaggcccgcaccgatcgcccttcccaacagttgcgcagcctgaatggcgaatggcgcctgatgcggtattttctccttacgcatctgtgcggtatttcacaccgcatatggtgcactctcagtacaatctgctctgatgccgcatagttaagccagccccgacacccgccaacacccgctgacgcgccctgacgggcttgtctgctcccggcatccgcttacagacaagctgtgacggtctcaCGCTcttcacgaagtggctcttcagtggacgaaagggcctcgtgatacgcctatttttataggttaatgtcatgataataatggtttcttagacgtcaggtggcacttttcggggaaatgtgcgcggaacccctatttgtttatttttctaaatacattcaaatatgtatccgctcatgagacaataaccctgataaatgcttcaataatattgaaaaaggaagagtatggctaaaatgagaatatcaccggaattgaaaaaactgatcgaaaaataccgctgcgtaaaagatacggaaggaatgtctcctgctaaggtatataagctggtgggagaaaatgaaaacctatatttaaaaatgacggacagccggtataaagggaccacctatgatgtggaacgggaaaaggacatgatgctatggctggaaggaaagctgcctgttccaaaggtcctgcactttgaacggcatgatggctggagcaatctgctcatgagtgaggccgatggcgtcctttgctcggaagagtatgaagatgaacaaagccctgaaaagattatcgagctgtatgcggagtgcatcaggctctttcactccatcgacatatcggattgtccctatacgaatagcttagacagccgcttagccgaattggattacttactgaataacgatctggccgatgtggattgcgaaaactgggaagaggacactccatttaaagatccgcgcgagctgtatgattttttaaagacggaaaagcccgaagaggaacttgtcttttcccacggcgacctgggagacagcaacatctttgtgaaagatggcaaagtaagtggctttattgatcttgggagaagcggcagggcggacaagtggtatgacattgccttctgcgtccggtcgatcagggaggatatcggggaagaacagtatgtcgagctattttttgacttactggggatcaagcctgattgggagaaaataaaatattatattttactggatgaattgttttagctgtcagaccaagtttactcatatatactttagattgatttaaaacttcatttttaatttaaaaggatctaggtgaagatcctttttgataatctcatgaccaaaatcccttaacgtgagttttcgttccactgagcgtcagaccccgtagaaaagatcaaaggatcttcttgagatcctttttttctgcgcgtaatctgctgcttgcaaacaaaaaaaccaccgctaccagcggtggtttgtttgccggatcaagagctaccaactctttttccgaaggtaactggcttcagcagagcgcagataccaaatactgtccttctagtgtagccgtagttaggccaccacttcaagaactctgtagcaccgcctacatacctcgctctgctaatcctgttaccagtggctgctgccagtggcgataagtcgtgtcttaccgggttggactcaagacgatagttaccggataaggcgcagcggtcgggctgaacggggggttcgtgcacacagcccagcttggagcgaacgacctacaccgaactgagatacctacagcgtgagctatgagaaagcgccacgcttcccgaagggagaaaggcggacaggtatccggtaagcggcagggtcggaacaggagagcgcacgagggagcttccagggggaaacgcctggtatctttatagtcctgtcgggtttcgccacctctgacttgagcgtcgatttttgtgatgctcgtcaggggggcggagcctatggaaaaacgccagcaacgcggcctttttacggttcctggccttttgctggccttttgctcacatgttctttcctgcgttatcccctgattctgtggataaccgtattaccgcctttgagtgagctgataccgctcgccgcagccgaacgaccgagcgcagcgagtcagtgagcgaggaagcggaagagcgcccaatacgcaaaccgcctctccccgcgcgttggccgattcattaatcactct"), true, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}

var pAGM1311_M = wtype.DNASequence{"pAGM1311_M", strings.ToUpper("gtgaagGGAGagagaccgcagctggcacgacaggtttcccgactggaaagcgggcagtgagcgcaacgcaattaatgtgagttagctcactcattaggcaccccaggctttacactttatgcttccggctcgtatgttgtgtggaattgtgagcggataacaatttcacacaggaaacagctatgaccatgattacgccaagcttgcatgcctgcaggtcgactctagaggatccccgggtaccgagctcgaattcactggccgtcgttttacaacgtcgtgactgggaaaaccctggcgttacccaacttaatcgccttgcagcacatccccctttcgccagctggcgtaatagcgaagaggcccgcaccgatcgcccttcccaacagttgcgcagcctgaatggcgaatggcgcctgatgcggtattttctccttacgcatctgtgcggtatttcacaccgcatatggtgcactctcagtacaatctgctctgatgccgcatagttaagccagccccgacacccgccaacacccgctgacgcgccctgacgggcttgtctgctcccggcatccgcttacagacaagctgtgacggtctcaCGCTcttcacgaagtggctcttcagtggacgaaagggcctcgtgatacgcctatttttataggttaatgtcatgataataatggtttcttagacgtcaggtggcacttttcggggaaatgtgcgcggaacccctatttgtttatttttctaaatacattcaaatatgtatccgctcatgagacaataaccctgataaatgcttcaataatattgaaaaaggaagagtatggctaaaatgagaatatcaccggaattgaaaaaactgatcgaaaaataccgctgcgtaaaagatacggaaggaatgtctcctgctaaggtatataagctggtgggagaaaatgaaaacctatatttaaaaatgacggacagccggtataaagggaccacctatgatgtggaacgggaaaaggacatgatgctatggctggaaggaaagctgcctgttccaaaggtcctgcactttgaacggcatgatggctggagcaatctgctcatgagtgaggccgatggcgtcctttgctcggaagagtatgaagatgaacaaagccctgaaaagattatcgagctgtatgcggagtgcatcaggctctttcactccatcgacatatcggattgtccctatacgaatagcttagacagccgcttagccgaattggattacttactgaataacgatctggccgatgtggattgcgaaaactgggaagaggacactccatttaaagatccgcgcgagctgtatgattttttaaagacggaaaagcccgaagaggaacttgtcttttcccacggcgacctgggagacagcaacatctttgtgaaagatggcaaagtaagtggctttattgatcttgggagaagcggcagggcggacaagtggtatgacattgccttctgcgtccggtcgatcagggaggatatcggggaagaacagtatgtcgagctattttttgacttactggggatcaagcctgattgggagaaaataaaatattatattttactggatgaattgttttagctgtcagaccaagtttactcatatatactttagattgatttaaaacttcatttttaatttaaaaggatctaggtgaagatcctttttgataatctcatgaccaaaatcccttaacgtgagttttcgttccactgagcgtcagaccccgtagaaaagatcaaaggatcttcttgagatcctttttttctgcgcgtaatctgctgcttgcaaacaaaaaaaccaccgctaccagcggtggtttgtttgccggatcaagagctaccaactctttttccgaaggtaactggcttcagcagagcgcagataccaaatactgtccttctagtgtagccgtagttaggccaccacttcaagaactctgtagcaccgcctacatacctcgctctgctaatcctgttaccagtggctgctgccagtggcgataagtcgtgtcttaccgggttggactcaagacgatagttaccggataaggcgcagcggtcgggctgaacggggggttcgtgcacacagcccagcttggagcgaacgacctacaccgaactgagatacctacagcgtgagctatgagaaagcgccacgcttcccgaagggagaaaggcggacaggtatccggtaagcggcagggtcggaacaggagagcgcacgagggagcttccagggggaaacgcctggtatctttatagtcctgtcgggtttcgccacctctgacttgagcgtcgatttttgtgatgctcgtcaggggggcggagcctatggaaaaacgccagcaacgcggcctttttacggttcctggccttttgctggccttttgctcacatgttctttcctgcgttatcccctgattctgtggataaccgtattaccgcctttgagtgagctgataccgctcgccgcagccgaacgaccgagcgcagcgagtcagtgagcgaggaagcggaagagcgcccaatacgcaaaccgcctctccccgcgcgttggccgattcattaatcactct"), true, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}
var pUC57 = wtype.DNASequence{"pUC57", strings.ToUpper("tcgcgcgtttcggtgatgacggtgaaaacctctgacacatgcagctcccggagacggtcacagcttgtctgtaagcggatgccgggagcagacaagcccgtcagggcgcgtcagcgggtgttggcgggtgtcggggctggcttaactatgcggcatcagagcagattgtactgagagtgcaccatatgcggtgtgaaataccgcacagatgcgtaaggagaaaataccgcatcaggcgccattcgccattcaggctgcgcaactgttgggaagggcgatcggtgcgggcctcttcgctattacgccagctggcgaaagggggatgtgctgcaaggcgattaagttgggtaacgccagggttttcccagtcacgacgttgtaaaacgacggccagtgaattcgagctcggtacctcgcgaatgcatctagatatcggatcccgggcccgtcgactgcagaggcctgcatgcaagcttggcgtaatcatggtcatagctgtttcctgtgtgaaattgttatccgctcacaattccacacaacatacgagccggaagcataaagtgtaaagcctggggtgcctaatgagtgagctaactcacattaattgcgttgcgctcactgcccgctttccagtcgggaaacctgtcgtgccagctgcattaatgaatcggccaacgcgcggggagaggcggtttgcgtattgggcgctcttccgcttcctcgctcactgactcgctgcgctcggtcgttcggctgcggcgagcggtatcagctcactcaaaggcggtaatacggttatccacagaatcaggggataacgcaggaaagaacatgtgagcaaaaggccagcaaaaggccaggaaccgtaaaaaggccgcgttgctggcgtttttccataggctccgcccccctgacgagcatcacaaaaatcgacgctcaagtcagaggtggcgaaacccgacaggactataaagataccaggcgtttccccctggaagctccctcgtgcgctctcctgttccgaccctgccgcttaccggatacctgtccgcctttctcccttcgggaagcgtggcgctttctcatagctcacgctgtaggtatctcagttcggtgtaggtcgttcgctccaagctgggctgtgtgcacgaaccccccgttcagcccgaccgctgcgccttatccggtaactatcgtcttgagtccaacccggtaagacacgacttatcgccactggcagcagccactggtaacaggattagcagagcgaggtatgtaggcggtgctacagagttcttgaagtggtggcctaactacggctacactagaagaacagtatttggtatctgcgctctgctgaagccagttaccttcggaaaaagagttggtagctcttgatccggcaaacaaaccaccgctggtagcggtggtttttttgtttgcaagcagcagattacgcgcagaaaaaaaggatctcaagaagatcctttgatcttttctacggggtctgacgctcagtggaacgaaaactcacgttaagggattttggtcatgagattatcaaaaaggatcttcacctagatccttttaaattaaaaatgaagttttaaatcaatctaaagtatatatgagtaaacttggtctgacagttaccaatgcttaatcagtgaggcacctatctcagcgatctgtctatttcgttcatccatagttgcctgactccccgtcgtgtagataactacgatacgggagggcttaccatctggccccagtgctgcaatgataccgcgagacccacgctcaccggctccagatttatcagcaataaaccagccagccggaagggccgagcgcagaagtggtcctgcaactttatccgcctccatccagtctattaattgttgccgggaagctagagtaagtagttcgccagttaatagtttgcgcaacgttgttgccattgctacaggcatcgtggtgtcacgctcgtcgtttggtatggcttcattcagctccggttcccaacgatcaaggcgagttacatgatcccccatgttgtgcaaaaaagcggttagctccttcggtcctccgatcgttgtcagaagtaagttggccgcagtgttatcactcatggttatggcagcactgcataattctcttactgtcatgccatccgtaagatgcttttctgtgactggtgagtactcaaccaagtcattctgagaatagtgtatgcggcgaccgagttgctcttgcccggcgtcaatacgggataataccgcgccacatagcagaactttaaaagtgctcatcattggaaaacgttcttcggggcgaaaactctcaaggatcttaccgctgttgagatccagttcgatgtaacccactcgtgcacccaactgatcttcagcatcttttactttcaccagcgtttctgggtgagcaaaaacaggaaggcaaaatgccgcaaaaaagggaataagggcgacacggaaatgttgaatactcatactcttcctttttcaatattattgaagcatttatcagggttattgtctcatgagcggatacatatttgaatgtatttagaaaaataaacaaataggggttccgcgcacatttccccgaaaagtgccacctgacgtctaagaaaccattattatcatgacattaacctataaaaa"), true, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}

var Promoter = wtype.DNASequence{"w", strings.ToUpper("ttgacggctagctcagtcctaggtacagtgctagc"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}
var Rbs = wtype.DNASequence{"x", strings.ToUpper("C"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}
var CDS = wtype.DNASequence{"y", strings.ToUpper("A"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}
var Term = wtype.DNASequence{"z", strings.ToUpper("T"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features}

var initialpartslist = map[string]wtype.DNASequence{
	"Vector1":       Vector1,
	"Part1":         Part1,
	"Part2":         Part2,
	"Promoter":      wtype.DNASequence{"AndersonProm", strings.ToUpper("ttgacggctagctcagtcctaggtacagtgctagc"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features},
	"Rbs":           wtype.DNASequence{"BBa_B0034", strings.ToUpper("aaagaggagaaa"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features},
	"CDS":           wtype.DNASequence{"GFP_BBa_E0040", strings.ToUpper("atgcgtaaaggagaagaacttttcactggagttgtcccaattcttgttgaattagatggtgatgttaatgggcacaaattttctgtcagtggagagggtgaaggtgatgcaacatacggaaaacttacccttaaatttatttgcactactggaaaactacctgttccatggccaacacttgtcactactttcggttatggtgttcaatgctttgcgagatacccagatcatatgaaacagcatgactttttcaagagtgccatgcccgaaggttatgtacaggaaagaactatatttttcaaagatgacgggaactacaagacacgtgctgaagtcaagtttgaaggtgatacccttgttaatagaatcgagttaaaaggtattgattttaaagaagatggaaacattcttggacacaaattggaatacaactataactcacacaatgtatacatcatggcagacaaacaaaagaatggaatcaaagttaacttcaaaattagacacaacattgaagatggaagcgttcaactagcagaccattatcaacaaaatactccaattggcgatggccctgtccttttaccagacaaccattacctgtccacacaatctgccctttcgaaagatcccaacgaaaagagagaccacatggtccttcttgagtttgtaacagctgctgggattacacatggcatggatgaactatacaaataataa"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features},
	"Term":          wtype.DNASequence{"Doubleterm_BBa_B0015", strings.ToUpper("ccaggcatcaaataaaacgaaaggctcagtcgaaagactgggcctttcgttttatctgttgtttgtcggtgaacgctctctactagagtcacactggctcaccttcgggtgggcctttctgcgtttata"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features},
	"pAGM1311_R":    wtype.DNASequence{"pAGM1311", strings.ToUpper("gtgaagGGAGagagaccgcagctggcacgacaggtttcccgactggaaagcgggcagtgagcgcaacgcaattaatgtgagttagctcactcattaggcaccccaggctttacactttatgcttccggctcgtatgttgtgtggaattgtgagcggataacaatttcacacaggaaacagctatgaccatgattacgccaagcttgcatgcctgcaggtcgactctagaggatccccgggtaccgagctcgaattcactggccgtcgttttacaacgtcgtgactgggaaaaccctggcgttacccaacttaatcgccttgcagcacatccccctttcgccagctggcgtaatagcgaagaggcccgcaccgatcgcccttcccaacagttgcgcagcctgaatggcgaatggcgcctgatgcggtattttctccttacgcatctgtgcggtatttcacaccgcatatggtgcactctcagtacaatctgctctgatgccgcatagttaagccagccccgacacccgccaacacccgctgacgcgccctgacgggcttgtctgctcccggcatccgcttacagacaagctgtgacggtctcaGCTTcttcacgaagtggctcttcagtggacgaaagggcctcgtgatacgcctatttttataggttaatgtcatgataataatggtttcttagacgtcaggtggcacttttcggggaaatgtgcgcggaacccctatttgtttatttttctaaatacattcaaatatgtatccgctcatgagacaataaccctgataaatgcttcaataatattgaaaaaggaagagtatggctaaaatgagaatatcaccggaattgaaaaaactgatcgaaaaataccgctgcgtaaaagatacggaaggaatgtctcctgctaaggtatataagctggtgggagaaaatgaaaacctatatttaaaaatgacggacagccggtataaagggaccacctatgatgtggaacgggaaaaggacatgatgctatggctggaaggaaagctgcctgttccaaaggtcctgcactttgaacggcatgatggctggagcaatctgctcatgagtgaggccgatggcgtcctttgctcggaagagtatgaagatgaacaaagccctgaaaagattatcgagctgtatgcggagtgcatcaggctctttcactccatcgacatatcggattgtccctatacgaatagcttagacagccgcttagccgaattggattacttactgaataacgatctggccgatgtggattgcgaaaactgggaagaggacactccatttaaagatccgcgcgagctgtatgattttttaaagacggaaaagcccgaagaggaacttgtcttttcccacggcgacctgggagacagcaacatctttgtgaaagatggcaaagtaagtggctttattgatcttgggagaagcggcagggcggacaagtggtatgacattgccttctgcgtccggtcgatcagggaggatatcggggaagaacagtatgtcgagctattttttgacttactggggatcaagcctgattgggagaaaataaaatattatattttactggatgaattgttttagctgtcagaccaagtttactcatatatactttagattgatttaaaacttcatttttaatttaaaaggatctaggtgaagatcctttttgataatctcatgaccaaaatcccttaacgtgagttttcgttccactgagcgtcagaccccgtagaaaagatcaaaggatcttcttgagatcctttttttctgcgcgtaatctgctgcttgcaaacaaaaaaaccaccgctaccagcggtggtttgtttgccggatcaagagctaccaactctttttccgaaggtaactggcttcagcagagcgcagataccaaatactgtccttctagtgtagccgtagttaggccaccacttcaagaactctgtagcaccgcctacatacctcgctctgctaatcctgttaccagtggctgctgccagtggcgataagtcgtgtcttaccgggttggactcaagacgatagttaccggataaggcgcagcggtcgggctgaacggggggttcgtgcacacagcccagcttggagcgaacgacctacaccgaactgagatacctacagcgtgagctatgagaaagcgccacgcttcccgaagggagaaaggcggacaggtatccggtaagcggcagggtcggaacaggagagcgcacgagggagcttccagggggaaacgcctggtatctttatagtcctgtcgggtttcgccacctctgacttgagcgtcgatttttgtgatgctcgtcaggggggcggagcctatggaaaaacgccagcaacgcggcctttttacggttcctggccttttgctggccttttgctcacatgttctttcctgcgttatcccctgattctgtggataaccgtattaccgcctttgagtgagctgataccgctcgccgcagccgaacgaccgagcgcagcgagtcagtgagcgaggaagcggaagagcgcccaatacgcaaaccgcctctccccgcgcgttggccgattcattaatcactct"), true, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features},
	"AmpR_promoter": wtype.DNASequence{"AmpR_promoter", strings.ToUpper("CGCGGAACCCCTATTTGTTTATTTTTCTAAATACATTCAAATATGTATCCGCTCATGAGACAATAACCCTGATAAATGCTTCAATAATATTGAAAAAGGAAGAGT"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features},
	"SmR":           wtype.DNASequence{"SmR", strings.ToUpper("atgagggaagcggtgatcgccgaagtatcgactcaactatcagaggtagttggcgtcatcgagcgccatctcgaaccgacgttgctggccgtacatttgtacggctccgcagtggatggcggcctgaagccacacagcgatattgatttgctggttacggtgaccgtaaggcttgatgaaacaacgcggcgagctttgatcaacgaccttttggaaacttcggcttcccctggagagagcgagattctccgcgctgtagaagtcaccattgttgtgcacgacgacatcattccgtggcgttatccagctaagcgcgaactgcaatttggagaatggcagcgcaatgacattcttgcaggtatcttcgagccagccacgatcgacattgatctggctatcttgctgacaaaagcaagagaacatagcgttgccttggtaggtccagcggcggaggaactctttgatccggttcctgaacaggatctatttgaggcgctaaatgaaaccttaacgctatggaactcgccgcccgactgggctggcgatgagcgaaatgtagtgcttacgttgtcccgcatttggtacagcgcagtaaccggcaaaatcgcgccgaaggatgtcgctgccgactgggcaatggagcgcctgccggcccagtatcagcccgtcatacttgaagctagacaggcttatcttggacaagaagaagatcgcttggcctcgcgcgcagatcagttggaagaatttgtccattacgtaaaaggcgagatcaccaaggtagtcggcaaataa"), false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", features},
	"pAGM1311_M":    pAGM1311_M,
}

func Partslist() (partslist map[string]wtype.DNASequence) {
	partslist, err := LookforParts()
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

func LookforParts() (partslist map[string]wtype.DNASequence, err error) {
	partslist = make(map[string]wtype.DNASequence)

	// first add PartsList above
	for key, value := range initialpartslist {

		// initialise parts into map
		partslist[key] = value
		/*if _, alreadyinmap := initialpartslist[key]; !alreadyinmap {
			partslist[key] = value
			fmt.Println("adding ", key, " to inventory")
		}*/
		// add by name too
		if _, alreadyinmap := initialpartslist[value.Nm]; !alreadyinmap {
			partslist[value.Nm] = value
			fmt.Println("adding ", value.Nm, " to inventory")
		} else {
			err = fmt.Errorf("cannot add " + key + " to inventory partslist as already found, change the name")
			//return
		}
	}

	files := make([]string, 0)
	// look for all parts in anthapath
	d, err := os.Open(anthapath.Path())
	if err != nil {
		return nil, err
	}
	defer d.Close()

	allfiles, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	//Determine if file extension is ".fasta"
	for _, file := range allfiles {
		if /*filepath.Ext(file.Name()) == ".gb" ||*/ filepath.Ext(file.Name()) == ".fasta" {
			files = append(files, file.Name())
		}
	}

	for _, filename := range files {
		file := filepath.Join(anthapath.Path(), filename)

		if filepath.Ext(file) == ".fasta" {
			sequences, _ := parser.FastatoDNASequences(file)

			for _, seq := range sequences {
				if _, alreadyinmap := partslist[seq.Nm]; !alreadyinmap {
					partslist[seq.Nm] = seq
				} else {
					err = fmt.Errorf("cannot add " + seq.Nm + " to inventory partslist as already found, change the name")
					//return
				}
			}
		} else if filepath.Ext(file) == ".gb" {
			seq, _ := parser.GenbanktoAnnotatedSeq(file)
			if _, alreadyinmap := partslist[seq.Nm]; !alreadyinmap {
				partslist[seq.Nm] = seq
			} else {
				err = fmt.Errorf("cannot add " + seq.Nm + " to inventory partslist as already found, change the name")
				//	return
			}

		}
	}

	return
}
