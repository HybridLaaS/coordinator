package lib

import (
	"log"
	"math/rand"
	"time"

	"gopkg.in/mail.v2"
)

func SendEmailTo(to, key string) {
	emailBody := "Please use the following code to log in: <code>" + key + "</code>"

	m := mail.NewMessage()
	m.SetHeader("From", "w1hlo.hamchat@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Login Code")
	m.SetBody("text/html", emailBody)

	d := mail.NewDialer("smtp.gmail.com", 587, "w1hlo.hamchat@gmail.com", "ckkm gaws caoe ries")

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Could not send email: %v", err)
	}
}

var emailTokens map[string]*EmailToken = make(map[string]*EmailToken)

type EmailToken struct {
	Email   string
	Token   string
	Expires time.Time
}

func (e *EmailToken) Expired() bool {
	return e.Expires.Before(time.Now())
}

func (e *EmailToken) Generate() string {
	var output string = ""

	for i := 0; i < 16; i++ {
		var char byte = byte(rand.Intn(26) + 97)
		output += string(char)

		if i%4 == 3 && i != 15 {
			output += "-"
		}
	}

	return output
}

func VerifyEmail(email string) *EmailToken {
	var token *EmailToken = new(EmailToken)
	token.Email = email
	token.Token = token.Generate()
	token.Expires = time.Now().Add(time.Minute * 3)

	emailTokens[email] = token

	SendEmailTo(email, token.Token)

	return token
}