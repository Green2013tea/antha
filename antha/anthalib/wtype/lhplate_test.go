package wtype

import (
	"fmt"
	"strings"
	"testing"
)

// platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64

func makeplatefortest() *LHPlate {
	swshp := NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell(nil, ZeroWellCoords(), "ul", 1000, 100, swshp, VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("testplate", "none", 8, 12, Coordinates{127.76, 85.48, 43.1}, welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
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

		if w.Crds.FormatA1() != crds {
			t.Fatal(fmt.Sprintf("ERROR: Well coords not consistent -- %s != %s", w.Crds.FormatA1(), crds))
		}

		if w.WContents.Loc == "" {
			t.Fatal(fmt.Sprintf("ERROR: Well contents do not have loc set"))
		}

		ltx := strings.Split(w.WContents.Loc, ":")

		if ltx[0] != plate.ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent -- %s != %s", ltx[0], plate.ID))
		}

		if ltx[0] != w.Plate.(*LHPlate).ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent with well -- %s != %s", ltx[0], w.Plate.(*LHPlate).ID))
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
	assertWellsEqual("Cols != WellCoords", ws3, ws4)

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
