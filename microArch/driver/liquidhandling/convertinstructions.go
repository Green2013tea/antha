// liquidhandling/convertinstructions.go Part of the Antha language
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

package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

//	this section aggregates instructions with the following constraints:
//
//	1) obey any requirement to do one sample at a time
//		-- bullet bitten: we cannot permit transfer to split up any multichannel instructions
//		   into singles here... this is a bit tricky but we must make it so
//		   some revision to how pragmas work may be needed: extend only to component type etc.
//
//	here is what a single sample assembled one thing at a time looks like
//	|
//	|		here is one sample assembled one component at a time looks like
//	|		|
//	i1(A)		i2(A B C)	--> the LHIVector contains these two, maxlen = 3, CmpAt (0) = [A A]
//	--										  CmpAt (1) = [  B]
//	i3(B) <------									  CmpAt (2) = [  C]
//	--          |-- these two are done separately (so they're boring)
//	i4(C) <------
//
// 	this should produce the output:
//	TFR(A A d1 d2), TFR(B d2), TFR(C d2), TFR(B d1), TFR(C d1)
//	iow it does i1 + first part of i2 in parallel, then the rest of i2 then i3 then i4

// 	issue is we cannot tolerate this situation
//
//	i1(A)		i2(A B)		i3(A C)
//	so we have to ensure the components line up

func ConvertInstructions(inssIn LHIVector, robot *LHProperties, carryvol wunit.Volume, channelprms *wtype.LHChannelParameter, multi int) (insOut []*TransferInstruction, err error) {
	insOut = make([]*TransferInstruction, 0, 1)

	for i := 0; i < inssIn.MaxLen(); i++ {
		comps := inssIn.CompsAt(i)
		lenToMake := 0
		// remove spaces between components
		cmpSquash := make([]*wtype.LHComponent, 0, lenToMake)
		for _, c := range comps {
			if c != nil {
				lenToMake += 1
				cmpSquash = append(cmpSquash, c)
			}
		}

		wh := make([]string, lenToMake)       // component types
		va := make([]wunit.Volume, lenToMake) // volumes
		// six parameters applying to the source
		// TODO --> this should create components if not already found
		fmt.Println(cmpSquash)
		fmt.Println(carryvol)
		fmt.Println(channelprms)
		fromPlateIDs, fromWells, fromvols, err := robot.GetComponents(cmpSquash, carryvol, channelprms.Orientation, multi, channelprms.Independent)

		// let's start making sense here

		if !allEqual(len(fromPlateIDs), len(fromWells), len(fromvols)) {
			panic("Lengths cannot differ here")
		}

		if err != nil {
			return nil, err
		}

		pf := make([]string, lenToMake)       // src plate positions
		wf := make([]string, lenToMake)       // src wells
		pfwx := make([]int, lenToMake)        // src plate X dim
		pfwy := make([]int, lenToMake)        // src plate Y dim
		vf := make([]wunit.Volume, lenToMake) // volumes
		ptf := make([]string, lenToMake)      // plate types

		// six parameters applying to the destination

		pt := make([]string, lenToMake)       // dest plate positions
		wt := make([]string, lenToMake)       // dest wells
		ptwx := make([]int, lenToMake)        // dimensions of plate pipetting to (X)
		ptwy := make([]int, lenToMake)        // dimensions of plate pipetting to (Y)
		vt := make([]wunit.Volume, lenToMake) // volume in well to
		ptt := make([]string, lenToMake)      // plate types

		ix := 0 // counts up cmpsquash

		for j, v := range comps {

			if comps[j] == nil {
				continue
			}

			var flhp, tlhp *wtype.LHPlate

			flhif := robot.PlateLookup[fromPlateIDs[ix][0]]

			if flhif != nil {
				flhp = flhif.(*wtype.LHPlate)
			} else {
				s := fmt.Sprint("NO SRC PLATE FOUND : ", ix, " ", fromPlateIDs[ix])
				err := wtype.LHError(wtype.LH_ERR_DIRE, s)

				return nil, err
			}

			tlhif := robot.PlateLookup[inssIn[j].PlateID()]

			if tlhif != nil {
				tlhp = tlhif.(*wtype.LHPlate)
			} else {
				s := fmt.Sprint("NO DST PLATE FOUND : ", ix, " ", inssIn[j].PlateID())
				err := wtype.LHError(wtype.LH_ERR_DIRE, s)

				return nil, err
			}

			wlt, ok := tlhp.WellAtString(inssIn[j].Welladdress)

			if !ok {
				err = wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("Well ", inssIn[j].Welladdress, " not found on dest plate ", inssIn[j].PlateID))
				return nil, err
			}

			v2 := wunit.NewVolume(v.Vol, v.Vunit)
			vt[ix] = wlt.CurrVolume()
			wh[ix] = v.TypeName()
			va[ix] = v2
			pt[ix] = robot.PlateIDLookup[inssIn[j].PlateID()]
			wt[ix] = inssIn[j].Welladdress
			ptwx[ix] = tlhp.WellsX()
			ptwy[ix] = tlhp.WellsY()
			ptt[ix] = tlhp.Type

			wlf, ok := flhp.WellAtString(fromWells[ix][0])

			if !ok {
				//logger.Fatal(fmt.Sprint("Well ", fromWells[ix], " not found on source plate ", fromPlateIDs[ix]))
				err = wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("Well ", fromWells[ix], " not found on source plate ", fromPlateIDs[ix]))
				return nil, err
			}

			vf[ix] = fromvols[j][0]
			//wlf.Remove(va[ix])

			pf[ix] = robot.PlateIDLookup[fromPlateIDs[ix][0]]
			wf[ix] = fromWells[ix][0]
			pfwx[ix] = flhp.WellsX()
			pfwy[ix] = flhp.WellsY()
			ptf[ix] = flhp.Type

			if v.Loc == "" {
				v.Loc = fromPlateIDs[ix][0] + ":" + fromWells[ix][0]
			}
			// add component to destination
			// need to ensure data are consistent
			vd := v.Dup()
			vd.ID = wlf.WContents.ID
			vd.ParentID = wlf.WContents.ParentID
			wlt.Add(vd)

			// add daughter ID to component in

			wlf.WContents.AddDaughterComponent(wlt.WContents)

			ix += 1
		}

		tfr := NewTransferInstruction(wh, pf, pt, wf, wt, ptf, ptt, va, vf, vt, pfwx, pfwy, ptwx, ptwy)
		insOut = append(insOut, tfr)
	}

	return insOut, nil
}
