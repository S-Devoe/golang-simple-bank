package util

import (
	"fmt"
	"math/rand"
)

var firstNames = []string{
	"John", "Emma", "Noah", "Olivia", "Liam", "Ava", "James", "Sophia", "William", "Isabella",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Martinez", "Hernandez",
}

// Function to generate random names
func GenerateRandomName() string {
	// Pick random first and last names
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]

	// Combine first and last names
	return fmt.Sprintf("%s %s", firstName, lastName)

}

// Function to generate random integer between min and max
func GenerateRandomInt(min, max int64) int64 {
	// Generate random integer
	return min + rand.Int63n(max-min+1)
}

// function to generate random string of length n
func GenerateRandomString(n int) string {
	// Generate random string
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomMoney() float64 {
	intValue := GenerateRandomInt(0, 1000)

	decimalPart := GenerateRandomInt(0, 99)
	return float64(intValue) + float64(decimalPart)/100
}

func GenerateRandomCurrency() string {
	currencies := []string{"USD", "NGN", "EUR", "GBP"}

	return currencies[rand.Intn(len(currencies))]
}
