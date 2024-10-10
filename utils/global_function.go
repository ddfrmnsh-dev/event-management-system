package utils

import "regexp"

// type GlobalFunction interface {
// 	IsEmail(input string) bool
// }

func IsEmail(input string) bool {
	// Regex sederhana untuk memeriksa format email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(input)
}
