package wtype

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

func A1ArrayFromWells(wells []*LHWell) []string {
	return A1ArrayFromWellCoords(WCArrayFromWells(wells))
}

func WCArrayFromWells(wells []*LHWell) []WellCoords {
	ret := make([]WellCoords, 0, len(wells))

	for _, w := range wells {
		if w == nil {
			continue
		}

		ret = append(ret, MakeWellCoords(w.Crds))
	}

	return ret
}

func WCArrayFromStrings(arr []string) []WellCoords {
	ret := make([]WellCoords, len(arr))

	for i, s := range arr {
		ret[i] = MakeWellCoords(s)
	}

	return ret
}

func A1ArrayFromWellCoords(arr []WellCoords) []string {
	ret := make([]string, len(arr))
	for i, v := range arr {
		ret[i] = v.FormatA1()
	}
	return ret
}

// make an array of these from an array of strings

func MakeWellCoordsArray(sa []string) []WellCoords {
	r := make([]WellCoords, len(sa))

	for i := 0; i < len(sa); i++ {
		r[i] = MakeWellCoords(sa[i])
	}

	return r
}

func WCArrayCols(wcA []WellCoords) []int {
	return squashedIntFromWCA(wcA, 0)
}

func WCArrayRows(wcA []WellCoords) []int {
	return squashedIntFromWCA(wcA, 1)
}

func containsInt(i int, ia []int) bool {
	for _, ii := range ia {
		if i == ii {
			return true
		}
	}
	return false
}

func squashedIntFromWCA(wcA []WellCoords, which int) []int {
	ret := make([]int, 0, len(wcA))
	for _, wc := range wcA {
		v := wc.X
		if which == 1 {
			v = wc.Y
		}

		// ignore nils

		if v == -1 {
			continue
		}

		if !containsInt(v, ret) {
			ret = append(ret, v)
		}
	}
	return ret
}

// convenience comparison operator

func CompareStringWellCoordsCol(sw1, sw2 string) int {
	w1 := MakeWellCoords(sw1)
	w2 := MakeWellCoords(sw2)
	return CompareWellCoordsCol(w1, w2)
}

func CompareWellCoordsCol(w1, w2 WellCoords) int {
	dx := w1.X - w2.X
	dy := w1.Y - w2.Y

	if dx < 0 {
		return -1
	} else if dx > 0 {
		return 1
	}

	if dy < 0 {
		return -1
	} else if dy > 0 {
		return 1
	} else {
		return 0
	}
}

func CompareStringWellCoordsRow(sw1, sw2 string) int {
	w1 := MakeWellCoords(sw1)
	w2 := MakeWellCoords(sw2)
	return CompareWellCoordsRow(w1, w2)
}

func CompareWellCoordsRow(w1, w2 WellCoords) int {
	dx := w1.X - w2.X
	dy := w1.Y - w2.Y

	if dy < 0 {
		return -1
	} else if dy > 0 {
		return 1
	}
	if dx < 0 {
		return -1
	} else if dx > 0 {
		return 1
	} else {
		return 0
	}
}

// convenience structure for handling well coordinates
type WellCoords struct {
	X int
	Y int
}

func ZeroWellCoords() WellCoords {
	return WellCoords{-1, -1}
}
func (wc WellCoords) IsZero() bool {
	if wc.Equals(ZeroWellCoords()) {
		return true
	}

	return false
}

func MatchString(s1, s2 string) bool {
	m, _ := regexp.MatchString(s1, s2)
	return m
}

func (wc WellCoords) Equals(w2 WellCoords) bool {
	if wc.X == w2.X && wc.Y == w2.Y {
		return true
	}

	return false
}

func MakeWellCoords(wc string) WellCoords {
	// try each one in turn

	r := MakeWellCoordsA1(wc)

	zero := WellCoords{-1, -1}

	if !r.Equals(zero) {
		return r
	}

	r = MakeWellCoords1A(wc)

	if !r.Equals(zero) {
		return r
	}

	r = MakeWellCoordsXY(wc)

	return r
}

// make well coordinates in the "A1" convention
func MakeWellCoordsA1(a1 string) WellCoords {
	re := regexp.MustCompile(`^([A-Z]{1,})([0-9]{1,2})$`)
	matches := re.FindStringSubmatch(a1)

	if matches == nil {
		return WellCoords{-1, -1}
	}
	/*
		re, _ := regexp.Compile("[A-Z]{1,}")
		ix := re.FindIndex([]byte(a1))
		endC := ix[1]
	*/

	X := wutil.ParseInt(matches[2]) - 1
	Y := wutil.AlphaToNum(matches[1]) - 1

	return WellCoords{X, Y}
}

// make well coordinates in the "1A" convention
func MakeWellCoords1A(a1 string) WellCoords {
	re := regexp.MustCompile(`^([0-9]{1,2})([A-Z]{1,})$`)
	matches := re.FindStringSubmatch(a1)

	if matches == nil {
		return WellCoords{-1, -1}
	}

	Y := wutil.AlphaToNum(matches[2]) - 1
	X := wutil.ParseInt(matches[1]) - 1
	return WellCoords{X, Y}
}

// make well coordinates in a manner compatble with "X1,Y1" etc.
func MakeWellCoordsXYsep(x, y string) WellCoords {
	r := WellCoords{wutil.ParseInt(y[1:len(y)]) - 1, wutil.ParseInt(x[1:len(x)]) - 1}

	if r.X < 0 || r.Y < 0 {
		return WellCoords{-1, -1}
	}

	return r
}

func MakeWellCoordsXY(xy string) WellCoords {
	tx := strings.Split(xy, "Y")
	if tx == nil || len(tx) != 2 || len(tx[0]) == 0 || len(tx[1]) == 0 {
		return WellCoords{-1, -1}
	}
	x := wutil.ParseInt(tx[0][1:len(tx[0])]) - 1
	y := wutil.ParseInt(tx[1]) - 1
	return WellCoords{x, y}
}

// return well coordinates in "X1Y1" format
func (wc WellCoords) FormatXY() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return "X" + strconv.Itoa(wc.X+1) + "Y" + strconv.Itoa(wc.Y+1)
}

func (wc WellCoords) Format1A() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return strconv.Itoa(wc.X+1) + wutil.NumToAlpha(wc.Y+1)
}

func (wc WellCoords) FormatA1() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return wutil.NumToAlpha(wc.Y+1) + strconv.Itoa(wc.X+1)
}

func (wc WellCoords) WellNumber() int {
	if wc.X < 0 || wc.Y < 0 {
		return -1
	}
	return (8*(wc.X-1) + wc.Y)
}

func (wc WellCoords) ColNumString() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return strconv.Itoa(wc.X + 1)
}

func (wc WellCoords) RowLettString() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return wutil.NumToAlpha(wc.Y + 1)
}

// comparison operators

func (wc WellCoords) RowLessThan(wc2 WellCoords) bool {
	if wc.Y == wc2.Y {
		return wc.X < wc2.Y
	}
	return wc.Y < wc2.Y
}

func (wc WellCoords) ColLessThan(wc2 WellCoords) bool {
	if wc.X == wc2.X {
		return wc.Y < wc2.Y
	}
	return wc.X < wc2.X
}

// convenience structure to allow sorting

type WellCoordArrayCol []WellCoords
type WellCoordArrayRow []WellCoords

func (wca WellCoordArrayCol) Len() int           { return len(wca) }
func (wca WellCoordArrayCol) Swap(i, j int)      { t := wca[i]; wca[i] = wca[j]; wca[j] = t }
func (wca WellCoordArrayCol) Less(i, j int) bool { return wca[i].RowLessThan(wca[j]) }

func (wca WellCoordArrayRow) Len() int           { return len(wca) }
func (wca WellCoordArrayRow) Swap(i, j int)      { t := wca[i]; wca[i] = wca[j]; wca[j] = t }
func (wca WellCoordArrayRow) Less(i, j int) bool { return wca[i].ColLessThan(wca[j]) }
