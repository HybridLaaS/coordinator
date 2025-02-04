package lib

import "regexp"

// Leaving the + in the first thing opens up to the possibility of mail bombing
// sample+1@gmail.com and sample+2@gmail.com point to the same mailbox, but are different emails
var emailRegex *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,4}$`)
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
