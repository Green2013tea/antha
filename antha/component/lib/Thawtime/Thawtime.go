//status = compiles and calculates; need to fill in correct parameters and check units
//currently using dummy values only so won't be accurate yet!
// Once working move from floats to antha types and units
package Thawtime

import (
	"fmt"                                                                 // we need this go library to print
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng" // all of our functions used here are in the Thaw.go file in the eng package which this points to
	//"github.com/montanaflynn/stats" // a rounding function is used from this third party library
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Many of the real parameters required will be looked up via the specific labware (platetype) and liquid type which are being used.

/* e.g. the sample volume as frozen by a previous storage protocol;
could be known or measured via liquid height detection on some liquid handlers */

// These should be captured via sensors just prior to execution

// This will be monitored via the thermometer in the freezer in which the sample was stored

/* This will offer another knob to tweak (in addition to the other parameters) as a means to improve
the correlation over time as we see how accurate the calculator is in practice */

func _requirements() {
}
func _setup(_ctx context.Context, _input *Input_) {
}
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	/*  Step 1. we need a mass for the following equations so we calculate this by looking up
	the liquid density and multiplying by the fill volume using this function from the engineering library */

	//fillvolume:= Fillvolume.SIValue()

	mass := eng.Massfromvolume(_input.Fillvolume, _input.Liquid)

	/*  Step 2. Required heat energy to melt the solid is calculated using the calculated mass along with the latent heat of melting
	which we find via a liquid class look up package which is not required for import here since it's imported from the engineering library */

	q := eng.Q(_input.Liquid, mass)

	/*  Step 3. Heat will be transferred via both convection through the air and conduction through the plate walls.
	Let's first work out the heat energy transferred via convection, this uses an empirical parameter,
	the convective heat transfer coefficient of air (HC_air), this is calculated via another function in the eng library.
	In future we will make this process slightly more sophisticated by adding conditions, since this empirical equation is
	only validated between air velocities 2 - 20 m/s. It could also be adjusted to calculate heat transfer if the sample
	is agitated on a shaker to speed up thawing. */

	hc_air := eng.Hc_air(_input.Airvelocity.SIValue())

	/*  Step 4. The rate of heat transfer by convection is then calculated using this value combined with the temperature differential
	(measured by the temp sensor) and surface area dictated by the plate type (another look up called from the eng library!)*/

	convection := eng.ConvectionPowertransferred(hc_air, _input.Platetype, _input.SurfaceTemp, _input.BulkTemp)

	/*  Step 5. We now estimate the heat transfer rate via conduction. For this we need to know the thermal conductivity of the plate material
	along with the wall thickness. As before, both of these are looked up via the labware library called by this function in the eng library */

	conduction := eng.ConductionPowertransferred(_input.Platetype, _input.SurfaceTemp, _input.BulkTemp)

	/*  Step 6. We're now ready to estimate the thawtime needed by simply dividing the estimated heat required to melt/thaw (i.e. q from step 2)
	by the combined rate of heat transfer estimated to occur via both convection and conduction */
	_output.Estimatedthawtime = eng.Thawtime(convection, conduction, q)

	/* Step 7. Since there're a lot of assumptions here (liquid behaves as water, no change in temperature gradient, no heat transferred via radiation,
	imprecision in the estimates and 	empirical formaulas) we'll multiply by a fudgefactor to be safer that we've definitely thawed,
	this (and all parameters!) can be adjusted over time as we see emprically how reliable this function is as more datapoints are collected */
	_output.Thawtimeused = wunit.NewTime(_output.Estimatedthawtime.SIValue()*_input.Fudgefactor, "s")

	_output.Status = fmt.Sprintln("For", mass.ToString(), "of", _input.Liquid, "in", _input.Platetype,
		"Thawtime required =", _output.Estimatedthawtime.ToString(),
		"Thawtime used =", _output.Thawtimeused.ToString(),
		"power required =", q, "J",
		"HC_air (convective heat transfer coefficient=", hc_air,
		"Convective power=", convection, "J/s",
		"conductive power=", conduction, "J/s")

}
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {

}

func _validation(_ctx context.Context, _input *Input_, _output *Output_) {

}

func _run(_ctx context.Context, value inject.Value) (inject.Value, error) {
	input := &Input_{}
	output := &Output_{}
	if err := inject.Assign(value, input); err != nil {
		return nil, err
	}
	_setup(_ctx, input)
	_steps(_ctx, input, output)
	_analysis(_ctx, input, output)
	_validation(_ctx, input, output)
	return inject.MakeValue(output), nil
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

func New() interface{} {
	return &Element_{
		inject.CheckedRunner{
			RunFunc: _run,
			In:      &Input_{},
			Out:     &Output_{},
		},
	}
}

type Element_ struct {
	inject.CheckedRunner
}

type Input_ struct {
	Airvelocity wunit.Velocity
	BulkTemp    wunit.Temperature
	Fillvolume  wunit.Volume
	Fudgefactor float64
	Liquid      string
	Platetype   string
	SurfaceTemp wunit.Temperature
}

type Output_ struct {
	Estimatedthawtime wunit.Time
	Status            string
	Thawtimeused      wunit.Time
}
