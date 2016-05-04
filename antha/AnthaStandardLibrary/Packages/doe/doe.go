// doe.go
package doe

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
)

type DOEPair struct {
	Factor string
	Levels []interface{}
}

func (pair DOEPair) LevelCount() (numberoflevels int) {
	numberoflevels = len(pair.Levels)
	return
}
func Pair(factordescription string, levels []interface{}) (doepair DOEPair) {
	doepair.Factor = factordescription
	doepair.Levels = levels
	return
}

type Run struct {
	RunNumber            int
	StdNumber            int
	Factordescriptors    []string
	Setpoints            []interface{}
	Responsedescriptors  []string
	ResponseValues       []interface{}
	AdditionalHeaders    []string // could represent a class e.g. Environment variable, processed, raw, location
	AdditionalSubheaders []string // e.g. well ID, Ambient Temp, order,
	AdditionalValues     []interface{}
}

func (run Run) AddResponseValue(responsedescriptor string, responsevalue interface{}) {

	for i, descriptor := range run.Responsedescriptors {
		if strings.ToUpper(descriptor) == strings.ToUpper(responsedescriptor) {
			run.ResponseValues[i] = responsevalue
		}
	}

}

// inexplicably not working
func (run Run) AddNewResponseFieldandValue(responsedescriptor string, responsevalue interface{}) {

	responsedescriptors := make([]string, len(run.Responsedescriptors))
	responsevalues := make([]interface{}, len(run.ResponseValues))

	responsedescriptors = run.Responsedescriptors
	responsevalues = run.ResponseValues

	responsedescriptors = append(responsedescriptors, responsedescriptor)
	responsevalues = append(responsevalues, responsevalue)

	run.Responsedescriptors = responsedescriptors
	run.ResponseValues = responsevalues
	fmt.Println(run)
}

func AddNewFactorFieldandValue(run Run, factordescriptor string, factorvalue interface{}) (newrun Run) {

	//var newrun Run

	factordescriptors := make([]string, len(run.Factordescriptors))
	factorvalues := make([]interface{}, len(run.Setpoints))

	factordescriptors = run.Factordescriptors
	factorvalues = run.Setpoints

	factordescriptors = append(factordescriptors, factordescriptor)
	factorvalues = append(factorvalues, factorvalue)

	//newrun = run

	newrun.Factordescriptors = factordescriptors
	newrun.Setpoints = factorvalues

	//run = newrun
	fmt.Println(run)

	return
}

func (run Run) AddAdditionalValue(additionalsubheader string, additionalvalue interface{}) {

	for i, descriptor := range run.AdditionalSubheaders {
		if strings.ToUpper(descriptor) == strings.ToUpper(additionalsubheader) {
			run.AdditionalValues[i] = additionalvalue
		}
	}

}

func (run Run) AddAdditionalHeaders(additionalheader string, additionalsubheader string) {

	headers := make([]string, len(run.AdditionalHeaders))

	for _, header := range run.AdditionalHeaders {
		headers = append(headers, header)
	}

	headers = append(headers, additionalheader)

	subheaders := make([]string, len(run.AdditionalSubheaders))

	for _, subheader := range run.AdditionalSubheaders {
		subheaders = append(subheaders, subheader)
	}

	subheaders = append(subheaders, additionalsubheader)

}

func (run Run) AddAdditionalHeaderandValue(additionalheader string, additionalsubheader string, additionalvalue interface{}) {
	run.AddAdditionalHeaders(additionalheader, additionalsubheader)
	run.AddAdditionalValue(additionalsubheader, additionalvalue)

}

func (run Run) CheckAdditionalInfo(subheader string, value interface{}) bool {

	for i, header := range run.AdditionalSubheaders {
		if strings.ToUpper(header) == strings.ToUpper(subheader) && run.AdditionalValues[i] == value {
			return true
		}
	}
	return false
}

func (run Run) GetAdditionalInfo(subheader string) (value interface{}, err error) {

	for i, header := range run.AdditionalSubheaders {
		if strings.ToUpper(header) == strings.ToUpper(subheader) {
			value = run.AdditionalValues[i]
			return value, err

		}
	}
	return value, fmt.Errorf("header not found")
}

func AddFixedFactors(runs []Run, fixedfactors []DOEPair) (runswithfixedfactors []Run) {

	if len(runs) > 0 {
		for _, run := range runs {
			descriptors := make([]string, len(run.Factordescriptors))
			setpoints := make([]interface{}, len(run.Setpoints))

			for i, descriptor := range run.Factordescriptors {
				descriptors[i] = descriptor
			}

			for i, setpoint := range run.Setpoints {
				setpoints[i] = setpoint
			}

			for _, fixed := range fixedfactors {

				descriptors = append(descriptors, fixed.Factor)
				setpoints = append(setpoints, fixed.Levels[0])

			}
			run.Factordescriptors = descriptors
			run.Setpoints = setpoints

		}

	} else {
		runs = RunsFromFixedFactors(fixedfactors)
	}

	runswithfixedfactors = runs

	return
}

func RunsFromFixedFactors(fixedfactors []DOEPair) (runswithfixedfactors []Run) {

	var run Run
	var descriptors = make([]string, 0)
	var setpoints = make([]interface{}, 0)

	for _, factor := range fixedfactors {

		descriptors = append(descriptors, factor.Factor)
		setpoints = append(setpoints, factor.Levels[0])

	}

	run.Factordescriptors = descriptors
	run.Setpoints = setpoints

	runswithfixedfactors = make([]Run, 1)

	runswithfixedfactors[0] = run

	return

}

func AllComboCount(pairs []DOEPair) (numberofuniquecombos int) {
	fmt.Println("In AllComboCount", "len(pairs)", len(pairs))
	var movingcount int
	movingcount = (pairs[0]).LevelCount()
	fmt.Println("levelcount", movingcount)
	fmt.Println("len(levels)", len(pairs[0].Levels))
	for i := 1; i < len(pairs); i++ {
		fmt.Println("levelcount", movingcount)
		movingcount = movingcount * (pairs[i]).LevelCount()
	}
	numberofuniquecombos = movingcount
	return
}

func FixedAndNonFixed(factors []DOEPair) (fixedfactors []DOEPair, nonfixed []DOEPair) {

	fixedfactors = make([]DOEPair, 0)
	nonfixed = make([]DOEPair, 0)

	for _, factor := range factors {
		if len(factor.Levels) == 1 {
			fixedfactors = append(fixedfactors, factor)
		} else if len(factor.Levels) > 1 {
			nonfixed = append(nonfixed, factor)
		}
	}
	return
}

func IsFixedFactor(factor DOEPair) (yesorno bool) {
	if len(factor.Levels) == 1 {
		yesorno = true
	}
	return
}

func AllCombinations(factors []DOEPair) (runs []Run) {

	fixed, nonfixed := FixedAndNonFixed(factors)

	fmt.Println("fixed:", fixed, "nonfixed: ", nonfixed)
	numberofruns := AllComboCount(factors)
	fmt.Println("Allcombocount: ", numberofruns, "Factors: ", factors)
	//numberoffactors := len(factors)

	runs = make([]Run, numberofruns)
	//setpoints := make([]interface{}, 0)
	//descriptors := make([]string, 0)

	//for i := 0; i < numberofruns; i++ {
	var swapevery int
	var numberofswaps int
	for i, factor := range factors {

		counter := 0
		runswitheachlevelforthisfactor := numberofruns / factor.LevelCount()

		if i == 0 {
			swapevery = runswitheachlevelforthisfactor
			numberofswaps = runswitheachlevelforthisfactor / swapevery
		} else {
			swapevery = swapevery / factor.LevelCount()
			numberofswaps = runswitheachlevelforthisfactor / swapevery
		}

		for j := 0; j < numberofswaps; j++ {
			for _, level := range factor.Levels {
				for k := 0; k < swapevery; k++ {

					runs[counter] = AddNewFactorFieldandValue(runs[counter], factor.Factor, level)
					runs[counter].RunNumber = counter + 1
					runs[counter].StdNumber = counter + 1
					counter++
					fmt.Println("counter: ", counter, "i: ", i, "Swapevery: ", swapevery, "numberofswaps: ", numberofswaps)
				}
			}
		}

	}

	//runs = AddFixedFactors(runs, fixed)
	//}
	return
}

/*
func AllCombinations(factors []DOEPair) (runs []Run) {

	fixed, nonfixed := FixedAndNonFixed(factors)

	//fmt.Println(factors)
	descriptors := make([]string, 0) //AllComboCount(factors))
	//fixedfactors := make([]DOEPair, 0)
	// range through factors and get annotation info first
	for _, factor := range factors {
		descriptors = append(descriptors, factor.Factor)

	}

	setpoints := make([]interface{}, 0)
	//runs = make([]Run, AllComboCount(factors))
	runs = make([]Run, 0)

	if AllComboCount(factors) == 0 {
		return
		//	panic("all combo count == 0")
		//	fmt.Println(factors)
	}

	for i, factor := range nonfixed {

		//setpoints = make([]interface{}, 0)

		fmt.Println("factor", i, "of", AllComboCount(factors), factor.Factor, factor.Levels)

		//descriptors = append(descriptors, factor.Factor)

		for j, level := range factor.Levels {

			fmt.Println("factor", i, "of", AllComboCount(factors), factor.Factor, "levels ", j, factor.Levels)

			var run Run

			setpoints = append(setpoints, level)

			//		}
			run.Factordescriptors = descriptors
			run.Setpoints = setpoints

			run.StdNumber = len(runs) + 1
			run.RunNumber = len(runs) + 1
			//runs[i+j] = run
			runs = append(runs, run)

		}

	}

	runs = AddFixedFactors(runs, fixed)

	return
}
*/
/*
// original
func AllCombinations(factors []DOEPair) (runs []Run) {
	fmt.Println(factors)

	fixed, nonfixed := FixedAndNonFixed(factors)

	descriptors := make([]string, 0)
	setpoints := make([]interface{}, 0)
	runs = make([]Run, AllComboCount(factors))
	var run Run
	for i, factor := range nonfixed {
		fmt.Println(factor, i, "of", AllComboCount(factors))
		for j, level := range factor.Levels {
			fmt.Println(factor, level, i, j, i+j)
			descriptors = append(descriptors, factor.Factor)
			setpoints = append(setpoints, level)
			run.Factordescriptors = descriptors
			run.Setpoints = setpoints
			runs[i+j] = run
		}
	}

	runs = AddFixedFactors(runs, fixed)

	return
}
*/
/*
func AllCombinations(factors []DOEPair) (runs []Run) {
	//fmt.Println(factors)
	descriptors := make([]string, AllComboCount(factors))
	setpoints := make([]interface{}, 0)
	runs = make([]Run, AllComboCount(factors))
	if AllComboCount(factors) == 0 {
		return
		//	panic("all combo count == 0")
		//	fmt.Println(factors)
	}
	var run Run
	var factorswithonelevel int
	for i, factor := range factors {

		fmt.Println("factor", i, "of", AllComboCount(factors), factor.Factor, factor.Levels)

		if len(factor.Levels) == 1 {
			factorswithonelevel = factorswithonelevel + 1
		}

		for j, level := range factor.Levels {

			//	fmt.Println("factor:", factor, i, j)

			if i-factorswithonelevel < 0 {

				//fmt.Println("factor:", factor, level, i, j) //, i+j)

			//			if len(descriptors) == len(factors) {
			//			descriptors[len(descriptors)-1] = factor.Factor
			//			setpoints[len(descriptors)-1] = level
			//		} else {
				descriptors[i] = factor.Factor
				setpoints = append(setpoints, level)
				//		}

				run.Factordescriptors = descriptors
				run.Setpoints = setpoints
				fmt.Println(i, j, factorswithonelevel, i+1-factorswithonelevel+j)
				runs[i+1-factorswithonelevel+j] = run
			} else {
				//		if len(descriptors) == len(factors) {
				//		descriptors[len(descriptors)-1] = factor.Factor
				//		setpoints[len(descriptors)-1] = level
				//	} else {
				descriptors[i] = factor.Factor
				setpoints = append(setpoints, level)
				//		}
				run.Factordescriptors = descriptors
				run.Setpoints = setpoints
				runs[i-factorswithonelevel+j] = run
			}
		}

	}
	return
}
*/

func ParseRunWellPair(pair string, nameappendage string) (runnumber int, well string, err error) {
	split := strings.Split(pair, ":")

	numberstring := strings.SplitAfter(split[0], nameappendage)

	//numberstring := split[0]
	fmt.Println("Pair", pair, "SPLIT", split /*string(numberstring[0])*/)
	fmt.Println("NUMBERSTRING!!", numberstring /*string(numberstring[0])*/)
	runnumber, err = strconv.Atoi(string(numberstring[1]))
	if err != nil {
		err = fmt.Errorf(err.Error(), "+ Failed at", pair, nameappendage)
	}
	well = split[1]
	return
}

func AddWelllocations(xlsxfile string, oldsheet int, runnumbertowellcombos []string, nameappendage string, pathtosave string, extracolumnheaders []string, extracolumnvalues []interface{}) error {

	var xlsxcell *xlsx.Cell

	file, err := spreadsheet.OpenFile(xlsxfile)
	if err != nil {
		return err
	}

	sheet := spreadsheet.Sheet(file, oldsheet)

	_ = file.AddSheet("hello")

	extracolumn := sheet.MaxCol + 1

	// add extra column headers first
	for _, extracolumnheader := range extracolumnheaders {
		xlsxcell = sheet.Rows[0].AddCell()

		xlsxcell.Value = "Extra column added"
		fmt.Println("CEllll added succesfully", sheet.Cell(0, extracolumn).String())
		xlsxcell = sheet.Rows[1].AddCell()
		xlsxcell.Value = extracolumnheader
	}

	// now add well position column
	xlsxcell = sheet.Rows[0].AddCell()

	xlsxcell.Value = "Location"
	fmt.Println("CEllll added succesfully", sheet.Cell(0, extracolumn).String())
	xlsxcell = sheet.Rows[1].AddCell()
	xlsxcell.Value = "Well ID"
	//	sheet.Cell(0, extracolumn).SetString("Location")
	//	sheet.Cell(1, extracolumn).SetString("Well")

	for i := 3; i < sheet.MaxRow; i++ {
		for _, pair := range runnumbertowellcombos {
			runnumber, well, err := ParseRunWellPair(pair, nameappendage)
			if err != nil {
				return err
			}
			xlsxrunmumber, err := sheet.Cell(i, 1).Int()
			if err != nil {
				return err
			}
			if xlsxrunmumber == runnumber {
				for _, extracolumnvalue := range extracolumnvalues {
					xlsxcell = sheet.Rows[i].AddCell()
					xlsxcell.SetValue(extracolumnvalue)
				}
				xlsxcell = sheet.Rows[i].AddCell()
				xlsxcell.Value = well

			}
		}
	}

	err = file.Save(pathtosave)

	return err
}

func RunsFromDXDesign(xlsx string, intfactors []string) (runs []Run, err error) {
	file, err := spreadsheet.OpenFile(xlsx)
	if err != nil {
		return runs, err
	}
	sheet := spreadsheet.Sheet(file, 0)

	runs = make([]Run, 0)
	var run Run

	var setpoint interface{}
	var descriptor string
	for i := 3; i < sheet.MaxRow; i++ {
		//maxfactorcol := 2
		factordescriptors := make([]string, 0)
		responsedescriptors := make([]string, 0)
		setpoints := make([]interface{}, 0)
		responsevalues := make([]interface{}, 0)
		otherheaders := make([]string, 0)
		othersubheaders := make([]string, 0)
		otherresponsevalues := make([]interface{}, 0)

		run.RunNumber, err = sheet.Cell(i, 1).Int()
		if err != nil {
			return runs, err
		}
		run.StdNumber, err = sheet.Cell(i, 0).Int()
		if err != nil {
			return runs, err
		}

		for j := 2; j < sheet.MaxCol; j++ {
			factororresponse := sheet.Cell(0, j).String()
			//fmt.Println(i, j, factororresponse)
			if strings.Contains(factororresponse, "Factor") {
				//	maxfactorcol = j
				descriptor = strings.Split(sheet.Cell(1, j).String(), ":")[1]
				factrodescriptor := descriptor
				fmt.Println(i, j, descriptor)

				cell := sheet.Cell(i, j)

				celltype := cell.Type()

				_, err := cell.Float()

				if err == nil || celltype == 1 {
					setpoint, _ = cell.Float()
					if search.InSlice(descriptor, intfactors) {
						setpoint, err = cell.Int()
						if err != nil {
							return runs, err
						}
					}
				} else {
					setpoint = cell.String()
				}

				if celltype == 3 {
					setpoint = cell.Bool()
				}
				factordescriptors = append(factordescriptors, factrodescriptor)
				setpoints = append(setpoints, setpoint)

			} else if strings.Contains(factororresponse, "Response") {
				descriptor = sheet.Cell(1, j).String()
				responsedescriptor := descriptor
				//fmt.Println("response", i, j, descriptor)
				responsedescriptors = append(responsedescriptors, responsedescriptor)

				cell := sheet.Cell(i, j)

				if cell == nil {

					break
				}

				celltype := cell.Type()

				if celltype == 1 {
					responsevalue, err := cell.Float()
					if err != nil {
						return runs, err
					}
					responsevalues = append(responsevalues, responsevalue)
				} else {
					responsevalue := cell.String()
					responsevalues = append(responsevalues, responsevalue)
				}

			} else {
				descriptor = sheet.Cell(1, j).String()
				responsedescriptor := descriptor
				//fmt.Println("Additional SubHeader", i, j, factororresponse)
				//fmt.Println("Additional SubHeader", i, j, descriptor)
				otherheaders = append(otherheaders, factororresponse)
				othersubheaders = append(othersubheaders, responsedescriptor)

				cell := sheet.Cell(i, j)

				if cell == nil {

					break
				}

				celltype := cell.Type()

				if celltype == 1 {
					responsevalue, err := cell.Float()
					if err != nil {
						return runs, err
					}
					otherresponsevalues = append(otherresponsevalues, responsevalue)
				} else {
					responsevalue := cell.String()
					otherresponsevalues = append(otherresponsevalues, responsevalue)
				}

			}
		}
		run.Factordescriptors = factordescriptors
		run.Responsedescriptors = responsedescriptors
		run.Setpoints = setpoints
		run.ResponseValues = responsevalues
		run.AdditionalHeaders = otherheaders
		run.AdditionalSubheaders = othersubheaders
		run.AdditionalValues = otherresponsevalues

		runs = append(runs, run)
		factordescriptors = make([]string, 0)
		responsedescriptors = make([]string, 0)

		// assuming this is necessary too
		otherheaders = make([]string, 0)
		othersubheaders = make([]string, 0)
	}

	return
}

func DXXLSXFilefromRuns(runs []Run, outputfilename string) (xlsxfile *xlsx.File) {

	// if output is a struct look for a sensible field to print

	//var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	xlsxfile = xlsx.NewFile()
	sheet = xlsxfile.AddSheet("Sheet1")

	// add headers
	row = sheet.AddRow()

	// 2 blank cells
	cell = row.AddCell()
	cell.Value = ""
	cell = row.AddCell()
	cell.Value = ""

	// take factor and run descriptors from first run (assuming they're all the same)
	for i, _ := range runs[0].Factordescriptors {
		cell = row.AddCell()
		cell.Value = "Factor " + strconv.Itoa(i+1)

	}
	for i, _ := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = "Response " + strconv.Itoa(i+1)

	}
	for _, additionalheader := range runs[0].AdditionalHeaders {
		cell = row.AddCell()
		cell.Value = additionalheader

	}
	// new row
	row = sheet.AddRow()

	// add Std and Run number headers
	cell = row.AddCell()
	cell.Value = "Std"
	cell = row.AddCell()
	cell.Value = "Run"

	// then add subheadings and descriptors
	for i, descriptor := range runs[0].Factordescriptors {

		letter := strings.ToUpper(wutil.Alphabet[i])

		cell = row.AddCell()
		cell.Value = letter + ":" + descriptor

	}
	for _, descriptor := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].AdditionalSubheaders {
		cell = row.AddCell()
		cell.Value = descriptor

	}

	// add blank row

	row = sheet.AddRow()

	//add data 1 row per run
	for _, run := range runs {

		row = sheet.AddRow()
		// Std
		cell = row.AddCell()
		cell.SetValue(run.StdNumber)

		// Run
		cell = row.AddCell()
		cell.SetValue(run.RunNumber)

		// factors
		for _, factor := range run.Setpoints {

			cell = row.AddCell()

			dna, amIdna := factor.(wtype.DNASequence)
			if amIdna {
				cell.SetValue(dna.Nm)
			} else {
				cell.SetValue(factor) //= factor.(string)
			}

		}

		// responses
		for _, response := range run.ResponseValues {
			cell = row.AddCell()
			cell.SetValue(response)
		}

		// additional
		for _, additional := range run.AdditionalValues {
			cell = row.AddCell()
			cell.SetValue(additional)
		}
	}
	err = xlsxfile.Save(outputfilename)
	if err != nil {
		fmt.Printf(err.Error())
	}
	return
}

// jmp

func RunsFromJMPDesign(xlsx string, factorcolumns []int, intfactors []string) (runs []Run, err error) {
	file, err := spreadsheet.OpenFile(xlsx)
	if err != nil {
		return runs, err
	}
	sheet := spreadsheet.Sheet(file, 0)

	runs = make([]Run, 0)
	var run Run

	var setpoint interface{}
	var descriptor string
	for i := 1; i < sheet.MaxRow; i++ {
		//maxfactorcol := 2
		factordescriptors := make([]string, 0)
		responsedescriptors := make([]string, 0)
		setpoints := make([]interface{}, 0)
		responsevalues := make([]interface{}, 0)
		otherheaders := make([]string, 0)
		othersubheaders := make([]string, 0)
		otherresponsevalues := make([]interface{}, 0)

		run.RunNumber = i //sheet.Cell(i, 1).Int()

		run.StdNumber = i //sheet.Cell(i, 0).Int()

		for j := 0; j < sheet.MaxCol; j++ {

			var factororresponse string

			if search.Contains(factorcolumns, j) {
				factororresponse = "Factor"
			} else {
				factororresponse = "Response"
			}
			//fmt.Println(i, j, factororresponse)
			if strings.Contains(factororresponse, "Factor") {
				//	maxfactorcol = j
				descriptor = sheet.Cell(0, j).String()
				factrodescriptor := descriptor
				fmt.Println(i, j, descriptor)

				cell := sheet.Cell(i, j)

				celltype := cell.Type()

				_, err := cell.Float()

				if err == nil || celltype == 1 {
					setpoint, _ = cell.Float()
					if search.InSlice(descriptor, intfactors) {
						setpoint, err = cell.Int()
						if err != nil {
							return runs, err
						}
					}
				} else {
					setpoint = cell.String()
				}

				if celltype == 3 {
					setpoint = cell.Bool()
				}
				factordescriptors = append(factordescriptors, factrodescriptor)
				setpoints = append(setpoints, setpoint)

			} else if strings.Contains(factororresponse, "Response") {
				descriptor = sheet.Cell(0, j).String()
				responsedescriptor := descriptor
				//fmt.Println("response", i, j, descriptor)
				responsedescriptors = append(responsedescriptors, responsedescriptor)

				cell := sheet.Cell(i, j)

				if cell == nil {

					break
				}

				celltype := cell.Type()

				if celltype == 1 {
					responsevalue, err := cell.Float()
					if err != nil {
						return runs, err
					}
					responsevalues = append(responsevalues, responsevalue)
				} else {
					responsevalue := cell.String()
					responsevalues = append(responsevalues, responsevalue)
				}

			} else {
				descriptor = sheet.Cell(0, j).String()
				responsedescriptor := descriptor
				//fmt.Println("Additional SubHeader", i, j, factororresponse)
				//fmt.Println("Additional SubHeader", i, j, descriptor)
				otherheaders = append(otherheaders, factororresponse)
				othersubheaders = append(othersubheaders, responsedescriptor)

				cell := sheet.Cell(i, j)

				if cell == nil {

					break
				}

				celltype := cell.Type()

				if celltype == 1 {
					responsevalue, err := cell.Float()
					if err != nil {
						return runs, err
					}
					otherresponsevalues = append(otherresponsevalues, responsevalue)
				} else {
					responsevalue := cell.String()
					otherresponsevalues = append(otherresponsevalues, responsevalue)
				}

			}
		}
		run.Factordescriptors = factordescriptors
		run.Responsedescriptors = responsedescriptors
		run.Setpoints = setpoints
		run.ResponseValues = responsevalues
		run.AdditionalHeaders = otherheaders
		run.AdditionalSubheaders = othersubheaders
		run.AdditionalValues = otherresponsevalues

		runs = append(runs, run)
		factordescriptors = make([]string, 0)
		responsedescriptors = make([]string, 0)

		// assuming this is necessary too
		otherheaders = make([]string, 0)
		othersubheaders = make([]string, 0)
	}

	return
}

func JMPXLSXFilefromRuns(runs []Run, outputfilename string) (xlsxfile *xlsx.File) {

	// if output is a struct look for a sensible field to print

	//var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	xlsxfile = xlsx.NewFile()
	sheet = xlsxfile.AddSheet("Sheet1")

	// new row
	row = sheet.AddRow()

	// then add subheadings and descriptors
	for _, descriptor := range runs[0].Factordescriptors {

		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].AdditionalSubheaders {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	//add data 1 row per run
	for _, run := range runs {

		row = sheet.AddRow()

		// factors
		for _, factor := range run.Setpoints {

			cell = row.AddCell()

			dna, amIdna := factor.(wtype.DNASequence)
			if amIdna {
				cell.SetValue(dna.Nm)
			} else {
				cell.SetValue(factor) //= factor.(string)
			}

		}

		// responses
		for _, response := range run.ResponseValues {
			cell = row.AddCell()
			cell.SetValue(response)
		}

		// additional
		for _, additional := range run.AdditionalValues {
			cell = row.AddCell()
			cell.SetValue(additional)
		}
	}
	err = xlsxfile.Save(outputfilename)
	if err != nil {
		fmt.Printf(err.Error())
	}
	return
}
