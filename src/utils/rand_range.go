package utils

import (
	"math/rand"
	"time"
)

/*
	Rand range funcs allow to get random elements in ranges
	The implementation is taken from here: https://stackoverflow.com/a/49747128/10352443

	Innokentiy Sokolov
	https://github.com/keshon

	2022-03-24
*/

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Array of random runes
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// --

// Array of random numbers in range: float64
func RandFloats64(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

// Array of random numbers in range: float32
func RandFloats32(min, max float32, n int) []float32 {
	res := make([]float32, n)
	for i := range res {
		res[i] = min + rand.Float32()*(max-min)
	}
	return res
}

// Array of random numbers in range: int64
func RandInts64(min, max int64, n int) []int64 {
	res := make([]int64, n)
	for i := range res {
		res[i] = min + rand.Int63()*(max-min)
	}
	return res
}

// Array of random numbers in range: int32
func RandInts32(min, max int32, n int) []int32 {
	res := make([]int32, n)
	for i := range res {
		res[i] = min + rand.Int31()*(max-min)
	}
	return res
}

// --

// Single random number in range: float32
func RandFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

// Single random number in range: float64
func RandFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Single random number in range: int
func RandInt(min, max int) int {
	return min + rand.Int()*(max-min)
}

// Single random number in range: int32
func RandInt32(min, max int32) int32 {
	return min + rand.Int31()*(max-min)
}

// Single random number in range: int64
func RandInt64(min, max int64) int64 {
	return min + rand.Int63()*(max-min)
}
