package wtype

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64

func makeplatefortest() *LHPlate {
	swshp := NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("DSW96", "", "", "ul", 200, 10, swshp, LHWBV, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("testplate", "none", 8, 12, 44.1, "mm", welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}
func make384platefortest() *LHPlate {
	swshp := NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("DSW384", "", "", "ul", 50, 5, swshp, LHWBV, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("testplate", "none", 16, 24, 44.1, "mm", welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}
func make1536platefortest() *LHPlate {
	swshp := NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("DSW1536", "", "", "ul", 15, 1, swshp, LHWBV, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("testplate", "none", 32, 48, 44.1, "mm", welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}
func make24platefortest() *LHPlate {
	swshp := NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("DSW24", "", "", "ul", 3000, 500, swshp, LHWBV, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("testplate", "none", 4, 6, 44.1, "mm", welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}
func make6platefortest() *LHPlate {
	swshp := NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("6wellplate", "", "", "ul", 3000, 500, swshp, LHWBV, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("testplate", "none", 2, 3, 44.1, "mm", welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}

func TestPlateCreation(t *testing.T) {
	p := makeplatefortest()
	validatePlate(t, p)
}

func TestPlateDup(t *testing.T) {
	p := makeplatefortest()
	d := p.Dup()
	validatePlate(t, d)
	for crds, w := range p.Wellcoords {
		w2 := d.Wellcoords[crds]

		if w.ID == w2.ID {
			t.Fatal(fmt.Sprintf("Error: coords %s has same IDs before / after dup", crds))
		}

		if w.WContents.Loc == w2.WContents.Loc {
			t.Fatal(fmt.Sprintf("Error: contents of wells at coords %s have same loc before and after regular Dup()", crds))
		}
	}
}

func TestPlateDupKeepIDs(t *testing.T) {
	p := makeplatefortest()
	d := p.DupKeepIDs()

	for crds, w := range p.Wellcoords {
		w2 := d.Wellcoords[crds]

		if w.ID != w2.ID {
			t.Fatal(fmt.Sprintf("Error: coords %s has different IDs", crds))
		}

		if w.WContents.ID != w2.WContents.ID {
			t.Fatal(fmt.Sprintf("Error: contents of wells at coords %s have different IDs", crds))

		}
		if w.WContents.Loc != w2.WContents.Loc {
			t.Fatal(fmt.Sprintf("Error: contents of wells at coords %s have different loc before and after DupKeepIDs()", crds))
		}
	}

}

func validatePlate(t *testing.T, plate *LHPlate) {
	assertWellsEqual := func(what string, as, bs []*LHWell) {
		seen := make(map[*LHWell]int)
		for _, w := range as {
			seen[w] += 1
		}
		for _, w := range bs {
			seen[w] += 1
		}
		for w, count := range seen {
			if count != 2 {
				t.Errorf("%s: no matching well found (%d != %d) for %p %s:%s", what, count, 2, w, w.ID, w.Crds)
			}
		}
	}

	var ws1, ws2, ws3, ws4 []*LHWell

	for _, w := range plate.HWells {
		ws1 = append(ws1, w)
	}
	for crds, w := range plate.Wellcoords {
		ws2 = append(ws2, w)

		if w.Crds != crds {
			t.Fatal(fmt.Sprintf("ERROR: Well coords not consistent -- %s != %s", w.Crds, crds))
		}

		if w.WContents.Loc == "" {
			t.Fatal(fmt.Sprintf("ERROR: Well contents do not have loc set"))
		}

		ltx := strings.Split(w.WContents.Loc, ":")

		if ltx[0] != plate.ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent -- %s != %s", ltx[0], plate.ID))
		}

		if ltx[0] != w.Plateid {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent with well -- %s != %s", ltx[0], w.Plateid))
		}

		if ltx[1] != crds {
			t.Fatal(fmt.Sprintf("ERROR: Coords for component not consistent: -- %s != %s", ltx[1], crds))
		}

	}

	for _, ws := range plate.Rows {
		for _, w := range ws {
			ws3 = append(ws3, w)
		}
	}
	for _, ws := range plate.Cols {
		for _, w := range ws {
			ws4 = append(ws4, w)
		}

	}
	assertWellsEqual("HWells != Rows", ws1, ws2)
	assertWellsEqual("Rows != Cols", ws2, ws3)
	assertWellsEqual("Cols != Wellcoords", ws3, ws4)

	// Check pointer-ID equality
	comp := make(map[string]*LHComponent)
	for _, w := range append(append(ws1, ws2...), ws3...) {
		c := w.WContents
		if c == nil || c.Vol == 0.0 {
			continue
		}
		if co, seen := comp[c.ID]; seen && co != c {
			t.Errorf("component %s duplicated as %+v and %+v", c.ID, c, co)
		} else if !seen {
			comp[c.ID] = c
		}
	}
}

func TestIsUserAllocated(t *testing.T) {
	p := makeplatefortest()

	if p.IsUserAllocated() {
		t.Fatal("Error: Plates must not start out user allocated")
	}
	p.Wellcoords["A1"].SetUserAllocated()

	if !p.IsUserAllocated() {
		t.Fatal("Error: Plates with at least one user allocated well must return true to IsUserAllocated()")
	}

	d := p.Dup()

	if !d.IsUserAllocated() {
		t.Fatal("Error: user allocation mark must survive Dup()lication")
	}

	d.Wellcoords["A1"].ClearUserAllocated()

	if d.IsUserAllocated() {
		t.Fatal("Error: user allocation mark not cleared")
	}

	if !p.IsUserAllocated() {
		t.Fatal("Error: UserAllocation mark must operate separately on Dup()licated plates")
	}
}

func TestMergeWith(t *testing.T) {
	p1 := makeplatefortest()
	p2 := makeplatefortest()

	c := NewLHComponent()

	c.CName = "Water1"
	c.Vol = 50.0
	c.Vunit = "ul"
	p1.Wellcoords["A1"].Add(c)
	p1.Wellcoords["A1"].SetUserAllocated()

	c = NewLHComponent()
	c.CName = "Butter"
	c.Vol = 80.0
	c.Vunit = "ul"
	p2.Wellcoords["A2"].Add(c)

	p1.MergeWith(p2)

	if !(p1.Wellcoords["A1"].WContents.CName == "Water1" && p1.Wellcoords["A1"].WContents.Vol == 50.0 && p1.Wellcoords["A1"].WContents.Vunit == "ul") {
		t.Fatal("Error: MergeWith should leave user allocated components alone")
	}

	if !(p1.Wellcoords["A2"].WContents.CName == "Butter" && p1.Wellcoords["A2"].WContents.Vol == 80.0 && p1.Wellcoords["A2"].WContents.Vunit == "ul") {
		t.Fatal("Error: MergeWith should add non user-allocated components to  plate merged with")
	}
}

func makeCV(name string, vol float64) ComponentVector {
	c := NewLHComponent()
	c.Type = LTWater
	c.CName = name
	c.Vol = vol
	CIDs := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}
	PIDs := []string{"Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1", "Plate1"}

	got := make([]*LHComponent, 8)

	for i := 0; i < 8; i++ {
		got[i] = c.Dup()
		got[i].Loc = PIDs[i] + ":" + CIDs[i]
	}

	return got
}

func makecomponent(cname string, vol float64) *LHComponent {
	c := NewLHComponent()
	c.Type = LTWater
	c.CName = cname
	c.Vol = vol
	c.Vunit = "ul"
	return c
}

/*
func TestFindCompMulti1(t *testing.T) {
	p := makeplatefortest()
	c := makecomponent("water", 1600.0)
	p.AddComponent(c, true)
	cv := makeCV("water", 50.0)

	pids, _, _, _ := p.FindComponentsMulti(cv, LHVChannel, 8, false)

	if len(pids) == 0 {
		t.Errorf("Didn't find a simple column of water... should have")
	}
}
*/

func TestLHPlateSerialize(t *testing.T) {
	p := makeplatefortest()
	c := NewLHComponent()
	c.CName = "Cthulhu"
	c.Type = LTWater
	c.Vol = 100.0
	_, err := p.AddComponent(c, false)

	if err != nil {
		t.Errorf(err.Error())
	}
	b, err := json.Marshal(p)

	if err != nil {
		t.Errorf(err.Error())
	}

	var p2 *LHPlate

	if err = json.Unmarshal(b, &p2); err != nil {
		t.Errorf(err.Error())
	}

	for i, w := range p.Wellcoords {
		w2 := p2.Wellcoords[i]

		if !reflect.DeepEqual(w.WContents, w2.WContents) {
			t.Errorf("%v =/= %v", w.WContents, w2.WContents)
		}
	}

	fMErr := func(s string) string {
		return s + " not maintained after marshal/unmarshal"
	}

	for i := 0; i < p2.WellsX(); i++ {
		for j := 0; j < p2.WellsY(); j++ {
			wc := WellCoords{X: i, Y: j}

			w := p2.Wellcoords[wc.FormatA1()]

			w.WContents.CName = wc.FormatA1()
			if p2.Rows[j][i].WContents.CName != wc.FormatA1() || p2.Cols[i][j].WContents.CName != wc.FormatA1() || p2.HWells[w.ID].WContents.CName != wc.FormatA1() {
				fmt.Println(p2.Cols[i][j].WContents.CName)
				fmt.Println(p2.Rows[j][i].WContents.CName)
				t.Errorf("Error: Wells inconsistent at position %s", wc.FormatA1())
			}

		}
	}

	// check extraneous parameters

	if p.ID != p2.ID {
		t.Errorf(fMErr("ID"))
	}

	if p.PlateName != p2.PlateName {
		t.Errorf(fMErr("Plate name"))
	}

	if p.Type != p2.Type {
		t.Errorf(fMErr("Type"))
	}

	if p.Mnfr != p2.Mnfr {
		t.Errorf(fMErr("Manufacturer"))
	}

	if p.Nwells != p2.Nwells {
		t.Errorf(fMErr("NWells"))
	}

	if p.Height != p2.Height {
		t.Errorf(fMErr("Height"))
	}

	if p.Hunit != p2.Hunit {
		t.Errorf(fMErr("Hunit"))
	}

	if p.WellXOffset != p2.WellXOffset {
		t.Errorf(fMErr("WellXOffset"))
	}

	if p.WellYOffset != p2.WellYOffset {
		t.Errorf(fMErr("WellYOffset"))
	}

	if p.WellXStart != p2.WellXStart {
		t.Errorf(fMErr("WellXStart"))
	}
	if p.WellYStart != p2.WellYStart {
		t.Errorf(fMErr("WellYStart"))
	}

	if p.WellZStart != p2.WellZStart {
		t.Errorf(fMErr("WellZStart"))
	}
}

func TestAddGetClearData(t *testing.T) {
	dat := []byte("3.5")

	t.Run("basic", func(t *testing.T) {
		p := makeplatefortest()

		if err := p.SetData("OD", dat); err != nil {
			t.Errorf(err.Error())
		}
		d, err := p.GetData("OD")
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reflect.DeepEqual(d, dat) {
			t.Errorf("Expected %v got %v", dat, d)
		}
	})

	t.Run("clear", func(t *testing.T) {
		p := makeplatefortest()

		if err := p.SetData("OD", dat); err != nil {
			t.Errorf(err.Error())
		}

		if err := p.ClearData("OD"); err != nil {
			t.Errorf(err.Error())
		}

		if _, err := p.GetData("OD"); err == nil {
			t.Errorf("ClearData should clear data but has not")
		}
	})

	t.Run("cannot update special", func(t *testing.T) {
		p := makeplatefortest()
		if err := p.SetData("IMSPECIAL", dat); err == nil {
			t.Errorf("Adding data with a reserved key should fail but does not")
		}
	})

}

func TestGetAllComponents(t *testing.T) {
	p := makeplatefortest()

	cmps := p.AllContents()

	if len(cmps) != p.WellsX()*p.WellsY() {
		t.Errorf("Expected %d components got %d", p.WellsX()*p.WellsY(), len(cmps))
	}
}
