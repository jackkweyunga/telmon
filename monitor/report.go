package monitor

import (
	"fmt"
	log "github.com/jackkweyunga/telmon/logging"
	"net/smtp"
	"os"
)

type Email struct {
	Password  string   `mapstructure:"password"`
	Sender    string   `mapstructure:"sender"`
	Receivers []string `mapstructure:"receivers"`
}

type Notify struct {
	Email `mapstructure:"email"`
}

// Sends notification as per configurations
// when triggered.
// @TODO: support more notification channels
func (n *Notify) report(msg string) {
	if n.Email.Sender != "" {
		//n.Email.sendMail(msg)
		return
	} else {
		log.Log.Fatalln("[TelnetClient] No notification configurations found.")
	}
}

func (e *Email) sendMail(msg string) {

	host := "smtp.gmail.com"
	port := "587"

	// We can't send strings directly in mail,
	// strings need to be converted into slice bytes
	body := []byte(msg)

	// PlainAuth uses the given username and password to
	// authenticate to host and act as identity.
	// Usually identity should be the empty string,
	// to act as username.
	auth := smtp.PlainAuth("", e.Sender, e.Password, host)

	// SendMail uses TLS connection to send the mail
	// The email is sent to all address in the toList,
	// the body should be of type bytes, not strings
	// This returns error if any occurred.
	err := smtp.SendMail(host+":"+port, auth, e.Sender, e.Receivers, body)

	// handling the errors
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Log.Println("[Telmon] Report Email sent successfully ... OK")

}
