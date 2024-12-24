package models

import (
	"math"
	"math/rand"
)

type Matrix struct {
	X    uint32
	P    uint32
	Data []uint32
}

func NewRandomMatrix(x, p uint32) Matrix {
	m := Matrix{
		X:    x,
		P:    p,
		Data: make([]uint32, int(math.Pow(float64(x), float64(p)))),
	}
	for i := range m.Data {
		m.Data[i] = uint32(rand.Intn(5))
	}
	return m
}

func (m *Matrix) Multiplication(lambda, mu uint32, other Matrix) Matrix {
	matrixResult := Matrix{
		X:    m.X,
		P:    (m.P - lambda - mu) + (other.P - lambda - mu) + lambda,
		Data: make([]uint32, int(math.Pow(float64(m.X), float64((m.P-lambda-mu)+(other.P-lambda-mu)+lambda)))),
	}

	indexLHS := make([]uint32, m.P)
	indexRHS := make([]uint32, other.P)
	indexResult := make([]uint32, matrixResult.P)

	for idx := 0; idx < len(matrixResult.Data); idx++ {
		tempValue := uint32(0)
		c := uint32(math.Pow(float64(m.X), float64(mu)))
		if mu > 0 {
			for i := uint32(0); i < c; i++ {
				tempValue += m.Data[calculateIndexFromArray(indexLHS, m.X)] * other.Data[calculateIndexFromArray(indexRHS, other.X)]
				incrementIndexVector(indexLHS, len(indexLHS)-1, m.X)
				incrementIndexVector(indexRHS, int(lambda+mu-1), other.X)
			}
		} else {
			tempValue += m.Data[calculateIndexFromArray(indexLHS, m.X)] * other.Data[calculateIndexFromArray(indexRHS, other.X)]
		}
		matrixResult.Data[idx] = tempValue
		incrementIndexVector(indexResult, len(indexResult)-1, matrixResult.X)
	}

	return matrixResult
}

func calculateIndexFromArray(arrayIndex []uint32, x uint32) int {
	var result uint32
	for idx := 0; idx < len(arrayIndex)-1; idx++ {
		result += arrayIndex[idx] * uint32(math.Pow(float64(x), float64(len(arrayIndex)-idx-1)))
	}
	result += arrayIndex[len(arrayIndex)-1]
	return int(result)
}

func incrementIndexVector(vec []uint32, idx int, x uint32) {
	for {
		if idx < 0 {
			return
		}
		if vec[idx] == x-1 {
			vec[idx] = 0
			idx--
		} else {
			vec[idx]++
			return
		}
	}
}
