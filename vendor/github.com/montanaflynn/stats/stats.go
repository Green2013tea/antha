package stats

import (
	"errors"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Min finds the lowest number in a slice
func Min(input []float64) (min float64, err error) {

	// Get the initial value or return an error
	if len(input) > 0 {
		min = input[0]
	} else {
		return 0, errors.New("Input must not be empty")
	}

	// Iterate until done checking for a lower value
	for i := 1; i < len(input); i++ {
		if input[i] < min {
			min = input[i]
		}
	}
	return min, nil
}

// Max finds the highest number in a slice
func Max(input []float64) (max float64, err error) {

	// Get the initial value
	if len(input) > 0 {
		max = input[0]
	} else {
		return 0, errors.New("Input must not be empty")
	}

	// Loop and replace higher values
	for i := 1; i < len(input); i++ {
		if input[i] > max {
			max = input[i]
		}
	}

	return max, nil
}

// Sum adds all the numbers of a slice together
func Sum(input []float64) (sum float64, err error) {

	// Add em up
	for _, n := range input {
		sum += float64(n)
	}

	return sum, nil
}

// Mean gets the average of a slice of numbers
func Mean(input []float64) (mean float64, err error) {

	if len(input) == 0 {
		return 0, errors.New("Input must not be empty")
	}

	sum, err := Sum(input)
	if err != nil {
		return 0, errors.New("Could not calculate sum")
	}

	return sum / float64(len(input)), nil
}

// Median gets the median number in a slice of numbers
func Median(input []float64) (median float64, err error) {

	// Start by sorting a copy of the slice
	c := copyslice(input)
	sort.Float64s(c)

	// No math is needed if there are no numbers
	// For even numbers we add the two middle numbers
	// and divide by two using the mean function above
	// For odd numbers we just use the middle number
	l := len(c)
	if l == 0 {
		return 0, errors.New("Input must not be empty")
	} else if l%2 == 0 {
		median, err = Mean(c[l/2-1 : l/2+1])
		if err != nil {
			return 0, errors.New("Could not calculate median using the mean func")
		}
	} else {
		median = float64(c[l/2])
	}

	return median, nil
}

// Mode gets the mode of a slice of numbers
func Mode(input []float64) (mode []float64, err error) {

	// Return the input if there's only one number
	l := len(input)
	if l == 1 {
		return input, nil
	} else if l == 0 {
		return nil, errors.New("Input must not be empty")
	}

	// Create a map with the counts for each number
	m := make(map[float64]int)
	for _, v := range input {
		m[v]++
	}

	// Find the highest counts to return as a slice
	// of ints to accomodate duplicate counts
	var current int
	for k, v := range m {

		// Depending if the count is lower, higher
		// or equal to the current numbers count
		// we return nothing, start a new mode or
		// append to the current mode
		switch {
		case v < current:
		case v > current:
			current = v
			mode = append(mode[:0], k)
		default:
			mode = append(mode, k)
		}
	}

	// Finally we check to see if there actually was
	// a mode by checking the length of the input and
	// mode against eachother
	lm := len(mode)
	if l == lm {
		return []float64{}, nil
	}

	return mode, nil
}

// Variance finds the variance for both population and sample data
func Variance(input []float64, sample int) (variance float64, err error) {

	if len(input) == 0 {
		return 0, errors.New("Input must not be empty")
	}

	// Sum the square of the mean subtracted from each number
	m, err := Mean(input)
	if err != nil {
		return 0, errors.New("Could not calculate mean")
	}

	for _, n := range input {
		variance += (float64(n) - m) * (float64(n) - m)
	}

	// When getting the mean of the squared differences
	// "sample" will allow us to know if it's a sample
	// or population and wether to subtract by one or not
	return variance / float64((len(input) - (1 * sample))), nil
}

// VarP finds the amount of variance within a population
func VarP(input []float64) (sdev float64, err error) {

	v, err := Variance(input, 0)
	if err != nil {
		return 0, err
	}

	return v, nil
}

// VarS finds the amount of variance within a sample
func VarS(input []float64) (sdev float64, err error) {

	v, err := Variance(input, 1)
	if err != nil {
		return 0, err
	}

	return v, nil
}

// StdDevP finds the amount of variation from the population
func StdDevP(input []float64) (sdev float64, err error) {

	if len(input) == 0 {
		return 0, errors.New("Input must not be empty")
	}

	// Get the population variance
	m, err := VarP(input)
	if err != nil {
		return 0, err
	}

	// Return the population standard deviation
	return math.Pow(m, 0.5), nil
}

// StdDevS finds the amount of variation from a sample
func StdDevS(input []float64) (sdev float64, err error) {

	if len(input) == 0 {
		return 0, errors.New("Input must not be empty")
	}

	// Get the sample variance
	sv, err := VarS(input)
	if err != nil {
		return 0, err
	}

	// Return the sample standard deviation
	return math.Pow(sv, 0.5), nil
}

// Round a float to a specific decimal place or precision
func Round(input float64, places int) (rounded float64, err error) {

	// If the float is not a number
	if math.IsNaN(input) {
		return 0.0, errors.New("Not a number")
	}

	// Find out the actual sign and correct the input for later
	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}

	// Use the places arg to get the amount of precision wanted
	precision := math.Pow(10, float64(places))

	// Find the decimal place we are looking to round
	digit := input * precision

	// Get the actual decimal number as a fraction to be compared
	_, decimal := math.Modf(digit)

	// If the decimal is less than .5 we round down otherwise up
	if decimal >= 0.5 {
		rounded = math.Ceil(digit)
	} else {
		rounded = math.Floor(digit)
	}

	// Finally we do the math to actually create a rounded number
	return rounded / precision * sign, nil
}

// Percentile finds the relative standing in a slice of floats
func Percentile(input []float64, percent float64) (percentile float64, err error) {

	if len(input) == 0 {
		return 0, errors.New("Input must not be empty")
	}

	// Start by sorting a copy of the slice
	c := copyslice(input)
	sort.Float64s(c)

	// Multiple percent by length of input
	index := (percent / 100) * float64(len(c))

	// Check if the index is a whole number
	if index == float64(int64(index)) {

		// Convert float to int
		i, err := Float64ToInt(index)
		if err != nil {
			return 0, errors.New("Could not turn float64 into int")
		}

		// Find the average of the index and following values
		percentile, err = Mean([]float64{c[i-1], c[i]})
		if err != nil {
			return 0, errors.New("Could not calculate percentile with the mean func")
		}

	} else {

		// Convert float to int
		i, err := Float64ToInt(index)
		if err != nil {
			return 0, errors.New("Could not turn float64 into int")
		}

		// Find the value at the index
		percentile = c[i-1]

	}

	return percentile, nil

}

// Float64ToInt rounds a float64 to an int
func Float64ToInt(input float64) (output int, err error) {

	// Round input to nearest whole number and convert to int
	r, err := Round(input, 0)
	if err != nil {
		return 0, err
	}

	return int(r), nil

}

// Coordinate holds the data in a series
type Coordinate struct {
	X, Y float64
}

// LinReg finds the least squares linear regression on data series
func LinReg(s []Coordinate) (regressions []Coordinate, err error) {

	if len(s) == 0 {
		return nil, errors.New("Input must not be empty")
	}

	// Placeholder for the math to be done
	var sum [5]float64

	// Loop over data keeping index in place
	i := 0
	for ; i < len(s); i++ {
		sum[0] += s[i].X
		sum[1] += s[i].Y
		sum[2] += s[i].X * s[i].X
		sum[3] += s[i].X * s[i].Y
		sum[4] += s[i].Y * s[i].Y
	}

	// Find gradient and intercept
	f := float64(i)
	gradient := (f*sum[3] - sum[0]*sum[1]) / (f*sum[2] - sum[0]*sum[0])
	intercept := (sum[1] / f) - (gradient * sum[0] / f)

	// Create the new regression series
	for j := 0; j < len(s); j++ {
		regressions = append(regressions, Coordinate{
			X: s[j].X,
			Y: s[j].X*gradient + intercept,
		})
	}

	return regressions, nil

}

// ExpReg returns an exponential regression on data series
func ExpReg(s []Coordinate) (regressions []Coordinate, err error) {

	if len(s) == 0 {
		return nil, errors.New("Input must not be empty")
	}

	var sum [6]float64

	for i := 0; i < len(s); i++ {
		sum[0] += s[i].X
		sum[1] += s[i].Y
		sum[2] += s[i].X * s[i].X * s[i].Y
		sum[3] += s[i].Y * math.Log(s[i].Y)
		sum[4] += s[i].X * s[i].Y * math.Log(s[i].Y)
		sum[5] += s[i].X * s[i].Y
	}

	denominator := (sum[1]*sum[2] - sum[5]*sum[5])
	a := math.Pow(math.E, (sum[2]*sum[3]-sum[5]*sum[4])/denominator)
	b := (sum[1]*sum[4] - sum[5]*sum[3]) / denominator

	for j := 0; j < len(s); j++ {
		regressions = append(regressions, Coordinate{
			X: s[j].X,
			Y: a * math.Pow(2.718281828459045, b*s[j].X),
		})
	}

	return regressions, nil

}

// LogReg returns an logarithmic regression on data series
func LogReg(s []Coordinate) (regressions []Coordinate, err error) {

	if len(s) == 0 {
		return nil, errors.New("Input must not be empty")
	}

	var sum [4]float64

	i := 0
	for ; i < len(s); i++ {
		sum[0] += math.Log(s[i].X)
		sum[1] += s[i].Y * math.Log(s[i].X)
		sum[2] += s[i].Y
		sum[3] += math.Pow(math.Log(s[i].X), 2)
	}

	f := float64(i)
	a := (f*sum[1] - sum[2]*sum[0]) / (f*sum[3] - sum[0]*sum[0])
	b := (sum[2] - a*sum[0]) / f

	for j := 0; j < len(s); j++ {
		regressions = append(regressions, Coordinate{
			X: s[j].X,
			Y: b + a*math.Log(s[j].X),
		})
	}

	return regressions, nil

}

// Sample returns sample from input with replacement or without
func Sample(input []float64, takenum int, replacement bool) ([]float64, error) {

	if len(input) == 0 {
		return nil, errors.New("Input must not be empty")
	}

	length := len(input)
	if replacement {

		result := []float64{}
		rand.Seed(unixnano())

		// In every step, randomly take the num for
		for i := 0; i < takenum; i++ {
			idx := rand.Intn(length)
			result = append(result, input[idx])
		}

		return result, nil

	} else if !replacement && takenum <= length {

		rand.Seed(unixnano())

		// Get permutation of number of indexies
		perm := rand.Perm(length)
		result := []float64{}

		// Get element of input by permutated index
		for _, idx := range perm[0:takenum] {
			result = append(result, input[idx])
		}

		return result, nil

	}

	return nil, errors.New("Number of taken elements must be less than length of input")
}

// unixnano returns nanoseconds from UTC epoch
func unixnano() int64 {
	return time.Now().UTC().UnixNano()
}

// copyslice copies a slice of float64s
func copyslice(input []float64) []float64 {
	s := make([]float64, len(input))
	copy(s, input)
	return s
}
