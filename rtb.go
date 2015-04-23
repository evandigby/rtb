// Package rtb defines the types for implementing a real time bidding application.
package rtb

import (
	"math"
)

// Are micro cents small enough to remove rounding error?
const MicroCentsConversionFactorInt64 = int64(1000000)
const MicroCentsConversionFactorFloat = float64(MicroCentsConversionFactorInt64)
const CPMConversionFactor = int64(1000)

// CpmToMicroCents converts from CPM (measured in dollars) to micro cents
func CpmToMicroCents(cpm float64) int64 {
	// 100 from dollars to cents, 1000 to microcents
	return DollarsToMicroCents(cpm)
}

// DollarsToMicroCents converts from dollars to micro cents
func DollarsToMicroCents(dollars float64) int64 {
	// 100 from dollars to cents, 1000 to microcents
	return int64(dollars * 100 * MicroCentsConversionFactorFloat)
}

// MicroCentsToCpm converts from micro cents to CPM (measured in dollars)
func MicroCentsToCpm(microCents int64) float64 {
	// 1000 from microcents to cents, 100 to dollars
	return MicroCentsToDollars(microCents)
}

// MicroCentsToDollars converts from micro cents to dollars
func MicroCentsToDollars(microCents int64) float64 {
	// 1000 from microcents to cents, 100 to dollars
	return (float64(microCents) / MicroCentsConversionFactorFloat) / 100.0
}

// MicroCentsPerImpression calculates the number of micro cents per impression based on a CPM value
func MicroCentsPerImpression(cpmInMicroCents int64) int64 {
	return cpmInMicroCents / CPMConversionFactor

}

func MicroCentsToDollarsRounded(microCents int64, decimals int) float64 {
	dollars := MicroCentsToDollars(microCents)
	factor := math.Pow(10, float64(decimals))
	return math.Floor((dollars*factor)+0.5) / factor
}
