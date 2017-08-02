package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
)

// Email ...
type Email struct {
	Client *smtp.Client
}

// Auth ... email authentication details
type Auth struct {
	Host string
	Port string
	User string
	Pass string
}

// Message ... email details
type Message struct {
	From    string
	To      string
	Subject string
	Body    string
}

// Init ... sets up the email client
func (e *Email) Init(auth Auth) error {

	if auth.Host == "" {
		return errors.New("Host missing")
	}
	if auth.User == "" {
		return errors.New("User missing")
	}
	if auth.Pass == "" {
		return errors.New("Pass missing")
	}

	emailServer := fmt.Sprintf("%s:%s", auth.Host, auth.Port)

	var client *smtp.Client
	var err error

	if auth.Port == "465" {

		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         auth.Host,
		}
		conn, err := tls.Dial("tcp", emailServer, tlsconfig)
		if err != nil {
			fmt.Printf("\t1 %v\n", err)
			return err
		}

		client, err = smtp.NewClient(conn, auth.Host)
		if err != nil {
			fmt.Printf("\t2 %v\n", err)
			return err
		}

		a := smtp.PlainAuth("", auth.User, auth.Pass, auth.Host)

		if err = client.Auth(a); err != nil {
			fmt.Printf("\t3 %v\n", err)
			return err
		}

	} else {
		client, err = smtp.Dial(emailServer)
		if err != nil {
			fmt.Printf("\t%4 %v\n", err)
			return err
		}
	}

	e.Client = client

	return nil
}

// Send ... sends email
func (e *Email) Send(m Message) error {

	from := mail.Address{"", m.From}
	to := mail.Address{"", m.To}
	subject := m.Subject
	body := m.Body

	// email headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	var err error

	if err = e.Client.Mail(from.Address); err != nil {
		return err
	}

	if err = e.Client.Rcpt(to.Address); err != nil {
		return err
	}

	w, err := e.Client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}