// antha/AnthaStandardLibrary/Packages/igem/igem.go: Part of the Antha language
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

package igem

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	//"/data"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	//"time"
	"path/filepath"
)

// http://parts.igem.org/Registry_API

/*
Input
FASTA Formatted Sequences
We will provide a daily update of part sequences, types, subparts, status, and short description for each part and for all parts. Go to http://parts.igem.org/fasta/parts/BBa_C0040 (substitute our desired part name for BBa_C0040) and you will receive a FASTA formatted file with the part's sequence. The header line has this format:
'>'[Part name] [First character of status] [Part Id Number] [Part type] [Short description]
Note: the short description has unusual characters converted to their two-digit hex value.
You can also get all of the parts in a single download (about 30 megabytes) as http://parts.igem.org/fasta/parts/All_Parts.
We are not yet updating these files on a daily basis. -- Randy May 30, 2009

XML Part Information
The database information about each part is now available as XML. You can get information about a part by entering a URL like this one. If your browser parses and displays XML in a formatted way, it will make some sense. The URL ends ...xml/part; you can follow with a list of part names separated by periods.
http://parts.igem.org/xml/part.BBa_B0034
The information for a part includes:
Part name, type, nickname, short description, status, rating, date entered, authors, quality
Lists of subparts (as specified by the designer, at the basic part level, and with scars)
Sequence
Features
Parameters
Categories
DNA Samples (not enabled not)
References (not enabled now)
Groups (not enabled now)
If you enter 'recursive' as the first part name, the returned XML will include details about all the subparts of this part.

*/

func MakeFastaURL(partname string) (Urlstring string) {
	// see comment above for structure
	//<domain> = substance | compound | assay | <other inputs>

	level1 := "http://parts.igem.org"
	// http://parts.igem.org/fasta/parts/all
	array := make([]string, 0)
	array = append(array, level1, "fasta", "parts", partname)

	Urlstring = strings.Join(array, "/")

	return Urlstring
}

func MakeXMLURL(partnames []string) (Urlstring string) {
	// see comment above for structure
	//<domain> = substance | compound | assay | <other inputs>

	parts := make([]string, 0)
	parts = append(parts, "part")
	for _, part := range partnames {
		parts = append(parts, part)
	}
	partconcat := strings.Join(parts, ".")

	level1 := "http://parts.igem.org"

	array := make([]string, 0)
	array = append(array, level1, "xml", partconcat)

	Urlstring = strings.Join(array, "/")

	return Urlstring
}

func SlurpOutput(Urlstring string) (output []byte) {
	fmt.Println("Slurping...", Urlstring)

	//var res *io.Reader

	res, err := http.Get(Urlstring)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("step 1")

	/*reader := bufio.NewReader(res.Body)

	//_, err = io.ReadFull(reader, output)

	output, err = reader.Peek(int(res.ContentLength)) returns negative count!

	fmt.Println("hello")
	fmt.Println(string(output))
	if err != nil {
		log.Fatal(err)
	}

	/*for{

	}
	*/

	//you should stream the resp.Body to disk by creating a file using os.Create() and io.Copy(file, resp.Body)
	//if _, err := os.Stat("Registry.fasta"); err != nil {
	//	if os.IsNotExist(err) {
	//		fmt.Printf("downloading registry.fasta")
	//f, _ := os.Create("temp.tmp")

	//err := ioutil.WriteFile("temp.tmp", res.Body, 0777)
	fmt.Println("step 2")
	//_, _ = io.Copy(f, res.Body)
	fmt.Println("step 3")
	//		return
	//	}
	//}

	fmt.Println("step 4")
	//output, err = ioutil.ReadFile("temp.tmp") // this is a very fast step!

	/*body := make([]byte, int64(res.ContentLength))
	b := bytes.NewBuffer(body)
	fmt.Println("step 2")
	_, err = b.ReadFrom(res.Body)
	fmt.Println("step 3")
	output = b.Bytes()*/
	output, err = ioutil.ReadAll(res.Body) // this is a slow step!
	fmt.Println("step 5")
	//res.Body.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(Urlstring, "=", string(output))
	return output
}

func UpdateRegistryfile() (err error) {
	fmt.Println("Getting latest registry database file from ...", "http://parts.igem.org/fasta/parts/All_Parts")
	res, err := http.Get("http://parts.igem.org/fasta/parts/All_Parts")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("step 1: creating registry file")

	/*for{

	}
	*/

	//you should stream the resp.Body to disk by creating a file using os.Create() and io.Copy(file, resp.Body)
	//if _, err := os.Stat("Registry.fasta"); err != nil {
	//	if os.IsNotExist(err) {
	//		fmt.Printf("downloading registry.fasta")
	f, _ := os.Create(filepath.Join(anthapath.Dirpath(), "iGem_registry.txt"))
	fmt.Println("step 2: copying registry... This could take a few minutes, don't go anywhere")
	_, err = io.Copy(f, res.Body) // takes just as long as ioutil.Readall()
	fmt.Println("step 3")
	//		return
	//	}
	//}

	fmt.Println("step 4")

	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(Urlstring, "=", string(output))
	return err
}

//[Part name] [First character of status] [Part Id Number] [Part type] [Short description]
/*
func ParseFastaOutput(fastaoutput []byte) (fastarecords []FastaPart) {

	fastaFh := bytes.NewReader(fastaoutput)

	fastarecords = FastaParse(fastaFh)

	return fastarecords
}
*/
type FastaPart struct {
	Part_id         string
	Desc            string
	Part_name       string
	Part_short_desc string
	Part_type       string
	Sample_status   string
	Seq_data        string
}

func Build_fasta(header string, seq bytes.Buffer) (Record FastaPart) {

	var record FastaPart

	fields := strings.SplitN(header, " ", 2)
	record.Desc = "`" + fields[1] + "`"

	fields = strings.SplitN(header, " ", 5)

	if len(fields) > 1 {
		record.Part_name = fields[0]
		record.Sample_status = fields[1]
		record.Part_id = fields[2]
		record.Part_type = fields[3]
		record.Part_short_desc = fields[4]

	} else {
		record.Part_name = fields[0]
		record.Part_short_desc = ""
	}

	record.Seq_data = seq.String()

	Record = record

	return Record
}

func FastaParse(fastaFh io.Reader) chan FastaPart {

	outputChannel := make(chan FastaPart)

	scanner := bufio.NewScanner(fastaFh)
	// scanner.Split(bufio.ScanLines)
	header := ""
	var seq bytes.Buffer

	go func() {
		// Loop over the letters in inputString
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) == 0 {
				continue
			}

			// line := scanner.Text()

			if line[0] == '>' {
				// If we stored a previous identifier, get the DNA string and map to the
				// identifier and clear the string
				if header != "" {
					// outputChannel <- build_fasta(header, seq.String())
					outputChannel <- Build_fasta(header, seq)
					// fmt.Println(record.id, len(record.seq))
					header = ""
					seq.Reset()
				}

				// Standard FASTA identifiers look like: ">id desc"
				header = line[1:]
			} else {
				// Append here since multi-line DNA strings are possible
				seq.WriteString(line)
			}

		}

		outputChannel <- Build_fasta(header, seq)

		// Close the output channel, so anything that loops over it
		// will know that it is finished.
		close(outputChannel)
	}()

	return outputChannel
}

func CountPartsinRegistryContaining(keystrings []string) (numberofparts int) {

	if anthapath.Anthafileexists("iGem_registry.txt") == false {
		err := UpdateRegistryfile()
		if err != nil {
			fmt.Println("error:", err)
		}
		//allparts := SlurpOutput("http://parts.igem.org/fasta/parts/All_Parts")
	}
	allparts, err := ioutil.ReadFile("iGem_registry.txt")
	if err != nil {
		fmt.Println("error:", err)
	}
	/*allparts, err := ioutil.ReadFile("iGem_registry.txt")
	if err != nil {
		fmt.Println("error:", err)
	}*/
	fastaFh := bytes.NewReader(allparts)

	//fasta := parser.FastaParse(allparts)
	records := make([][]string, 0)
	seq := make([]string, 0)
	//seq = []string{"#Name", "Sequence", "Plasmid?", "Seq Type", "Class"}
	records = append(records, seq)
	for record := range FastaParse(fastaFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if enzymes.Containsallthings(record.Desc, keystrings) {
			numberofparts = numberofparts + 1
		}
		/*	if strings.Contains(record.Desc, "Amino acid") || strings.Contains(record.Id, "aa") {
				seqtype = "AA"
			}

			if strings.Contains(record.Desc, "Class:") {
				uptoclass := strings.Index(record.Desc, "Class:")
				prefix := uptoclass + len("class:")
				class = record.Desc[prefix:]
			}
			seq = []string{record.Id, record.Seq, plasmidstatus, seqtype, class}*/
		/*records = append(records, seq)*/
	}
	return numberofparts
}

func Exists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func FilterRegistry(keystrings []string) (listofpartIDs []string) {
	if anthapath.Anthafileexists("iGem_registry.txt") == false {
		err := UpdateRegistryfile()
		if err != nil {
			fmt.Println("error:", err)
		}
		//allparts := SlurpOutput("http://parts.igem.org/fasta/parts/All_Parts")
	}
	allparts, err := ioutil.ReadFile("iGem_registry.txt")
	if err != nil {
		fmt.Println("error:", err)
	}

	/*f, _ := filepath.Abs("$GOPATH/src/github.com/antha-lang/antha/antha/AnthaStandardLibrary/Localdata")
	fmt.Println(f)

	err := os.Chdir("$GOPATH/src/github.com/antha-lang/antha/antha/AnthaStandardLibrary/Localdata")
	if err != nil {
		fmt.Println("error:", err)
	}
	*/
	//allparts, err := ioutil.ReadFile("iGem_registry.txt")
	/*if err != nil {
		fmt.Println("error:", err)
	}

	//err = os.Chdir("cd -")

	if err != nil {
		fmt.Println("error:", err)
	}
	*/
	fmt.Println("slurped registry")
	fastaFh := bytes.NewReader(allparts)

	//fasta := parser.FastaParse(allparts)
	records := make([][]string, 0)
	seq := make([]string, 0)
	//seq = []string{"#Name", "Sequence", "Plasmid?", "Seq Type", "Class"}
	records = append(records, seq)

	listofpartIDs = make([]string, 0)

	for record := range FastaParse(fastaFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if enzymes.Containsallthings(record.Desc, keystrings) && record.Seq_data != "" {
			fmt.Println(record.Part_name)
			listofpartIDs = append(listofpartIDs, record.Part_name)
		}
		/*	if strings.Contains(record.Desc, "Amino acid") || strings.Contains(record.Id, "aa") {
				seqtype = "AA"
			}

			if strings.Contains(record.Desc, "Class:") {
				uptoclass := strings.Index(record.Desc, "Class:")
				prefix := uptoclass + len("class:")
				class = record.Desc[prefix:]
			}
			seq = []string{record.Id, record.Seq, plasmidstatus, seqtype, class}*/
		/*records = append(records, seq)*/
	}
	return listofpartIDs
}

/*
func PartSummary(parts []string) [][]string {

	parsedxml := PartProperties(parts)


	for _,part := range parsedxml.Partlist.Parts {
		name := part.Part_name
		seq := "Seq:"+ part.Sequencelist[0].Seq_data
		part_type := "type:" + part.Part_type
		part.Part_results
		part.Part_rating
		part.Part_short_desc
		for _, category := range part.Categories {

		}

	}
return[][]string
}
*/
func ParseOutput(xmldata []byte) (parsedxml Rsbpml) {

	err := xml.Unmarshal(xmldata, &parsedxml)
	if err != nil {
		fmt.Println("error:", err)
	}

	//fmt.Println(parsedxml)

	return parsedxml
}

func Partpropertiesmini(parts []string) (parsedxml Rsbpml) {
	url := MakeXMLURL(parts)

	urloutput := SlurpOutput(url)

	parsedxml = ParseOutput(urloutput)

	return parsedxml
}

/*

func PartPropertiesChan(parts []string) chan Rsbpml {

	outputChannel := make(chan Fasta)

	fmt.Println("len(parts =", len(parts))
	if len(parts) > 10 {

		partslice := parts[0:10]
		fmt.Println("len(partslice) =", len(partslice))

		go func() {
		for parts []string {


			}

		outputChannel <- Partpropertiesmini(partslice)

		url := MakeXMLURL(partslice)

		urloutput := SlurpOutput(url)

		parsedxml = ParseOutput(urloutput)

		//var parts []Part
		newparsedxml := make([]Part, 0)
		fmt.Println("len(parsedxml.Partlist[0].Parts) =", len(parsedxml.Partlist[0].Parts))
		for _, part := range parsedxml.Partlist[0].Parts {
			newparsedxml = append(newparsedxml, part)
		}

		var parsedxml Rsbpml
		//n := 0
		partsleft := (len(parts) - len(partslice))
		/*if partsleft > 10 {
			n = 10
		} else {
			n = len(parts) - 10
		}

		for i := 10; i < len(parts); i = i + 10 {
			partslice = parts[i : i+10]
			url := MakeXMLURL(partslice)

			urloutput := SlurpOutput(url)

			parsedxml = ParseOutput(urloutput)
			fmt.Println("len(parsedxml.Partlist[0].Parts) =", len(parsedxml.Partlist[0].Parts))
			for _, part := range parsedxml.Partlist[0].Parts {
				newparsedxml = append(newparsedxml, part)
				fmt.Println("len(newparsedxml)", len(newparsedxml))
			}
			var parsedxml Rsbpml
			partsleft = partsleft - len(partslice)
			fmt.Println("parts left = ", partsleft)
			if partsleft < 10 {
				//for i := 0; i < len(parts); i = i + 10 {
				partslice = parts[len(parts)-partsleft : len(parts)]
				url = MakeXMLURL(partslice)

				urloutput = SlurpOutput(url)

				parsedxml = ParseOutput(urloutput)

				for _, part := range parsedxml.Partlist[0].Parts {
					newparsedxml = append(newparsedxml, part)
					fmt.Println("len(newparsedxml)", len(newparsedxml))
					parsedxml.Partlist[0].Parts = newparsedxml
					fmt.Println("newparsedxml", newparsedxml)

				}
				{
					break
				}
			}
		}
	} else {

		url := MakeXMLURL(parts)

		urloutput := SlurpOutput(url)

		parsedxml = ParseOutput(urloutput)
	}
	return parsedxml
}
*/

func PartProperties(parts []string) (parsedxml Rsbpml) {
	fmt.Println("len(parts =", len(parts))
	if len(parts) > 14 {

		partslice := parts[0:14]
		fmt.Println("len(partslice) =", len(partslice))

		parsedxml = Partpropertiesmini(partslice)

		//var parts []Part
		newparsedxml := make([]Part, 0)
		fmt.Println("len(parsedxml.Partlist[0].Parts) =", len(parsedxml.Partlist[0].Parts))
		for _, part := range parsedxml.Partlist[0].Parts {
			newparsedxml = append(newparsedxml, part)
		}

		var parsedxml Rsbpml
		//n := 0
		partsleft := (len(parts) - len(partslice))
		/*if partsleft > 10 {
			n = 10
		} else {
			n = len(parts) - 10
		}
		*/
		for i := 10; i < len(parts); i = i + 14 {
			partslice = parts[i : i+14]
			parsedxml = Partpropertiesmini(partslice)
			fmt.Println("len(parsedxml.Partlist[0].Parts) =", len(parsedxml.Partlist[0].Parts))
			for _, part := range parsedxml.Partlist[0].Parts {
				newparsedxml = append(newparsedxml, part)
				fmt.Println("len(newparsedxml)", len(newparsedxml))
			}
			var parsedxml Rsbpml
			partsleft = partsleft - len(partslice)
			fmt.Println("parts left = ", partsleft)
			if partsleft < 14 {
				//for i := 0; i < len(parts); i = i + 10 {
				partslice = parts[len(parts)-partsleft : len(parts)]
				parsedxml = Partpropertiesmini(partslice)

				for _, part := range parsedxml.Partlist[0].Parts {
					newparsedxml = append(newparsedxml, part)
					fmt.Println("len(newparsedxml)", len(newparsedxml))
					parsedxml.Partlist[0].Parts = newparsedxml
					fmt.Println("newparsedxml", newparsedxml)

				}
				{
					break
				}
			}
		}
	} else {

		parsedxml = Partpropertiesmini(parts)
	}
	return parsedxml
}

// Add Get funcs to get data from Rsbpml? Would be much faster

func GetSequence(partname string) (sequence string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml := ParseOutput(urloutput)

	sequence = parsedxml.Partlist[0].Parts[0].Sequencelist[0].Seq_data // [0].Seq_data

	sequence = strings.Replace(sequence, "\n", "", -1)

	return sequence
}

func GetType(partname string) (parttype string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml := ParseOutput(urloutput)

	parttype = parsedxml.Partlist[0].Parts[0].Part_type // [0].Seq_data

	return parttype
}

func GetCategories(partname string) (categories Categories) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml := ParseOutput(urloutput)

	categories = parsedxml.Partlist[0].Parts[0].Categories // [0].Seq_data

	return categories
}

func GetResults(partname string) (results string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml := ParseOutput(urloutput)

	results = parsedxml.Partlist[0].Parts[0].Part_results // [0].Seq_data

	return results
}

func GetResultsfromSubset(partname string, parsedxml Rsbpml) (results string) {

	/*parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml = ParseOutput(urloutput)*/
	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			results = part.Part_results
		}
	}

	results = strings.ToUpper(results)

	//desc = parsedxml.Partlist[0].Parts[0].Part_short_desc // [0].Seq_data

	return results
}

func GetRating(partname string) (rating string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml := ParseOutput(urloutput)

	rating = parsedxml.Partlist[0].Parts[0].Part_rating // [0].Seq_data

	return rating
}

func GetDescription(partname string) (desc string) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml := ParseOutput(urloutput)

	desc = parsedxml.Partlist[0].Parts[0].Part_short_desc // [0].Seq_data

	return desc
}

func GetDescriptionfromSubset(partname string, parsedxml Rsbpml) (desc string) {

	/*parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml = ParseOutput(urloutput)*/
	for _, part := range parsedxml.Partlist[0].Parts {
		if part.Part_name == partname {
			desc = part.Part_short_desc
		}
	}

	desc = strings.ToUpper(desc)

	//desc = parsedxml.Partlist[0].Parts[0].Part_short_desc // [0].Seq_data

	return desc
}

func GetPart(partname string) (partproperties Part) {

	parts := make([]string, 0)
	parts = append(parts, partname)
	url := MakeXMLURL(parts)
	urloutput := SlurpOutput(url)
	parsedxml := ParseOutput(urloutput)

	partproperties = parsedxml.Partlist[0].Parts[0] // [0].Seq_data
	return partproperties
}

type Rsbpml struct {
	Partlist []Part_list `xml:"part_list"`
}

type Part_list struct {
	Parts []Part `xml:"part"`
}

type Part struct {
	Part_id            string         `xml:"part_id"`
	Part_name          string         `xml:"part_name"`
	Part_short_name    string         `xml:"part_short_name"`
	Part_short_desc    string         `xml:"part_short_desc"`
	Part_type          string         `xml:"part_type"`
	Release_status     string         `xml:"release_status"`
	Sample_status      string         `xml:"sample_status"`
	Part_results       string         `xml:"part_results"`
	Part_nickname      string         `xml:"part_nickname"`
	Part_rating        string         `xml:"part_rating"`
	Part_url           string         `xml:"part_url"`
	Part_entered       string         `xml:"part_entered"`
	Part_author        string         `xml:"part_author"`
	Deep_subparts      Subparts       `xml:"deep_subparts"`
	Specified_subparts Subparts       `xml:"specified_subparts"`
	Specified_subscars Subscars       `xml:"specified_subscars"`
	Sequencelist       []Sequence     `xml:"sequences"`
	Features           IgemFeatures   `xml:"features"`
	Parameters         IgemParameters `xml:"parameters"`
	Categories         Categories     `xml:"categories"`
	Twins              Twins          `xml:"twins"`
}

type Subparts struct {
	Subparts []Subpart `xml:"subpart"`
}

type Subpart struct {
	Part_id         string `xml:"part_id"`
	Part_name       string `xml:"part_name"`
	Part_short_desc string `xml:"part_short_desc"`
	Part_type       string `xml:"part_type"`
	Part_nickname   string `xml:"part_nickname"`
}

type Subscars struct {
	Subparts []Subpart `xml:"subpart"`
	Scars    []Scar    `xml:"scar"`
}

type Scar struct {
	Scar_id       string `xml:"scar_id"`
	Scar_standard string `xml:"scar_standard"`
	Scar_type     string `xml:"scar_type"`
	Scar_name     string `xml:"scar_name"`
	Scar_nickname string `xml:"scar_nickname"`
	Scar_comments string `xml:"scar_comments"`
	Scar_sequence string `xml:"scar_sequence"`
}

/*type Sequences struct {
	Sequencelist []Sequence `xml:"sequences"`
}
*/
type Sequence struct {
	Seq_data string `xml:"seq_data"`
}

type IgemFeatures struct {
	Features []IgemFeature `xml:"feature"`
}

type IgemFeature struct {
	Id        string `xml:"id"`
	Title     string `xml:"title"`
	Type      string `xml:"type"`
	Direction string `xml:"direction"`
	Startpos  string `xml:"startpos"`
	Endpos    string `xml:"endpos"`
}

type IgemParameters struct {
	Parameters []IgemParameter `xml:"parameter"`
}

type IgemParameter struct {
	Name      string `xml:"name"`
	Value     string `xml:"value"`
	Units     string `xml:"units"`
	Url       string `xml:"url"`
	Id        string `xml:"id"`
	M_date    string `xml:"m_date"`
	User_id   string `xml:"user_id"`
	User_name string `xml:"user_name"`
}

type Categories struct {
	Categories []string `xml:"category"`
}

type Twins struct {
	Twins []string `xml:"twin"`
}
