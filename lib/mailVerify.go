package lib

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gopkg.in/mail.v2"
)

func SendEmailTo(to, key string) {
	emailBody := "Please use the following code to log in: <code>" + key + "</code>"

	m := mail.NewMessage()
	m.SetHeader("From", Config.SmtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Login Code")
	m.SetBody("text/html", emailBody)

	fmt.Println(Config.SmtpHost, Config.SmtpPort, Config.SmtpUser, Config.SmtpPassword)

	d := mail.NewDialer(Config.SmtpHost, Config.SmtpPort, Config.SmtpUser, Config.SmtpPassword)

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
	token.Expires = time.Now().Add(time.Minute * 10)

	emailTokens[email] = token

	SendEmailTo(email, token.Token)

	return token
}

type PendingAccount struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`

	emailToken      EmailToken
	deletionTimeout *time.Timer
}

var pendingAccounts map[string]PendingAccount = make(map[string]PendingAccount)

func InitializeAccountCreation(acc PendingAccount) {
	if _, ok := pendingAccounts[acc.Email]; ok {
		pendingAccounts[acc.Email].deletionTimeout.Stop()
		delete(pendingAccounts, acc.Email)
	}

	acc.emailToken = *VerifyEmail(acc.Email)
	acc.deletionTimeout = time.NewTimer(time.Minute * 10)
	pendingAccounts[acc.Email] = acc

	go func() {
		<-acc.deletionTimeout.C
		delete(pendingAccounts, acc.Email)
	}()
}

func CompleteAccountCreation(email, token string) *PendingAccount {
	if acc, ok := pendingAccounts[email]; ok {
		if acc.emailToken.Token == token && !acc.emailToken.Expired() {
			delete(pendingAccounts, email)
			acc.deletionTimeout.Stop()
			return &acc
		}
	}

	return nil
}
