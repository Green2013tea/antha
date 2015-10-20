package stats

import (
	"math"
	"reflect"
	"sort"
	"testing"
)

func TestMin(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out float64
	}{
		{[]float64{1.1, 2, 3, 4, 5}, 1.1},
		{[]float64{10.534, 3, 5, 7, 9}, 3.0},
		{[]float64{-5, 1, 5}, -5.0},
		{[]float64{5}, 5},
	} {
		got, err := Min(c.in)
		if err != nil {
			t.Errorf("Returned an error")
		}
		if got != c.out {
			t.Errorf("Min(%.1f) => %.1f != %.1f", c.in, c.out, got)
		}
	}
}

func TestMax(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 5.0},
		{[]float64{10.5, 3, 5, 7, 9}, 10.5},
		{[]float64{-20, -1, -5.5}, -1.0},
		{[]float64{-1.0}, -1.0},
	} {
		got, err := Max(c.in)
		if err != nil {
			t.Errorf("Returned an error")
		}
		if got != c.out {
			t.Errorf("Max(%.1f) => %.1f != %.1f", c.in, c.out, got)
		}
	}
}

func TestMean(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 3.0},
		{[]float64{1, 2, 3, 4, 5, 6}, 3.5},
		{[]float64{1}, 1.0},
	} {
		got, _ := Mean(c.in)
		if got != c.out {
			t.Errorf("Mean(%.1f) => %.1f != %.1f", c.in, c.out, got)
		}
	}
	_, err := Mean([]float64{})
	if err == nil {
		t.Errorf("Should have returned an error")
	}
}

func TestMedian(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out float64
	}{
		{[]float64{5, 3, 4, 2, 1}, 3.0},
		{[]float64{6, 3, 2, 4, 5, 1}, 3.5},
		{[]float64{1}, 1.0},
	} {
		got, _ := Median(c.in)
		if got != c.out {
			t.Errorf("Median(%.1f) => %.1f != %.1f", c.in, c.out, got)
		}
	}
	_, err := Median([]float64{})
	if err == nil {
		t.Errorf("Should have returned an error")
	}
}

func TestMedianSortSideEffects(t *testing.T) {
	s := []float64{0.1, 0.3, 0.2, 0.4, 0.5}
	a := []float64{0.1, 0.3, 0.2, 0.4, 0.5}
	Median(s)
	if !reflect.DeepEqual(s, a) {
		t.Errorf("%.1f != %.1f", s, a)
	}
}

func TestMode(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out []float64
	}{
		{[]float64{5, 3, 4, 2, 1}, []float64{}},
		{[]float64{5, 5, 3, 4, 2, 1}, []float64{5}},
		{[]float64{5, 5, 3, 3, 4, 2, 1}, []float64{3, 5}},
		{[]float64{1}, []float64{1}},
	} {
		got, err := Mode(c.in)
		if err != nil {
			t.Errorf("Returned an error")
		}
		sort.Float64s(got)
		if !reflect.DeepEqual(c.out, got) {
			t.Errorf("Mode(%.1f) => %.1f != %.1f", c.in, got, c.out)
		}
	}
}

func TestSum(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out float64
	}{
		{[]float64{1, 2, 3}, 6},
		{[]float64{1.0, 1.1, 1.2, 2.2}, 5.5},
		{[]float64{1, -1, 2, -3}, -1},
		{[]float64{}, 0},
	} {
		got, err := Sum(c.in)
		if err != nil {
			t.Errorf("Returned an error")
		}
		if !reflect.DeepEqual(c.out, got) {
			t.Errorf("Sum(%.1f) => %.1f != %.1f", c.in, got, c.out)
		}
	}
}

func TestVariance(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		pop int
		out float64
	}{
		{[]float64{}, 0, 0.0},
		{[]float64{}, 1, 0.0},
		{[]float64{1, 2, 3}, 0, 0.7},
		{[]float64{1, 2, 3}, 1, 1.0},
	} {
		v, _ := Variance(c.in, c.pop)
		got, err := Round(v, 1)
		if err != nil {
			t.Errorf("Returned an error")
		}
		if got != c.out {
			t.Errorf("Variance(%.1f) => %.1f != %.1f", c.in, c.out, got)
		}
	}
}

func TestVarP(t *testing.T) {
	m, _ := VarP([]float64{})
	if m != 0.0 {
		t.Errorf("%.1f != %.1f", m, 0.0)
	}
	m, _ = VarP([]float64{1, 2, 3})
	m, err := Round(m, 1)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 0.7 {
		t.Errorf("%.1f != %.1f", m, 0.7)
	}
}

func TestVarS(t *testing.T) {
	m, _ := VarS([]float64{})
	if m != 0.0 {
		t.Errorf("%.1f != %.1f", m, 0.0)
	}
	m, _ = VarS([]float64{1, 2, 3})
	if m != 1.0 {
		t.Errorf("%.1f != %.1f", m, 1.0)
	}
}

func TestStdDevP(t *testing.T) {
	s, _ := StdDevP([]float64{1, 2, 3})
	m, err := Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 0.82 {
		t.Errorf("%.10f != %.10f", m, 0.82)
	}
	s, _ = StdDevP([]float64{-1, -2, -3.3})
	m, err = Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 0.94 {
		t.Errorf("%.10f != %.10f", m, 0.94)
	}

	m, _ = StdDevP([]float64{})
	if m != 0.0 {
		t.Errorf("%.1f != %.1f", m, 0.0)
	}
}

func TestStdDevS(t *testing.T) {
	s, _ := StdDevS([]float64{1, 2, 3})
	m, err := Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 1.0 {
		t.Errorf("%.10f != %.10f", m, 1.0)
	}
	s, _ = StdDevS([]float64{-1, -2, -3.3})
	m, err = Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 1.15 {
		t.Errorf("%.10f != %.10f", m, 1.15)
	}

	m, _ = StdDevS([]float64{})
	if m != 0.0 {
		t.Errorf("%.1f != %.1f", m, 0.0)
	}
}

func TestRound(t *testing.T) {
	m, err := Round(0.1111, 1)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 0.1 {
		t.Errorf("%.1f != %.1f", m, 0.1)
	}

	m, err = Round(-0.1111, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != -0.11 {
		t.Errorf("%.1f != %.1f", m, -0.11)
	}

	m, err = Round(5.3253, 3)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 5.325 {
		t.Errorf("%.1f != %.1f", m, 5.325)
	}

	m, err = Round(5.3253, 0)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 5.0 {
		t.Errorf("%.1f != %.1f", m, 5.0)
	}

	m, err = Round(math.NaN(), 2)
	if err == nil {
		t.Errorf("Round should error on NaN")
	}
}

func TestPercentile(t *testing.T) {
	m, _ := Percentile([]float64{43, 54, 56, 61, 62, 66}, 90)
	if m != 62.0 {
		t.Errorf("%.1f != %.1f", m, 62.0)
	}
	m, _ = Percentile([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 50)
	if m != 5.5 {
		t.Errorf("%.1f != %.1f", m, 5.5)
	}
	m, _ = Percentile([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 99.9)
	if m != 10.0 {
		t.Errorf("%.1f != %.1f", m, 10.0)
	}
}

func TestPercentileSortSideEffects(t *testing.T) {
	s := []float64{43, 54, 56, 44, 62, 66}
	a := []float64{43, 54, 56, 44, 62, 66}
	Percentile(s, 90)
	if !reflect.DeepEqual(s, a) {
		t.Errorf("%.1f != %.1f", s, a)
	}
}

func TestFloat64ToInt(t *testing.T) {
	m, _ := Float64ToInt(234.0234)
	if m != 234 {
		t.Errorf("%x != %x", m, 234)
	}
	m, _ = Float64ToInt(-234.0234)
	if m != -234 {
		t.Errorf("%x != %x", m, -234)
	}
	m, _ = Float64ToInt(1)
	if m != 1 {
		t.Errorf("%x != %x", m, 1)
	}
}

func TestLinReg(t *testing.T) {
	data := []Coordinate{
		{1, 2.3},
		{2, 3.3},
		{3, 3.7},
		{4, 4.3},
		{5, 5.3},
	}

	r, _ := LinReg(data)
	a := 2.3800000000000026
	if r[0].Y != a {
		t.Errorf("%v != %v", r, a)
	}
	a = 3.0800000000000014
	if r[1].Y != a {
		t.Errorf("%v != %v", r, a)
	}
	a = 3.7800000000000002
	if r[2].Y != a {
		t.Errorf("%v != %v", r, a)
	}
	a = 4.479999999999999
	if r[3].Y != a {
		t.Errorf("%v != %v", r, a)
	}
	a = 5.179999999999998
	if r[4].Y != a {
		t.Errorf("%v != %v", r, a)
	}
}

func TestExpReg(t *testing.T) {
	data := []Coordinate{
		{1, 2.3},
		{2, 3.3},
		{3, 3.7},
		{4, 4.3},
		{5, 5.3},
	}

	r, _ := ExpReg(data)
	a, _ := Round(r[0].Y, 3)
	if a != 2.515 {
		t.Errorf("%v != %v", r, 2.515)
	}
	a, _ = Round(r[1].Y, 3)
	if a != 3.032 {
		t.Errorf("%v != %v", r, 3.032)
	}
	a, _ = Round(r[2].Y, 3)
	if a != 3.655 {
		t.Errorf("%v != %v", r, 3.655)
	}
	a, _ = Round(r[3].Y, 3)
	if a != 4.407 {
		t.Errorf("%v != %v", r, 4.407)
	}
	a, _ = Round(r[4].Y, 3)
	if a != 5.313 {
		t.Errorf("%v != %v", r, 5.313)
	}
}

func TestLogReg(t *testing.T) {
	data := []Coordinate{
		{1, 2.3},
		{2, 3.3},
		{3, 3.7},
		{4, 4.3},
		{5, 5.3},
	}

	r, _ := LogReg(data)
	a := 2.1520822363811702
	if r[0].Y != a {
		t.Errorf("%v != %v", r, a)
	}
	a = 3.3305559222492214
	if r[1].Y != a {
		t.Errorf("%v != %v", r, a)
	}
	a = 4.019918836568674
	if r[2].Y != a {
		t.Errorf("%v != %v", r, a)
	}
	a = 4.509029608117273
	if r[3].Y != a {
		t.Errorf("%v != %v", r, a)
	}
	a = 4.888413396683663
	if r[4].Y != a {
		t.Errorf("%v != %v", r, a)
	}
}

func TestSample(t *testing.T) {
	_, err := Sample([]float64{}, 10, false)
	if err == nil {
		t.Errorf("Returned an error")
	}

	_, err2 := Sample([]float64{0.1, 0.2}, 10, false)
	if err2 == nil {
		t.Errorf("Returned an error")
	}
}

func TestSampleWithoutReplacement(t *testing.T) {
	arr := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	result, _ := Sample(arr, 5, false)
	checks := map[float64]bool{}
	for _, res := range result {
		_, ok := checks[res]
		if ok {
			t.Errorf("%v already seen", res)
		}
		checks[res] = true
	}
}

func TestSampleWithReplacement(t *testing.T) {
	arr := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	numsamples := 100
	result, _ := Sample(arr, numsamples, true)
	if len(result) != numsamples {
		t.Errorf("%v != %v", len(result), numsamples)
	}
}
