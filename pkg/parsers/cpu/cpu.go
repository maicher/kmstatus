package cpu

// Freq holds CPU frequency in kHz.
type Freq int

// Load holds the average CPU load, calculated between the previous and current stats.
// The result is combined for all cores and is expressed in percentages.
type Load float64

// Temp holds temperature values of the CPUs.
// Each value is for one CPU socket and is expressed in Celsius degrees.
type Temp []int

type CPU struct {
	Freq Freq
	Load Load
	Temp Temp
}
