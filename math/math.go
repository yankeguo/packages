// Copyright (c) The Cryptorium Authors.
// Licensed under the MIT License.

package math

import (
	"math"
	"math/big"
)

func PercentageDiff(old, new float64) float64 {
	diff := new - old

	if old == 0 {
		if new > 0 {
			return math.MaxFloat64
		}
		return -math.MaxFloat64
	}
	if new == 0 {
		if old > 0 {
			return math.MaxFloat64
		}
		return -math.MaxFloat64
	}

	if old > new {
		return -math.Abs((diff / new) * 100)
	}
	return math.Abs((diff / old) * 100)

}

func BigIntToFloat(input *big.Int) float64 {
	fl, _ := big.NewFloat(0).SetInt(input).Float64()
	return fl
}

func BigIntToFloatDiv(input *big.Int, devider float64) float64 {
	fl := BigIntToFloat(input)
	if devider == 1 {
		return fl
	}
	f := 0.0
	if input != nil {
		return fl / devider
	}
	return f
}

func FloatToBigIntMul(input float64, multiplier float64) *big.Int {
	if input == 0 {
		return big.NewInt(0)
	}
	bigE18 := big.NewFloat(0).Mul(big.NewFloat(input), big.NewFloat(multiplier))
	result, _ := bigE18.Int(nil)
	return result
}

func FloatToBigInt(input float64) *big.Int {
	result, _ := big.NewFloat(input).Int(nil)
	return result
}
