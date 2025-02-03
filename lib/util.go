package lib

import "regexp"

var emailRegex *regexp.Regexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
var nameRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9_ ']+$`)

func IsEmailValid(email string) bool {
	return emailRegex.MatchString(email)
}

func IsPasswordValid(password string) bool {
	return len(password) >= 8 && len(password) <= 64
}

func IsNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 48 && nameRegex.MatchString(name)
}

func CleanFromSQLInjection(s string) string {
	return s
}