package utils

import (
	"math/rand"
)

var alphabets = "abcdefghijklmnopqrstuvwxyz"
var numbers = "123456789"

func RandomString(r int) string {
	wholeLetter := []rune{}
	k := len(alphabets)

	for i := 0; i < r; i++ {
		index := rand.Intn(k)
		wholeLetter = append(wholeLetter, rune(alphabets[index]))
	}
	return string(wholeLetter)
}

func RandIntegers(r int) string {
	wholeFigure := []rune{}
	k := len(numbers)

	for i := 0; i < r; i++ {
		index := rand.Intn(k)
		wholeFigure = append(wholeFigure, rune(numbers[index]))
	}
	return string(wholeFigure)
}

func randomInteger(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

// func randomFloat(min, max float64) float64 {
// 	return min + rand.Float64()*(max)
// }

////////////////////////////////////////////////////////////////////////////////

func RandomEmail() string {
	return RandomString(7) + "@testing.com"
}

func RandomPhone() string {
	return RandIntegers(11)
}

func RandomName() string {
	return RandomString(5)
}

func RandomAddress() string {
	return RandomString(30)
}

func RandomText() string {
	return RandomString(100)
}

func RandomPrice() string {
	return RandIntegers(5) + "." + RandIntegers(2)
}

func RandomQty() int32 {
	return randomInteger(1, 2000)
}
