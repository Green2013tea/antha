package wunit

import (
	"fmt"
	"testing"
)

type concConversionTest struct {
	StockConc    Concentration
	TargetConc   Concentration
	TotalVolume  Volume
	VolumeNeeded Volume
}

type massConversionTest struct {
	Conc  Concentration
	Vol   Volume
	Mass  Mass
	Error bool
}

type densityConversionTest struct {
	Density Density
	Vol     Volume
	Mass    Mass
}

var (
	// Some Concentrations
	nilConc Concentration
	uMPer0  = NewConcentration(0, "uM")

	x2   = NewConcentration(2, "X")
	x10  = NewConcentration(10, "X")
	x100 = NewConcentration(100, "X")

	kgPerL01 = NewConcentration(0.1, "kg/L")
	gPerL1   = NewConcentration(1, "g/L")
	gPerL10  = NewConcentration(10, "g/L")
	gPerL100 = NewConcentration(100, "g/L")

	mPerL01  = NewConcentration(0.1, "M/L")
	mMPerL1  = NewConcentration(1, "mM/L")
	uMPerL10 = NewConcentration(10, "uM")

	// Some volumes
	nilVol Volume
	l0     = NewVolume(0, "l")
	l01    = NewVolume(0.1, "l")
	ul1    = NewVolume(1, "ul")
	ul10   = NewVolume(10, "ul")
	ul100  = NewVolume(100, "ul")
	ml1    = NewVolume(1, "ml")
)

var tests1 = []concConversionTest{
	concConversionTest{StockConc: kgPerL01, TargetConc: gPerL1, TotalVolume: ul100, VolumeNeeded: ul1},
	concConversionTest{StockConc: gPerL100, TargetConc: gPerL1, TotalVolume: ul100, VolumeNeeded: ul1},
	concConversionTest{StockConc: x100, TargetConc: x2, TotalVolume: ul100, VolumeNeeded: NewVolume(2.0, "ul")},
	concConversionTest{StockConc: mMPerL1, TargetConc: uMPerL10, TotalVolume: ul100, VolumeNeeded: NewVolume(1.0, "ul")},
	concConversionTest{StockConc: mPerL01, TargetConc: mMPerL1, TotalVolume: ul100, VolumeNeeded: NewVolume(1.0, "ul")},
	concConversionTest{StockConc: mPerL01, TargetConc: uMPer0, TotalVolume: ul100, VolumeNeeded: l0},
}

var tests2 = []massConversionTest{
	massConversionTest{Conc: NewConcentration(1.0, "g/L"), Mass: NewMass(1000.0, "mg"), Vol: NewVolume(1.0, "l"), Error: false},
	massConversionTest{Conc: NewConcentration(1.0, "kg/L"), Mass: NewMass(1.0, "kg"), Vol: NewVolume(1.0, "l"), Error: false},
	massConversionTest{Conc: NewConcentration(0.1, "mg/L"), Mass: NewMass(0.1, "mg"), Vol: NewVolume(1.0, "l"), Error: false},
	massConversionTest{Conc: NewConcentration(100, "ng/ul"), Mass: NewMass(100, "ng"), Vol: NewVolume(1.0, "ul"), Error: false},
	massConversionTest{Conc: NewConcentration(0, "g/l"), Mass: NewMass(0, "g"), Vol: NewVolume(1.0, "l"), Error: true},
}

var tests3 = []densityConversionTest{
	densityConversionTest{Density: NewDensity(1.0, "kg/m^3"), Mass: NewMass(1.0, "kg"), Vol: NewVolume(1000, "l")},
	densityConversionTest{Density: NewDensity(1000.0, "kg/m^3"), Mass: NewMass(1.0, "g"), Vol: NewVolume(1, "ml")},
	densityConversionTest{Density: NewDensity(1000.0, "kg/m^3"), Mass: NewMass(0.0, "g"), Vol: NewVolume(0, "l")},
}

func TestVolumeForTargetConcentration(t *testing.T) {

	for _, test := range tests1 {

		vol, err := VolumeForTargetConcentration(test.TargetConc, test.StockConc, test.TotalVolume)

		if err != nil {
			t.Error(
				"for", test, "\n",
				"got error:", err.Error(), "\n",
			)
		}

		if !vol.EqualTo(test.VolumeNeeded) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected Vol:", test.VolumeNeeded.ToString(), "\n",
				"Got Vol:", vol.ToString(), "\n",
			)
		}

	}
}

func TestMassForTargetConcentration(t *testing.T) {

	for _, test := range tests2 {

		mass, err := MassForTargetConcentration(test.Conc, test.Vol)

		if err != nil {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"got error:", err.Error(), "\n",
			)
		}

		if !mass.EqualTo(test.Mass) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected mass:", test.Mass.ToString(), "\n",
				"Got mass:", mass.ToString(), "\n",
			)
		}

	}
}

func TestVolumeForTargetMass(t *testing.T) {

	for _, test := range tests2 {

		vol, err := VolumeForTargetMass(test.Mass, test.Conc)

		if err != nil && !test.Error {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"got error:", err.Error(), "\n",
			)
		} else if err == nil && test.Error {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"expected error but got none.", "\n",
			)
		}

		if !vol.EqualTo(test.Vol) && !test.Error {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected vol:", test.Vol.ToString(), "\n",
				"Got vol:", vol.ToString(), "\n",
			)
		}

	}
}

func TestMasstoVolume(t *testing.T) {

	for _, test := range tests3 {

		vol := MasstoVolume(test.Mass, test.Density)

		if !vol.EqualTo(test.Vol) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected vol:", test.Vol.ToString(), "\n",
				"Got vol:", vol.ToString(), "\n",
			)
		}

	}
}

func TestVolumetoMass(t *testing.T) {

	for _, test := range tests3 {

		mass := VolumetoMass(test.Vol, test.Density)

		if !mass.EqualTo(test.Mass) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected mass:", test.Mass.ToString(), "\n",
				"Got mass:", mass.ToString(), "\n",
			)
		}

	}
}
