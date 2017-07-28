package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var emailHost string = os.Getenv("SMTP_host")
var emailPort string = os.Getenv("SMTP_port")
var emailPass string = os.Getenv("SMTP_pass")
var emailFrom string = "noreply@guitar-lessons-stalbans.co.uk"

var emailRecipient string = "xxxxx@gmail.com"

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type IntroData struct {
	Message string
	Error   bool
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, nil)
}

func main() {

	if emailPass == "" {
		log.Fatal("ERROR: email authentication (password) must be supplied for SMTP BTInternet server and xxxxx@btinternet.com")
	}
	if emailHost == "" {
		emailHost = "mail.btinternet.com"
	}
	if emailPort == "" {
		emailPort = "465"
	}
	emailServer := strings.Join([]string{emailHost, emailPort}, ":")

	// static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// show the main page
	http.Handle("/", &templateHandler{filename: "intro.html"})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("in contact handler")

		t := templateHandler{filename: "intro.html"}

		t.once.Do(func() {
			t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
		})
		err := r.ParseForm()

		if err != nil {
			log.Println(err)

			successMsg := IntroData{Message: "Form was not successfully submitted. Please try again.", Error: true}
			jsonData, _ := json.Marshal(successMsg)
			w.Write(jsonData)

			return
		}

		email := r.PostFormValue("email")

		if email == "" {
			http.Redirect(w, r, "/", 301)
			return
		}

		firstName := r.PostFormValue("first_name")
		lastName := r.PostFormValue("last_name")
		note := r.PostFormValue("note")

		fmt.Printf("\t\tForm data : firstName=%s, lastName=%s,email=%s,note=%s\n", firstName, lastName, email, note)

		subject := fmt.Sprintf("Guitar lessons interest (%s %s - %s)", firstName, lastName, email)

		message := "First name: %s \r\nLast name: %s \r\n\r\nEmail: %s \r\n\r\nNote: %s\r\n"
		emailBody := fmt.Sprintf(message, firstName, lastName, email, note)

		from := mail.Address{"", emailFrom}
		to1 := mail.Address{"", emailRecipient}
		emailSubject := subject

		// email headers
		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = to1.String()
		headers["Subject"] = emailSubject

		// Setup message
		emailMessage := ""
		for k, v := range headers {
			emailMessage += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		emailMessage += "\r\n" + emailBody

		// send an email about the contact
		auth := smtp.PlainAuth("", "xxxxx@btinternet.com", emailPass, "mail.btinternet.com")

		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         emailHost,
		}

		// need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls) (https://gist.github.com/chrisgillis/10888032)
		conn, err := tls.Dial("tcp", emailServer, tlsconfig)
		if err != nil {
			log.Panic(err)
		}

		c, err := smtp.NewClient(conn, emailHost)
		if err != nil {
			log.Panic(err)
		}

		// Auth
		if err = c.Auth(auth); err != nil {
			log.Panic(err)
		}

		// To && From
		if err = c.Mail(from.Address); err != nil {
			log.Panic(err)
		}

		if err = c.Rcpt(to1.Address); err != nil {
			log.Panic(err)
		}

		// Data
		ew, err := c.Data()
		if err != nil {
			log.Panic(err)
		}

		_, err = ew.Write([]byte(emailMessage))
		if err != nil {
			log.Panic(err)
		}

		err = ew.Close()
		if err != nil {
			log.Panic(err)
		}

		c.Quit()

		// TODO: store in database

		// show the successful message
		w.Header().Set("Content-Type", "application/json")

		successMsg := IntroData{Message: "Thank you. Form was successfully submitted. I shall contact you shortly.", Error: false}
		jsonData, err := json.Marshal(successMsg)

		if err != nil {
			fmt.Printf("contact handler error: %v \n\n", err)
		}
		w.Write(jsonData)

	})

	// ------- WEB SERVER -------
	fmt.Println("Starting web server ...")
	if err := http.ListenAndServe("192.168.1.91:80", nil); err != nil {
		log.Fatal("Server problems: ", err)
	}
}
