package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func sendMail(msg string, from string, password string, toList []string) {

	host := "smtp.gmail.com"
	port := "587"

	// We can't send strings directly in mail,
	// strings need to be converted into slice bytes
	body := []byte(msg)

	// PlainAuth uses the given username and password to
	// authenticate to host and act as identity.
	// Usually identity should be the empty string,
	// to act as username.
	auth := smtp.PlainAuth("", from, password, host)

	// SendMail uses TLS connection to send the mail
	// The email is sent to all address in the toList,
	// the body should be of type bytes, not strings
	// This returns error if any occurred.
	err := smtp.SendMail(host+":"+port, auth, from, toList, body)

	// handling the errors
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
