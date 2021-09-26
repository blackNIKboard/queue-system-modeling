package poisson

import "math"

//qPrForAll = [(a ** i) / (math.factorial(i)) * (math.exp(-1 * a)) for i in range(20)]

func GeneratePoissonProcess(a float64) []float64 {
	pp := []float64{}

	for i := 0; i < 20; i++ {
		pp = append(pp, (math.Pow(a, float64(i)))/(Factorial(float64(i)))*(math.Exp(-1*a)))
	}

	return pp
}

func Factorial(n float64) (result float64) {
	if n > 0 {
		result = n * Factorial(n-1)

		return result
	}

	return 1
}
